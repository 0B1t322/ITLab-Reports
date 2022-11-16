package app

import (
	drafts "github.com/RTUITLab/ITLab-Reports/internal/domain/drafts/repository"
	reports "github.com/RTUITLab/ITLab-Reports/internal/domain/reports/repository"
	daldrafts "github.com/RTUITLab/ITLab-Reports/internal/infrastructure/dal/mongo/drafts"
	dalreports "github.com/RTUITLab/ITLab-Reports/internal/infrastructure/dal/mongo/reports"
	"github.com/samber/do"
)

func (a *App) configureRepositories() {
	// Reports repository
	do.Provide(
		a.injector,
		func(i *do.Injector) (reports.ReportRepository, error) {
			return dalreports.NewMongoReportsRepositoryFrom(i)
		},
	)

	// Drafts repository
	do.Provide(
		a.injector,
		func(i *do.Injector) (drafts.DraftRepository, error) {
			return daldrafts.NewMongoDraftRepositoryFrom(i)
		},
	)
}
