package main

import (
	"github.com/RTUITLab/ITLab-Reports/app"
	"github.com/RTUITLab/ITLab-Reports/config"
)

// @title ITLab-Reports API
// @version 1.0
// @description This is a server to work with reports
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	cfg := config.GetConfig()
	app := app.New(
		cfg,
	)
	app.StartHTTP()
}