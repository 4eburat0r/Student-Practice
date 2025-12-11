package app

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"vacancy-service/internal/config"
	"vacancy-service/internal/handler"
	"vacancy-service/internal/server"
	"vacancy-service/internal/service"
	"vacancy-service/internal/repository/postgres"
	"vacancy-service/pkg/logger"
	"vacancy-service/pkg/validator"
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
		log.Error("failed to connect to db", zap.Error(err))
		return nil, err
	}

	v := validator.New()

	// DI: repo -> service -> handler
	vacRepo := postgres.NewVacancyRepository(db)
	vacService := service.NewVacancyService(vacRepo)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	h := handler.NewHandler(vacService, log, v)
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
