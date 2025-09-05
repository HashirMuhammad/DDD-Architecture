package config

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	URI        string
	Database   string
	Collection string
}

func NewMongoConfig() *MongoConfig {
	return &MongoConfig{
		URI:        "mongodb://localhost:27017/",
		Database:   "UserServiceDB",
		Collection: "users",
	}
}

func (c *MongoConfig) Connect() (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(c.URI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(c.Database)
	return db, nil
}

func (c *MongoConfig) GetConnectionURI() string {
	return c.URI
}
