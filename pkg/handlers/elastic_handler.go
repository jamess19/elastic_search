package handlers

import (
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
	r.MustBind(&req); 

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
