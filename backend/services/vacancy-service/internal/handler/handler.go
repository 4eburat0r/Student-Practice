package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"vacancy-service/internal/handler/middleware"
	"vacancy-service/internal/service"
	"vacancy-service/pkg/validator"
)

type Handler struct {
	vacancyService service.VacancyService
	logger         *zap.Logger
	validator      *validator.Validator
}

func NewHandler(vacancyService service.VacancyService, logger *zap.Logger, v *validator.Validator) *Handler {
	return &Handler{
		vacancyService: vacancyService,
		logger:         logger,
		validator:      v,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	// глобальные middleware (по желанию)
	router.Use(
		middleware.CORS(),
		middleware.Logger(h.logger),
		middleware.Recovery(),
	)

	api := router.Group("/vacancies")
	{
		// общедоступные (для студентов/гостей)
		api.GET("/feed", h.GetFeed)
		api.GET("/:id", h.GetVacancyByID)

		// ручки работодателя (нужен авторизованный employer)
		employer := api.Group("")
		employer.Use(
			middleware.AuthRequired(),
			middleware.RequireEmployer(),
		)
		{
			employer.GET("/me", h.GetMyVacancies)
			employer.POST("", h.CreateVacancy)
			employer.PUT("/:id", h.UpdateVacancy)
			employer.POST("/:id/publish", h.PublishVacancy)
			employer.DELETE("/:id", h.DeleteVacancy)
		}
	}
}
