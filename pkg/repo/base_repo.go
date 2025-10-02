package repo

import (
	"business/pkg/model"
	"business/pkg/utils"
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"runtime/debug"
	"time"

	"gitlab.com/goxp/cloud0/logger"

	"github.com/google/uuid"
	"gitlab.com/goxp/cloud0/ginext"
	"gorm.io/gorm"
)

const (
	StateNew byte = iota + 1 // starts from 1
	StateDoing
	StateDone

	generalQueryTimeout         = 60 * time.Second
	generalQueryTimeout2Minutes = 120 * time.Second
	defaultPageSize             = 30
	maxPageSize                 = 1000
)

func NewPGRepo(db *gorm.DB) PGInterface {
	return &RepoPG{db: db}
}

type PGInterface interface {
	// DB
	DBWithTimeout(ctx context.Context) (*gorm.DB, context.CancelFunc)
	DB() (db *gorm.DB)
	Transaction(ctx context.Context, f func(rp PGInterface) error) error

	// Business methods
	CreateBusiness(ctx context.Context, business *model.Business, tx *gorm.DB) error
	CreateBusiness_v2(ctx context.Context, business_chan <- chan model.Business, worker_name string, done chan <-bool, tx *gorm.DB) error
	GetListBusiness(ctx context.Context, req *model.GetListBusinessRequest, tx *gorm.DB) (rs model.GetListBusinessResponse, err error)
	GetListBusiness_v2(ctx context.Context, req *model.GetListBusinessRequest, tx *gorm.DB) (rs model.GetListBusinessResponse, err error)
	GetOneBusiness(ctx context.Context, businessID uuid.UUID, tx *gorm.DB) (rs *model.Business, err error)
	GetOneBusiness_v2(ctx context.Context, businessId uuid.UUID, tx *gorm.DB) (rs *model.Business, err error)
	UpdateBusiness(ctx context.Context, business *model.Business, tx *gorm.DB) error
	DeleteBusiness(ctx context.Context, business *model.Business, tx *gorm.DB) error

	// Staff methods
	CreateStaff(ctx context.Context, staff *model.Staff, tx *gorm.DB) error
	GetOneStaff(ctx context.Context, staffID uuid.UUID, tx *gorm.DB) (*model.Staff, error)
	GetListStaff(ctx context.Context, req *model.GetListStaffRequest, tx *gorm.DB) (model.GetListStaffResponse, error)
	GetListStaffWithPaging(ctx context.Context, req *model.GetListStaffRequest, tx *gorm.DB) (model.GetListStaffResponse, error)
	UpdateStaff(ctx context.Context, staff *model.Staff, tx *gorm.DB) error
	DeleteStaff(ctx context.Context, staff *model.Staff, tx *gorm.DB) error
	GetStaffByBusinessID(ctx context.Context, businessID uuid.UUID, tx *gorm.DB) (model.GetListStaffResponse, error)
}

type RepoPG struct {
	db    *gorm.DB
	debug bool
}

func (r *RepoPG) Transaction(ctx context.Context, f func(rp PGInterface) error) (err error) {
	log := logger.WithCtx(ctx, "RepoPG.Transaction")
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()
	// create new instance to run the transaction
	repo := *r
	tx = tx.Begin()
	repo.db = tx
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = errors.New(fmt.Sprint(r))
			log.WithError(err).Error("error_500: Panic when run Transaction")
			debug.PrintStack()
			return
		}
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
	err = f(&repo)
	if err != nil {
		log.WithError(err).Error("error_500: Error when run Transaction")
		return err
	}
	return nil
}

func (r *RepoPG) DB() *gorm.DB {
	return r.db
}

func (r *RepoPG) DBWithTimeout(ctx context.Context) (*gorm.DB, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(ctx, generalQueryTimeout)
	return r.db.WithContext(ctx), cancel
}

func (r *RepoPG) DBWithTimeout2Minutes(ctx context.Context) (*gorm.DB, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(ctx, generalQueryTimeout2Minutes)
	return r.db.WithContext(ctx), cancel
}

func (r *RepoPG) GetPage(page int) int {
	if page == 0 {
		return 1
	}
	return page
}

func (r *RepoPG) GetOffset(page int, pageSize int) int {
	return (page - 1) * pageSize
}

func (r *RepoPG) GetPageSize(pageSize int) int {
	if pageSize == 0 {
		return defaultPageSize
	}
	if pageSize > maxPageSize {
		return maxPageSize
	}
	return pageSize
}

func (r *RepoPG) GetTotalPages(totalRows, pageSize int) int {
	return int(math.Ceil(float64(totalRows) / float64(pageSize)))
}

func (r *RepoPG) GetOrder(sort string) string {
	if sort == "" {
		sort = "created_at desc"
	}
	return sort
}

func (r *RepoPG) GetOrderBy(sort string) string {
	if sort == "" {
		sort = "created_at desc"
	}
	return sort
}

func (r *RepoPG) GetPaginationInfo(query string, tx *gorm.DB, totalRow, page, pageSize int) (rs ginext.BodyMeta, err error) {
	tm := struct {
		Count int `json:"count"`
	}{}
	if query != "" {
		if err = tx.Raw(query).Scan(&tm).Error; err != nil {
			return nil, err
		}
		totalRow = tm.Count
	}

	return ginext.BodyMeta{
		"page":        page,
		"page_size":   pageSize,
		"total_pages": r.GetTotalPages(totalRow, pageSize),
		"total_rows":  totalRow,
	}, nil
}

func (r *RepoPG) ReturnErrorInGetFuncV2(ctx context.Context, logStr string, err error, key string, value interface{}) error {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0)).WithField(key, value)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.WithError(err).Error(fmt.Sprintf("error_404: %s - RepoPG", logStr))
		return ginext.NewError(http.StatusNotFound, utils.MessageError()[http.StatusNotFound])
	}
	log.WithError(err).Error(fmt.Sprintf("error_500: %s - RepoPG", logStr))
	return ginext.NewError(http.StatusInternalServerError, err.Error())
}
