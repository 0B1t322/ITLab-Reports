package migrations

import (
	"context"

	stage "github.com/0B1t322/MongoBuilder/aggregation"
	"github.com/0B1t322/MongoBuilder/operators/aggregation"
	"github.com/0B1t322/MongoBuilder/operators/query"
	"github.com/0B1t322/MongoBuilder/operators/field"

	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	migrate.Register(
		func(db *mongo.Database) error {
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
		func(db *mongo.Database) error {
			// We don't need to down it
			return nil
		},
	)
}