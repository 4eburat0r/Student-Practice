package handler

import (
	"fmt"
	"io"
	"net/http"

	"api-gateway/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateVacancy(c *gin.Context) {
	token := c.GetString("token")
	role := c.GetString("role")

	if role != "employer" {
		response.Error(c, http.StatusForbidden, "Only employers can create vacancies")
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	respBody, statusCode, err := h.deps.VacancyClient.ProxyRequest(
		c.Request.Context(),
		"POST",
		"/api/vacancies",
		body,
		token,
	)
	if err != nil {
		response.Error(c, statusCode, "Failed to create vacancy")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}

func (h *Handler) GetMyVacancies(c *gin.Context) {
	token := c.GetString("token")
	role := c.GetString("role")

	if role != "employer" {
		response.Error(c, http.StatusForbidden, "Only employers can access this endpoint")
		return
	}

	respBody, statusCode, err := h.deps.VacancyClient.ProxyRequest(
		c.Request.Context(),
		"GET",
		"/api/vacancies/me",
		nil,
		token,
	)
	if err != nil {
		response.Error(c, statusCode, "Failed to get vacancies")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}

func (h *Handler) GetVacanciesFeed(c *gin.Context) {
	token := c.GetString("token")

	queryString := c.Request.URL.RawQuery
	path := "/api/vacancies/feed"
	if queryString != "" {
		path += "?" + queryString
	}

	respBody, statusCode, err := h.deps.VacancyClient.ProxyRequest(
		c.Request.Context(),
		"GET",
		path,
		nil,
		token,
	)
	if err != nil {
		response.Error(c, statusCode, "Failed to get vacancies feed")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}

func (h *Handler) GetVacancyByID(c *gin.Context) {
	id := c.Param("id")
	token := c.GetString("token")

	path := fmt.Sprintf("/api/vacancies/%s", id)

	respBody, statusCode, err := h.deps.VacancyClient.ProxyRequest(
		c.Request.Context(),
		"GET",
		path,
		nil,
		token,
	)
	if err != nil {
		response.Error(c, statusCode, "Failed to get vacancy")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}

func (h *Handler) UpdateVacancy(c *gin.Context) {
	id := c.Param("id")
	token := c.GetString("token")
	role := c.GetString("role")

	if role != "employer" {
		response.Error(c, http.StatusForbidden, "Only employers can update vacancies")
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	path := fmt.Sprintf("/api/vacancies/%s", id)

	respBody, statusCode, err := h.deps.VacancyClient.ProxyRequest(
		c.Request.Context(),
		"PUT",
		path,
		body,
		token,
	)
	if err != nil {
		response.Error(c, statusCode, "Failed to update vacancy")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}

func (h *Handler) DeleteVacancy(c *gin.Context) {
	id := c.Param("id")
	token := c.GetString("token")
	role := c.GetString("role")

	if role != "employer" {
		response.Error(c, http.StatusForbidden, "Only employers can delete vacancies")
		return
	}

	path := fmt.Sprintf("/api/vacancies/%s", id)

	respBody, statusCode, err := h.deps.VacancyClient.ProxyRequest(
		c.Request.Context(),
		"DELETE",
		path,
		nil,
		token,
	)
	if err != nil {
		response.Error(c, statusCode, "Failed to delete vacancy")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}

func (h *Handler) PublishVacancy(c *gin.Context) {
	id := c.Param("id")
	token := c.GetString("token")
	role := c.GetString("role")

	if role != "employer" {
		response.Error(c, http.StatusForbidden, "Only employers can publish vacancies")
		return
	}

	path := fmt.Sprintf("/api/vacancies/%s/publish", id)

	respBody, statusCode, err := h.deps.VacancyClient.ProxyRequest(
		c.Request.Context(),
		"POST",
		path,
		nil,
		token,
	)
	if err != nil {
		response.Error(c, statusCode, "Failed to publish vacancy")
		return
	}

	c.Data(statusCode, "application/json", respBody)
}
