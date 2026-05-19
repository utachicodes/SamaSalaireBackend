package models

import "go.mongodb.org/mongo-driver/v2/bson"

type SalaryComponent struct {
	ID         bson.ObjectID `bson:"_id,omitempty" json:"id"`
	EmployeeID bson.ObjectID `bson:"employee_id" json:"employee_id"`
	Name       string        `bson:"name" json:"name"`
	Type       string        `bson:"type" json:"type"` // "earning" | "deduction"
	Amount     float64       `bson:"amount" json:"amount"`
}
