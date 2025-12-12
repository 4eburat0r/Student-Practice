package handler

import (
	"net/http"

	"api-gateway/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HealthCheck(c *gin.Context) {
	response.Success(c, http.StatusOK, gin.H{
		"status":  "ok",
		"service": "api-gateway",
	})
}
