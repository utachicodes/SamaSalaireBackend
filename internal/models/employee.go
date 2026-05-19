package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Employee struct {
	ID         bson.ObjectID  `bson:"_id,omitempty" json:"id"`
	Name       string         `bson:"name" json:"name"`
	Email      string         `bson:"email" json:"email"`
	HireDate   time.Time      `bson:"hire_date" json:"hire_date"`
	JobTitle   string         `bson:"job_title" json:"job_title"`
	Department string         `bson:"department" json:"department"`
	ManagerID  *bson.ObjectID `bson:"manager_id,omitempty" json:"manager_id,omitempty"`
	Role       string         `bson:"role" json:"role"`
	IsActive   bool           `bson:"is_active" json:"is_active"`
}
