package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"response-service/internal/app"
	"response-service/internal/config"
)

func main() {
	// загружаем конфиг
	cfg := config.Load()

	// контекст с graceful shutdown
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	// инициализация приложения
	application, err := app.New(ctx, cfg)
	if err != nil {
		os.Exit(1)
	}
	defer application.Close()

	// запуск HTTP-сервера
	if err := application.Run(ctx); err != nil {
		os.Exit(1)
	}
}
