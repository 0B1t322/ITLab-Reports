package app

import (
	"context"
	"time"

	_ "github.com/RTUITLab/ITLab-Reports/migrations/mongo"
	"github.com/sirupsen/logrus"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"

	stage "github.com/0B1t322/MongoBuilder/aggregation"
	"github.com/0B1t322/MongoBuilder/operators/aggregation"
	"github.com/0B1t322/MongoBuilder/operators/field"
	"github.com/0B1t322/MongoBuilder/operators/query"
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
	if err := migrate.Up(migrate.AllAvailable); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"from": "NewReportsRepository",
				"err":  err,
			},
		).Panic("Failed to migrate")
	}

	return nil
}

func (a *App) registerMigrations(db *mongo.Database) *migrate.Migrate {
	return migrate.NewMigrate(
		db,
		a.changeDateFromStringMigration(),
		a.addNameField(),
	)
}

func (a *App) changeDateFromStringMigration() migrate.Migration {
	return migrate.Migration{
		Version:     3,
		Description: "change-date-from-string-to-date",
		Up: func(db *mongo.Database) error {
			pipeline := []interface{}{
				bson.M{
					"$set": bson.M{
						"date": bson.M{
							"$cond": bson.A{
								bson.M{
									"$eq": bson.A{
										bson.M{
											"$type": "$date",
										},
										"string",
									},
								},
								bson.M{
									"$toDate": bson.M{
										"$concat": bson.A{
											"$date",
											"Z",
										},
									},
								},
								"$date",
							},
						},
					},
				},
				bson.M{
					"$out": "reports",
				},
			}
			ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
			_, err := db.Collection("reports").Aggregate(ctx, pipeline)
			if err != nil {
				return err
			}

			return nil
		},
		Down: func(db *mongo.Database) error {
			// We don't need to down it
			return nil
		},
	}
}

func (a *App) addNameField() migrate.Migration {
	return migrate.Migration{
		Version:     4,
		Description: "add_name_field",
		Up: func(db *mongo.Database) error {
			_, err := db.Collection("reports").Aggregate(
				context.Background(),
				[]bson.M{
					stage.Match(
						query.Regex(
							"text",
							`@\n\t\n@`,
						),
					),
					stage.AddFields(
						stage.AddFieldArg().
							AddField(
								"text",
								aggregation.Last(
									aggregation.Split(
										field.Field("text"),
										"@\n\t\n@",
									),
								),
							).
							AddField(
								"name",
								aggregation.First(
									aggregation.Split(
										field.Field("text"),
										"@\n\t\n@",
									),
								),
							),
					),
					stage.Merge(
						stage.MergeArg().
							IntoCollection("reports").
							WhenMatched(stage.WhenMatchedArg().Merge()),
					),
				},
			)
			return err
		},
		Down: func(db *mongo.Database) error {
			// We don't need to down it
			return nil
		},
	}
}

func (a *App) addState() migrate.Migration {
	return migrate.Migration{
		Version:     5,
		Description: "Add reports state field",
		Up: func(db *mongo.Database) error {
			return nil
		},
		Down: func(db *mongo.Database) error {
			return nil
		},
	}
}
