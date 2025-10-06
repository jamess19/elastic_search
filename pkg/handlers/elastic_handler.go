package handlers

import (
	"business/pkg/es"
	"business/pkg/model"
	"business/pkg/service"
	"net/http"
	"gitlab.com/goxp/cloud0/ginext"
	"gitlab.com/goxp/cloud0/logger"
)

type ElasticHandlers struct {
	service service.EsInterface
}

func NewElasticHandlers(service *service.EsService) *ElasticHandlers {
	return &ElasticHandlers{service: service}
}

// PushToElastic
// @Tags Elastic
// @Security ApiKeyAuth
// @Summary push business in postgre to elastic
// @Description  push business in postgre to elastic
// @ID PushToElastic
// @Accept  json
// @Produce  json
// @Success 200 {object} []model.Business
// @Router /api/v1/elastic/push-to-elastic [post]
func (h *ElasticHandlers) PushToElastic(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.GinCtx, "PushToElastic")

	var req model.GetListBusinessRequest
	r.MustBind(&req)

	result, err := h.service.PushToEs(r.Context(), &req)
	if err != nil {
		log.WithError(err).Error("Failed to push data to Elastic")
		return nil, err
	}

	return &ginext.Response{
		Code: http.StatusOK,
		GeneralBody: &ginext.GeneralBody{
			Data: result,
		},
	}, nil
}

// SearchByField
// @Summary Search businesses by filters
// @Description Search documents with multiple filters and pagination
// @Tags Elastic
// @ID SearchByField
// @Accept json
// @Produce json
// @Param request body es.SearchRequest true "Search request"
// @Success 200 {object} []es.SearchResult
// @Router /api/v1/elastic/search-by-field [post]
func (h *ElasticHandlers) SearchByField(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.GinCtx, "SearchByField")

	var req es.SearchRequest
	r.MustBind(&req)

	// Set defaults to prevent panics
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	if req.Index == "" {
		req.Index = "business"
	}

	result, err := h.service.SearchWithField(r.Context(), req)
	if err != nil {
		log.WithError(err).Error("Error when get list business")
		return nil, err
	}

	return &ginext.Response{
		Code: http.StatusOK,
		GeneralBody: &ginext.GeneralBody{
			Data: result,
			Meta: map[string]interface{}{
				"page":      req.Page,
				"page_size": req.Size,
			},
		},
	}, nil
}

// FullTextSearch
// @Summary Full-text search documents
// @Description Perform multi-field full-text search with filters, pagination, and sort
// @Tags Elastic
// @ID FullTextSearch
// @Accept json
// @Produce json
// @Param request body es.SearchRequest true "Search request"
// @Success 200 {object} []es.SearchResult
// @Router /api/v1/elastic/fulltext-search [get]
func (h *ElasticHandlers) FullTextSearch(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.GinCtx, "FullTextSearch")

	var req es.SearchRequest
	r.MustBind(&req)

	// Set defaults to prevent panics
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	if req.Index == "" {
		req.Index = "business"
	}

	result, err := h.service.FullTextSearch(r.GinCtx, req)
	if err != nil {
		log.WithError(err).Error("Failed to perform full-text search")
		return nil, err
	}

	var total int64
	if result != nil {
		total = result.Hits.Total.Value
	}

	return &ginext.Response{
		Code: http.StatusOK,
		GeneralBody: &ginext.GeneralBody{
			Data: result,
			Meta: map[string]interface{}{
				"page": req.Page,
				"size": req.Size,
				"total": total,
			},
		},
	}, nil
}
