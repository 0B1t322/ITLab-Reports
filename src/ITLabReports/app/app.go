package app

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/RTUITLab/ITLab-Reports/config"
	_ "github.com/RTUITLab/ITLab-Reports/docs"
	"github.com/RTUITLab/ITLab-Reports/pkg/adapters/toidchecker"
	"github.com/RTUITLab/ITLab-Reports/service/idvalidator"
	"github.com/RTUITLab/ITLab-Reports/service/reports"
	"github.com/RTUITLab/ITLab-Reports/service/salary"
	"github.com/RTUITLab/ITLab-Reports/transport/middlewares"
	"github.com/RTUITLab/ITLab-Reports/transport/report"
	pb "github.com/RTUITLab/ITLab/proto/reports/v1"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	swag "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
)

type App struct {
	Router        *mux.Router
	GRPCServer    *grpc.Server
	Auther        middlewares.Auther
	IdChecker     middlewares.IdChecker
	SalaryService salary.SalaryService

	ReportService reports.Service
	DraftService  reports.Service

	ReportEndpoints report.Endpoints
	DraftEndpoints  report.Endpoints

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
	// Build salary service
	app.buildSalaryService()
	// Build services
	if err := app.BuildReportService(); err != nil {
		log.Panicf("Failed to build app: %v", err)
	}

	if err := app.BuildDraftService(); err != nil {
		log.Panicf("Failed to build app: %v", err)
	}

	app.BuildReportEndpoints()
	app.BuildDraftEndpoints()

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

func (a *App) buildSalaryService() {
	if !a.cfg.App.TestMode {
		a.SalaryService = salary.NewExternalRestSalaryService(
			a.cfg.App.ITLabURL,
			nil,
		)
	} else {
		a.SalaryService = salary.NewTestModeSalaryService()
	}
}

func (a *App) BuildHTTP() {
	draftHTTPEndpoints := a.BuildDraftsHTTPV1(a.DraftEndpoints)

	a.BuildReportsHTTPV1(a.ReportEndpoints, ToDraftService(draftHTTPEndpoints))
	a.BuildReportsHTTPV2(a.ReportEndpoints)

	docs := a.Router.PathPrefix("/reports/swagger")
	docs.Handler(
		swag.WrapHandler,
	)
}

func (a *App) BuildGRPC() {
	grpcServer := grpc.NewServer()

	pb.RegisterReportsServer(grpcServer, a.BuildReportsGRPCV1(a.ReportEndpoints))

	a.GRPCServer = grpcServer
}

func (a *App) StartGRPC() {
	a.BuildGRPC()
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", a.cfg.App.GrpcAppPort))
	if err != nil {
		log.Panicf("Failed to start grpc application %v", err)
	}

	log.Infof("Starting grpc Application is port %s", a.cfg.App.GrpcAppPort)
	if err := a.GRPCServer.Serve(lis); err != nil {
		log.Panicf("Failed to start grpc application %v", err)
	}
}

func (a *App) StartHTTP() {
	a.BuildHTTP()

	log.Infof("Starting http Application is port %s", a.cfg.App.AppPort)
	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", a.cfg.App.AppPort),
		Handler:        a.Router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		IdleTimeout:    2 * time.Second,
	}
	if err := s.ListenAndServe(); err != nil {
		log.Panicf("Failed to start http application %v", err)
	}

}
