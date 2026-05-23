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

type SalaryHandler struct {
	db *mongo.Database
}

func NewSalaryHandler(db *mongo.Database) *SalaryHandler {
	return &SalaryHandler{db: db}
}

func (h *SalaryHandler) ListByEmployee(c *gin.Context) {
	empID, err := bson.ObjectIDFromHex(c.Param("employeeId"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid employee id")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := h.db.Collection(database.ColSalaryComponents).Find(ctx, bson.M{"employee_id": empID})
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to fetch salary components")
		return
	}
	var components []models.SalaryComponent
	if err := cursor.All(ctx, &components); err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to decode components")
		return
	}
	if components == nil {
		components = []models.SalaryComponent{}
	}
	RespondOK(c, components)
}

func (h *SalaryHandler) Create(c *gin.Context) {
	var comp models.SalaryComponent
	if err := c.ShouldBindJSON(&comp); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	comp.ID = bson.NewObjectID()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := h.db.Collection(database.ColSalaryComponents).InsertOne(ctx, comp); err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to create salary component")
		return
	}
	RespondCreated(c, comp)
}

func (h *SalaryHandler) Update(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid salary component id")
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

	result, err := h.db.Collection(database.ColSalaryComponents).UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": body})
	if err != nil || result.MatchedCount == 0 {
		RespondError(c, http.StatusNotFound, "salary component not found")
		return
	}

	var updated models.SalaryComponent
	h.db.Collection(database.ColSalaryComponents).FindOne(ctx, bson.M{"_id": id}).Decode(&updated)
	RespondOK(c, updated)
}

func (h *SalaryHandler) Delete(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := h.db.Collection(database.ColSalaryComponents).DeleteOne(ctx, bson.M{"_id": id})
	if err != nil || result.DeletedCount == 0 {
		RespondError(c, http.StatusNotFound, "salary component not found")
		return
	}
	RespondOK(c, gin.H{"message": "deleted"})
}
