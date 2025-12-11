package main

import (
	"context"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/gin-gonic/gin"

	"backend/services/user-service/internal/config"
	"backend/services/user-service/internal/server"
)

func main() {
	_ = godotenv.Load() // optional

	cfg := config.LoadFromEnv()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()

	r := gin.Default()
	srv := server.NewServer(r, pool, cfg)

	if err := srv.Run(); err != nil {
		log.Fatalf("server run: %v", err)
	}
}
