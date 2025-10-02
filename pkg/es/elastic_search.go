package es

import "time"

type BulkDocument struct {
	ID   string
	Data interface{}
}

// SearchResult represents search response
type SearchResult struct {
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
		Hits []struct {
			ID     string                 `json:"_id"`
			Source map[string]interface{} `json:"_source"`
			Score  float64                `json:"_score"`
		} `json:"hits"`
	} `json:"hits"`
}

// Config holds ElasticSearch configuration
type Config struct {
	Addresses []string
	Username  string
	Password  string
	CloudID   string
	APIKey    string
	Timeout   time.Duration
}
