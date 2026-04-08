package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func NewMongo(ctx context.Context, uri string, dbName string) (*Mongo, error) {

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	ctxPing, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Ping(ctxPing, nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)

	return &Mongo{
		Client: client,
		DB:     db,
	}, nil
}
