package config

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

func MongoClient() (*mongo.Database, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return &mongo.Database{}, err
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return &mongo.Database{}, err
	}

	database := client.Database("testing")
	return database, err
}
