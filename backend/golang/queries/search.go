package queries

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SearchQueries struct {
	db *pgxpool.Pool
}

func NewSearchQueries(db *pgxpool.Pool) *SearchQueries {
	return &SearchQueries{db: db}
}

func (sq *SearchQueries) BuildStatement(query string) (pgx.Rows, error) {
	sql := `
		SELECT id, quote, author, category, tags, 
		       paradedb.score(id) as relevance,
		       paradedb.snippet(quote) as highlighted_quote
		FROM quotes 
		WHERE quote @@@ $1 OR author @@@ $1 OR tags @@@ $1 OR category @@@ $1
		ORDER BY paradedb.score(id) DESC
		LIMIT 20
	`
	
	return sq.db.Query(context.Background(), sql, query)
}

func (sq *SearchQueries) BuildStatementWithFilters(query string, params BrowseParams) (pgx.Rows, error) {
	// Build BM25 boolean query parts
	var booleanParts []string
	var args []interface{}
	argIndex := 1

	// Main search query
	booleanParts = append(booleanParts, fmt.Sprintf("paradedb.match('quote', $%d)", argIndex))
	args = append(args, query)
	argIndex++

	// Add BM25-indexed filters to boolean query
	for _, category := range params.Categories {
		booleanParts = append(booleanParts, fmt.Sprintf("paradedb.term('category', $%d)", argIndex))
		args = append(args, category)
		argIndex++
	}

	// Build WHERE clause parts for non-BM25 fields
	var whereClauses []string

	// Tags filter (exact array matching)
	if len(params.Tags) > 0 {
		placeholders := make([]string, len(params.Tags))
		for i, tag := range params.Tags {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, tag)
			argIndex++
		}
		whereClauses = append(whereClauses, fmt.Sprintf("tags @> ARRAY[%s]", strings.Join(placeholders, ",")))
	}

	// Popularity range filters
	if params.PopularityMin != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("popularity >= $%d", argIndex))
		args = append(args, *params.PopularityMin)
		argIndex++
	}
	if params.PopularityMax != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("popularity <= $%d", argIndex))
		args = append(args, *params.PopularityMax)
		argIndex++
	}

	// Date range filters
	if params.DateFrom != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, *params.DateFrom)
		argIndex++
	}
	if params.DateTo != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, *params.DateTo)
		argIndex++
	}

	// Calculate OFFSET
	offset := (params.Page - 1) * params.Limit

	// Build the complete SQL query
	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = " AND " + strings.Join(whereClauses, " AND ")
	}

	sql := fmt.Sprintf(`
		SELECT id, quote, author, category, tags, popularity, created_at,
		       paradedb.score(id) as relevance,
		       paradedb.snippet(quote) as highlighted_quote
		FROM quotes 
		WHERE quotes @@@ paradedb.with_index('quotes_search_idx', 
			paradedb.boolean(must => ARRAY[%s])
		)%s
		ORDER BY paradedb.score(id) DESC
		LIMIT $%d OFFSET $%d
	`, strings.Join(booleanParts, ","), whereClause, argIndex, argIndex+1)

	args = append(args, params.Limit, offset)

	return sq.db.Query(context.Background(), sql, args...)
}

func (sq *SearchQueries) BuildResponse(rows pgx.Rows, query string) (SearchResponse, error) {
	var quotes []Quote
	
	for rows.Next() {
		var q Quote
		var tagsArray []string
		
		err := rows.Scan(&q.ID, &q.Quote, &q.Author, &q.Category, &tagsArray, &q.Relevance, &q.HighlightedQuote)
		if err != nil {
			return SearchResponse{}, err
		}
		
		q.Tags = tagsArray
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

func (sq *SearchQueries) BuildResponseWithFilters(rows pgx.Rows, query string, params BrowseParams) (BrowseResponse, error) {
	var quotes []Quote
	
	// Parse quotes
	for rows.Next() {
		var q Quote
		var tagsArray []string
		var createdAt interface{}
		
		var relevance *float64
		var highlightedQuote *string
		err := rows.Scan(&q.ID, &q.Quote, &q.Author, &q.Category, &tagsArray, &q.Popularity, &createdAt, &relevance, &highlightedQuote)
		if err != nil {
			return BrowseResponse{}, err
		}
		
		q.Tags = tagsArray
		// Handle created_at conversion if needed
		if createdAt != nil {
			if t, ok := createdAt.(time.Time); ok {
				createdAtStr := t.Format(time.RFC3339)
				q.CreatedAt = &createdAtStr
			}
		}
		// Handle relevance and highlighted quote
		if relevance != nil {
			q.Relevance = *relevance
		}
		if highlightedQuote != nil {
			q.HighlightedQuote = highlightedQuote
		}
		quotes = append(quotes, q)
	}

	if err := rows.Err(); err != nil {
		return BrowseResponse{}, err
	}

	// Get total count for search with filters
	totalCount, err := sq.getTotalCountWithFilters(query, params)
	if err != nil {
		return BrowseResponse{}, err
	}

	// Build pagination
	pagination := sq.buildPagination(params.Page, params.Limit, totalCount)

	// Build active filters
	activeFilters := sq.buildActiveFilters(params)

	response := BrowseResponse{
		Quotes:        quotes,
		Pagination:    pagination,
		ActiveFilters: activeFilters,
	}

	// Add facets if requested
	if params.IncludeFacets {
		facets, err := sq.buildFacetsWithSearch(query, params)
		if err != nil {
			return BrowseResponse{}, err
		}
		response.Facets = facets
	}

	return response, nil
}

func (sq *SearchQueries) getTotalCountWithFilters(query string, params BrowseParams) (int, error) {
	// Build same query as search but with COUNT
	var booleanParts []string
	var args []interface{}
	argIndex := 1

	// Main search query
	booleanParts = append(booleanParts, fmt.Sprintf("paradedb.match('quote', $%d)", argIndex))
	args = append(args, query)
	argIndex++

	// Add BM25-indexed filters to boolean query
	for _, category := range params.Categories {
		booleanParts = append(booleanParts, fmt.Sprintf("paradedb.term('category', $%d)", argIndex))
		args = append(args, category)
		argIndex++
	}

	// Build WHERE clause parts for non-BM25 fields
	var whereClauses []string

	// Tags filter (exact array matching)
	if len(params.Tags) > 0 {
		placeholders := make([]string, len(params.Tags))
		for i, tag := range params.Tags {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, tag)
			argIndex++
		}
		whereClauses = append(whereClauses, fmt.Sprintf("tags @> ARRAY[%s]", strings.Join(placeholders, ",")))
	}

	// Popularity range filters
	if params.PopularityMin != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("popularity >= $%d", argIndex))
		args = append(args, *params.PopularityMin)
		argIndex++
	}
	if params.PopularityMax != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("popularity <= $%d", argIndex))
		args = append(args, *params.PopularityMax)
		argIndex++
	}

	// Date range filters
	if params.DateFrom != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, *params.DateFrom)
		argIndex++
	}
	if params.DateTo != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, *params.DateTo)
		argIndex++
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = " AND " + strings.Join(whereClauses, " AND ")
	}

	sql := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM quotes 
		WHERE quotes @@@ paradedb.with_index('quotes_search_idx', 
			paradedb.boolean(must => ARRAY[%s])
		)%s
	`, strings.Join(booleanParts, ","), whereClause)

	var count int
	err := sq.db.QueryRow(context.Background(), sql, args...).Scan(&count)
	return count, err
}

func (sq *SearchQueries) buildPagination(page, limit, totalCount int) Pagination {
	totalPages := (totalCount + limit - 1) / limit // ceiling division
	
	return Pagination{
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		TotalCount: totalCount,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

func (sq *SearchQueries) buildActiveFilters(params BrowseParams) ActiveFilters {
	return ActiveFilters{
		Categories:    params.Categories,
		Tags:          params.Tags,
		PopularityMin: params.PopularityMin,
		PopularityMax: params.PopularityMax,
		DateFrom:      params.DateFrom,
		DateTo:        params.DateTo,
	}
}

func (sq *SearchQueries) buildFacetsWithSearch(query string, params BrowseParams) (*Facets, error) {
	// For search with filters, we'll generate facets based on the search results
	// This is a simplified version - you could make this more sophisticated
	facets := &Facets{}

	// Get category facets for search results
	categoryFacets, err := sq.getCategoryFacetsWithSearch(query, params)
	if err != nil {
		return nil, err
	}
	facets.Categories = categoryFacets

	// Get tag facets for search results
	tagFacets, err := sq.getTagFacetsWithSearch(query, params)
	if err != nil {
		return nil, err
	}
	facets.Tags = tagFacets

	// Get popularity range for search results
	popularityRange, err := sq.getPopularityRangeWithSearch(query, params)
	if err != nil {
		return nil, err
	}
	facets.PopularityRange = popularityRange

	return facets, nil
}

func (sq *SearchQueries) getCategoryFacetsWithSearch(query string, params BrowseParams) ([]FacetItem, error) {
	// Build facet query without category filter to show all categories in search results
	facetParams := params
	facetParams.Categories = nil
	
	var booleanParts []string
	var args []interface{}
	argIndex := 1

	// Main search query
	booleanParts = append(booleanParts, fmt.Sprintf("paradedb.match('quote', $%d)", argIndex))
	args = append(args, query)
	argIndex++

	// Build WHERE clause for non-category filters
	var whereClauses []string
	if len(facetParams.Tags) > 0 {
		placeholders := make([]string, len(facetParams.Tags))
		for i, tag := range facetParams.Tags {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, tag)
			argIndex++
		}
		whereClauses = append(whereClauses, fmt.Sprintf("tags @> ARRAY[%s]", strings.Join(placeholders, ",")))
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = " AND " + strings.Join(whereClauses, " AND ")
	}

	whereClause += " AND category IS NOT NULL"

	sql := fmt.Sprintf(`
		SELECT category, COUNT(*) as count
		FROM quotes
		WHERE quotes @@@ paradedb.with_index('quotes_search_idx', 
			paradedb.boolean(must => ARRAY[%s])
		)%s
		GROUP BY category
		ORDER BY count DESC
		LIMIT $%d
	`, strings.Join(booleanParts, ","), whereClause, argIndex)

	args = append(args, params.FacetLimit)

	rows, err := sq.db.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var facets []FacetItem
	for rows.Next() {
		var item FacetItem
		err := rows.Scan(&item.Value, &item.Count)
		if err != nil {
			return nil, err
		}
		facets = append(facets, item)
	}

	return facets, rows.Err()
}

func (sq *SearchQueries) getTagFacetsWithSearch(query string, params BrowseParams) ([]FacetItem, error) {
	// Similar to category facets but for tags
	facetParams := params
	facetParams.Tags = nil

	var booleanParts []string
	var args []interface{}
	argIndex := 1

	// Main search query
	booleanParts = append(booleanParts, fmt.Sprintf("paradedb.match('quote', $%d)", argIndex))
	args = append(args, query)
	argIndex++

	// Add category filters if any
	for _, category := range facetParams.Categories {
		booleanParts = append(booleanParts, fmt.Sprintf("paradedb.term('category', $%d)", argIndex))
		args = append(args, category)
		argIndex++
	}

	sql := fmt.Sprintf(`
		SELECT unnest(tags) as tag, COUNT(*) as count
		FROM quotes
		WHERE quotes @@@ paradedb.with_index('quotes_search_idx', 
			paradedb.boolean(must => ARRAY[%s])
		)
		GROUP BY tag
		ORDER BY count DESC
		LIMIT $%d
	`, strings.Join(booleanParts, ","), argIndex)

	args = append(args, params.FacetLimit)

	rows, err := sq.db.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var facets []FacetItem
	for rows.Next() {
		var item FacetItem
		err := rows.Scan(&item.Value, &item.Count)
		if err != nil {
			return nil, err
		}
		facets = append(facets, item)
	}

	return facets, rows.Err()
}

func (sq *SearchQueries) getPopularityRangeWithSearch(query string, params BrowseParams) (*PopularityRange, error) {
	var booleanParts []string
	var args []interface{}
	argIndex := 1

	// Main search query
	booleanParts = append(booleanParts, fmt.Sprintf("paradedb.match('quote', $%d)", argIndex))
	args = append(args, query)
	argIndex++

	// Add category filters if any
	for _, category := range params.Categories {
		booleanParts = append(booleanParts, fmt.Sprintf("paradedb.term('category', $%d)", argIndex))
		args = append(args, category)
		argIndex++
	}

	sql := fmt.Sprintf(`
		SELECT MIN(popularity), MAX(popularity)
		FROM quotes
		WHERE quotes @@@ paradedb.with_index('quotes_search_idx', 
			paradedb.boolean(must => ARRAY[%s])
		) AND popularity IS NOT NULL
	`, strings.Join(booleanParts, ","))

	var min, max *float64
	err := sq.db.QueryRow(context.Background(), sql, args...).Scan(&min, &max)
	if err != nil {
		return nil, err
	}

	if min == nil || max == nil {
		return nil, nil
	}

	return &PopularityRange{
		Min: *min,
		Max: *max,
	}, nil
}