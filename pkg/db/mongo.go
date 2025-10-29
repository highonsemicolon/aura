package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBRegistry map[string]*mongo.Database

func InitMongoRegistry(ctx context.Context, uri string, dbs map[string]string) (DBRegistry, *mongo.Client, error) {
	opts := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	registry := DBRegistry{}

	for alias, name := range dbs {
		registry[alias] = client.Database(name)
	}

	return registry, client, nil
}
