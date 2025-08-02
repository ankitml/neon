package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handlers struct {
	db *pgxpool.Pool
}

type Quote struct {
	ID                int      `json:"id"`
	Quote             string   `json:"quote"`
	Author            string   `json:"author"`
	Category          string   `json:"category"`
	Tags              []string `json:"tags"`
	Relevance         float64  `json:"relevance"`
	HighlightedQuote  string   `json:"highlighted_quote,omitempty"`
}

type SearchResponse struct {
	Results []Quote `json:"results"`
	Count   int     `json:"count"`
	Query   string  `json:"query"`
}

type HealthResponse struct {
	Status string `json:"status"`
	Message string `json:"message"`
}

func (h *Handlers) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	response := HealthResponse{
		Status:  "ok",
		Message: "Quotes API is running",
	}
	
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) SearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Get search query from URL parameter
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		http.Error(w, `{"error": "q parameter is required"}`, http.StatusBadRequest)
		return
	}

	// Execute database query
	rows, err := h.BuildQuery(query)
	if err != nil {
		http.Error(w, `{"error": "Database query failed"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Build and send response
	response, err := h.BuildResponse(rows, query)
	if err != nil {
		http.Error(w, `{"error": "Failed to parse results"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) BuildQuery(query string) (pgx.Rows, error) {
	sql := `
		SELECT id, quote, author, category, tags, 
		       paradedb.score(id) as relevance,
		       paradedb.snippet(quote) as highlighted_quote
		FROM quotes 
		WHERE quote @@@ $1 OR author @@@ $1 OR tags @@@ $1 OR category @@@ $1
		ORDER BY paradedb.score(id) DESC
		LIMIT 20
	`
	
	return h.db.Query(context.Background(), sql, query)
}

func (h *Handlers) BuildResponse(rows pgx.Rows, query string) (SearchResponse, error) {
	var quotes []Quote
	
	for rows.Next() {
		var q Quote
		err := rows.Scan(&q.ID, &q.Quote, &q.Author, &q.Category, &q.Tags, &q.Relevance, &q.HighlightedQuote)
		if err != nil {
			return SearchResponse{}, err
		}
		quotes = append(quotes, q)
	}

	if err := rows.Err(); err != nil {
		return SearchResponse{}, err
	}

	return SearchResponse{
		Results: quotes,
		Count:   len(quotes),
		Query:   query,
	}, nil
}