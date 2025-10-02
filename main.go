package main

import (
	"business/conf"
	_ "business/docs"
	"business/pkg/route"
	"business/pkg/utils"
	"context"
	"os"

	"gitlab.com/goxp/cloud0/logger"
)

const (
	APPNAME = "Template"
)

// @title template API
// @version 1.0
// @description This is Template api docs.
// @termsOfService http://swagger.io/terms/

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3333
// @BasePath  
func main() {
	conf.SetEnv()
	logger.Init(APPNAME)
	utils.LoadMessageError()
	app := route.NewService()
	ctx := context.Background()
	err := app.Start(ctx)
	if err != nil {
		logger.Tag("main").Error(err)
	}
	os.Clearenv()
}
