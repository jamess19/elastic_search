package handlers

import (
	"business/pkg/model"
	"business/pkg/service"
	"business/pkg/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/praslar/lib/common"
	"gitlab.com/goxp/cloud0/ginext"
	"gitlab.com/goxp/cloud0/logger"
)

type BusinessHandlers struct {
	service service.BusinessInterface
}

func NewBusinessHandlers(service service.BusinessInterface) *BusinessHandlers {
	return &BusinessHandlers{service: service}
}

// CreateBusiness
// @Tags Business
// @Security ApiKeyAuth
// @Summary Create new Business
// @Description create a new Business for system
// @ID CreateBusiness
// @Accept  json
// @Produce  json
// @Param data body model.BusinessRequest true "body data"
// @Success 200 {object} model.Business
// @Router /api/v1/business/create [post]
func (h *BusinessHandlers) CreateBusiness(r *ginext.Request) (*ginext.Response, error) {
	req := model.BusinessRequest{}
	r.MustBind(&req)

	if err := common.CheckRequireValid(req); err != nil {
		return nil, ginext.NewError(http.StatusBadRequest, err.Error())
	}

	rs, err := h.service.CreateBusiness(r.GinCtx, req)
	if err != nil {
		return nil, err
	}

	return ginext.NewResponseData(http.StatusCreated, rs), nil
}

// CreateBusiness_v2
// @Tags Business
// @Security ApiKeyAuth
// @Summary Create new 10.000 businesses
// @Description create new 10.000 businesses for system
// @ID CreateBusiness_v2
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Business
// @Router /api/v1/business/create-v2 [post]
func (h *BusinessHandlers) CreateBusiness_v2(r *ginext.Request) (*ginext.Response, error) {
	// req := []model.BusinessRequest{}
	// r.MustBind(&req)

	// if err := common.CheckRequireValid(req); err != nil {
	// 	return nil, ginext.NewError(http.StatusBadRequest, err.Error())
	// }

	rs, err := h.service.CreateBusiness_v2(r.GinCtx)
	if err != nil {
		return nil, err
	}
	

	return ginext.NewResponseData(http.StatusCreated, rs), nil
}

// UpdateBusiness
// @Tags Business
// @Security ApiKeyAuth
// @Summary update Business
// @Description update Business
// @ID UpdateBusiness
// @Accept  json
// @Produce  json
// @Param data body model.BusinessRequest true "body data"
// @Success 200 {object} model.Business
// @Router /api/v1/business/update/{id} [put]
func (h *BusinessHandlers) UpdateBusiness(r *ginext.Request) (*ginext.Response, error) {

	req := model.BusinessRequest{}
	// if req.ID = utils.ParseIDFromUri(r.GinCtx); req.ID == nil {
	// 	return nil, ginext.NewError(http.StatusForbidden, "Wrong ID")
	// }

	r.MustBind(&req)

	rs, err := h.service.UpdateBusiness(r.GinCtx, req)
	if err != nil {
		return nil, err
	}
	return ginext.NewResponseData(http.StatusOK, rs), nil
}

// ListBusiness
// @Tags Business
// @Security ApiKeyAuth
// @Summary List Businesss
// @Description Get a list of Businesss
// @ID ListBusiness
// @Accept  json
// @Produce  json
// @Success 200 {object} []model.Business
// @Router /api/v1/business/get-list [get]
func (h *BusinessHandlers) ListBusiness(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.GinCtx, "ListBusiness")

	var req model.GetListBusinessRequest
	r.MustBind(&req)

	rs, err := h.service.GetListBusiness(r.Context(), &req)
	if err != nil {
		log.WithError(err).Error("Error when get list business")
		return nil, err
	}

	return &ginext.Response{
		Code: http.StatusOK,
		GeneralBody: &ginext.GeneralBody{
			Data: rs.Data,
			Meta: rs.Meta,
		},
	}, nil
}

// ListBusiness_v2
// @Tags Business
// @Security ApiKeyAuth
// @Summary List Businesss
// @Description Get a list of Businesss
// @ID ListBusiness-v2
// @Accept  json
// @Produce  json
// @Success 200 {object} []model.Business
// @Router /api/v1/business/get-list-v2 [get]
func (h *BusinessHandlers) ListBusiness_v2(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.GinCtx, "ListBusiness_v2")

	var req model.GetListBusinessRequest
	r.MustBind(&req)

	rs, err := h.service.GetListBusiness_v2(r.Context(), &req)
	if err != nil {
		log.WithError(err).Error("Error when get list business")
		return nil, err
	}

	return &ginext.Response{
		Code: http.StatusOK,
		GeneralBody: &ginext.GeneralBody{
			Data: rs.Data,
			Meta: rs.Meta,
		},
	}, nil
}

// GetOneBusiness
// @Tags Business
// @Security ApiKeyAuth
// @Summary Get one Business
// @Description Get details of a specific Business by ID then get staff with id
// @ID GetOneBusiness
// @Accept  json
// @Produce  json
// @Param id path string true "Business ID"
// @Success 200 {object} model.Business
// @Router /api/v1/business/get-one/{id} [get]
func (h *BusinessHandlers) GetOneBusiness(r *ginext.Request) (*ginext.Response, error) {
	ID := &uuid.UUID{}
	if ID = utils.ParseIDFromUri(r.GinCtx); ID == nil {
		return nil, ginext.NewError(http.StatusForbidden, "Wrong ID")
	}

	Business, err := h.service.GetOneBusiness(r.Context(), *ID)
	if err != nil {
		return nil, err
	}

	return ginext.NewResponseData(http.StatusOK, Business), nil
}

// GetOneBusiness_v2
// @Tags Business
// @Security ApiKeyAuth
// @Summary Get one Business
// @Description Get details of a specific Business by ID  and staffs with preloading
// @ID GetOneBusiness_v2
// @Accept  json
// @Produce  json
// @Param id path string true "Business ID"
// @Success 200 {object} model.Business
// @Router /api/v1/business/get-one-v2/{id} [get]
func (h *BusinessHandlers) GetOneBusiness_v2(r *ginext.Request) (*ginext.Response, error) {
	ID := &uuid.UUID{}
	if ID = utils.ParseIDFromUri(r.GinCtx); ID == nil {
		return nil, ginext.NewError(http.StatusForbidden, "Wrong ID")
	}

	Business, err := h.service.GetOneBusiness_v2(r.Context(), *ID)
	if err != nil {
		return nil, err
	}

	return ginext.NewResponseData(http.StatusOK, Business), nil
}

// DeleteBusiness
// @Tags Business
// @Security ApiKeyAuth
// @Summary Delete a Business
// @Description Delete a specific Business by ID
// @ID DeleteBusiness
// @Accept  json
// @Produce  json
// @Param id path string true "Business ID"
// @Router /api/v1/business/delete/{id} [delete]
func (h *BusinessHandlers) DeleteBusiness(r *ginext.Request) (*ginext.Response, error) {

	businessID := &uuid.UUID{}
	if businessID = utils.ParseIDFromUri(r.GinCtx); businessID == nil {
		return nil, ginext.NewError(http.StatusForbidden, "Wrong ID")
	}

	if err := h.service.DeleteBusiness(r.Context(), *businessID); err != nil {
		return nil, err
	}

	return ginext.NewResponse(http.StatusOK), nil
}
