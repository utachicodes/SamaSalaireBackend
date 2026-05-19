package services

import (
	"context"
	"fmt"
	"time"

	"samasalaire-backend/internal/database"
	"samasalaire-backend/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func RunPayroll(db *mongo.Database, periodID bson.ObjectID) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var period models.PayrollPeriod
	if err := db.Collection(database.ColPayrollPeriods).FindOne(ctx, bson.M{"_id": periodID}).Decode(&period); err != nil {
		return 0, fmt.Errorf("payroll period not found")
	}
	if period.Status != "open" {
		return 0, fmt.Errorf("payroll period is not open")
	}

	cursor, err := db.Collection(database.ColEmployees).Find(ctx, bson.M{"is_active": true})
	if err != nil {
		return 0, fmt.Errorf("failed to fetch employees")
	}
	var employees []models.Employee
	if err := cursor.All(ctx, &employees); err != nil {
		return 0, fmt.Errorf("failed to decode employees")
	}

	count := 0
	for _, emp := range employees {
		compCursor, err := db.Collection(database.ColSalaryComponents).Find(ctx, bson.M{"employee_id": emp.ID})
		if err != nil {
			continue
		}
		var components []models.SalaryComponent
		compCursor.All(ctx, &components)

		var gross, deductions float64
		for _, comp := range components {
			if comp.Type == "earning" {
				gross += comp.Amount
			} else {
				deductions += comp.Amount
			}
		}

		payslip := models.Payslip{
			ID:          bson.NewObjectID(),
			EmployeeID:  emp.ID,
			PeriodID:    periodID,
			GrossPay:    gross,
			Deductions:  deductions,
			NetPay:      gross - deductions,
			GeneratedAt: time.Now(),
		}

		// Upsert: remove existing payslip for this period+employee, then insert
		db.Collection(database.ColPayslips).DeleteOne(ctx, bson.M{
			"employee_id": emp.ID,
			"period_id":   periodID,
		})
		db.Collection(database.ColPayslips).InsertOne(ctx, payslip)
		count++
	}

	return count, nil
}
