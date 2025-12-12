package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"response-service/internal/handler/middleware"
	"response-service/internal/service"
	"response-service/pkg/validator"
)


type Handler struct {
	responseService service.ResponseService
	logger          *zap.Logger
	validator       *validator.Validator
}

func NewHandler(responseService service.ResponseService, logger *zap.Logger, v *validator.Validator) *Handler {
	return &Handler{
		responseService: responseService,
		logger:          logger,
		validator:       v,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.Use(
		middleware.CORS(),
		middleware.Logger(h.logger),
		middleware.Recovery(),
	)

	// Все ручки откликов — только для студента
	api := router.Group("/response")
	api.Use(
		middleware.AuthRequired(),
		middleware.RequireStudent(),
	)
	{
		api.POST("", h.CreateResponse)
		api.GET("/me", h.ListMyResponses)
		api.GET("/:id", h.GetResponseByID)
		api.DELETE("/:id", h.DeleteResponse)
	}
}
