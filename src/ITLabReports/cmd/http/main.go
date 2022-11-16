package main

import (
	"github.com/RTUITLab/ITLab-Reports/internal/app"
	"github.com/RTUITLab/ITLab-Reports/internal/config"
)

// @title ITLab-Reports API
// @version 2.0
// @description This is a server to work with reports
// @description.markdown
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	config.InitGlobalConfig()

	app := app.NewApp(config.GlobalConfig)

	app.RunHTTP()
}
