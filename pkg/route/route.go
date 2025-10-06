package route

import (
	"business/pkg/handlers"
	"business/pkg/middleware"
	"business/pkg/repo"
	"business/pkg/es"
	service2 "business/pkg/service"

	"github.com/caarlos0/env/v6"
	swaggerFiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"

	"gitlab.com/goxp/cloud0/ginext"
	"gitlab.com/goxp/cloud0/service"
)

type extraSetting struct {
	DbDebugEnable bool `env:"DB_DEBUG_ENABLE" envDefault:"true"`
}

type Service struct {
	*service.BaseApp
	setting *extraSetting
}

func NewService() *Service {
	s := &Service{
		service.NewApp("MVT Adapter", "v1.0"),
		&extraSetting{},
	}

	// repo
	_ = env.Parse(s.setting)
	db := s.GetDB()
	if s.setting.DbDebugEnable {
		db = db.Debug()
	}
	repoPG := repo.NewPGRepo(db)
	esConfig := es.Config{
		Addresses: []string{"http://localhost:9200"},
		Username: "elastic",
		Password: "VAMQ+bNBKdC5bu3Ee-1y",
	}
	client, err := es.NewClient(esConfig)
	if err != nil {
		panic(err)
	}
	// service
	businessService := service2.NewBusinessService(repoPG)
	staffService := service2.NewStaffService(repoPG)
	esService := service2.NewEsService(repoPG, client)
	// handle
	businessHandle := handlers.NewBusinessHandlers(businessService)
	staffHandle := handlers.NewStaffHandler(staffService)
	esHandle := handlers.NewElasticHandlers(esService)

	// Áp dụng CORS middleware cho toàn bộ router
	s.Router.Use(middleware.CORSMiddleware())

	v1Api := s.Router.Group("/api/v1")
	swaggerApi := s.Router.Group("/")

	// swagger
	swaggerApi.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))

	// Khởi tạo rate limiter
	// rateLimiter := middleware.NewRateLimiter(redisClient, 5*time.Second, 100)

	// Áp dụng rate limiter cho tất cả các route
	// v1Api.Use(rateLimiter.RateLimit())

	v1Api.POST("/business/create", middleware.LoggingRequest(), ginext.WrapHandler(businessHandle.CreateBusiness)) // only admin portal
	v1Api.POST("/business/create-v2", middleware.LoggingRequest(), ginext.WrapHandler(businessHandle.CreateBusiness_v2)) // only admin portal
	v1Api.GET("/business/get-one/:id", ginext.WrapHandler(businessHandle.GetOneBusiness))
	v1Api.GET("/business/get-one-v2/:id", ginext.WrapHandler(businessHandle.GetOneBusiness_v2))
	v1Api.GET("/business/get-list", ginext.WrapHandler(businessHandle.ListBusiness))
	v1Api.GET("/business/get-list-v2", ginext.WrapHandler(businessHandle.ListBusiness_v2))
	v1Api.PUT("/business/update/:id", ginext.WrapHandler(businessHandle.UpdateBusiness))    // only admin portal
	v1Api.DELETE("/business/delete/:id", ginext.WrapHandler(businessHandle.DeleteBusiness)) // only admin portal

	v1Api.POST("/staff/create", middleware.LoggingRequest(), ginext.WrapHandler(staffHandle.CreateStaff)) // only admin portal
	v1Api.GET("/staff/get-one/:id", ginext.WrapHandler(staffHandle.GetOneStaff))
	v1Api.GET("/staff/get-list", ginext.WrapHandler(staffHandle.ListStaff))
	v1Api.PUT("/staff/update/:id", ginext.WrapHandler(staffHandle.UpdateStaff))    // only admin portal
	v1Api.DELETE("/staff/delete/:id", ginext.WrapHandler(staffHandle.DeleteStaff)) // only admin portal
	v1Api.GET("/staff/get-list-paging", ginext.WrapHandler(staffHandle.ListStaffWithPaging))

	v1Api.POST("/elastic/push-to-elastic", ginext.WrapHandler(esHandle.PushToElastic))
	v1Api.POST("/elastic/search-by-field", ginext.WrapHandler(esHandle.SearchByField))
	v1Api.POST("/elastic/fulltext-search", ginext.WrapHandler(esHandle.FullTextSearch))

	
	// Migrate
	migrateHandler := handlers.NewMigrationHandler(db)
	s.Router.POST("/internal/migrate", migrateHandler.Migrate)
	return s
}
