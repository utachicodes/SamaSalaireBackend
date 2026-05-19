package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type LeaveType struct {
	ID                bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name              string        `bson:"name" json:"name"`
	AnnualEntitlement float64       `bson:"annual_entitlement" json:"annual_entitlement"`
	AllowsNegative    bool          `bson:"allows_negative" json:"allows_negative"`
}

type LeaveBalance struct {
	ID            bson.ObjectID `bson:"_id,omitempty" json:"id"`
	EmployeeID    bson.ObjectID `bson:"employee_id" json:"employee_id"`
	LeaveTypeID   bson.ObjectID `bson:"leave_type_id" json:"leave_type_id"`
	RemainingDays float64       `bson:"remaining_days" json:"remaining_days"`
}

type LeaveRequest struct {
	ID          bson.ObjectID  `bson:"_id,omitempty" json:"id"`
	EmployeeID  bson.ObjectID  `bson:"employee_id" json:"employee_id"`
	LeaveTypeID bson.ObjectID  `bson:"leave_type_id" json:"leave_type_id"`
	StartDate   time.Time      `bson:"start_date" json:"start_date"`
	EndDate     time.Time      `bson:"end_date" json:"end_date"`
	Status      string         `bson:"status" json:"status"` // "draft"|"pending"|"approved"|"rejected"
	Reason      string         `bson:"reason" json:"reason"`
	DecidedBy   *bson.ObjectID `bson:"decided_by,omitempty" json:"decided_by,omitempty"`
	DecidedAt   *time.Time     `bson:"decided_at,omitempty" json:"decided_at,omitempty"`
}
