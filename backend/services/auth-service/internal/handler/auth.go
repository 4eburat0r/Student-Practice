package handler

import (
    "log"
    "net/http"

    "auth-service/internal/domain/entity"
    "auth-service/internal/domain/errors"
    "auth-service/internal/domain/models/request"
    "auth-service/internal/service"

    "github.com/gin-gonic/gin"
)

type Handler struct {
    authService service.AuthService
}

func NewHandler(authService service.AuthService) *Handler {
    return &Handler{
        authService: authService,
    }
}

func (h *Handler) RegisterStudent(c *gin.Context) {
    var req request.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request body",
        })
        return
    }

    req.Role = entity.RoleStudent 

    if err := h.authService.Register(c.Request.Context(), &req); err != nil {
        h.handleError(c, err)
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Student registered successfully",
    })
}

func (h *Handler) RegisterEmployer(c *gin.Context) {
    var req request.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request body",
        })
        return
    }

    req.Role = entity.RoleEmployer

    if err := h.authService.Register(c.Request.Context(), &req); err != nil {
        h.handleError(c, err)
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Employer registered successfully",
    })
}

func (h *Handler) Login(c *gin.Context) {
    var req request.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request body",
        })
        return
    }

    resp, err := h.authService.Login(c.Request.Context(), &req)
    if err != nil {
        h.handleError(c, err)
        return
    }

    c.JSON(http.StatusOK, resp)
}

func (h *Handler) RefreshToken(c *gin.Context) {
    var req request.RefreshTokenRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request body",
        })
        return
    }

    resp, err := h.authService.RefreshToken(c.Request.Context(), &req)
    if err != nil {
        h.handleError(c, err)
        return
    }

    c.JSON(http.StatusOK, resp)
}

func (h *Handler) ValidateToken(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "Authorization header required",
        })
        return
    }

    token := authHeader[len("Bearer "):]
    claims, err := h.authService.ValidateToken(c.Request.Context(), token)
    if err != nil {
        h.handleError(c, err)
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "valid":   true,
        "user_id": claims.UserID,
        "email":   claims.Email,
        "role":    claims.Role,
    })
}

func (h *Handler) handleError(c *gin.Context, err error) {
    switch err {
    case errors.ErrInvalidCredentials:
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
    case errors.ErrUserAlreadyExists:
        c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
    case errors.ErrUserNotFound:
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
    case errors.ErrUnauthorized:
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
    case errors.ErrInvalidToken:
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
    default:
        log.Printf("Internal error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Internal server error",
        })
    }
}

func (h *Handler) GetRoles(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "roles": []gin.H{
            {
                "value": "student",
                "label": "Студент",
                "description": "Ищу стажировку",
            },
            {
                "value": "employer",
                "label": "Работодатель",
                "description": "Предлагаю стажировку",
            },
        },
    })
}
