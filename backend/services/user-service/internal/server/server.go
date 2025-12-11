package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"backend/services/user-service/internal/config"
	"backend/services/user-service/internal/handler"
	"backend/services/user-service/internal/repository/postgres"
	"backend/services/user-service/internal/service"
)

type Server struct {
	router *gin.Engine
	pool   *pgxpool.Pool
	cfg    *config.Config
}

func NewServer(r *gin.Engine, pool *pgxpool.Pool, cfg *config.Config) *Server {
	return &Server{router: r, pool: pool, cfg: cfg}
}

func (s *Server) Run() error {
	// repositories
	userRepo := postgres.NewUserRepo(s.pool)
	studentRepo := postgres.NewStudentRepo(s.pool)
	employerRepo := postgres.NewEmployerRepo(s.pool)

	// services
	userSvc := service.NewUserService(userRepo)
	studentSvc := service.NewStudentService(studentRepo, userRepo)
	employerSvc := service.NewEmployerService(employerRepo, userRepo)

	// handlers
	h := handler.NewHandler(userSvc, studentSvc, employerSvc, s.cfg.AuthIntrospect)
	api := s.router.Group("/users")
	h.RegisterRoutes(api)

	// health
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	addr := fmt.Sprintf(":%s", s.cfg.Port)
	return s.router.Run(addr)
}
