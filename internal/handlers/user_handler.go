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
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	db *mongo.Database
}

func NewUserHandler(db *mongo.Database) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) List(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, _ := h.db.Collection(database.ColUsers).Find(ctx, bson.M{})
	var users []models.User
	cursor.All(ctx, &users)
	if users == nil {
		users = []models.User{}
	}
	RespondOK(c, users)
}

type createUserBody struct {
	EmployeeID string `json:"employee_id" binding:"required"`
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Role       string `json:"role" binding:"required"`
}

func (h *UserHandler) Create(c *gin.Context) {
	var body createUserBody
	if err := c.ShouldBindJSON(&body); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	empID, err := bson.ObjectIDFromHex(body.EmployeeID)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid employee_id")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to hash password")
		return
	}

	user := models.User{
		ID:           bson.NewObjectID(),
		EmployeeID:   empID,
		Username:     body.Username,
		PasswordHash: string(hash),
		Role:         body.Role,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := h.db.Collection(database.ColUsers).InsertOne(ctx, user); err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to create user (username may already exist)")
		return
	}
	RespondCreated(c, user)
}

type updateUserBody struct {
	Role     string `json:"role"`
	Password string `json:"password"`
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var body updateUserBody
	if err := c.ShouldBindJSON(&body); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	update := bson.M{}
	if body.Role != "" {
		update["role"] = body.Role
	}
	if body.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			RespondError(c, http.StatusInternalServerError, "failed to hash password")
			return
		}
		update["password_hash"] = string(hash)
	}

	if len(update) == 0 {
		RespondError(c, http.StatusBadRequest, "no fields to update")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h.db.Collection(database.ColUsers).UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	var user models.User
	h.db.Collection(database.ColUsers).FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	RespondOK(c, user)
}
