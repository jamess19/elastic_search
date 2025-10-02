package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// Ping checks connection to ElasticSearch
func (c *esClient) Ping(ctx context.Context) error {
	res, err := c.client.Ping(c.client.Ping.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("ping error: %s", res.String())
	}

	return nil
}

// CreateIndex creates a new index with mapping
func (c *esClient) CreateIndex(ctx context.Context, indexName string, mapping interface{}) error {
	var body []byte
	var err error

	if mapping != nil {
		body, err = json.Marshal(mapping)
		if err != nil {
			return fmt.Errorf("failed to marshal mapping: %w", err)
		}
	}

	res, err := c.client.Indices.Create(
		indexName,
		c.client.Indices.Create.WithContext(ctx),
		c.client.Indices.Create.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("create index error: %s", res.String())
	}

	return nil
}

// DeleteIndex deletes an index
func (c *esClient) DeleteIndex(ctx context.Context, indexName string) error {
	res, err := c.client.Indices.Delete(
		[]string{indexName},
		c.client.Indices.Delete.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("failed to delete index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("delete index error: %s", res.String())
	}

	return nil
}

// IndexExists checks if an index exists
func (c *esClient) IndexExists(ctx context.Context, indexName string) (bool, error) {
	res, err := c.client.Indices.Exists(
		[]string{indexName},
		c.client.Indices.Exists.WithContext(ctx),
	)
	if err != nil {
		return false, fmt.Errorf("failed to check index existence: %w", err)
	}
	defer res.Body.Close()

	return res.StatusCode == 200, nil
}

// IndexDocument indexes a single document
func (c *esClient) IndexDocument(ctx context.Context, indexName, docID string, doc interface{}) error {
	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: docID,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("failed to index document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("index document error: %s", res.String())
	}

	return nil
}

// GetDocument retrieves a document by ID
func (c *esClient) GetDocument(ctx context.Context, indexName, docID string, result interface{}) error {
	res, err := c.client.Get(
		indexName,
		docID,
		c.client.Get.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("failed to get document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("get document error: %s", res.String())
	}

	var response struct {
		Source json.RawMessage `json:"_source"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if err := json.Unmarshal(response.Source, result); err != nil {
		return fmt.Errorf("failed to unmarshal source: %w", err)
	}

	return nil
}

// UpdateDocument updates a document
func (c *esClient) UpdateDocument(ctx context.Context, indexName, docID string, doc interface{}) error {
	update := map[string]interface{}{
		"doc": doc,
	}

	data, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("failed to marshal update: %w", err)
	}

	req := esapi.UpdateRequest{
		Index:      indexName,
		DocumentID: docID,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("update document error: %s", res.String())
	}

	return nil
}

// DeleteDocument deletes a document
func (c *esClient) DeleteDocument(ctx context.Context, indexName, docID string) error {
	req := esapi.DeleteRequest{
		Index:      indexName,
		DocumentID: docID,
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("delete document error: %s", res.String())
	}

	return nil
}

// BulkIndex performs bulk indexing
func (c *esClient) BulkIndex(ctx context.Context, indexName string, docs []BulkDocument) error {
	if len(docs) == 0 {
		return nil
	}

	var buf bytes.Buffer

	for _, doc := range docs {
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": indexName,
				"_id":    doc.ID,
			},
		}

		metaJSON, err := json.Marshal(meta)
		if err != nil {
			return fmt.Errorf("failed to marshal meta: %w", err)
		}

		buf.Write(metaJSON)
		buf.WriteByte('\n')

		docJSON, err := json.Marshal(doc.Data)
		if err != nil {
			return fmt.Errorf("failed to marshal document: %w", err)
		}

		buf.Write(docJSON)
		buf.WriteByte('\n')
	}

	res, err := c.client.Bulk(
		bytes.NewReader(buf.Bytes()),
		c.client.Bulk.WithContext(ctx),
		c.client.Bulk.WithRefresh("true"),
	)
	if err != nil {
		return fmt.Errorf("bulk request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk error: %s", res.String())
	}

	return nil
}

// Search performs a search query
func (c *esClient) Search(ctx context.Context, indexName string, query interface{}) (*SearchResult, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("failed to encode query: %w", err)
	}

	res, err := c.client.Search(
		c.client.Search.WithContext(ctx),
		c.client.Search.WithIndex(indexName),
		c.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	var result SearchResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}