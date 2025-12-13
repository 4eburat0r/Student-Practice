package handler

import (
	"api-gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	router.GET("/health", h.HealthCheck)

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.Register)
			auth.POST("/login", h.Login)
			auth.POST("/refresh", h.RefreshToken)

			student := auth.Group("/student")
			{
				student.POST("/register", h.RegisterStudent)
				student.POST("/login", h.LoginStudent)
			}

			employer := auth.Group("/employer")
			{
				employer.POST("/register", h.RegisterEmployer)
				employer.POST("/login", h.LoginEmployer)
			}
		}

		resumes := api.Group("/resumes")
		{
			resumes.GET("/feed", h.GetResumesFeed)
			resumes.GET("/:id", h.GetResumeByID)
		}

		vacancies := api.Group("/vacancies")
		{
			vacancies.GET("/feed", h.GetVacanciesFeed)
			vacancies.GET("/:id", h.GetVacancyByID)
		}

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(h.deps.AuthClient))
		{
			users := protected.Group("/users")
			{
				students := users.Group("/students")
				{
					students.GET("/me", h.GetMyStudentProfile)
					students.PUT("/me", h.UpdateMyStudentProfile)
					students.DELETE("/me", h.DeleteMyStudentProfile)
				}

				employers := users.Group("/employers")
				{
					employers.GET("/me", h.GetMyEmployerProfile)
					employers.PUT("/me", h.UpdateMyEmployerProfile)
					employers.DELETE("/me", h.DeleteMyEmployerProfile)
				}
			}

			resumes := protected.Group("/resumes")
			{
				resumes.POST("", h.CreateResume)
				resumes.GET("/me", h.GetMyResumes)
				resumes.PUT("/:id", h.UpdateResume)
				resumes.DELETE("/:id", h.DeleteResume)
				resumes.POST("/:id/publish", h.PublishResume)
			}

			vacancies := protected.Group("/vacancies")
			{
				vacancies.POST("", h.CreateVacancy)
				vacancies.GET("/me", h.GetMyVacancies)
				vacancies.PUT("/:id", h.UpdateVacancy)
				vacancies.DELETE("/:id", h.DeleteVacancy)
				vacancies.POST("/:id/publish", h.PublishVacancy)
			}

			responses := protected.Group("/responses")
			{
				responses.POST("", h.CreateResponse)
				responses.GET("/me", h.GetMyResponses)
				responses.GET("/:id", h.GetResponseByID)
				responses.PUT("/:id/status", h.UpdateResponseStatus)
				responses.DELETE("/:id", h.DeleteResponse)
			}

			protected.GET("/profile/me", h.GetFullProfile)
			protected.GET("/dashboard", h.GetDashboard)
		}
	}

	return router
}
