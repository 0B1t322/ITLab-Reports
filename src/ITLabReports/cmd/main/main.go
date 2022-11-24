package main

import (
	"github.com/RTUITLab/ITLab-Reports/internal/app"
	"github.com/RTUITLab/ITLab-Reports/internal/config"
)

func main() {
	config.InitGlobalConfig()

	app := app.NewApp(config.GlobalConfig)
	app.ConfigureDependencies()
	app.ConfigureSharedControllersOptions()
	app.ConfigureGRPCControllerOptions()
	app.ConfigureHTTPControllerOptions()

	go app.RunGRPC()
	go app.RunHTTP()
	select {}
}
