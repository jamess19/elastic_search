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

type StaffHandlers struct {
	service service.StaffInterface
}

func NewStaffHandler(service service.StaffInterface) *StaffHandlers {
	return &StaffHandlers{service: service}
}

// CreateStaff
// @Tags Staff
// @Security ApiKeyAuth
// @Summary Create new staff
// @Description create a new staff for system
// @ID CreateStaff
// @Accept  json
// @Produce  json
// @Param data body model.StaffRequest true "body data"
// @Success 200 {object} model.Staff
// @Router /api/v1/staff/create [post]
func (h *StaffHandlers) CreateStaff(r *ginext.Request) (*ginext.Response, error) {
	req := model.StaffRequest{}
	r.MustBind(&req)

	if err := common.CheckRequireValid(req); err != nil {
		return nil, ginext.NewError(http.StatusBadRequest, err.Error())
	}

	rs, err := h.service.CreateStaff(r.GinCtx, req)
	if err != nil {
		return nil, err
	}

	return ginext.NewResponseData(http.StatusCreated, rs), nil
}

// UpdateStaff
// @Tags Staff
// @Security ApiKeyAuth
// @Summary update Staff
// @Description update Staff
// @ID UpdateStaff
// @Accept  json
// @Produce  json
// @Param data body model.StaffRequest true "body data"
// @Success 200 {object} model.Staff
// @Router /api/v1/staff/update/{id} [put]
func (h *StaffHandlers) UpdateStaff(r *ginext.Request) (*ginext.Response, error) {

	req := model.StaffRequest{}
	// if req.ID = utils.ParseIDFromUri(r.GinCtx); req.ID == nil {
	// 	return nil, ginext.NewError(http.StatusForbidden, "Wrong ID")
	// }

	r.MustBind(&req)

	rs, err := h.service.UpdateStaff(r.GinCtx, req)
	if err != nil {
		return nil, err
	}
	return ginext.NewResponseData(http.StatusOK, rs), nil
}

// ListStaff
// @Tags Staff
// @Security ApiKeyAuth
// @Summary List Staffs
// @Description Get a list of Staffs
// @ID ListStaff
// @Accept  json
// @Produce  json
// @Success 200 {object} []model.Staff
// @Router /api/v1/staff/get-list [get]
func (h *StaffHandlers) ListStaff(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.GinCtx, "ListStaff")

	var req model.GetListStaffRequest
	r.MustBind(&req)

	rs, err := h.service.GetListStaff(r.Context(), &req)
	if err != nil {
		log.WithError(err).Error("Error when get list Staff")
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

// GetOneStaff
// @Tags Staff
// @Security ApiKeyAuth
// @Summary Get one Staff
// @Description Get details of a specific Staff by ID
// @ID GetOneStaff
// @Accept  json
// @Produce  json
// @Param id path string true "Staff ID"
// @Success 200 {object} model.Staff
// @Router /api/v1/staff/get-one/{id} [get]
func (h *StaffHandlers) GetOneStaff(r *ginext.Request) (*ginext.Response, error) {
	ID := &uuid.UUID{}
	if ID = utils.ParseIDFromUri(r.GinCtx); ID == nil {
		return nil, ginext.NewError(http.StatusForbidden, "Wrong ID")
	}

	Staff, err := h.service.GetOneStaff(r.Context(), *ID)
	if err != nil {
		return nil, err
	}

	return ginext.NewResponseData(http.StatusOK, Staff), nil
}

// DeleteStaff
// @Tags Staff
// @Security ApiKeyAuth
// @Summary Delete a Staff
// @Description Delete a specific Staff by ID
// @ID DeleteStaff
// @Accept  json
// @Produce  json
// @Param id path string true "Staff ID"
// @Router /api/v1/staff/delete/{id} [delete]
func (h *StaffHandlers) DeleteStaff(r *ginext.Request) (*ginext.Response, error) {

	StaffID := &uuid.UUID{}
	if StaffID = utils.ParseIDFromUri(r.GinCtx); StaffID == nil {
		return nil, ginext.NewError(http.StatusForbidden, "Wrong ID")
	}

	if err := h.service.DeleteStaff(r.Context(), *StaffID); err != nil {
		return nil, err
	}

	return ginext.NewResponse(http.StatusOK), nil
}

// ListStaffWithPaging
// @Tags Staff
// @Security ApiKeyAuth
// @Summary Get List with paging
// Description Get a list of staff with paging
// @ID ListStaffWithPaging
// Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param sort query string false "Sort by column, e.g. 'staff.created_at desc'" 
// @Param keyword query string false "search by name,..." 
// @Router /api/v1/staff/get-list-paging [get]
func (h *StaffHandlers) ListStaffWithPaging(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.GinCtx, "GetListWithPaging")
	var req = model.GetListStaffRequest{}
	r.MustBind(&req)
	
	res, err := h.service.GetListStaffWithPaging(r.Context(),&req)
	if(err!=nil) {
		log.WithError(err).Error("Error when get list staff")
		return nil, err
	}
	return &ginext.Response{
		Code: http.StatusOK,
		GeneralBody: &ginext.GeneralBody{
			Data: res.Data,
			Meta: res.Meta,
		},
	}, nil
}

