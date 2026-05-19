package middleware

import (
	"context"
	"log"
	"time"

	"samasalaire-backend/internal/database"
	"samasalaire-backend/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func AuditLog(db *mongo.Database, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Status() >= 400 {
			return
		}

		userIDStr, _ := c.Get("userID")
		userIDRaw, _ := userIDStr.(string)
		userID, err := bson.ObjectIDFromHex(userIDRaw)
		if err != nil {
			return
		}

		log := &models.AuditLog{
			UserID:    userID,
			Action:    action,
			Target:    c.Request.URL.Path,
			Timestamp: time.Now(),
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		db.Collection(database.ColAuditLogs).InsertOne(ctx, log)
	}
}

// avoid name clash with standard log package
var _ = log.Println
