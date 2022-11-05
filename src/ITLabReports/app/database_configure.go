package app

import (
	"context"

	migrations "github.com/RTUITLab/ITLab-Reports/migrations/mongo/v2"
	"github.com/sirupsen/logrus"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

func (a *App) configureMongoDatabase() error {
	conn, err := mongo.Connect(
		context.Background(),
		options.Client().
			ApplyURI(a.cfg.MongoDB.URI),
	)
	if err != nil {
		return err
	}

	defer conn.Disconnect(context.Background())

	connStr, err := connstring.Parse(a.cfg.MongoDB.URI)
	if err != nil {
		return err
	}

	db := conn.Database(connStr.Database)
	migrate.SetDatabase(db)

	if err := a.registerMigrations(db).
		Up(migrate.AllAvailable); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"from": "NewReportsRepository",
				"err":  err,
			},
		).Panic("Failed to migrate")
	}

	return nil
}

type Migration interface {
	Version() uint64
	Description() string
	Up(db *mongo.Database) error
	Down(db *mongo.Database) error
}

func ToMigration(m Migration) migrate.Migration {
	return migrate.Migration{
		Version:     m.Version(),
		Description: m.Description(),
		Up: func(db *mongo.Database) error {
			return m.Up(db)
		},
		Down: func(db *mongo.Database) error {
			return m.Down(db)
		},
	}
}

func (a *App) registerMigrations(db *mongo.Database) *migrate.Migrate {
	return migrate.NewMigrate(
		db,
		ToMigration(&migrations.ChangeDateFromStringToDate{}),
		ToMigration(&migrations.AddNameFieldMigration{}),
		ToMigration(
			migrations.NewAddStateMigration(
				a.SalaryService,
				a.TokenService,
			),
		),
	)
}
