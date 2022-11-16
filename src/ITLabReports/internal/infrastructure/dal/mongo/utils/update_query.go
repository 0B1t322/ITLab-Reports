package utils

import (
	stage "github.com/0B1t322/MongoBuilder/aggregation"
	"github.com/0B1t322/MongoBuilder/operators/aggregation"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateQuery(model any) (query []bson.M) {
	query = append(
		query,
		stage.ReplaceRoot(
			aggregation.MergeObjects(
				"$ROOT",
				model,
			),
		),
	)
	return
}
