package repo

import (
	"business/pkg/model"
	"business/pkg/utils"
	"context"
	"net/http"

	"github.com/google/uuid"
	"gitlab.com/goxp/cloud0/ginext"
	"gitlab.com/goxp/cloud0/logger"
	"gorm.io/gorm"
)

func (r *RepoPG) CreateStaff(ctx context.Context, staff *model.Staff, tx *gorm.DB) error {
	log := logger.WithCtx(ctx, "RepoPG.CreateStaff")

	db := r.db
	if tx != nil {
		db = tx
	} else {
		var cancel context.CancelFunc
		db, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}

	if err := db.Create(staff).Error; err != nil {
		log.WithError(err).WithField("staff", staff).Error("Error when create staff")
		return ginext.NewError(http.StatusInternalServerError, utils.MessageError()[http.StatusInternalServerError])
	}

	return nil
}

func (r *RepoPG) GetOneStaff(ctx context.Context, staffID uuid.UUID, tx *gorm.DB) (*model.Staff, error) {
	log := logger.WithCtx(ctx, "RepoPG.GetOneStaff")

	db := r.db
	if tx != nil {
		db = tx
	} else {
		var cancel context.CancelFunc
		db, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}

	staff := &model.Staff{}
	if err := db.Where("id = ?", staffID).First(staff).Error; err != nil {
		log.WithError(err).WithField("staffID", staffID).Error("Error when get one staff")
		return nil, r.ReturnErrorInGetFuncV2(ctx, "GetOneStaff", err, "staffID", staffID)
	}

	return staff, nil
}

func (r *RepoPG) GetListStaff(ctx context.Context, req *model.GetListStaffRequest, tx *gorm.DB) (model.GetListStaffResponse, error) {
	log := logger.WithCtx(ctx, "RepoPG.GetListStaff")

	db := r.db
	if tx != nil {
		db = tx
	} else {
		var cancel context.CancelFunc
		db, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}

	var staffs []model.Staff
	var rs model.GetListStaffResponse

	if err := db.Find(&staffs).Error; err != nil {
		log.WithError(err).Error("Error when get list staff")
		return rs, ginext.NewError(http.StatusInternalServerError, utils.MessageError()[http.StatusInternalServerError])
	}

	rs.Data = staffs
	rs.Meta = nil

	return rs, nil
}

func (r *RepoPG) GetListStaffWithPaging(ctx context.Context, req *model.GetListStaffRequest, tx *gorm.DB) (rs model.GetListStaffResponse, err error) {
	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}
	page := r.GetPage(req.Page)
	pageSize := r.GetPageSize(req.PageSize)

	if req.Sort == "" {
		req.Sort = "staff.created_at desc"
	}

	tx = tx.WithContext(ctx).Model(&model.Staff{})

	var total int64

	// query va sort
	if err := tx.Count(&total).Select("staff.*").Where("fullname ILIKE ?", "%"+req.Keyword+"%").Limit(pageSize).Offset(r.GetOffset(page, pageSize)).
		Order(r.GetOrderBy(req.Sort)).Find(&rs.Data).Error; err != nil {
		return rs, err
	}

	// pagination
	if rs.Meta, err = r.GetPaginationInfo("", tx, int(total), page, pageSize); err != nil {
		return rs, err
	}
	return rs, nil
}

func (r *RepoPG) UpdateStaff(ctx context.Context, staff *model.Staff, tx *gorm.DB) error {
	log := logger.WithCtx(ctx, "RepoPG.UpdateStaff")

	db := r.db
	if tx != nil {
		db = tx
	} else {
		var cancel context.CancelFunc
		db, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}

	if err := db.Save(staff).Error; err != nil {
		log.WithError(err).WithField("staff", staff).Error("Error when update staff")
		return ginext.NewError(http.StatusInternalServerError, utils.MessageError()[http.StatusInternalServerError])
	}

	return nil
}

func (r *RepoPG) DeleteStaff(ctx context.Context, staff *model.Staff, tx *gorm.DB) error {
	log := logger.WithCtx(ctx, "RepoPG.DeleteStaff")

	db := r.db
	if tx != nil {
		db = tx
	} else {
		var cancel context.CancelFunc
		db, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}

	if err := db.Delete(staff).Error; err != nil {
		log.WithError(err).WithField("staff", staff).Error("Error when delete staff")
		return ginext.NewError(http.StatusInternalServerError, utils.MessageError()[http.StatusInternalServerError])
	}

	return nil
}

func (r *RepoPG) GetStaffByBusinessID(ctx context.Context, businessID uuid.UUID, tx *gorm.DB) (model.GetListStaffResponse, error) {
	log := logger.WithCtx(ctx, "RepoPG.GetStaffByBusinessID")

	db := r.db
	if tx != nil {
		db = tx
	} else {
		var cancel context.CancelFunc
		db, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}

	var staffs []model.Staff
	var rs model.GetListStaffResponse
	if err := db.Where("business_id = ?", businessID).Find(&staffs).Error; err != nil {
		log.WithError(err).Error("Error when get list staff with ID")
		return rs, ginext.NewError(http.StatusInternalServerError, utils.MessageError()[http.StatusInternalServerError])
	}

	rs.Data = staffs
	rs.Meta = nil

	return rs, nil
}
