package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	HeaderUserID = "X-User-Id"
	HeaderRole   = "X-User-Role"

	RoleStudent  = "student"
	RoleEmployer = "employer"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.GetHeader(HeaderUserID)
		role := c.GetHeader(HeaderRole)
		if idStr == "" || role == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
			return
		}
		c.Set("user_id", id)
		c.Set("user_role", role)
		c.Next()
	}
}

func RequireStudent() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, ok := c.Get("user_role")
		if !ok || roleVal != RoleStudent {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}
