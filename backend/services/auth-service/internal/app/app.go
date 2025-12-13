package app

import (
    "auth-service/internal/config"
    "auth-service/internal/handler"
    "auth-service/internal/handler/middleware"

    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
    Config *config.Config
    DB     *pgxpool.Pool
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

    handler.RegisterRoutes(router, deps.AuthService, deps.JWTManager, cfg.JWT.Secret)

    return &App{
        Config: cfg,
        DB:     deps.DB,
        Router: router,
    }, nil
}

func (a *App) Close() {
    if a.DB != nil {
        a.DB.Close()
    }
}