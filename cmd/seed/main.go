package main

import (
	"context"
	"log"
	"samasalaire-backend/internal/config"
	"samasalaire-backend/internal/database"
	"samasalaire-backend/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.Load()
	client := database.Connect(cfg)
	db := database.GetDB(client, cfg.DBName)
	ctx := context.Background()

	// Drop existing data
	for _, col := range []string{
		database.ColEmployees, database.ColUsers, database.ColSalaryComponents,
		database.ColPayrollPeriods, database.ColPayslips, database.ColLeaveTypes,
		database.ColLeaveBalances, database.ColLeaveRequests, database.ColAuditLogs,
	} {
		db.Collection(col).Drop(ctx)
	}

	database.CreateIndexes(db)

	// --- Leave types ---
	annual := models.LeaveType{ID: bson.NewObjectID(), Name: "Annual Leave", AnnualEntitlement: 21, AllowsNegative: false}
	sick := models.LeaveType{ID: bson.NewObjectID(), Name: "Sick Leave", AnnualEntitlement: 10, AllowsNegative: false}
	personal := models.LeaveType{ID: bson.NewObjectID(), Name: "Personal Leave", AnnualEntitlement: 5, AllowsNegative: false}
	insertMany(ctx, db.Collection(database.ColLeaveTypes), []interface{}{annual, sick, personal})

	// --- Employees ---
	hire := func(y, m, d int) time.Time { return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC) }

	adminEmp := employee("Admin User", "admin@samasalaire.sn", hire(2020, 1, 1), "System Administrator", "IT", nil, "admin")
	hrEmp := employee("HR Officer", "hr@samasalaire.sn", hire(2021, 3, 15), "HR Manager", "Human Resources", nil, "hr")
	mgrEmp := employee("Manager One", "manager@samasalaire.sn", hire(2019, 6, 1), "Team Manager", "Operations", nil, "manager")
	emp1 := employee("Alice Diallo", "alice@samasalaire.sn", hire(2022, 1, 10), "Software Engineer", "IT", &mgrEmp.ID, "employee")
	emp2 := employee("Bob Ndiaye", "bob@samasalaire.sn", hire(2022, 4, 5), "Accountant", "Finance", &hrEmp.ID, "employee")
	emp3 := employee("Claire Seck", "claire@samasalaire.sn", hire(2023, 7, 20), "Operations Analyst", "Operations", &mgrEmp.ID, "employee")

	employees := []interface{}{adminEmp, hrEmp, mgrEmp, emp1, emp2, emp3}
	insertMany(ctx, db.Collection(database.ColEmployees), employees)

	// --- Users ---
	users := []interface{}{
		user(adminEmp.ID, "admin", "admin123", "admin"),
		user(hrEmp.ID, "hr", "hr123", "hr"),
		user(mgrEmp.ID, "manager", "manager123", "manager"),
		user(emp1.ID, "alice", "alice123", "employee"),
		user(emp2.ID, "bob", "bob123", "employee"),
		user(emp3.ID, "claire", "claire123", "employee"),
	}
	insertMany(ctx, db.Collection(database.ColUsers), users)

	// --- Salary components ---
	components := []interface{}{
		salaryComp(emp1.ID, "Base Salary", "earning", 450000),
		salaryComp(emp1.ID, "Housing Allowance", "earning", 75000),
		salaryComp(emp1.ID, "Income Tax", "deduction", 45000),
		salaryComp(emp1.ID, "Social Security", "deduction", 22500),
		salaryComp(emp2.ID, "Base Salary", "earning", 380000),
		salaryComp(emp2.ID, "Transport Allowance", "earning", 30000),
		salaryComp(emp2.ID, "Income Tax", "deduction", 38000),
		salaryComp(emp2.ID, "Social Security", "deduction", 19000),
		salaryComp(emp3.ID, "Base Salary", "earning", 320000),
		salaryComp(emp3.ID, "Meal Allowance", "earning", 25000),
		salaryComp(emp3.ID, "Income Tax", "deduction", 30000),
		salaryComp(emp3.ID, "Social Security", "deduction", 16000),
		salaryComp(mgrEmp.ID, "Base Salary", "earning", 650000),
		salaryComp(mgrEmp.ID, "Management Bonus", "earning", 100000),
		salaryComp(mgrEmp.ID, "Income Tax", "deduction", 75000),
		salaryComp(mgrEmp.ID, "Social Security", "deduction", 32500),
		salaryComp(hrEmp.ID, "Base Salary", "earning", 500000),
		salaryComp(hrEmp.ID, "Housing Allowance", "earning", 60000),
		salaryComp(hrEmp.ID, "Income Tax", "deduction", 56000),
		salaryComp(hrEmp.ID, "Social Security", "deduction", 28000),
	}
	insertMany(ctx, db.Collection(database.ColSalaryComponents), components)

	// --- Leave balances ---
	allEmployees := []models.Employee{adminEmp, hrEmp, mgrEmp, emp1, emp2, emp3}
	leaveTypes := []models.LeaveType{annual, sick, personal}
	var balances []interface{}
	for _, emp := range allEmployees {
		for _, lt := range leaveTypes {
			balances = append(balances, models.LeaveBalance{
				ID:            bson.NewObjectID(),
				EmployeeID:    emp.ID,
				LeaveTypeID:   lt.ID,
				RemainingDays: lt.AnnualEntitlement,
			})
		}
	}
	insertMany(ctx, db.Collection(database.ColLeaveBalances), balances)

	log.Println("Seed complete!")
	log.Println("Test credentials:")
	log.Println("  admin / admin123")
	log.Println("  hr    / hr123")
	log.Println("  manager / manager123")
	log.Println("  alice   / alice123")
	log.Println("  bob     / bob123")
	log.Println("  claire  / claire123")
}

func employee(name, email string, hireDate time.Time, jobTitle, dept string, managerID *bson.ObjectID, role string) models.Employee {
	return models.Employee{
		ID:         bson.NewObjectID(),
		Name:       name,
		Email:      email,
		HireDate:   hireDate,
		JobTitle:   jobTitle,
		Department: dept,
		ManagerID:  managerID,
		Role:       role,
		IsActive:   true,
	}
}

func user(employeeID bson.ObjectID, username, password, role string) models.User {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("bcrypt error: %v", err)
	}
	return models.User{
		ID:           bson.NewObjectID(),
		EmployeeID:   employeeID,
		Username:     username,
		PasswordHash: string(hash),
		Role:         role,
	}
}

func salaryComp(employeeID bson.ObjectID, name, compType string, amount float64) models.SalaryComponent {
	return models.SalaryComponent{
		ID:         bson.NewObjectID(),
		EmployeeID: employeeID,
		Name:       name,
		Type:       compType,
		Amount:     amount,
	}
}

func insertMany(ctx context.Context, col *mongo.Collection, docs []interface{}) {
	if _, err := col.InsertMany(ctx, docs); err != nil {
		log.Printf("InsertMany error on %s: %v", col.Name(), err)
	}
}
