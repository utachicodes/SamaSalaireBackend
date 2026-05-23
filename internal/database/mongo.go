package database

import (
	"context"
	"log"
	"time"

	"samasalaire-backend/internal/config"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect(cfg *config.Config) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("mongo: connect failed: %v", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalf("mongo: ping failed: %v", err)
	}

	log.Println("Connected to MongoDB successfully")
	return client
}

func GetDB(client *mongo.Client, dbName string) *mongo.Database {
	return client.Database(dbName)
}

func GetCollection(db *mongo.Database, name string) *mongo.Collection {
	return db.Collection(name)
}
