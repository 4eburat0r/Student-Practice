package handler

import (
	"io"
	"net/http"

	"api-gateway/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetMyStudentProfile(c *gin.Context) {
	token := c.GetString("token")

	respBody, statusCode, err := h.deps.UserClient.GetStudentProfile(c.Request.Context(), token)
	if err != nil {
		response.Error(c, statusCode, "Failed to get profile")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}

func (h *Handler) UpdateMyStudentProfile(c *gin.Context) {
	token := c.GetString("token")

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	respBody, statusCode, err := h.deps.UserClient.UpdateStudentProfile(c.Request.Context(), body, token)
	if err != nil {
		response.Error(c, statusCode, "Failed to update profile")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}

func (h *Handler) DeleteMyStudentProfile(c *gin.Context) {
	token := c.GetString("token")

	respBody, statusCode, err := h.deps.UserClient.DeleteStudentProfile(c.Request.Context(), token)
	if err != nil {
		response.Error(c, statusCode, "Failed to delete profile")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}

func (h *Handler) GetMyEmployerProfile(c *gin.Context) {
	token := c.GetString("token")

	respBody, statusCode, err := h.deps.UserClient.GetEmployerProfile(c.Request.Context(), token)
	if err != nil {
		response.Error(c, statusCode, "Failed to get profile")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}

func (h *Handler) UpdateMyEmployerProfile(c *gin.Context) {
	token := c.GetString("token")

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	respBody, statusCode, err := h.deps.UserClient.UpdateEmployerProfile(c.Request.Context(), body, token)
	if err != nil {
		response.Error(c, statusCode, "Failed to update profile")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}

func (h *Handler) DeleteMyEmployerProfile(c *gin.Context) {
	token := c.GetString("token")

	respBody, statusCode, err := h.deps.UserClient.DeleteEmployerProfile(c.Request.Context(), token)
	if err != nil {
		response.Error(c, statusCode, "Failed to delete profile")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}
