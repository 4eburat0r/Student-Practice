package response

import (
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func Error(c *gin.Context, statusCode int, err interface{}) {
	var errorMessage string

	switch v := err.(type) {
	case error:
		errorMessage = v.Error()
	case string:
		errorMessage = v
	case gin.H:
		c.JSON(statusCode, gin.H{
			"success": false,
			"error":   v,
		})
		return
	default:
		errorMessage = "Unknown error"
	}

	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Error:   errorMessage,
	})
}
