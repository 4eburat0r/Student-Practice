package app

import (
	"log"
	"resume-service/internal/config"
	"resume-service/internal/handler"
	"resume-service/internal/handler/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type App struct {
	Config *config.Config
	DB     *pgxpool.Pool
	Redis  *redis.Client
	Router *gin.Engine
}

func NewApp(cfg *config.Config) (*App, error) {
	deps, err := NewDependencies(cfg)
	if err != nil {
		return nil, err
	}

	router := gin.New()

	router.Use(middleware.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "resume-service"})
	})

	h := handler.NewHandler(deps.ResumeService, deps.JWTManager)

	handler.RegisterRoutes(router, h, cfg.JWT.Secret)

	log.Println("✅ Resume Service initialized successfully")

	return &App{
		Config: cfg,
		DB:     deps.DB,
		Redis:  deps.RedisClient,
		Router: router,
	}, nil
}

func (a *App) Close() {
	if a.DB != nil {
		log.Println("Closing database connection...")
		a.DB.Close()
	}
	if a.Redis != nil {
		log.Println("Closing Redis connection...")
		_ = a.Redis.Close()
	}
}
