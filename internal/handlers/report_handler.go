package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"samasalaire-backend/internal/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ReportHandler struct {
	db *mongo.Database
}

func NewReportHandler(db *mongo.Database) *ReportHandler {
	return &ReportHandler{db: db}
}

func (h *ReportHandler) PayrollSummary(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	filter := bson.M{}
	if periodID := c.Query("period_id"); periodID != "" {
		pid, err := bson.ObjectIDFromHex(periodID)
		if err == nil {
			filter["period_id"] = pid
		}
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$period_id"},
			{Key: "total_gross", Value: bson.M{"$sum": "$gross_pay"}},
			{Key: "total_deductions", Value: bson.M{"$sum": "$deductions"}},
			{Key: "total_net", Value: bson.M{"$sum": "$net_pay"}},
			{Key: "employee_count", Value: bson.M{"$sum": 1}},
		}}},
	}

	cursor, err := h.db.Collection(database.ColPayslips).Aggregate(ctx, pipeline)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "aggregation failed")
		return
	}
	var results []bson.M
	cursor.All(ctx, &results)
	if results == nil {
		results = []bson.M{}
	}
	RespondOK(c, results)
}

func (h *ReportHandler) LeaveSummary(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	filter := bson.M{}
	if yearStr := c.Query("year"); yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err == nil {
			start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
			end := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
			filter["start_date"] = bson.M{"$gte": start, "$lt": end}
		}
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "leave_type_id", Value: "$leave_type_id"},
				{Key: "status", Value: "$status"},
			}},
			{Key: "count", Value: bson.M{"$sum": 1}},
		}}},
	}

	cursor, err := h.db.Collection(database.ColLeaveRequests).Aggregate(ctx, pipeline)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "aggregation failed")
		return
	}
	var results []bson.M
	cursor.All(ctx, &results)
	if results == nil {
		results = []bson.M{}
	}
	RespondOK(c, results)
}
