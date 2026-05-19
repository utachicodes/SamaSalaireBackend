package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AuditLog struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    bson.ObjectID `bson:"user_id" json:"user_id"`
	Action    string        `bson:"action" json:"action"`
	Target    string        `bson:"target" json:"target"`
	Timestamp time.Time     `bson:"timestamp" json:"timestamp"`
}
