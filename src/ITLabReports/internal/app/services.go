package app

import (
	"time"

	reportsadapters "github.com/RTUITLab/ITLab-Reports/internal/adapters/reports/service"
	drafts "github.com/RTUITLab/ITLab-Reports/internal/domain/drafts/service"
	reports "github.com/RTUITLab/ITLab-Reports/internal/domain/reports/service"
	user "github.com/RTUITLab/ITLab-Reports/internal/domain/user/service"
	"github.com/samber/do"
)

func (a *App) configureServices() {
	a.configureDomainServices()
}

func (a *App) configureDomainServices() {
	// User service
	do.Provide(
		a.injector,
		func(i *do.Injector) (user.UserService, error) {
			return user.NewJWKUserService(
				user.JWKAuthServiceConfig{
					UserRole:       a.cfg.Auth.Roles.User,
					AdminRole:      a.cfg.Auth.Roles.Admin,
					SuperAdminRole: a.cfg.Auth.Roles.SuperAdmin,
					RoleClaim:      a.cfg.Auth.Audience,
					ScopeClaim:     "scope",
					UserIDClaim:    "sub",
					KeyURL:         a.cfg.Auth.KeyURL,
					RequiredScope:  a.cfg.Auth.Scope,
					RefreshTime:    time.Minute,
				},
			), nil
		},
	)

	// Drafts service
	do.Provide(
		a.injector,
		func(i *do.Injector) (drafts.DraftsService, error) {
			return drafts.NewDraftsServiceImplFrom(i)
		},
	)

	// Reports service

	// Provide reports adapters
	do.Provide(
		a.injector,
		func(i *do.Injector) (reports.DraftService, error) {
			return reportsadapters.NewReportsDraftServiceFrom(i)
		},
	)

	// Provide reports service
	do.Provide(
		a.injector,
		func(i *do.Injector) (reports.ReportsService, error) {
			return reports.NewReportsServiceImplFrom(i)
		},
	)
}
