package handlers

import (
	"context"
	"net/http"
	"time"

	"samasalaire-backend/internal/config"
	"samasalaire-backend/internal/database"
	"samasalaire-backend/internal/middleware"
	"samasalaire-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db *mongo.Database
}

func NewAuthHandler(db *mongo.Database) *AuthHandler {
	return &AuthHandler{db: db}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "username and password must be provided")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := h.db.Collection(database.ColUsers).
		FindOne(ctx, bson.M{"username": req.Username}).
		Decode(&user)
	if err != nil {
		RespondError(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		RespondError(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	cfg := config.Load()
	claims := middleware.UserClaims{
		UserID:     user.ID.Hex(),
		EmployeeID: user.EmployeeID.Hex(),
		Role:       user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.JWTExpiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to generate token")
		return
	}

	RespondOK(c, gin.H{
		"token": signed,
		"user": gin.H{
			"id":          user.ID.Hex(),
			"role":        user.Role,
			"employee_id": user.EmployeeID.Hex(),
			"username":    user.Username,
		},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	RespondOK(c, gin.H{"message": "logged out"})
}
