package service

import (
	"context"
	"business/pkg/repo"
	"business/pkg/model"
	"business/pkg/utils"
	"github.com/google/uuid"
	"gitlab.com/goxp/cloud0/ginext"
	"gitlab.com/goxp/cloud0/logger"
	"gorm.io/gorm"
	"net/http"
	"github.com/jinzhu/copier"

)

type StaffService struct {
	repo repo.PGInterface
}

func NewStaffService(repo repo.PGInterface) StaffInterface {
	return &StaffService{repo: repo}
}

type StaffInterface interface {
	CreateStaff(ctx context.Context, req model.StaffRequest) (*model.Staff, error)
	UpdateStaff(ctx context.Context, req model.StaffRequest) (*model.Staff, error)
	GetOneStaff(ctx context.Context, StaffID uuid.UUID) (*model.Staff, error)
	GetListStaff(ctx context.Context, req *model.GetListStaffRequest) (model.GetListStaffResponse, error)
	GetListStaffWithPaging(ctx context.Context, req *model.GetListStaffRequest) (model.GetListStaffResponse, error)
	DeleteStaff(ctx context.Context, StaffID uuid.UUID) error
}

func (s *StaffService) CreateStaff(ctx context.Context, req model.StaffRequest) (*model.Staff, error) {
	Staff := &model.Staff{}

	copier.Copy(Staff,req)

	if err:= s.repo.CreateStaff(ctx,Staff,nil); err!=nil {
		return nil, err
	}
	
	return Staff, nil
}

func (s *StaffService) 	UpdateStaff(ctx context.Context, req model.StaffRequest) (*model.Staff, error) {
	log := logger.WithCtx(ctx, "BusinessService.UpdateBusiness")
	Staff ,err := s.repo.GetOneStaff(ctx, req.ID, nil)
	
	if err != nil {
		log.WithError(err).WithField("req", req).Error("Error get Staff for updating")
		return nil, ginext.NewError(http.StatusForbidden, "Error get Business for updating")
	}

	copier.Copy(Staff,req)

	if err:= s.repo.UpdateStaff(ctx, Staff, nil); err!=nil {
		log.WithError(err).WithField("req",req).Error("Error update Staff")
	}
	return Staff, nil
}

func (s *StaffService) GetOneStaff(ctx context.Context, StaffID uuid.UUID) (*model.Staff, error) {
	log := logger.WithCtx(ctx, "BusinessService.GetOneStaff")

	Staff, err := s.repo.GetOneStaff(ctx, StaffID, nil)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			{
				log.WithError(err).WithField("StaffID", StaffID).Error("Error when call func GetOneStaff")
				return nil, ginext.NewError(http.StatusNotFound, utils.MessageError()[http.StatusNotFound])
			}
		default:
			{
				log.WithError(err).WithField("BusinessID", StaffID).Error("Error when call func GetOneBusiness")
				return nil, ginext.NewError(http.StatusInternalServerError, utils.MessageError()[http.StatusInternalServerError])
			}
		}
	}

	return Staff, nil

}

func (s *StaffService) GetListStaff(ctx context.Context, req *model.GetListStaffRequest) (model.GetListStaffResponse, error) {
	log := logger.WithCtx(ctx, "StaffService.GetListStaff")

	res, err:= s.repo.GetListStaff(ctx,req,nil); 
	if (err!=nil) {
		log.WithError(err).Error("Error when call getListStaff")
		return res, ginext.NewError(http.StatusInternalServerError, utils.MessageError()[http.StatusInternalServerError])
	}

	return res, nil

}

func (s *StaffService) GetListStaffWithPaging(ctx context.Context, req *model.GetListStaffRequest) (model.GetListStaffResponse, error) {
	log := logger.WithCtx(ctx, "StaffService.GetListStaffWithPaging")

	rs, err := s.repo.GetListStaffWithPaging(ctx, req, nil)
	if(err!= nil) {
		log.WithError(err).Error("Error when call GetListStaffWithPaging")
		return rs, ginext.NewError(http.StatusInternalServerError, utils.MessageError()[http.StatusInternalServerError])
	}

	return rs, nil
}
func (s *StaffService) DeleteStaff(ctx context.Context, StaffID uuid.UUID) error {
	log:= logger.WithCtx(ctx, "StaffService.DeleteStaff")

	Staff, err := s.repo.GetOneStaff(ctx, StaffID, nil);
	if(err!=nil) {
		log.WithError(err).WithField("StaffID", StaffID).Error("Error when con fun GetOneStaff")
		return ginext.NewError(http.StatusNotFound, err.Error())
	}

	if err = s.repo.DeleteStaff(ctx, Staff, nil); err != nil {
		log.WithError(err).WithField("Staff", Staff).Error("Error when call func DeleteBusiness")
		return ginext.NewError(http.StatusInternalServerError, utils.MessageError()[http.StatusInternalServerError])
	}

	return nil
}
