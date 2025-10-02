package service

import (
	"business/pkg/model"
	"business/pkg/repo"
	"business/pkg/es"
	"context"

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
}

func (e *EsService) PushToEs(ctx context.Context, req *model.GetListBusinessRequest) (*model.Business, error) {
	// Lấy toàn bộ business từ repo
	businesses, err := e.repo.GetListBusiness(ctx,req,nil)
	if err != nil {
		return nil, err
	}

	// Tạo index nếu chưa có
	
	indexName := "business"
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
		return nil, err
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