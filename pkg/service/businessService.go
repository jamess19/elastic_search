package service

import (
	"business/pkg/model"
	"business/pkg/repo"
	"business/pkg/utils"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gitlab.com/goxp/cloud0/ginext"
	"gitlab.com/goxp/cloud0/logger"
	"gorm.io/gorm"
)

type BusinessService struct {
	repo repo.PGInterface
}

func NewBusinessService(repo repo.PGInterface) BusinessInterface {
	return &BusinessService{repo: repo}
}

type BusinessInterface interface {
	CreateBusiness(ctx context.Context, req model.BusinessRequest) (*model.Business, error)
	CreateBusiness_v2(ctx context.Context) ([]model.Business, error)
	UpdateBusiness(ctx context.Context, req model.BusinessRequest) (*model.Business, error)
	GetOneBusiness(ctx context.Context, BusinessID uuid.UUID) (*model.Business, error)
	GetOneBusiness_v2(ctx context.Context, BusinessID uuid.UUID) (*model.Business, error)
	GetListBusiness(ctx context.Context, req *model.GetListBusinessRequest) (model.GetListBusinessResponse, error)
	GetListBusiness_v2(ctx context.Context, req *model.GetListBusinessRequest) (model.GetListBusinessResponse, error)
	DeleteBusiness(ctx context.Context, BusinessID uuid.UUID) error
}

func (s *BusinessService) CreateBusiness(ctx context.Context, req model.BusinessRequest) (*model.Business, error) {

	Business := &model.Business{}

	copier.Copy(Business, req)

	if err := s.repo.CreateBusiness(ctx, Business, nil); err != nil {
		return nil, err
	}

	return Business, nil
}

func (s *BusinessService) CreateBusiness_v2(ctx context.Context) ([]model.Business, error) {
	log := logger.WithCtx(ctx, "BusinessService.CreateBusiness_v2")
	BusinessList := make([]model.Business, 0, 10000)

	// Tạo 10.000 business ngẫu nhiên trực tiếp
	for i := 0; i < 10000; i++ {
		b := model.Business{
			Name:         randomdata.SillyName(),
			Description:  randomdata.Paragraph(),
			Address:      randomdata.Address(),
			BusinessType: randomdata.StringSample("type1", "type2", "type3"),
		}
		BusinessList = append(BusinessList, b)
	}

	business_chan := make(chan model.Business, 10000)
	done := make(chan bool, 5)
	start := time.Now().UnixMilli();
	for w := 1; w <= 20; w++ {
		worker_name := "worker" + fmt.Sprint(w)
		go func(worker string) {
			if err := s.repo.CreateBusiness_v2(ctx, business_chan, worker, done, nil); err != nil {
				logger.WithCtx(ctx, "BusinessService.CreateBusiness_v2").
					WithError(err).
					WithField("worker", worker).
					Error("Error in CreateBusiness_v2")
			}
			done <- true
		}(worker_name)
	}

	for _, b := range BusinessList {
		business_chan <- b
	}
	close(business_chan)

	// Chờ tất cả worker xong
	for i := 0; i < 20; i++ {
		<-done
	}
	duration := time.Now().UnixMilli() - start
	log.Info("Execution time: %d ms", duration)
	return BusinessList, nil
}

func (s *BusinessService) UpdateBusiness(ctx context.Context, req model.BusinessRequest) (*model.Business, error) {
	log := logger.WithCtx(ctx, "BusinessService.UpdateBusiness")

	Business, err := s.repo.GetOneBusiness(ctx, req.ID, nil)
	if err != nil {
		log.WithError(err).WithField("req", req).Error("Error get Business for updating")
		return nil, ginext.NewError(http.StatusForbidden, "Error get Business for updating")
	}

	copier.Copy(Business, req)

	if err := s.repo.UpdateBusiness(ctx, Business, nil); err != nil {
		log.WithError(err).WithField("req", req).Error("Error update Business")
		return nil, err
	}

	return Business, nil
}

func (s *BusinessService) GetListBusiness(ctx context.Context, req *model.GetListBusinessRequest) (model.GetListBusinessResponse, error) {
	log := logger.WithCtx(ctx, "BusinessService.GetListBusiness")

	res, err := s.repo.GetListBusiness(ctx, req, nil)
	if err != nil {
		log.WithError(err).Error("Error when call func GetListBusiness")
		return res, ginext.NewError(http.StatusInternalServerError, utils.MessageError()[http.StatusInternalServerError])
	}

	for i := range res.Data {
		staffs, err := s.repo.GetStaffByBusinessID(ctx, res.Data[i].ID, nil)
		if err != nil {
			log.WithError(err).Error("Error when call GetStaffByBusinessID")
			return model.GetListBusinessResponse{},
				ginext.NewError(http.StatusInternalServerError, utils.MessageError()[http.StatusInternalServerError])
		}
		res.Data[i].Staffs = staffs.Data
	}

	return res, nil
}

func (s *BusinessService) GetListBusiness_v2(ctx context.Context, req *model.GetListBusinessRequest) (model.GetListBusinessResponse, error) {
	log := logger.WithCtx(ctx, "BusinessService.GetListBusiness")

	res, err := s.repo.GetListBusiness_v2(ctx, req, nil)
	if err != nil {
		log.WithError(err).Error("Error when call func GetListBusiness")
		return res, ginext.NewError(http.StatusInternalServerError, utils.MessageError()[http.StatusInternalServerError])
	}

	return res, nil
}

func (s *BusinessService) DeleteBusiness(ctx context.Context, BusinessID uuid.UUID) error {
	log := logger.WithCtx(ctx, "BusinessService.DeleteBusiness")

	Business, err := s.repo.GetOneBusiness(ctx, BusinessID, nil)
	if err != nil {
		log.WithError(err).WithField("BusinessID", BusinessID).Error("Error when call func GetOneBusiness")
		return ginext.NewError(http.StatusNotFound, err.Error())
	}

	if err = s.repo.DeleteBusiness(ctx, Business, nil); err != nil {
		log.WithError(err).WithField("Business", Business).Error("Error when call func DeleteBusiness")
		return ginext.NewError(http.StatusInternalServerError, utils.MessageError()[http.StatusInternalServerError])
	}

	return nil
}

func (s *BusinessService) GetOneBusiness(ctx context.Context, BusinessID uuid.UUID) (*model.Business, error) {
	log := logger.WithCtx(ctx, "BusinessService.GetOneBusiness")

	Business, err := s.repo.GetOneBusiness(ctx, BusinessID, nil)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			{
				log.WithError(err).WithField("BusinessID", BusinessID).Error("Error when call func GetOneBusiness")
				return nil, ginext.NewError(http.StatusNotFound, utils.MessageError()[http.StatusNotFound])
			}
		default:
			{
				log.WithError(err).WithField("BusinessID", BusinessID).Error("Error when call func GetOneBusiness")
				return nil, ginext.NewError(http.StatusInternalServerError, utils.MessageError()[http.StatusInternalServerError])
			}
		}
	}

	res, err := s.repo.GetStaffByBusinessID(ctx, Business.ID, nil)
	Business.Staffs = res.Data

	return Business, nil
}

// Get with preloading
func (s *BusinessService) GetOneBusiness_v2(ctx context.Context, BusinessID uuid.UUID) (*model.Business, error) {
	log := logger.WithCtx(ctx, "BusinessService.GetOneBusiness_v2")

	Business, err := s.repo.GetOneBusiness_v2(ctx, BusinessID, nil)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			{
				log.WithError(err).WithField("BusinessID", BusinessID).Error("Error when call func GetOneBusiness_v2")
				return nil, ginext.NewError(http.StatusNotFound, utils.MessageError()[http.StatusNotFound])
			}
		default:
			{
				log.WithError(err).WithField("BusinessID", BusinessID).Error("Error when call func GetOneBusiness_v2")
				return nil, ginext.NewError(http.StatusInternalServerError, utils.MessageError()[http.StatusInternalServerError])
			}
		}
	}

	return Business, nil
}
