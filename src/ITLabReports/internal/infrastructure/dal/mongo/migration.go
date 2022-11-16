package dal

import (
	"github.com/samber/lo"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	AllAvailable = migrate.AllAvailable
)

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

func RegisterMigrations(db *mongo.Database, migrations ...Migration) *migrate.Migrate {
	return migrate.NewMigrate(
		db,
		lo.Map(
			migrations,
			func(m Migration, _ int) migrate.Migration {
				return ToMigration(m)
			},
		)...,
	)
}
