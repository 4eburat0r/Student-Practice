package app

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"response-service/internal/config"
	"response-service/internal/handler"
	"response-service/internal/repository/postgres"
	"response-service/internal/server"
	"response-service/internal/service"
	"response-service/pkg/logger"
	"response-service/pkg/validator"
)

type App struct {
	cfg    *config.Config
	logger *zap.Logger
	db     *pgxpool.Pool
	server *server.Server
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	log, err := logger.New()
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.New(ctx, cfg.DBURL)
	if err != nil {
		log.Error("failed to connect db", zap.Error(err))
		return nil, err
	}

	v := validator.New()

	respRepo := postgres.NewResponseRepository(db)
	respService := service.NewResponseService(respRepo)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	h := handler.NewHandler(respService, log, v)
	h.RegisterRoutes(router)

	srv := server.New(router, cfg.ServerPort, log)

	return &App{
		cfg:    cfg,
		logger: log,
		db:     db,
		server: srv,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.server.Run(ctx)
}

func (a *App) Close() {
	a.db.Close()
	_ = a.logger.Sync()
}
