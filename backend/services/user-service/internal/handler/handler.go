package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"backend/services/user-service/internal/domain/entity"
	"backend/services/user-service/internal/service"
)

type Handler struct {
	userSvc     *service.UserService
	studentSvc  *service.StudentService
	employerSvc *service.EmployerService
	authURL     string
	httpClient  *http.Client
}

func NewHandler(us *service.UserService, ss *service.StudentService, es *service.EmployerService, authURL string) *Handler {
	return &Handler{
		userSvc:     us,
		studentSvc:  ss,
		employerSvc: es,
		authURL:     authURL,
		httpClient:  &http.Client{Timeout: 5 * time.Second},
	}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	// users
	rg.POST("", h.CreateUser)
	rg.GET("/:id", h.GetUserByID)

	// students - protected
	sg := rg.Group("/students")
	sg.Use(h.authMiddleware())
	{
		sg.GET("/me", h.GetStudentMe)
		sg.PUT("/me", h.UpdateStudentMe)
		sg.DELETE("/me", h.DeleteStudentMe)
	}

	// employers - protected
	eg := rg.Group("/employers")
	eg.Use(h.authMiddleware())
	{
		eg.GET("/me", h.GetEmployerMe)
		eg.PUT("/me", h.UpdateEmployerMe)
		eg.DELETE("/me", h.DeleteEmployerMe)
	}
}

// --- DTOs
type createUserReq struct {
	ID       *int64 `json:"id"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=student employer admin"`
}

type updateStudentReq struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	MiddleName  string `json:"middle_name"`
	BirthDate   string `json:"birth_date"`
	PhotoPath   string `json:"photo_path"`
	Facility    string `json:"facility"`
	Speciality  string `json:"speciality"`
	Course      int    `json:"course"`
	Description string `json:"description"`
}

type updateEmployerReq struct {
	CompanyName string `json:"company_name" binding:"required"`
	Description string `json:"description"`
	PhotoPath   string `json:"photo_path"`
	WebsiteURL  string `json:"website_url"`
	Requisites  string `json:"requisites"`
}

// ---- handlers
func (h *Handler) CreateUser(c *gin.Context) {
	var req createUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.userSvc.CreateUser(c.Request.Context(), req.Email, req.Password, req.Role, req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user_id": id})
}

func (h *Handler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	u, err := h.userSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, u)
}

// validate response struct
type validateResp struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	Valid  bool   `json:"valid"`
}

// auth middleware (delegates to auth-service). Also contains a dev fallback.
func (h *Handler) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// DEV shortcut: accept devtoken => user_id=1
		if c.GetHeader("Authorization") == "Bearer devtoken" {
			c.Set("user_id", int64(1))
			c.Set("role", "student")
			c.Next()
			return
		}

		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()
		req, _ := http.NewRequestWithContext(ctx, "POST", h.authURL, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := h.httpClient.Do(req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "auth service error"})
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(resp.Body)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token", "auth_resp": string(b)})
			return
		}
		var vr validateResp
		if err := json.NewDecoder(resp.Body).Decode(&vr); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "bad validate response"})
			return
		}
		if !vr.Valid || vr.UserID == 0 || vr.Role == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("user_id", vr.UserID)
		c.Set("role", vr.Role)
		c.Next()
	}
}

// student handlers
func (h *Handler) GetStudentMe(c *gin.Context) {
	raw, _ := c.Get("user_id")
	uid := raw.(int64)
	st, err := h.studentSvc.GetProfile(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}
	c.JSON(http.StatusOK, st)
}

func (h *Handler) UpdateStudentMe(c *gin.Context) {
	raw, _ := c.Get("user_id")
	uid := raw.(int64)

	var req updateStudentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	st := &entity.Student{
		UserID:      uid,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		MiddleName:  req.MiddleName,
		BirthDate:   req.BirthDate,
		PhotoPath:   req.PhotoPath,
		Facility:    req.Facility,
		Speciality:  req.Speciality,
		Course:      req.Course,
		Description: req.Description,
	}
	if err := h.studentSvc.UpdateProfile(c.Request.Context(), st); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *Handler) DeleteStudentMe(c *gin.Context) {
	raw, _ := c.Get("user_id")
	uid := raw.(int64)
	if err := h.studentSvc.DeleteAccount(c.Request.Context(), uid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// employer handlers
func (h *Handler) GetEmployerMe(c *gin.Context) {
	raw, _ := c.Get("user_id")
	uid := raw.(int64)
	e, err := h.employerSvc.GetProfile(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}
	c.JSON(http.StatusOK, e)
}

func (h *Handler) UpdateEmployerMe(c *gin.Context) {
	raw, _ := c.Get("user_id")
	uid := raw.(int64)

	var req updateEmployerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	e := &entity.Employer{
		UserID:      uid,
		CompanyName: req.CompanyName,
		Description: req.Description,
		PhotoPath:   req.PhotoPath,
		WebsiteURL:  req.WebsiteURL,
		Requisites:  req.Requisites,
	}
	if err := h.employerSvc.UpdateProfile(c.Request.Context(), e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *Handler) DeleteEmployerMe(c *gin.Context) {
	raw, _ := c.Get("user_id")
	uid := raw.(int64)
	if err := h.employerSvc.DeleteAccount(c.Request.Context(), uid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
