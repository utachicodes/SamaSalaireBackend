package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type PayrollPeriod struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	StartDate time.Time     `bson:"start_date" json:"start_date"`
	EndDate   time.Time     `bson:"end_date" json:"end_date"`
	Status    string        `bson:"status" json:"status"` // "open" | "closed"
}

type Payslip struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	EmployeeID  bson.ObjectID `bson:"employee_id" json:"employee_id"`
	PeriodID    bson.ObjectID `bson:"period_id" json:"period_id"`
	GrossPay    float64       `bson:"gross_pay" json:"gross_pay"`
	Deductions  float64       `bson:"deductions" json:"deductions"`
	NetPay      float64       `bson:"net_pay" json:"net_pay"`
	GeneratedAt time.Time     `bson:"generated_at" json:"generated_at"`
}
