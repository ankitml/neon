package queries

import (
	"context"

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