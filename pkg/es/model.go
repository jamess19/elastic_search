package es

import (
	"time"
)

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
}

type SearchRequest struct {
    Index   string      `json:"index" example:"business"`      // Tên index
    Page    int          `json:"page" example:"1"`              // Số trang (bắt đầu từ 1)
    Size    int         `json:"size" example:"10"`             // Số document mỗi trang
    Sort    string      `json:"sort,omitempty" example:"{\"created_at\":\"desc\"}"` // key=field, value="asc|desc"
    Filters BusinessFilter `json:"filters,omitempty"` // các field cần filter
    Source  []string               `json:"_source,omitempty"`             // chọn field nào trả về (optional)
}

type BusinessFilter struct {
	Name        string    `json:"name"`
	Description string    `gorm:"type:text"`
	Address     string    `json:"address"`
	BusinessType        string    `json:"type"`
	Status      string    `json:"status"`
	CreateAt    *time.Time `gorm:"column:created_at"`
}

