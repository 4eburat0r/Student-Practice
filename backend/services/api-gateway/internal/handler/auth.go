package handler

import (
	"io"
	"net/http"

	"api-gateway/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Register(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	respBody, statusCode, err := h.deps.AuthClient.Register(c.Request.Context(), body)
	if err != nil {
		response.Error(c, statusCode, "Registration failed")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}

func (h *Handler) Login(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	respBody, statusCode, err := h.deps.AuthClient.Login(c.Request.Context(), body)
	if err != nil {
		response.Error(c, statusCode, "Login failed")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}

func (h *Handler) RefreshToken(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	respBody, statusCode, err := h.deps.AuthClient.RefreshToken(c.Request.Context(), body)
	if err != nil {
		response.Error(c, statusCode, "Token refresh failed")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}
