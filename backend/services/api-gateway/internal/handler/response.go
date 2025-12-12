package handler

import (
	"fmt"
	"io"
	"net/http"

	"api-gateway/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateResponse(c *gin.Context) {
	token := c.GetString("token")
	role := c.GetString("role")

	if role != "student" {
		response.Error(c, http.StatusForbidden, "Only students can create responses")
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	respBody, statusCode, err := h.deps.ResponseClient.ProxyRequest(
		c.Request.Context(),
		"POST",
		"/api/responses",
		body,
		token,
	)
	if err != nil {
		response.Error(c, statusCode, "Failed to create response")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}

func (h *Handler) GetMyResponses(c *gin.Context) {
	token := c.GetString("token")

	respBody, statusCode, err := h.deps.ResponseClient.ProxyRequest(
		c.Request.Context(),
		"GET",
		"/api/responses/me",
		nil,
		token,
	)
	if err != nil {
		response.Error(c, statusCode, "Failed to get responses")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}

func (h *Handler) GetResponseByID(c *gin.Context) {
	id := c.Param("id")
	token := c.GetString("token")

	path := fmt.Sprintf("/api/responses/%s", id)

	respBody, statusCode, err := h.deps.ResponseClient.ProxyRequest(
		c.Request.Context(),
		"GET",
		path,
		nil,
		token,
	)
	if err != nil {
		response.Error(c, statusCode, "Failed to get response")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}

func (h *Handler) UpdateResponseStatus(c *gin.Context) {
	id := c.Param("id")
	token := c.GetString("token")
	role := c.GetString("role")

	if role != "employer" {
		response.Error(c, http.StatusForbidden, "Only employers can update response status")
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	path := fmt.Sprintf("/api/responses/%s/status", id)

	respBody, statusCode, err := h.deps.ResponseClient.ProxyRequest(
		c.Request.Context(),
		"PUT",
		path,
		body,
		token,
	)
	if err != nil {
		response.Error(c, statusCode, "Failed to update response status")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}

func (h *Handler) DeleteResponse(c *gin.Context) {
	id := c.Param("id")
	token := c.GetString("token")

	path := fmt.Sprintf("/api/responses/%s", id)

	respBody, statusCode, err := h.deps.ResponseClient.ProxyRequest(
		c.Request.Context(),
		"DELETE",
		path,
		nil,
		token,
	)
	if err != nil {
		response.Error(c, statusCode, "Failed to delete response")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}
