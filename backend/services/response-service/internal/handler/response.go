package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	reqModel "response-service/internal/domain/models/request"
	"response-service/pkg/response"
)

func (h *Handler) CreateResponse(c *gin.Context) {
	studentID := c.GetInt64("user_id")

	var req reqModel.CreateResponseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	resp, err := h.responseService.Create(c.Request.Context(), studentID, req)
	if err != nil {
		response.HandleError(c, err, h.logger)
		return
	}
	response.SuccessCreated(c, resp)
}

func (h *Handler) ListMyResponses(c *gin.Context) {
	studentID := c.GetInt64("user_id")

	var q reqModel.ResponseListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.ValidationError(c, err)
		return
	}

	list, err := h.responseService.ListByStudent(c.Request.Context(), studentID, q)
	if err != nil {
		response.HandleError(c, err, h.logger)
		return
	}
	response.Success(c, list)
}

func (h *Handler) GetResponseByID(c *gin.Context) {
	studentID := c.GetInt64("user_id")

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	resp, err := h.responseService.GetByID(c.Request.Context(), studentID, id)
	if err != nil {
		response.HandleError(c, err, h.logger)
		return
	}
	response.Success(c, resp)
}

func (h *Handler) DeleteResponse(c *gin.Context) {
	studentID := c.GetInt64("user_id")

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.responseService.Delete(c.Request.Context(), studentID, id); err != nil {
		response.HandleError(c, err, h.logger)
		return
	}
	c.Status(204)
}
