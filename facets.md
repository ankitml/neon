# Search and Facets Implementation

## Overview
Implementation of combined search and filtering functionality using ParadeDB BM25 search with hybrid filtering approach. The browse API with facets has been completed. This document focuses on adding filtering support to the search API.

## Database Schema and BM25 Index
ParadeDB version: **0.15.26**
BM25 Index: `quotes_search_idx` with fields:
- `id` (I64) - Primary key, fast field
- `quote` (Str) - Full-text searchable content  
- `author` (Str) - Author names, searchable
- `category` (Str) - Quote categories, filterable
- `tags` (Str) - Tags array, searchable
- `ctid` (U64) - Row identifier, fast field

Non-indexed fields (require WHERE clauses):
- `popularity` (DECIMAL) - Numerical score for ranking
- `created_at` (TIMESTAMP) - Date-based filtering

## Search + Filters Implementation Status

### ✅ Completed: Browse API with Facets
- Full browse functionality with faceting
- Category, tag, popularity, and date filtering
- Pagination and sorting
- Frontend integration with refinements UI

### 🎯 Current Goal: Hybrid Search + Filters

## Supported Filter Types with ParadeDB BM25

### ✅ **Boolean Query Method** (for BM25-indexed fields):

**Category Filtering:**
```sql
paradedb.boolean(
  must => ARRAY[
    paradedb.match('quote', 'search_term'),
    paradedb.term('category', 'exact_category_value')
  ]
)
```

**Author Search (fuzzy matching):**
```sql
paradedb.boolean(
  must => ARRAY[
    paradedb.match('quote', 'search_term'),
    paradedb.match('author', 'author_name_part')
  ]
)
```

**Tag Filtering (text search within tags):**
```sql
paradedb.boolean(
  must => ARRAY[
    paradedb.match('quote', 'search_term'),
    paradedb.match('tags', 'tag_value')
  ]
)
```

### ✅ **WHERE Clause Method** (for non-indexed fields):

**Popularity Range:**
```sql
WHERE quotes @@@ paradedb.with_index('quotes_search_idx', paradedb.match('quote', 'term'))
  AND popularity BETWEEN 0.01 AND 0.15
```

**Date Range:**
```sql
WHERE quotes @@@ paradedb.with_index('quotes_search_idx', paradedb.match('quote', 'term'))
  AND created_at >= '2020-01-01'
  AND created_at <= '2024-12-31'
```

**Exact Tag Matching (array membership):**
```sql
WHERE quotes @@@ paradedb.with_index('quotes_search_idx', paradedb.match('quote', 'term'))
  AND 'Love' = ANY(tags)
```

**Multiple Tag AND Logic:**
```sql
WHERE quotes @@@ paradedb.with_index('quotes_search_idx', paradedb.match('quote', 'term'))
  AND tags @> ARRAY['tag1', 'tag2']
```

## Implementation Plan

### Phase 1: Hybrid Search Handler ✅ Ready to Implement
1. **Extend Search Handler** (`handlers.go`)
   - Add filter parameter parsing (reuse from browse API)
   - Implement hybrid query building:
     - Boolean query for BM25-indexed fields
     - WHERE clauses for non-indexed fields
   - Maintain search relevance scoring

2. **Query Builder Strategy**
   ```go
   // Pseudo-code approach
   func buildSearchWithFilters(searchTerm string, filters BrowseParams) string {
     // Build boolean query parts for BM25 fields
     booleanMust := []string{
       paradedb.match('quote', searchTerm)
     }
     
     if len(filters.Categories) > 0 {
       for _, cat := range filters.Categories {
         booleanMust = append(booleanMust, paradedb.term('category', cat))
       }
     }
     
     // Build WHERE clause parts for non-BM25 fields  
     var whereClauses []string
     if filters.PopularityMin != nil {
       whereClauses = append(whereClauses, "popularity >= $X")
     }
     
     // Combine: BM25 query + WHERE filters
     sql := fmt.Sprintf(`
       SELECT id, quote, author, category, tags, popularity, created_at, 
              paradedb.score(id) as bm25_score
       FROM quotes 
       WHERE quotes @@@ paradedb.with_index('quotes_search_idx', 
         paradedb.boolean(must => ARRAY[%s])
       ) %s
       ORDER BY paradedb.score(id) DESC
       LIMIT $Y OFFSET $Z
     `, strings.Join(booleanMust, ","), strings.Join(whereClauses, " AND "))
   }
   ```

### Phase 2: Frontend Integration
1. **Update Search Component**
   - Add same filter UI used in browse mode
   - Show "Search + Filters" indicator when both active
   - Maintain search relevance scores in results

2. **Unified API Response**
   - Include facets in search responses (based on current filters)
   - Show active filters with remove buttons
   - Display search relevance scores alongside filtered results

## Filter Support Matrix

| Filter Type | Browse API | Search API (Planned) | Method |
|-------------|------------|---------------------|---------|
| Categories | ✅ | ✅ | Boolean query |
| Tags (exact) | ✅ | ✅ | WHERE clause |
| Tags (fuzzy) | ❌ | ✅ | Boolean query |
| Author search | ❌ | ✅ | Boolean query |
| Popularity range | ✅ | ✅ | WHERE clause |
| Date range | ✅ | ✅ | WHERE clause |
| Sorting | ✅ | ⚠️ | Score-based only |

## Technical Notes

- **Hybrid Approach**: Combine ParadeDB boolean queries with PostgreSQL WHERE clauses
- **Performance**: BM25 filters integrate with search scoring; WHERE filters applied post-search
- **Scoring**: Maintain search relevance while filtering
- **Facets**: Generate facets based on search results + current filters  
- **Compatibility**: All existing browse functionality remains unchanged

