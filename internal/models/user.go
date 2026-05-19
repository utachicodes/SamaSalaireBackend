package models

import "go.mongodb.org/mongo-driver/v2/bson"

type User struct {
	ID           bson.ObjectID `bson:"_id,omitempty" json:"id"`
	EmployeeID   bson.ObjectID `bson:"employee_id" json:"employee_id"`
	Username     string        `bson:"username" json:"username"`
	PasswordHash string        `bson:"password_hash" json:"-"`
	Role         string        `bson:"role" json:"role"`
}
