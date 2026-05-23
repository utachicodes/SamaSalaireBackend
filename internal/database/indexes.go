package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func CreateIndexes(db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	createUniqueIndex(ctx, db.Collection(ColUsers), bson.D{{Key: "username", Value: 1}})
	createUniqueIndex(ctx, db.Collection(ColEmployees), bson.D{{Key: "email", Value: 1}})
	createUniqueIndex(ctx, db.Collection(ColLeaveBalances), bson.D{{Key: "employee_id", Value: 1}, {Key: "leave_type_id", Value: 1}})
	createUniqueIndex(ctx, db.Collection(ColPayslips), bson.D{{Key: "employee_id", Value: 1}, {Key: "period_id", Value: 1}})

	log.Println("mongo: indexes ensured")
}

func createUniqueIndex(ctx context.Context, col *mongo.Collection, keys bson.D) {
	_, err := col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    keys,
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		log.Printf("Warning: could not create index on %s: %v", col.Name(), err)
	}
}
