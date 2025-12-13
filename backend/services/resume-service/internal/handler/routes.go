package handler

import (
	"resume-service/internal/handler/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, h *Handler, jwtSecret string) {
	v1 := router.Group("/api")
	{
		resumes := v1.Group("/resumes")
		{
			resumes.GET("/feed", h.GetResumeFeed)
			resumes.GET("/:id", h.GetResumeByID)
		}

		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(jwtSecret))
		{
			protected.POST("/resumes", h.CreateResume)
			protected.GET("/resumes/me", h.GetMyResumes)
			protected.PUT("/resumes/:id", h.UpdateResume)
			protected.DELETE("/resumes/:id", h.DeleteResume)
			protected.POST("/resumes/:id/publish", h.PublishResume)
		}
	}
}
