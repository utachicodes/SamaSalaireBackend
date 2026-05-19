package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RespondOK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"data": data, "error": nil})
}

func RespondCreated(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, gin.H{"data": data, "error": nil})
}

func RespondError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"data": nil, "error": message})
}
