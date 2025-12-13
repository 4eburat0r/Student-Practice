package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"api-gateway/pkg/response"
	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=student employer"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=student employer"`
}

func (h *Handler) Register(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	var req RegisterRequest
	if len(body) > 0 {
		if err := json.Unmarshal(body, &req); err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid JSON format")
			return
		}
	}

	if roleOverride, ok := c.Get("role_override"); ok {
		if roleStr, ok := roleOverride.(string); ok {
			req.Role = roleStr
		}
	}

	if req.Email == "" || req.Password == "" || req.Role == "" {
		response.Error(c, http.StatusBadRequest, "Email, password and role are required")
		return
	}

	if req.Role != "student" && req.Role != "employer" {
		response.Error(c, http.StatusBadRequest, "Role must be 'student' or 'employer'")
		return
	}

	safeBody, err := json.Marshal(req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to process request")
		return
	}

	respBody, statusCode, err := h.deps.AuthClient.Register(c.Request.Context(), safeBody)
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

	var req LoginRequest
	if len(body) > 0 {
		if err := json.Unmarshal(body, &req); err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid JSON format")
			return
		}
	}

	if roleOverride, ok := c.Get("role_override"); ok {
		if roleStr, ok := roleOverride.(string); ok {
			req.Role = roleStr
		}
	}

	if req.Email == "" || req.Password == "" || req.Role == "" {
		response.Error(c, http.StatusBadRequest, "Email, password and role are required")
		return
	}

	if req.Role != "student" && req.Role != "employer" {
		response.Error(c, http.StatusBadRequest, "Role must be 'student' or 'employer'")
		return
	}

	safeBody, err := json.Marshal(req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to process request")
		return
	}

	respBody, statusCode, err := h.deps.AuthClient.Login(c.Request.Context(), safeBody)
	if err != nil {
		response.Error(c, statusCode, "Login failed")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}

func (h *Handler) RegisterStudent(c *gin.Context) {
	c.Set("role_override", "student")
	h.Register(c)
}

func (h *Handler) RegisterEmployer(c *gin.Context) {
	c.Set("role_override", "employer")
	h.Register(c)
}

func (h *Handler) LoginStudent(c *gin.Context) {
	c.Set("role_override", "student")
	h.Login(c)
}

func (h *Handler) LoginEmployer(c *gin.Context) {
	c.Set("role_override", "employer")
	h.Login(c)
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
