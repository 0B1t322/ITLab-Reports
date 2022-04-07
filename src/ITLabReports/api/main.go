package main

import (
	"ITLabReports/config"
	"ITLabReports/server"
	"fmt"
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
	app := &server.App{}
	app.Init(cfg)
	app.Run(":"+cfg.App.AppPort)
	fmt.Scanln()
}
