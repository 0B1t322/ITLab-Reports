package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/RTUITLab/ITLab-Reports/config"
	_ "github.com/RTUITLab/ITLab-Reports/docs"
	"github.com/RTUITLab/ITLab-Reports/pkg/adapters/toidchecker"
	"github.com/RTUITLab/ITLab-Reports/service/idvalidator"
	"github.com/RTUITLab/ITLab-Reports/transport/middlewares"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	swag "github.com/swaggo/http-swagger"
)

type App struct {
	Router    *mux.Router
	Auther    middlewares.Auther
	IdChecker middlewares.IdChecker

	cfg *config.Config
}

func New(cfg *config.Config) *App {
	app := &App{
		cfg: cfg,
	}

	app.Router = mux.NewRouter().PathPrefix("/api").Subrouter()

	// Set log level
	if !cfg.App.TestMode {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}

	// Build auth middleware
	app.buildAuthMiddleware()

	// Build idChecker
	app.buildIdChecker()

	return app
}

func (a *App) buildAuthMiddleware() {
	if !a.cfg.App.TestMode {
		a.Auther = middlewares.NewJWKSAuth(
			middlewares.WithAdminRole(a.cfg.Auth.Roles.Admin),
			middlewares.WithUserRole(a.cfg.Auth.Roles.User),
			middlewares.WithSuperAdminRole(a.cfg.Auth.Roles.SuperAdmin),
			middlewares.WithJWKSUrl(a.cfg.Auth.KeyURL),
			middlewares.WithRoleClaim(a.cfg.Auth.Audience),
		)
	} else {
		a.Auther = middlewares.NewTestAuth(
			middlewares.WithAdminRole(a.cfg.Auth.Roles.Admin),
			middlewares.WithUserRole(a.cfg.Auth.Roles.User),
			middlewares.WithSuperAdminRole(a.cfg.Auth.Roles.SuperAdmin),
			middlewares.WithJWKSUrl(a.cfg.Auth.KeyURL),
			middlewares.WithRoleClaim(a.cfg.Auth.Audience),
		)
	}
}

func (a *App) buildIdChecker() {
	if !a.cfg.App.TestMode {
		a.IdChecker = toidchecker.ToIdChecker(
			idvalidator.New(
				idvalidator.ExternalRestIDValidator(
					a.cfg.App.ITLabURL,
					nil, // can be nil
				),
			),
		)
	} else {
		a.IdChecker = toidchecker.ToIdChecker(
			idvalidator.New(
				idvalidator.AlwaysTrueIdValidator(),
			),
		)
	}
}

func (a *App) BuildDraftHTTP() (DraftEndpoints, error) {
	s, err := a.BuildDraftService()
	if err != nil {
		return DraftEndpoints{}, err
	}

	e := a.BuildReportsEndpoints(s)

	draftEnds := a.BuildDraftsHTTPV1(e)

	return draftEnds, nil
}

func (a *App) BuildReportsHTTP(d DraftEndpoints) error {
	s, err := a.BuildReportService()
	if err != nil {
		return err
	}

	e := a.BuildReportsEndpoints(s)

	a.BuildReportsHTTPV1(e, ToDraftService(d))
	a.BuildReportsHTTPV2(e)
	return nil
}

func (a *App) BuildHTTP() error {
	draft, err := a.BuildDraftHTTP()
	if err != nil {
		return err
	}

	if err := a.BuildReportsHTTP(draft); err != nil {
		return err
	}

	docs := a.Router.PathPrefix("/reports/swagger")
	docs.Handler(
		swag.WrapHandler,
	)

	return nil
}

func (a *App) StartHTTP() {
	if err := a.BuildHTTP(); err != nil {
		log.Panicf("Failed to start application %v", err)
	}

	log.Infof("Starting Application is port %s", a.cfg.App.AppPort)
	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", a.cfg.App.AppPort),
		Handler:        a.Router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		IdleTimeout:    2 * time.Second,
	}
	if err := s.ListenAndServe(); err != nil {
		log.Panicf("Failed to start application %v", err)
	}

}
