package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"quotes-api/queries"
)

type Handlers struct {
	db            *pgxpool.Pool
	searchQueries *queries.SearchQueries
	browseQueries *queries.BrowseQueries
}

func NewHandlers(db *pgxpool.Pool) *Handlers {
	return &Handlers{
		db:            db,
		searchQueries: queries.NewSearchQueries(db),
		browseQueries: queries.NewBrowseQueries(db),
	}
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
	
	// Get search query from URL parameter (can be empty for browse mode)
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	
	// Parse all browse parameters (including filters)
	params, err := h.parseBrowseParams(r)
	if err != nil {
		log.Printf("Invalid search parameters: %v", err)
		http.Error(w, `{"error": "Invalid parameters"}`, http.StatusBadRequest)
		return
	}

	if query == "" {
		// Browse mode: no search query, use browse logic
		rows, err := h.browseQueries.BuildStatement(params)
		if err != nil {
			log.Printf("Browse query failed: %v", err)
			http.Error(w, `{"error": "Database query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Build and send browse response
		response, err := h.browseQueries.BuildResponse(rows, params)
		if err != nil {
			log.Printf("Browse response failed: %v", err)
			http.Error(w, `{"error": "Failed to parse results"}`, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(response)
	} else {
		// Search mode: use search with filters
		rows, err := h.searchQueries.BuildStatementWithFilters(query, params)
		if err != nil {
			log.Printf("Search with filters query failed: %v", err)
			http.Error(w, `{"error": "Database query failed"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Build and send search response with filters
		response, err := h.searchQueries.BuildResponseWithFilters(rows, query, params)
		if err != nil {
			log.Printf("Search with filters response failed: %v", err)
			http.Error(w, `{"error": "Failed to parse results"}`, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(response)
	}
}

func (h *Handlers) BrowseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Parse query parameters into BrowseParams
	params, err := h.parseBrowseParams(r)
	if err != nil {
		log.Printf("Invalid browse parameters: %v", err)
		http.Error(w, `{"error": "Invalid parameters"}`, http.StatusBadRequest)
		return
	}

	// Execute database query
	rows, err := h.browseQueries.BuildStatement(params)
	if err != nil {
		log.Printf("Browse query failed: %v", err)
		http.Error(w, `{"error": "Database query failed"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Build and send response
	response, err := h.browseQueries.BuildResponse(rows, params)
	if err != nil {
		log.Printf("Browse response failed: %v", err)
		http.Error(w, `{"error": "Failed to parse results"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) parseBrowseParams(r *http.Request) (queries.BrowseParams, error) {
	params := queries.BrowseParams{
		Page:          1,
		Limit:         20,
		Sort:          "popularity",
		Order:         "desc",
		IncludeFacets: true,
		FacetLimit:    10,
	}

	// Parse page
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		}
	}

	// Parse limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			params.Limit = limit
		}
	}

	// Parse sort
	if sort := r.URL.Query().Get("sort"); sort != "" {
		params.Sort = sort
	}

	// Parse order
	if order := r.URL.Query().Get("order"); order != "" {
		params.Order = order
	}

	// Parse categories (array parameter)
	params.Categories = r.URL.Query()["categories[]"]

	// Parse tags (array parameter)
	params.Tags = r.URL.Query()["tags[]"]

	// Parse popularity range
	if minStr := r.URL.Query().Get("popularity_min"); minStr != "" {
		if min, err := strconv.ParseFloat(minStr, 64); err == nil {
			params.PopularityMin = &min
		}
	}
	if maxStr := r.URL.Query().Get("popularity_max"); maxStr != "" {
		if max, err := strconv.ParseFloat(maxStr, 64); err == nil {
			params.PopularityMax = &max
		}
	}

	// Parse date range
	if dateFrom := r.URL.Query().Get("date_from"); dateFrom != "" {
		params.DateFrom = &dateFrom
	}
	if dateTo := r.URL.Query().Get("date_to"); dateTo != "" {
		params.DateTo = &dateTo
	}

	// Parse facets flag
	if facetsStr := r.URL.Query().Get("facets"); facetsStr != "" {
		if facets, err := strconv.ParseBool(facetsStr); err == nil {
			params.IncludeFacets = facets
		}
	}

	// Parse facet limit
	if facetLimitStr := r.URL.Query().Get("facet_limit"); facetLimitStr != "" {
		if facetLimit, err := strconv.Atoi(facetLimitStr); err == nil && facetLimit > 0 {
			params.FacetLimit = facetLimit
		}
	}

	return params, nil
}