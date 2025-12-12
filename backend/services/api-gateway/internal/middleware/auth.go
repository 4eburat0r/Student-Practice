package middleware

import (
	"net/http"
	"strings"

	"api-gateway/internal/client"
	"api-gateway/pkg/logger"
	"api-gateway/pkg/response"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authClient *client.AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("Missing authorization header",
				logger.String("path", c.Request.URL.Path),
				logger.String("ip", c.ClientIP()),
			)
			response.Error(c, http.StatusUnauthorized, "Missing authorization header")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warn("Invalid authorization format",
				logger.String("header", authHeader),
			)
			response.Error(c, http.StatusUnauthorized, "Invalid authorization format")
			c.Abort()
			return
		}

		token := parts[1]

		claims, err := authClient.ValidateToken(c.Request.Context(), token)
		if err != nil {
			logger.Warn("Token validation failed",
				logger.Err(err),
				logger.String("ip", c.ClientIP()),
			)
			response.Error(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "Invalid token claims")
			c.Abort()
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "Invalid token claims")
			c.Abort()
			return
		}

		c.Set("user_id", int(userID))
		c.Set("role", role)
		c.Set("token", token)

		c.Next()
	}
}
