# Browse API with Facets and Pagination

## Overview
Design and implement a browse API that allows users to explore quotes with filtering, faceting, and pagination capabilities. The facets system should be generic and reusable for both browse and search functionality.

## Database Schema Analysis
From the `quotes` table, we have these facetable fields:
- `author` (VARCHAR(255)) - Individual quote authors
- `category` (VARCHAR(100)) - Quote categories (arts, books, death, education, etc.)
- `tags` (TEXT[]) - Array of tags per quote (Life, Good, humor, Love, etc.)
- `popularity` (DECIMAL) - Numerical score for ranking
- `created_at` (TIMESTAMP) - Date-based filtering

## API Design

### Endpoint: `/api/browse`

**Query Parameters:**
- `page` (int, default: 1) - Page number for pagination
- `limit` (int, default: 20, max: 100) - Items per page
- `sort` (string, default: "popularity") - Sort field: "popularity", "created_at", "random"
- `order` (string, default: "desc") - Sort order: "asc", "desc"

**Filter Parameters:**
- `categories[]` (string array) - Filter by categories
- `tags[]` (string array) - Filter by tags (AND logic - quote must have ALL specified tags)
- `popularity_min` (float) - Minimum popularity score
- `popularity_max` (float) - Maximum popularity score
- `date_from` (ISO date) - Created after date
- `date_to` (ISO date) - Created before date

**Facet Parameters:**
- `facets` (boolean, default: true) - Include facet counts in response
- `facet_limit` (int, default: 10) - Max items per facet

### Response Format

```json
{
  "quotes": [
    {
      "id": 1,
      "quote": "Life is what happens...",
      "author": "John Lennon",
      "category": "life", 
      "tags": ["Life", "inspirational"],
      "popularity": 0.85,
      "created_at": "2024-01-15T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total_pages": 250,
    "total_count": 5000,
    "has_next": true,
    "has_prev": false
  },
  "facets": {
    "categories": [
      {"value": "life", "count": 1200},
      {"value": "love", "count": 890}
    ],
    "tags": [
      {"value": "Life", "count": 380},
      {"value": "inspirational", "count": 192}
    ],
    "popularity_range": {
      "min": 0.01,
      "max": 0.99
    }
  },
  "active_filters": {
    "categories": ["life"],
    "tags": ["inspirational"]
  }
}
```

## Implementation Plan

### Phase 1: Core Browse API
1. **Create Browse Handler** (`handlers.go`)
   - New `BrowseQuotesHandler` function
   - Parse query parameters with validation
   - Build dynamic SQL with WHERE clauses
   - Implement pagination logic

2. **Add Route** (`main.go`)
   - Register `/api/browse` endpoint
   - Add CORS headers

3. **Database Queries**
   - Main quotes query with filters and pagination
   - Count query for total results
   - Facet aggregation queries

### Phase 2: Generic Faceting System
1. **Facet Types**
   - `TermFacet`: For categorical fields (category)
   - `ArrayFacet`: For array fields (tags) 
   - `RangeFacet`: For numerical fields (popularity)
   - `DateRangeFacet`: For timestamp fields (future use)

2. **Facet Interface**
   ```go
   type Facet interface {
       GetName() string
       BuildQuery(filters map[string]interface{}) string
       ParseResults(rows *sql.Rows) ([]FacetItem, error)
   }
   
   type FacetItem struct {
       Value string `json:"value"`
       Count int    `json:"count"`
   }
   ```

3. **Facet Manager**
   - Registry of available facets
   - Execute facet queries in parallel
   - Combine results

### Phase 3: Search Integration
1. **Extend Search API**
   - Add same filter parameters to `/api/search`
   - Include facets in search responses
   - Maintain search relevance + filtering

2. **Unified Filter System**
   - Shared filter parsing logic
   - Common SQL builder functions
   - Consistent response format

### Phase 4: Performance Optimization
1. **Database Indexes**
   - Composite indexes for common filter combinations
   - Array indexes for tags (GIN index)
   - Consider materialized views for heavy facet queries

2. **Caching Strategy**
   - Redis cache for facet counts
   - Cache invalidation on data updates
   - Response compression

## Decisions Made

1. **Tag Filtering Logic**: ✅ AND logic (quote must have ALL specified tags)
2. **Author Faceting**: ✅ Skip for now (will implement search-as-you-type later)
3. **Default Sorting**: ✅ `popularity DESC` as default, but sorting parameter required
4. **Performance Priority**: ✅ Accuracy over speed (real-time facet counts)
5. **UI Integration**: ✅ Same UI with layered functionality (details below)

## UI Integration Plan

### Layout Structure
```
┌─────────────────────────────────────┐
│         Search Box                  │ ← Search API
├─────────────────────────────────────┤
│    Filter Row (Categories, Tags,    │ ← Browse API
│    Sort, Date Range)                │
├─────────────────────────────────────┤
│                                     │
│         Quotes List                 │ ← Results from active API
│         (Paginated)                 │
│                                     │
└─────────────────────────────────────┘
```

### API Usage Logic
- **Search Mode**: When user types in search box → Use `/api/search`
- **Browse Mode**: When user uses filters/sorting only → Use `/api/browse` 
- **Future**: Combine search + filters → Enhanced `/api/search` with filter support

### Frontend State Management
```javascript
const state = {
  mode: 'browse', // 'browse' | 'search'
  searchQuery: '',
  filters: {
    categories: [],
    tags: [],
    sort: 'popularity',
    order: 'desc',
    popularity_min: null,
    popularity_max: null
  },
  results: [],
  facets: {},
  pagination: {}
}
```

### Component Updates Required

1. **SearchBox Component**
   - Add `onInput` handler to switch to search mode
   - Clear filters when switching to search mode (optional)
   - Debounce search queries (300ms)

2. **New FilterRow Component**
   - Category multi-select dropdown (from facets)
   - Tag multi-select with popular tags (from facets, AND logic)
   - Sort dropdown (popularity, created_at, random)
   - Sort order toggle (asc/desc) 
   - Clear all filters button

3. **Updated QuotesList Component**
   - Handle both search and browse results
   - Show active filters as removable chips
   - Loading states for both APIs
   - Pagination controls

4. **New Facets Component**
   - Display facet counts next to filter options
   - Update when filters change
   - Show "applied" state for active filters

### API Switching Logic
```javascript
const fetchQuotes = async () => {
  const isSearchMode = searchQuery.trim().length > 0;
  
  if (isSearchMode) {
    // Use search API
    const params = new URLSearchParams({
      q: searchQuery,
      page: pagination.page,
      limit: pagination.limit
    });
    const response = await fetch(`/api/search?${params}`);
  } else {
    // Use browse API with filters
    const params = new URLSearchParams({
      page: pagination.page,
      limit: pagination.limit,
      sort: filters.sort,
      order: filters.order,
      facets: true
    });
    
    // Add array filters
    filters.categories.forEach(cat => params.append('categories[]', cat));
    filters.tags.forEach(tag => params.append('tags[]', tag));
    
    const response = await fetch(`/api/browse?${params}`);
  }
};
```

### Phase Integration with Existing Work
- **Phase 1**: Implement browse API backend (reuse existing Go structure)
- **Phase 2**: Update frontend components (modify existing SvelteKit app) 
- **Phase 3**: Add faceting system to search API
- **Phase 4**: Performance optimization

## Technical Notes

- Use pgx driver for efficient bulk operations
- Implement query builder pattern for dynamic WHERE clauses
- Consider using PostgreSQL's array operators for tag filtering
- Use LIMIT/OFFSET for pagination (may switch to cursor-based for very large datasets)
- Add request timeout handling (10s max)
- Include query performance metrics in response headers