package response

import "github.com/gin-gonic/gin"

func Success(c *gin.Context, data any) {
	c.JSON(200, gin.H{"data": data})
}

func SuccessCreated(c *gin.Context, data any) {
	c.JSON(201, gin.H{"data": data})
}

func BadRequest(c *gin.Context, msg string) {
	c.JSON(400, gin.H{"error": msg})
}
