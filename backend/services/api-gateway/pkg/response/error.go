package response

import (
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Error:   message,
	})
}
