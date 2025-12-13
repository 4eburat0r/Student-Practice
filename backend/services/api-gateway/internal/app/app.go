package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-gateway/internal/config"
	"api-gateway/internal/deps"
	"api-gateway/internal/handler"
	"api-gateway/internal/server"
	"api-gateway/pkg/logger"

	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) error {
	if err := logger.InitLogger(cfg.Env); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	defer logger.Sync()

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	dependencies := deps.NewDependencies(&deps.Config{
		AuthServiceURL:     cfg.AuthServiceURL,
		UserServiceURL:     cfg.UserServiceURL,
		ResumeServiceURL:   cfg.ResumeServiceURL,
		VacancyServiceURL:  cfg.VacancyServiceURL,
		ResponseServiceURL: cfg.ResponseServiceURL,
	})

	h := handler.NewHandler(dependencies)

	router := h.InitRoutes()
	srv := server.NewServer(cfg.ServerPort, router)

	go func() {
		logger.Info("Starting API Gateway",
			logger.String("port", fmt.Sprintf("%d", cfg.ServerPort)),
			logger.String("env", cfg.Env),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", logger.String("error", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	logger.Info("Server exited")
	return nil
}
