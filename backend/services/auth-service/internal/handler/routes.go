package handler

import (
	"auth-service/internal/handler/middleware"
	"auth-service/internal/service"
	"auth-service/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, authService service.AuthService, jwtManager *jwt.JWTManager, jwtSecret string) {
	handler := NewHandler(authService)

	auth := r.Group("/api/auth")
	{
		auth.POST("", handler.GetRoles)

		student := auth.Group("/student")
		{
			student.POST("/register", handler.RegisterStudent)
			student.POST("/login", handler.Login)
		}

		employer := auth.Group("/employer")
		{
			employer.POST("/register", handler.RegisterEmployer)
			employer.POST("/login", handler.Login)
		}

		auth.POST("/refresh", handler.RefreshToken)
		auth.POST("/validate", handler.ValidateToken)
	}

	protected := r.Group("/api/auth")
	protected.Use(middleware.Auth(jwtManager))
	{
		protected.GET("/me", handler.GetCurrentUser)
	}
}

func (h *Handler) GetCurrentUser(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	email, _ := middleware.GetEmail(c)
	role, _ := middleware.GetRole(c)

	c.JSON(200, gin.H{
		"user_id": userID,
		"email":   email,
		"role":    role,
	})
}
