package migrations

import (
	"context"

	stage "github.com/0B1t322/MongoBuilder/aggregation"
	"github.com/0B1t322/MongoBuilder/operators/aggregation"
	"github.com/0B1t322/MongoBuilder/operators/field"
	"github.com/0B1t322/MongoBuilder/operators/query"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AddNameFieldMigration struct{}

func (a *AddNameFieldMigration) Version() uint64 {
	return 4
}

func (a *AddNameFieldMigration) Description() string {
	return "add_name_field"
}

func (m *AddNameFieldMigration) Up(db *mongo.Database) error {
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
}

func (m *AddNameFieldMigration) Down(db *mongo.Database) error {
	// We don't need to down it
	return nil
}
