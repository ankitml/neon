# ParadeDB Index Plan for Themed Quote Collections

Based on ParadeDB documentation analysis and project requirements for a full-text search application, this document outlines the recommended search indexes for the quotes table.

## Table Structure
```sql
CREATE TABLE quotes (
    id SERIAL PRIMARY KEY,
    quote TEXT NOT NULL,
    author VARCHAR(255) NOT NULL,
    tags TEXT[] NOT NULL DEFAULT '{}',
    popularity DECIMAL(10, 8),
    category VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Recommended Indexes

### 1. Comprehensive BM25 Search Index (All Fields)
```sql
CREATE INDEX quotes_search_idx ON quotes
USING bm25 (id, quote, author, category, tags)
WITH (key_field='id');
```

**Purpose:** Complete search functionality for theme-based quote discovery
**Key Limitation:** ParadeDB allows only **one BM25 index per table**, so all searchable fields must be in a single index
**Features:**
- Searches across quote text, author names, categories, and tags
- BM25 relevance scoring for optimal ranking
- Supports phrase search, fuzzy matching, and boolean queries

## Search Capabilities Enabled

### Core Search Features
1. **Theme-based Search**: `WHERE quote @@@ 'courage innovation'`
2. **Author Filtering**: `WHERE quote @@@ 'perseverance' AND author @@@ 'Roosevelt'`
3. **Phrase Search**: `WHERE quote @@@ '"never give up"'`
4. **Fuzzy Search**: `WHERE quote @@@ paradedb.fuzzy_term('quote', 'perserverance', 1)`
5. **Tag Filtering**: `WHERE tags @@@ 'inspirational wisdom'`
6. **Boolean Queries**: Complex theme combinations with AND/OR logic

### Advanced Features
- **BM25 Relevance Scoring** for ranking quote relevance: `ORDER BY paradedb.score(id) DESC`
- **Highlighting** to show matching terms: `paradedb.snippet(quote)`
- **Fast Field Sorting** by author/category for secondary ordering
- **Filtering** by popularity ranges: `WHERE popularity @@@ '>0.5'`

## Query Examples

### Basic Theme Search
```sql
SELECT id, quote, author, paradedb.score(id) as relevance
FROM quotes 
WHERE quote @@@ 'courage perseverance'
ORDER BY paradedb.score(id) DESC
LIMIT 10;
```

### Advanced Search with Highlighting
```sql
SELECT id, 
       paradedb.snippet(quote) as highlighted_quote,
       author, 
       category,
       paradedb.score(id) as relevance
FROM quotes 
WHERE quote @@@ 'innovation technology' 
   OR tags @@@ 'future progress'
ORDER BY paradedb.score(id) DESC, popularity DESC
LIMIT 20;
```

### Phrase Search with Author Filter
```sql
SELECT quote, author, category
FROM quotes 
WHERE quote @@@ '"give up"~1' 
  AND author @@@ 'Churchill Roosevelt Lincoln'
ORDER BY paradedb.score(id) DESC;
```

## Performance Considerations
- Primary index supports all core search operations
- Fast fields enable efficient sorting and filtering
- BM25 scoring provides relevance-based ranking
- Position recording enables phrase and proximity searches
- Tag index allows independent tag-based queries

## Implementation Priority
1. **Primary BM25 Search Index** - Essential for core functionality
2. **Tags Array Search Index** - Important for enhanced theme discovery

This configuration provides comprehensive full-text search capabilities while maintaining optimal performance for the themed quote collections application.