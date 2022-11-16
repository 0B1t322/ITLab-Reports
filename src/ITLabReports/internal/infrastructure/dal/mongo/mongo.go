package dal

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

func ConnectToMongoDB(dsn string) (*mongo.Database, error) {
	c, err := mongo.NewClient(
		options.Client().
			ApplyURI(dsn),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create mongo client: %v", err)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Minute,
	)
	defer cancel()
	if err := c.Connect(ctx); err != nil {
		return nil, fmt.Errorf("Failed to connect to mongo: %v", err)
	}

	constr, err := connstring.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse mongo uri: %v", err)
	}

	db := c.Database(constr.Database)

	return db, nil
}
