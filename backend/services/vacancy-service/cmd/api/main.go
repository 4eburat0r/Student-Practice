package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"vacancy-service/internal/app"
	"vacancy-service/internal/config"
)

func main() {
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	application, err := app.New(ctx, cfg)
	if err != nil {
		os.Exit(1)
	}
	defer application.Close()

	if err := application.Run(ctx); err != nil {
		os.Exit(1)
	}
}
