package handlers

import (
	"context"
	"net/http"
	"time"

	"samasalaire-backend/internal/database"
	"samasalaire-backend/internal/models"
	"samasalaire-backend/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PayrollHandler struct {
	db *mongo.Database
}

func NewPayrollHandler(db *mongo.Database) *PayrollHandler {
	return &PayrollHandler{db: db}
}

func (h *PayrollHandler) ListPeriods(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := h.db.Collection(database.ColPayrollPeriods).Find(ctx, bson.M{})
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to fetch periods")
		return
	}
	var periods []models.PayrollPeriod
	cursor.All(ctx, &periods)
	if periods == nil {
		periods = []models.PayrollPeriod{}
	}
	RespondOK(c, periods)
}

func (h *PayrollHandler) CreatePeriod(c *gin.Context) {
	var period models.PayrollPeriod
	if err := c.ShouldBindJSON(&period); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if !period.StartDate.Before(period.EndDate) {
		RespondError(c, http.StatusBadRequest, "start_date must be strictly before end_date")
		return
	}
	period.ID = bson.NewObjectID()
	period.Status = "open"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := h.db.Collection(database.ColPayrollPeriods).InsertOne(ctx, period); err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to create period")
		return
	}
	RespondCreated(c, period)
}

func (h *PayrollHandler) RunPayroll(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	count, err := services.RunPayroll(h.db, id)
	if err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	RespondOK(c, gin.H{
		"message":            "payroll run completed",
		"payslips_generated": count,
	})
}

func (h *PayrollHandler) FinalizePeriod(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := h.db.Collection(database.ColPayrollPeriods).UpdateOne(ctx,
		bson.M{"_id": id, "status": "open"},
		bson.M{"$set": bson.M{"status": "closed"}},
	)
	if err != nil || result.MatchedCount == 0 {
		RespondError(c, http.StatusBadRequest, "payroll period not found or already finalized")
		return
	}
	RespondOK(c, gin.H{"message": "payroll period finalized"})
}

func (h *PayrollHandler) ListPayslips(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}
	role, _ := c.Get("role")
	if role == "employee" {
		empIDStr, _ := c.Get("employeeID")
		empID, _ := bson.ObjectIDFromHex(empIDStr.(string))
		filter["employee_id"] = empID
	}
	if periodID := c.Query("period_id"); periodID != "" {
		pid, err := bson.ObjectIDFromHex(periodID)
		if err == nil {
			filter["period_id"] = pid
		}
	}

	cursor, err := h.db.Collection(database.ColPayslips).Find(ctx, filter)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to fetch payslips")
		return
	}
	var payslips []models.Payslip
	cursor.All(ctx, &payslips)
	if payslips == nil {
		payslips = []models.Payslip{}
	}
	RespondOK(c, payslips)
}

func (h *PayrollHandler) GetPayslip(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	role, _ := c.Get("role")
	if role == "employee" {
		empIDStr, _ := c.Get("employeeID")
		empID, _ := bson.ObjectIDFromHex(empIDStr.(string))
		filter["employee_id"] = empID
	}

	var payslip models.Payslip
	if err := h.db.Collection(database.ColPayslips).FindOne(ctx, filter).Decode(&payslip); err != nil {
		RespondError(c, http.StatusNotFound, "payslip not found")
		return
	}
	RespondOK(c, payslip)
}
