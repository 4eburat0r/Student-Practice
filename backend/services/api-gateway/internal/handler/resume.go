package handler

import (
	"fmt"
	"io"
	"net/http"

	"api-gateway/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateResume(c *gin.Context) {
	token := c.GetString("token")
	role := c.GetString("role")

	if role != "student" {
		response.Error(c, http.StatusForbidden, "Only students can create resumes")
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	respBody, statusCode, err := h.deps.ResumeClient.ProxyRequest(
		c.Request.Context(),
		"POST",
		"/api/resumes",
		body,
		token,
	)
	if err != nil {
		response.Error(c, statusCode, "Failed to create resume")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}

func (h *Handler) GetMyResumes(c *gin.Context) {
	token := c.GetString("token")
	role := c.GetString("role")

	if role != "student" {
		response.Error(c, http.StatusForbidden, "Only students can access this endpoint")
		return
	}

	respBody, statusCode, err := h.deps.ResumeClient.ProxyRequest(
		c.Request.Context(),
		"GET",
		"/api/resumes/me",
		nil,
		token,
	)
	if err != nil {
		response.Error(c, statusCode, "Failed to get resumes")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}

func (h *Handler) GetResumesFeed(c *gin.Context) {
	token := c.GetString("token")

	queryString := c.Request.URL.RawQuery
	path := "/api/resumes/feed"
	if queryString != "" {
		path += "?" + queryString
	}

	respBody, statusCode, err := h.deps.ResumeClient.ProxyRequest(
		c.Request.Context(),
		"GET",
		path,
		nil,
		token,
	)
	if err != nil {
		response.Error(c, statusCode, "Failed to get resumes feed")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}

func (h *Handler) GetResumeByID(c *gin.Context) {
	id := c.Param("id")
	token := c.GetString("token")

	path := fmt.Sprintf("/api/resumes/%s", id)

	respBody, statusCode, err := h.deps.ResumeClient.ProxyRequest(
		c.Request.Context(),
		"GET",
		path,
		nil,
		token,
	)
	if err != nil {
		response.Error(c, statusCode, "Failed to get resume")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}

func (h *Handler) UpdateResume(c *gin.Context) {
	id := c.Param("id")
	token := c.GetString("token")
	role := c.GetString("role")

	if role != "student" {
		response.Error(c, http.StatusForbidden, "Only students can update resumes")
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	path := fmt.Sprintf("/api/resumes/%s", id)

	respBody, statusCode, err := h.deps.ResumeClient.ProxyRequest(
		c.Request.Context(),
		"PUT",
		path,
		body,
		token,
	)
	if err != nil {
		response.Error(c, statusCode, "Failed to update resume")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}

func (h *Handler) DeleteResume(c *gin.Context) {
	id := c.Param("id")
	token := c.GetString("token")
	role := c.GetString("role")

	if role != "student" {
		response.Error(c, http.StatusForbidden, "Only students can delete resumes")
		return
	}

	path := fmt.Sprintf("/api/resumes/%s", id)

	respBody, statusCode, err := h.deps.ResumeClient.ProxyRequest(
		c.Request.Context(),
		"DELETE",
		path,
		nil,
		token,
	)
	if err != nil {
		response.Error(c, statusCode, "Failed to delete resume")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}

func (h *Handler) PublishResume(c *gin.Context) {
	id := c.Param("id")
	token := c.GetString("token")
	role := c.GetString("role")

	if role != "student" {
		response.Error(c, http.StatusForbidden, "Only students can publish resumes")
		return
	}

	path := fmt.Sprintf("/api/resumes/%s/publish", id)

	respBody, statusCode, err := h.deps.ResumeClient.ProxyRequest(
		c.Request.Context(),
		"POST",
		path,
		nil,
		token,
	)
	if err != nil {
		response.Error(c, statusCode, "Failed to publish resume")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}
