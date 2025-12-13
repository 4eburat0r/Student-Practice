package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	reqModel "vacancy-service/internal/domain/models/request"
	"vacancy-service/pkg/response"
)

func (h *Handler) CreateVacancy(c *gin.Context) {
	employerID := c.GetInt64("user_id")

	var req reqModel.CreateVacancyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	vac, err := h.vacancyService.Create(c.Request.Context(), employerID, req)
	if err != nil {
		response.HandleError(c, err, h.logger)
		return
	}
	response.SuccessCreated(c, vac)
}

func (h *Handler) GetVacancyByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	vac, err := h.vacancyService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.HandleError(c, err, h.logger)
		return
	}
	response.Success(c, vac)
}

func (h *Handler) GetMyVacancies(c *gin.Context) {
	employerID := c.GetInt64("user_id")

	var q reqModel.VacancyFeedQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.ValidationError(c, err)
		return
	}
	list, err := h.vacancyService.GetByEmployer(c.Request.Context(), employerID, q.Limit, q.Offset)
	if err != nil {
		response.HandleError(c, err, h.logger)
		return
	}
	response.Success(c, list)
}

func (h *Handler) GetFeed(c *gin.Context) {
	var q reqModel.VacancyFeedQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.ValidationError(c, err)
		return
	}
	list, err := h.vacancyService.GetFeed(c.Request.Context(), q)
	if err != nil {
		response.HandleError(c, err, h.logger)
		return
	}
	response.Success(c, list)
}

func (h *Handler) UpdateVacancy(c *gin.Context) {
	employerID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	var req reqModel.UpdateVacancyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	vac, err := h.vacancyService.Update(c.Request.Context(), employerID, id, req)
	if err != nil {
		response.HandleError(c, err, h.logger)
		return
	}
	response.Success(c, vac)
}

func (h *Handler) DeleteVacancy(c *gin.Context) {
	employerID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	if err := h.vacancyService.Delete(c.Request.Context(), employerID, id); err != nil {
		response.HandleError(c, err, h.logger)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) PublishVacancy(c *gin.Context) {
	employerID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	if err := h.vacancyService.Publish(c.Request.Context(), employerID, id); err != nil {
		response.HandleError(c, err, h.logger)
		return
	}
	c.Status(http.StatusNoContent)
}
