package migrations

import (
	"context"
	"time"

	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	migrate.Register(
		func(db *mongo.Database) error {
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
		func(db *mongo.Database) error {
			// We don't need to down it

			return nil
		},
	)
}
