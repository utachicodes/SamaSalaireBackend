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

type EmployeeHandler struct {
	db *mongo.Database
}

func NewEmployeeHandler(db *mongo.Database) *EmployeeHandler {
	return &EmployeeHandler{db: db}
}

func (h *EmployeeHandler) List(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}
	if dep := c.Query("department"); dep != "" {
		filter["department"] = dep
	}
	if active := c.Query("is_active"); active == "false" {
		filter["is_active"] = false
	} else {
		filter["is_active"] = true
	}

	role, _ := c.Get("role")
	if role == "employee" {
		empIDStr, _ := c.Get("employeeID")
		empID, err := bson.ObjectIDFromHex(empIDStr.(string))
		if err != nil {
			RespondError(c, http.StatusBadRequest, "invalid employee id")
			return
		}
		filter["_id"] = empID
	}

	cursor, err := h.db.Collection(database.ColEmployees).Find(ctx, filter)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to fetch employees")
		return
	}
	var employees []models.Employee
	if err := cursor.All(ctx, &employees); err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to decode employees")
		return
	}
	if employees == nil {
		employees = []models.Employee{}
	}
	RespondOK(c, employees)
}

func (h *EmployeeHandler) Create(c *gin.Context) {
	var emp models.Employee
	if err := c.ShouldBindJSON(&emp); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	emp.ID = bson.NewObjectID()
	emp.IsActive = true

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := h.db.Collection(database.ColEmployees).InsertOne(ctx, emp); err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to create employee")
		return
	}
	RespondCreated(c, emp)
}

func (h *EmployeeHandler) Get(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid employee id")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var emp models.Employee
	if err := h.db.Collection(database.ColEmployees).FindOne(ctx, bson.M{"_id": id}).Decode(&emp); err != nil {
		RespondError(c, http.StatusNotFound, "employee not found")
		return
	}
	RespondOK(c, emp)
}

func (h *EmployeeHandler) Update(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid employee id")
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

	result, err := h.db.Collection(database.ColEmployees).UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": body})
	if err != nil || result.MatchedCount == 0 {
		RespondError(c, http.StatusNotFound, "employee not found")
		return
	}

	var updated models.Employee
	h.db.Collection(database.ColEmployees).FindOne(ctx, bson.M{"_id": id}).Decode(&updated)
	RespondOK(c, updated)
}

func (h *EmployeeHandler) Delete(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := h.db.Collection(database.ColEmployees).UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"is_active": false}})
	if err != nil || result.MatchedCount == 0 {
		RespondError(c, http.StatusNotFound, "employee not found")
		return
	}
	RespondOK(c, gin.H{"message": "employee deactivated"})
}
