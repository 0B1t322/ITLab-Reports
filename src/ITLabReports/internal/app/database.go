package app

import (
	dal "github.com/RTUITLab/ITLab-Reports/internal/infrastructure/dal/mongo"
	"github.com/RTUITLab/ITLab-Reports/internal/infrastructure/dal/mongo/migrations"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"
)

func (a *App) configureDatabase() {
	a.configureMongo()
}

func (a *App) configureMongo() {
	db, err := dal.ConnectToMongoDB(a.cfg.MongoDB.URI)
	if err != nil {
		logrus.Panic("Failed to connect to database: ", err)
	}

	if err := dal.RegisterMigrations(
		db,
		&migrations.ChangeDateFromStringToDate{},
		&migrations.AddNameFieldMigration{},
		// TODO: Add this migration after salary service will be ready
		// migrations.NewAddStateMigration(
		// 	do.MustInvoke[salary.SalaryService](a.injector),
		// 	do.MustInvoke[token.TokenService](a.injector),
		// ),
	).Up(dal.AllAvailable); err != nil {
		logrus.Panic("Failed to migrate database: ", err)
	}

	do.ProvideValue(
		a.injector,
		db,
	)
}
