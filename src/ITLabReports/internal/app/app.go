package app

import (
	"fmt"
	"net/http"

	_ "github.com/RTUITLab/ITLab-Reports/docs"
	"github.com/RTUITLab/ITLab-Reports/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	cfg config.Config

	injector *do.Injector
}

func NewApp(cfg config.Config) *App {
	return &App{
		cfg:      cfg,
		injector: do.New(),
	}
}

func (a *App) configureDependencies() {
	a.configureExternalServices()
	a.configureDatabase()
	a.configureRepositories()
	a.configureServices()
}

func (a *App) RunHTTP() {
	a.configureDependencies()

	app := gin.New()

	if !a.cfg.App.TestMode {
		gin.SetMode(gin.ReleaseMode)
	}

	app.Use(
		gin.Recovery(),
	)

	root := app.Group("/api/reports")

	// init swagger handler
	root.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// init redirect to swagger
	root.GET(
		"/swagger",
		func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/api/reports/swagger/index.html")
		},
	)

	for _, controller := range a.configureHTTPControllers() {
		controller.Build(root)
	}

	logrus.Infof("Start HTTP application on port :%s", a.cfg.App.AppPort)
	if err := app.Run(fmt.Sprintf(":%s", a.cfg.App.AppPort)); err != nil {
		logrus.Panic("Failed to start HTTP application: ", err)
	}
}
