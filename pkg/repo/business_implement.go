package repo

import (
	"business/pkg/model"
	"context"
	"net/http"

	"github.com/google/uuid"
	"gitlab.com/goxp/cloud0/ginext"
	"gitlab.com/goxp/cloud0/logger"
	"gorm.io/gorm"
)

func (r *RepoPG) CreateBusiness(ctx context.Context, business *model.Business, tx *gorm.DB) error {
	log := logger.WithCtx(ctx, "RepoPG.CreateBusiness")

	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}

	if err := tx.Create(business).Error; err != nil {
		log.WithError(err).Error("Error when call func CreateBusiness")
		return ginext.NewError(http.StatusInternalServerError, "Error when run query create Business")
	}

	return nil
}


func (r *RepoPG) CreateBusiness_v2(ctx context.Context, business_chan <- chan model.Business, worker_name string, done chan <-bool, tx *gorm.DB) error {
	log := logger.WithCtx(ctx, "RepoPG.CreateBusiness_v2")

	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}
	batch := []model.Business{}

	for bus := range business_chan {
		bus.WorkerName = worker_name
		batch = append(batch, bus)
		if len(batch) == 10 {
			if err := tx.Create(batch).Error; err != nil {
				log.WithError(err).Error("Error when call func CreateBusiness")
				return ginext.NewError(http.StatusInternalServerError, "Error when run query create Business")
			}
			batch = batch[:0]
		}
	} 
	if len(batch) > 0 {
        if err := tx.Create(batch).Error; err != nil {
				log.WithError(err).Error("Error when call func CreateBusiness")
				return ginext.NewError(http.StatusInternalServerError, "Error when run query create Business")
		}
    }
    done <- true

	return nil
}

// Get one business 
func (r *RepoPG) GetOneBusiness(ctx context.Context, businessId uuid.UUID, tx *gorm.DB) (rs *model.Business, err error) {
	var cancel context.CancelFunc

	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}
	if err := tx.First(&rs, businessId).Error; err != nil {
		return rs, err
	}


	return rs, nil
}

// Get one business
func (r *RepoPG) GetOneBusiness_v2(ctx context.Context, businessId uuid.UUID, tx *gorm.DB) (rs *model.Business, err error) {
	var cancel context.CancelFunc

	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}

	if err := tx.Preload("Staffs").First(&rs, businessId).Error; err != nil {
		return rs, err
	}
	return rs, nil
}

func (r RepoPG) UpdateBusiness(ctx context.Context, business *model.Business, tx *gorm.DB) error {
	log := logger.WithCtx(ctx, "RepoPG.UpdateBusiness")

	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}

	if err := tx.WithContext(ctx).Where("id = ?", business.ID).Save(&business).Error; err != nil {
		log.WithError(err).Errorf("Error when call func UpdateBusiness")
		return ginext.NewError(http.StatusInternalServerError, "Error when run query update Business")
	}

	return nil
}

func (r *RepoPG) DeleteBusiness(ctx context.Context, business *model.Business, tx *gorm.DB) error {
	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}
	return tx.WithContext(ctx).Delete(&business).Error
}

func (r *RepoPG) GetListBusiness(ctx context.Context, req *model.GetListBusinessRequest, tx *gorm.DB) (rs model.GetListBusinessResponse, err error) {
	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}

	page := r.GetPage(req.Page)

	pageSize := r.GetPageSize(req.PageSize)
	if req.Sort == "" {
		req.Sort = "Business.created_at desc"
	}

	tx = tx.WithContext(ctx).Model(&model.Business{})

	if req.Name != nil {
		tx = tx.Where("name = ?", req.Name)
	}

	if req.ManagerID != nil {
		tx = tx.Where("manager_id = ?", req.ManagerID)
	}

	if req.Address != nil {
		tx = tx.Where("address = ?", req.Address)
	}

	var total int64

	// Get list bussiness
	if err := tx.Count(&total).Select("business.*").Limit(pageSize).Offset(r.GetOffset(page, pageSize)).
		Order(r.GetOrderBy(req.Sort)).Find(&rs.Data).Error; err != nil {
		return rs, err
	}

	// Pagination
	if rs.Meta, err = r.GetPaginationInfo("", tx, int(total), page, pageSize); err != nil {
		return rs, err
	}

	return rs, nil
}

func (r *RepoPG) GetListBusiness_v2(ctx context.Context, req *model.GetListBusinessRequest, tx *gorm.DB) (rs model.GetListBusinessResponse, err error) {
	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}

	page := r.GetPage(req.Page)

	pageSize := r.GetPageSize(req.PageSize)
	if req.Sort == "" {
		req.Sort = "Business.created_at desc"
	}

	tx = tx.WithContext(ctx).Model(&model.Business{})

	if req.Name != nil {
		tx = tx.Where("name = ?", req.Name)
	}

	if req.ManagerID != nil {
		tx = tx.Where("manager_id = ?", req.ManagerID)
	}

	if req.Address != nil {
		tx = tx.Where("address = ?", req.Address)
	}

	var total int64

	// Get list bussiness
	if err := tx.Count(&total).Select("business.*").Limit(pageSize).Offset(r.GetOffset(page, pageSize)).
		Order(r.GetOrderBy(req.Sort)).Preload("Staffs").Find(&rs.Data).Error; err != nil {
		return rs, err
	}

	// Pagination
	if rs.Meta, err = r.GetPaginationInfo("", tx, int(total), page, pageSize); err != nil {
		return rs, err
	}

	return rs, nil
}
