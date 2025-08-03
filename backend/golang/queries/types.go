package queries

// Quote represents a quote from the database
type Quote struct {
	ID               int      `json:"id"`
	Quote            string   `json:"quote"`
	Author           string   `json:"author"`
	Category         *string  `json:"category,omitempty"`
	Tags             []string `json:"tags"`
	Relevance        float64  `json:"relevance,omitempty"`
	HighlightedQuote *string  `json:"highlighted_quote,omitempty"`
	Popularity       *float64 `json:"popularity,omitempty"`
	CreatedAt        *string  `json:"created_at,omitempty"`
}

// SearchResponse represents the response for search API
type SearchResponse struct {
	Results []Quote `json:"results"`
	Count   int     `json:"count"`
	Query   string  `json:"query"`
}

// Pagination represents pagination information
type Pagination struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	TotalPages int  `json:"total_pages"`
	TotalCount int  `json:"total_count"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// FacetItem represents a single facet value with count
type FacetItem struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

// PopularityRange represents min/max popularity values
type PopularityRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// Facets represents all available facets
type Facets struct {
	Categories      []FacetItem      `json:"categories"`
	Tags            []FacetItem      `json:"tags"`
	PopularityRange *PopularityRange `json:"popularity_range,omitempty"`
}

// ActiveFilters represents currently applied filters
type ActiveFilters struct {
	Categories    []string `json:"categories,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	PopularityMin *float64 `json:"popularity_min,omitempty"`
	PopularityMax *float64 `json:"popularity_max,omitempty"`
	DateFrom      *string  `json:"date_from,omitempty"`
	DateTo        *string  `json:"date_to,omitempty"`
}

// BrowseResponse represents the response for browse API
type BrowseResponse struct {
	Quotes        []Quote       `json:"quotes"`
	Pagination    Pagination    `json:"pagination"`
	Facets        *Facets       `json:"facets,omitempty"`
	ActiveFilters ActiveFilters `json:"active_filters"`
}

// BrowseParams represents parameters for browse queries
type BrowseParams struct {
	Page          int      `json:"page"`
	Limit         int      `json:"limit"`
	Sort          string   `json:"sort"`
	Order         string   `json:"order"`
	Categories    []string `json:"categories"`
	Tags          []string `json:"tags"`
	PopularityMin *float64 `json:"popularity_min"`
	PopularityMax *float64 `json:"popularity_max"`
	DateFrom      *string  `json:"date_from"`
	DateTo        *string  `json:"date_to"`
	IncludeFacets bool     `json:"include_facets"`
	FacetLimit    int      `json:"facet_limit"`
}