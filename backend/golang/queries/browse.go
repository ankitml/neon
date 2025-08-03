package queries

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BrowseQueries struct {
	db *pgxpool.Pool
}

func NewBrowseQueries(db *pgxpool.Pool) *BrowseQueries {
	return &BrowseQueries{db: db}
}

func (bq *BrowseQueries) BuildStatement(params BrowseParams) (pgx.Rows, error) {
	// Build WHERE clause
	whereClause, args := bq.buildWhereClause(params)
	
	// Build ORDER BY clause
	orderBy := bq.buildOrderBy(params.Sort, params.Order)
	
	// Calculate OFFSET
	offset := (params.Page - 1) * params.Limit
	
	// Build main query
	sql := fmt.Sprintf(`
		SELECT id, quote, author, category, tags, popularity, created_at
		FROM quotes
		%s
		%s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderBy, len(args)+1, len(args)+2)
	
	// Add limit and offset to args
	args = append(args, params.Limit, offset)
	
	return bq.db.Query(context.Background(), sql, args...)
}

func (bq *BrowseQueries) BuildResponse(rows pgx.Rows, params BrowseParams) (BrowseResponse, error) {
	var quotes []Quote
	
	// Parse quotes
	for rows.Next() {
		var q Quote
		var tagsArray []string
		var createdAt interface{}
		
		err := rows.Scan(&q.ID, &q.Quote, &q.Author, &q.Category, &tagsArray, &q.Popularity, &createdAt)
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
		quotes = append(quotes, q)
	}

	if err := rows.Err(); err != nil {
		return BrowseResponse{}, err
	}

	// Get total count
	totalCount, err := bq.getTotalCount(params)
	if err != nil {
		return BrowseResponse{}, err
	}

	// Build pagination
	pagination := bq.buildPagination(params.Page, params.Limit, totalCount)

	// Build active filters
	activeFilters := bq.buildActiveFilters(params)

	response := BrowseResponse{
		Quotes:        quotes,
		Pagination:    pagination,
		ActiveFilters: activeFilters,
	}

	// Add facets if requested
	if params.IncludeFacets {
		facets, err := bq.buildFacets(params)
		if err != nil {
			return BrowseResponse{}, err
		}
		response.Facets = facets
	}

	return response, nil
}

func (bq *BrowseQueries) buildWhereClause(params BrowseParams) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Category filter
	if len(params.Categories) > 0 {
		placeholders := make([]string, len(params.Categories))
		for i, category := range params.Categories {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, category)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("category = ANY(ARRAY[%s])", strings.Join(placeholders, ",")))
	}

	// Tags filter (AND logic - must have ALL specified tags)
	if len(params.Tags) > 0 {
		placeholders := make([]string, len(params.Tags))
		for i, tag := range params.Tags {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, tag)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("tags @> ARRAY[%s]", strings.Join(placeholders, ",")))
	}

	// Popularity range filter
	if params.PopularityMin != nil {
		conditions = append(conditions, fmt.Sprintf("popularity >= $%d", argIndex))
		args = append(args, *params.PopularityMin)
		argIndex++
	}
	if params.PopularityMax != nil {
		conditions = append(conditions, fmt.Sprintf("popularity <= $%d", argIndex))
		args = append(args, *params.PopularityMax)
		argIndex++
	}

	// Date range filter
	if params.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, *params.DateFrom)
		argIndex++
	}
	if params.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, *params.DateTo)
		argIndex++
	}

	if len(conditions) == 0 {
		return "", args
	}

	return "WHERE " + strings.Join(conditions, " AND "), args
}

func (bq *BrowseQueries) buildOrderBy(sort, order string) string {
	// Validate sort field
	validSorts := map[string]string{
		"popularity": "popularity",
		"created_at": "created_at",
		"random":     "RANDOM()",
	}
	
	sortField, exists := validSorts[sort]
	if !exists {
		sortField = "popularity" // default
	}

	// Validate order
	if order != "asc" && order != "desc" {
		order = "desc" // default
	}

	// Special case for random
	if sort == "random" {
		return "ORDER BY RANDOM()"
	}

	return fmt.Sprintf("ORDER BY %s %s", sortField, strings.ToUpper(order))
}

func (bq *BrowseQueries) getTotalCount(params BrowseParams) (int, error) {
	whereClause, args := bq.buildWhereClause(params)
	
	sql := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM quotes
		%s
	`, whereClause)
	
	var count int
	err := bq.db.QueryRow(context.Background(), sql, args...).Scan(&count)
	return count, err
}

func (bq *BrowseQueries) buildPagination(page, limit, totalCount int) Pagination {
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

func (bq *BrowseQueries) buildActiveFilters(params BrowseParams) ActiveFilters {
	return ActiveFilters{
		Categories:    params.Categories,
		Tags:          params.Tags,
		PopularityMin: params.PopularityMin,
		PopularityMax: params.PopularityMax,
		DateFrom:      params.DateFrom,
		DateTo:        params.DateTo,
	}
}

func (bq *BrowseQueries) buildFacets(params BrowseParams) (*Facets, error) {
	facets := &Facets{}

	// Get category facets
	categoryFacets, err := bq.getCategoryFacets(params)
	if err != nil {
		return nil, err
	}
	facets.Categories = categoryFacets

	// Get tag facets
	tagFacets, err := bq.getTagFacets(params)
	if err != nil {
		return nil, err
	}
	facets.Tags = tagFacets

	// Get popularity range
	popularityRange, err := bq.getPopularityRange(params)
	if err != nil {
		return nil, err
	}
	facets.PopularityRange = popularityRange

	return facets, nil
}

func (bq *BrowseQueries) getCategoryFacets(params BrowseParams) ([]FacetItem, error) {
	// Build facet query without category filter to show all categories
	facetParams := params
	facetParams.Categories = nil
	whereClause, args := bq.buildWhereClause(facetParams)
	
	// Handle the WHERE clause properly
	if whereClause == "" {
		whereClause = "WHERE category IS NOT NULL"
	} else {
		whereClause += " AND category IS NOT NULL"
	}
	
	sql := fmt.Sprintf(`
		SELECT category, COUNT(*) as count
		FROM quotes
		%s
		GROUP BY category
		ORDER BY count DESC
		LIMIT $%d
	`, whereClause, len(args)+1)
	
	args = append(args, params.FacetLimit)
	
	rows, err := bq.db.Query(context.Background(), sql, args...)
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

func (bq *BrowseQueries) getTagFacets(params BrowseParams) ([]FacetItem, error) {
	// Build facet query without tag filter to show all tags
	facetParams := params
	facetParams.Tags = nil
	whereClause, args := bq.buildWhereClause(facetParams)
	
	sql := fmt.Sprintf(`
		SELECT unnest(tags) as tag, COUNT(*) as count
		FROM quotes
		%s
		GROUP BY tag
		ORDER BY count DESC
		LIMIT $%d
	`, whereClause, len(args)+1)
	
	args = append(args, params.FacetLimit)
	
	rows, err := bq.db.Query(context.Background(), sql, args...)
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

func (bq *BrowseQueries) getPopularityRange(params BrowseParams) (*PopularityRange, error) {
	// Build facet query without popularity filter
	facetParams := params
	facetParams.PopularityMin = nil
	facetParams.PopularityMax = nil
	whereClause, args := bq.buildWhereClause(facetParams)
	
	// Handle the WHERE clause properly
	if whereClause == "" {
		whereClause = "WHERE popularity IS NOT NULL"
	} else {
		whereClause += " AND popularity IS NOT NULL"
	}
	
	sql := fmt.Sprintf(`
		SELECT MIN(popularity), MAX(popularity)
		FROM quotes
		%s
	`, whereClause)
	
	var min, max *float64
	err := bq.db.QueryRow(context.Background(), sql, args...).Scan(&min, &max)
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