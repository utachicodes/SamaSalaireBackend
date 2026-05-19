package handlers

import (
	"context"
	"net/http"
	"time"

	"samasalaire-backend/internal/database"
	"samasalaire-backend/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type LeaveHandler struct {
	db *mongo.Database
}

func NewLeaveHandler(db *mongo.Database) *LeaveHandler {
	return &LeaveHandler{db: db}
}

// ---- Leave Types ----

func (h *LeaveHandler) ListLeaveTypes(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, _ := h.db.Collection(database.ColLeaveTypes).Find(ctx, bson.M{})
	var types []models.LeaveType
	cursor.All(ctx, &types)
	if types == nil {
		types = []models.LeaveType{}
	}
	RespondOK(c, types)
}

func (h *LeaveHandler) CreateLeaveType(c *gin.Context) {
	var lt models.LeaveType
	if err := c.ShouldBindJSON(&lt); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	lt.ID = bson.NewObjectID()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := h.db.Collection(database.ColLeaveTypes).InsertOne(ctx, lt); err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to create leave type")
		return
	}
	RespondCreated(c, lt)
}

func (h *LeaveHandler) UpdateLeaveType(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var body bson.M
	if err := c.ShouldBindJSON(&body); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	delete(body, "_id")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h.db.Collection(database.ColLeaveTypes).UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": body})
	var updated models.LeaveType
	h.db.Collection(database.ColLeaveTypes).FindOne(ctx, bson.M{"_id": id}).Decode(&updated)
	RespondOK(c, updated)
}

// ---- Leave Balances ----

func (h *LeaveHandler) GetBalance(c *gin.Context) {
	empID, err := bson.ObjectIDFromHex(c.Param("employeeId"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid employee id")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, _ := h.db.Collection(database.ColLeaveBalances).Find(ctx, bson.M{"employee_id": empID})
	var balances []models.LeaveBalance
	cursor.All(ctx, &balances)
	if balances == nil {
		balances = []models.LeaveBalance{}
	}
	RespondOK(c, balances)
}

// ---- Leave Requests ----

func (h *LeaveHandler) ListRequests(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}
	role, _ := c.Get("role")
	empIDStr, _ := c.Get("employeeID")
	empID, _ := bson.ObjectIDFromHex(empIDStr.(string))

	switch role {
	case "employee":
		filter["employee_id"] = empID
	case "manager":
		// find all direct reports of this manager
		reportCursor, err := h.db.Collection(database.ColEmployees).Find(ctx, bson.M{"manager_id": empID})
		if err != nil {
			RespondError(c, http.StatusInternalServerError, "failed to fetch reports")
			return
		}
		var reports []models.Employee
		reportCursor.All(ctx, &reports)
		reportIDs := make([]bson.ObjectID, len(reports))
		for i, r := range reports {
			reportIDs[i] = r.ID
		}
		filter["employee_id"] = bson.M{"$in": reportIDs}
	// hr and admin see all — no filter
	}

	if status := c.Query("status"); status != "" {
		filter["status"] = status
	}

	cursor, _ := h.db.Collection(database.ColLeaveRequests).Find(ctx, filter)
	var requests []models.LeaveRequest
	cursor.All(ctx, &requests)
	if requests == nil {
		requests = []models.LeaveRequest{}
	}
	RespondOK(c, requests)
}

func (h *LeaveHandler) CreateRequest(c *gin.Context) {
	var req models.LeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	empIDStr, _ := c.Get("employeeID")
	empID, _ := bson.ObjectIDFromHex(empIDStr.(string))
	req.ID = bson.NewObjectID()
	req.EmployeeID = empID
	req.Status = "pending"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check balance
	days := req.EndDate.Sub(req.StartDate).Hours()/24 + 1

	var balance models.LeaveBalance
	err := h.db.Collection(database.ColLeaveBalances).FindOne(ctx, bson.M{
		"employee_id":   empID,
		"leave_type_id": req.LeaveTypeID,
	}).Decode(&balance)

	if err == nil {
		var lt models.LeaveType
		h.db.Collection(database.ColLeaveTypes).FindOne(ctx, bson.M{"_id": req.LeaveTypeID}).Decode(&lt)
		if !lt.AllowsNegative && balance.RemainingDays < days {
			RespondError(c, http.StatusBadRequest, "insufficient leave balance")
			return
		}
	}

	if _, err := h.db.Collection(database.ColLeaveRequests).InsertOne(ctx, req); err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to create leave request")
		return
	}
	RespondCreated(c, req)
}

type decideBody struct {
	Status string `json:"status" binding:"required"`
	Note   string `json:"note"`
}

func (h *LeaveHandler) DecideRequest(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var body decideBody
	if err := c.ShouldBindJSON(&body); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if body.Status != "approved" && body.Status != "rejected" {
		RespondError(c, http.StatusBadRequest, "status must be approved or rejected")
		return
	}

	deciderIDStr, _ := c.Get("employeeID")
	deciderID, _ := bson.ObjectIDFromHex(deciderIDStr.(string))
	now := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var req models.LeaveRequest
	if err := h.db.Collection(database.ColLeaveRequests).FindOne(ctx, bson.M{"_id": id}).Decode(&req); err != nil {
		RespondError(c, http.StatusNotFound, "leave request not found")
		return
	}
	if req.Status != "pending" {
		RespondError(c, http.StatusBadRequest, "request is not pending")
		return
	}

	update := bson.M{
		"$set": bson.M{
			"status":     body.Status,
			"decided_by": deciderID,
			"decided_at": now,
		},
	}
	h.db.Collection(database.ColLeaveRequests).UpdateOne(ctx, bson.M{"_id": id}, update)

	if body.Status == "approved" {
		days := req.EndDate.Sub(req.StartDate).Hours()/24 + 1
		h.db.Collection(database.ColLeaveBalances).UpdateOne(ctx, bson.M{
			"employee_id":   req.EmployeeID,
			"leave_type_id": req.LeaveTypeID,
		}, bson.M{"$inc": bson.M{"remaining_days": -days}})
	}

	h.db.Collection(database.ColLeaveRequests).FindOne(ctx, bson.M{"_id": id}).Decode(&req)
	RespondOK(c, req)
}
