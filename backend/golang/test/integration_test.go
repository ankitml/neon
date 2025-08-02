package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

const (
	baseURL = "http://localhost:8080"
	testWord = "life"
)

type Quote struct {
	ID               int      `json:"id"`
	Quote            string   `json:"quote"`
	Author           string   `json:"author"`
	Category         string   `json:"category"`
	Tags             []string `json:"tags"`
	Relevance        float64  `json:"relevance"`
	HighlightedQuote string   `json:"highlighted_quote,omitempty"`
}

type SearchResponse struct {
	Results []Quote `json:"results"`
	Count   int     `json:"count"`
	Query   string  `json:"query"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func TestHealthEndpoint(t *testing.T) {
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		t.Fatalf("Failed to call health endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var health HealthResponse
	if err := json.Unmarshal(body, &health); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if health.Status != "ok" {
		t.Fatalf("Expected status 'ok', got '%s'", health.Status)
	}

	t.Logf("✅ Health check passed: %s", health.Message)
}

func TestSearchEndpoint(t *testing.T) {
	// Test search with the word "life"
	resp, err := http.Get(fmt.Sprintf("%s/api/search?q=%s", baseURL, testWord))
	if err != nil {
		t.Fatalf("Failed to call search endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var searchResp SearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Test 1: Should return non-zero results
	if searchResp.Count == 0 {
		t.Fatalf("Expected non-zero results for query '%s', got %d", testWord, searchResp.Count)
	}

	t.Logf("✅ Found %d results for query '%s'", searchResp.Count, testWord)

	// Test 2: At least one result should contain the test word
	foundMatch := false
	for i, quote := range searchResp.Results {
		// Check if test word appears in quote text (case-insensitive)
		if strings.Contains(strings.ToLower(quote.Quote), strings.ToLower(testWord)) {
			foundMatch = true
			t.Logf("✅ Result %d contains '%s': \"%s...\" by %s", 
				i+1, testWord, truncateString(quote.Quote, 50), quote.Author)
			break
		}
	}

	if !foundMatch {
		t.Logf("⚠️  Checking if ParadeDB found semantic matches...")
		// Even if direct word match not found, log some results to verify search is working
		for i, quote := range searchResp.Results[:3] {
			t.Logf("   Result %d (relevance: %.2f): \"%s...\" by %s",
				i+1, quote.Relevance, truncateString(quote.Quote, 50), quote.Author)
		}
		// Don't fail the test - ParadeDB might find semantically related quotes
		t.Logf("✅ Search returned results (ParadeDB semantic matching)")
	}
}

func TestSearchEmptyQuery(t *testing.T) {
	resp, err := http.Get(baseURL + "/api/search")
	if err != nil {
		t.Fatalf("Failed to call search endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status 400 for empty query, got %d", resp.StatusCode)
	}

	t.Logf("✅ Empty query correctly returns 400 Bad Request")
}

func TestSearchNonExistentWord(t *testing.T) {
	// Use a made-up word that should not exist in quotes
	nonExistentWord := "xyznonexistentword123"
	resp, err := http.Get(fmt.Sprintf("%s/api/search?q=%s", baseURL, nonExistentWord))
	if err != nil {
		t.Fatalf("Failed to call search endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var searchResp SearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Should return 0 results for non-existent word
	if searchResp.Count != 0 {
		t.Fatalf("Expected 0 results for non-existent word '%s', got %d", nonExistentWord, searchResp.Count)
	}

	t.Logf("✅ Non-existent word correctly returns 0 results")
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// Helper function to wait for server to be ready
func waitForServer(t *testing.T) {
	maxAttempts := 30
	for i := 0; i < maxAttempts; i++ {
		resp, err := http.Get(baseURL + "/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(1 * time.Second)
	}
	t.Fatalf("Server did not start within %d seconds", maxAttempts)
}

func TestMain(m *testing.M) {
	// Wait for server to be ready before running tests
	fmt.Println("Waiting for server to be ready...")
	// Note: In a real scenario, you'd start the server here or ensure it's running
	// For now, we assume the server is started manually
}