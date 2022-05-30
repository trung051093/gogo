package elasticsearchmodel

// SearchResults wraps the Elasticsearch search response.
type SearchResults struct {
	Total int                      `json:"total"`
	Hits  []map[string]interface{} `json:"hits"`
}
