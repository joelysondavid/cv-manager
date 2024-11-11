package db

import (
	"context"

	"github.com/joelysondavid/cv-manager/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New(databaseName string) (database *mongo.Database, err error) {
	env := config.GetEnv()
	clientOptions := options.Client().ApplyURI(env.DBUrl)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		return nil, err
	}

	db := *client.Database(env.DBName)
	return &db, nil
}
