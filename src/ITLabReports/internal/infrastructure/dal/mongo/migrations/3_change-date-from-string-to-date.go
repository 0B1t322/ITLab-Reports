package migrations

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChangeDateFromStringToDate struct{}

func (a *ChangeDateFromStringToDate) Version() uint64 {
	return 3
}

func (a *ChangeDateFromStringToDate) Description() string {
	return "change-date-from-string-date"
}

func (a *ChangeDateFromStringToDate) Up(db *mongo.Database) error {
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
}

func (a *ChangeDateFromStringToDate) Down(db *mongo.Database) error {
	// We don't need to down it

	return nil
}
