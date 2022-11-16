package app

import (
	"github.com/RTUITLab/ITLab-Reports/internal/services/salary"
	"github.com/RTUITLab/ITLab-Reports/internal/services/token"
	salaryGrpc "github.com/RTUITLab/ITLab/proto/salary/v1"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func (a *App) configureExternalServices() {
	// Configure external salary service
	a.configureSalaryService()
}

func (a *App) configureSalaryService() {
	if a.cfg.App.TestMode {
		do.Provide(
			a.injector,
			func(i *do.Injector) (salary.SalaryService, error) {
				return salary.NewTestModeSalaryService(), nil
			},
		)
	} else {
		conn, err := grpc.Dial(a.cfg.App.SalaryGRPCAddr, grpc.WithInsecure())
		if err != nil {
			logrus.Panicf("Failed to connect to salary service: %v", err)
		}

		client := salaryGrpc.NewApprovedReportsSalaryClient(conn)

		do.Provide(
			a.injector,
			func(i *do.Injector) (salary.SalaryService, error) {
				return salary.NewExternalGRPCSalaryService(client), nil
			},
		)
	}
}

func (a *App) configureTokenService() {
	if a.cfg.App.TestMode {
		do.Provide(
			a.injector,
			func(i *do.Injector) (token.TokenService, error) {
				return token.NewTestTokenService(), nil
			},
		)
	} else {
		do.Provide(
			a.injector,
			func(i *do.Injector) (token.TokenService, error) {
				return token.NewExternalTokenService(
					a.cfg.RemoteApi.ClientID,
					a.cfg.RemoteApi.ClientSecret,
					a.cfg.RemoteApi.TokenURL,
				)
			},
		)
	}
}
