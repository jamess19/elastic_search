package service

import (
	"business/pkg/es"
	"business/pkg/model"
	"business/pkg/repo"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gitlab.com/goxp/cloud0/logger"
)

type EsService struct {
	client es.Client
	repo repo.PGInterface
}

func NewEsService( repo repo.PGInterface, client es.Client) *EsService {
	return &EsService{repo: repo, client: client}
}

type EsInterface interface {
	PushToEs(ctx context.Context, req *model.GetListBusinessRequest) (*model.Business, error) 
	SearchWithField(ctx context.Context, req es.SearchRequest) (*model.GetListBusinessResponse, error)
	FullTextSearch(ctx context.Context, req es.SearchRequest) (*es.SearchResult, error)

}

func (e *EsService) PushToEs(ctx context.Context, req *model.GetListBusinessRequest) (*model.Business, error) {
	// Lấy toàn bộ business từ repo
	businesses, err := e.repo.GetListBusiness(ctx,req,nil)
	if err != nil {
		return nil, err
	}

	// Tạo index nếu chưa có
	indexName := "business"
	
	// Check if index exists first
	exists, err := e.client.IndexExists(ctx, indexName)
	if err != nil {
		return nil, fmt.Errorf("failed to check index existence: %w", err)
	}
	
	if !exists {
		mapping := map[string]interface{}{
			"mappings": map[string]interface{}{
				"properties": map[string]interface{}{
					"id": map[string]string{
						"type": "keyword",
					},
					"name": map[string]string{
						"type": "text",
					},
					"description": map[string]string{
						"type": "text",
					},
					"address": map[string]string{
						"type": "text",
					},
					"businessType": map[string]string{
						"type": "keyword",
					},
					"status": map[string]string{
						"type": "keyword",
					},
					"createAt": map[string]string{
						"type": "date",
					},
					"workerName": map[string]string{
						"type": "text",
					},
					"staffs": map[string]interface{}{
						"type": "nested",
						"properties": map[string]interface{}{
							"id": map[string]string{
								"type": "keyword",
							},
							"name": map[string]string{
								"type": "text",
							},
							"role": map[string]string{
								"type": "keyword",
							},
						},
					},
				},
			},
		}
		if err := e.client.CreateIndex(ctx, indexName, mapping); err != nil {
			return nil, fmt.Errorf("failed to create index: %w", err)
		}
	}

	for _, b := range businesses.Data {
		if err := e.client.IndexDocument(ctx, indexName, b.ID.String(), b); err != nil {
			return nil, err
		}
	}

	if len(businesses.Data) > 0 {
		return &businesses.Data[0], nil
	}
	return nil, nil
}

func (e* EsService) SearchWithField(ctx context.Context, req es.SearchRequest) (*model.GetListBusinessResponse, error) {
	log:= logger.WithCtx(ctx, "esService.SearchWithField")
	from := (req.Page - 1) * req.Size
	data, _ := json.Marshal(req.Filters)
	var filterMap map[string]interface{}
	err := json.Unmarshal(data, &filterMap)
	if(err!= nil) {
		log.WithError(err).Error("error in marshal")
	}	
	mustQueries := make([]map[string]interface{}, 0)
	for field, value := range filterMap {
		if v, ok := value.(string); ok && strings.TrimSpace(v) != "" {
			mustQueries = append(mustQueries, map[string]interface{}{
				"match": map[string]interface{}{
					field: v,
				},
			})
		}
	}

	query := map[string]interface{}{
		"from": from,
		"size": req.Size,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": mustQueries,
			},
		},
	}

	result, err := e.client.Search(ctx, req.Index, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search with filters: %w", err)
	}

	businesses := make([]model.Business, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		var b model.Business
		src, _ := json.Marshal(hit.Source)
		if err := json.Unmarshal(src, &b); err == nil {
			if parsedID, err := uuid.Parse(hit.ID); err == nil {
				b.ID = parsedID
			}
			businesses = append(businesses, b)
		}
	}

	resp := &model.GetListBusinessResponse{
		Data: businesses,
		Meta: map[string]interface{}{
			"total": result.Hits.Total.Value,
			"page":  req.Page,
			"size":  req.Size,
		},
	}

	return resp, nil
}

func (e *EsService) FullTextSearch(ctx context.Context, req es.SearchRequest) (*es.SearchResult, error) {
    from := (req.Page - 1) * req.Size
    if from < 0 {
        from = 0
    }

    // Convert Filters (struct) → map
    data, _ := json.Marshal(req.Filters)
    var filterMap map[string]interface{}
    _ = json.Unmarshal(data, &filterMap)

    mustQueries := make([]map[string]interface{}, 0)

    // Exact match filters
    for field, value := range filterMap {
        if str, ok := value.(string); ok && strings.TrimSpace(str) != "" {
            mustQueries = append(mustQueries, map[string]interface{}{
                "match": map[string]interface{}{
                    field: str,
                },
            })
        }
    }

    // Multi-field full-text search (tất cả string filters đều dùng multi_match)
    multiFields := []string{}
    multiQuery := ""
    for field, value := range filterMap {
        if str, ok := value.(string); ok && strings.TrimSpace(str) != "" {
            multiFields = append(multiFields, field)
            multiQuery = str
        }
    }
    if len(multiFields) > 0 && multiQuery != "" {
        mustQueries = append(mustQueries, map[string]interface{}{
            "multi_match": map[string]interface{}{
                "query":  multiQuery,
                "fields": multiFields,
            },
        })
    }

    // Build sort
    sortQuery := []map[string]interface{}{}
    if req.Sort != "" {
        parts := strings.Split(req.Sort, ":")
        if len(parts) == 2 {
            sortQuery = append(sortQuery, map[string]interface{}{
                parts[0]: map[string]string{"order": parts[1]},
            })
        }
    }

    // Build final query
    query := map[string]interface{}{
        "from": from,
        "size": req.Size,
        "query": map[string]interface{}{
            "bool": map[string]interface{}{
                "must": mustQueries,
            },
        },
    }
    if len(sortQuery) > 0 {
        query["sort"] = sortQuery
    }
    if len(req.Source) > 0 {
        query["_source"] = req.Source
    }

    // Call ES
    var buf bytes.Buffer
    if err := json.NewEncoder(&buf).Encode(query); err != nil {
        return nil, fmt.Errorf("encode query failed: %w", err)
    }

	result, err := e.client.Search(ctx, req.Index, query)
	if err != nil {
		return nil, fmt.Errorf("full-text search failed: %w", err)
	}

	return result, nil
}
