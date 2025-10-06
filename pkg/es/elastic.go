package es

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
)

func NewClient(cfg Config) (Client, error) {

	esCfg := elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
	}

	client, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	return &esClient{client: client}, nil
}

type Client interface {
	// Index operations
	CreateIndex(ctx context.Context, indexName string, mapping interface{}) error
	IndexExists(ctx context.Context, indexName string) (bool, error)
	
	// Document operations
	IndexDocument(ctx context.Context, indexName, docID string, doc interface{}) error
	
	// Bulk operations
	BulkIndex(ctx context.Context, indexName string, docs []BulkDocument) error
	
	// Search operations
	Search(ctx context.Context, indexName string, query interface{}) (*SearchResult, error)
	
	// Health check
	Ping(ctx context.Context) error
}

type esClient struct {
	client *elasticsearch.Client
}

