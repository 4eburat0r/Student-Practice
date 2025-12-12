package handler

import (
	"net/http"
	"resume-service/internal/domain/errors"
	"resume-service/internal/domain/models/request"
	"resume-service/internal/service"
	"resume-service/pkg/jwt"
	"resume-service/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	resumeService service.Service
	jwtManager    *jwt.JWTManager
}

func NewHandler(resumeService service.Service, jwtManager *jwt.JWTManager) *Handler {
	return &Handler{
		resumeService: resumeService,
		jwtManager:    jwtManager,
	}
}

func (h *Handler) CreateResume(c *gin.Context) {

	studentID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, errors.ErrUnauthorized)
		return
	}

	role, _ := c.Get("role")
	if role != "student" {
		response.Error(c, http.StatusForbidden, errors.ErrForbidden)
		return
	}

	var req request.CreateResumeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	resume, err := h.resumeService.Create(c.Request.Context(), studentID.(int), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err)
		return
	}

	response.Success(c, http.StatusCreated, resume)
}

func (h *Handler) GetResumeByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	resume, err := h.resumeService.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == errors.ErrResumeNotFound {
			response.Error(c, http.StatusNotFound, err)
			return
		}
		response.Error(c, http.StatusInternalServerError, err)
		return
	}

	response.Success(c, http.StatusOK, resume)
}

func (h *Handler) GetMyResumes(c *gin.Context) {
	studentID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, errors.ErrUnauthorized)
		return
	}

	role, _ := c.Get("role")
	if role != "student" {
		response.Error(c, http.StatusForbidden, errors.ErrForbidden)
		return
	}

	resumes, err := h.resumeService.GetMyResumes(c.Request.Context(), studentID.(int))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err)
		return
	}

	response.Success(c, http.StatusOK, resumes)
}

func (h *Handler) GetResumeFeed(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	resumes, err := h.resumeService.GetFeed(c.Request.Context(), limit, offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err)
		return
	}

	response.Success(c, http.StatusOK, resumes)
}

func (h *Handler) UpdateResume(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	studentID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, errors.ErrUnauthorized)
		return
	}

	role, _ := c.Get("role")
	if role != "student" {
		response.Error(c, http.StatusForbidden, errors.ErrForbidden)
		return
	}

	var req request.UpdateResumeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	resume, err := h.resumeService.Update(c.Request.Context(), id, studentID.(int), &req)
	if err != nil {
		if err == errors.ErrResumeNotFound {
			response.Error(c, http.StatusNotFound, err)
			return
		}
		if err == errors.ErrForbidden {
			response.Error(c, http.StatusForbidden, err)
			return
		}
		response.Error(c, http.StatusInternalServerError, err)
		return
	}

	response.Success(c, http.StatusOK, resume)
}

func (h *Handler) DeleteResume(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	studentID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, errors.ErrUnauthorized)
		return
	}

	role, _ := c.Get("role")
	if role != "student" {
		response.Error(c, http.StatusForbidden, errors.ErrForbidden)
		return
	}

	err = h.resumeService.Delete(c.Request.Context(), id, studentID.(int))
	if err != nil {
		if err == errors.ErrResumeNotFound {
			response.Error(c, http.StatusNotFound, err)
			return
		}
		if err == errors.ErrForbidden {
			response.Error(c, http.StatusForbidden, err)
			return
		}
		response.Error(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) PublishResume(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	studentID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, errors.ErrUnauthorized)
		return
	}

	role, _ := c.Get("role")
	if role != "student" {
		response.Error(c, http.StatusForbidden, errors.ErrForbidden)
		return
	}

	err = h.resumeService.Publish(c.Request.Context(), id, studentID.(int))
	if err != nil {
		if err == errors.ErrResumeNotFound {
			response.Error(c, http.StatusNotFound, err)
			return
		}
		if err == errors.ErrForbidden {
			response.Error(c, http.StatusForbidden, err)
			return
		}
		response.Error(c, http.StatusInternalServerError, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Resume published successfully"})
}
