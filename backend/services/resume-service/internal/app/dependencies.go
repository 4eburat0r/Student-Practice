package app

import (
	"context"
	"fmt"
	"log"

	"resume-service/internal/config"
	"resume-service/internal/repository/postgres"
	"resume-service/internal/repository/redis"
	"resume-service/internal/service"
	"resume-service/pkg/jwt"

	"github.com/jackc/pgx/v5/pgxpool"
	redisClient "github.com/redis/go-redis/v9"
)

type Dependencies struct {
	DB            *pgxpool.Pool
	RedisClient   *redisClient.Client
	JWTManager    *jwt.JWTManager
	ResumeService service.Service
}

func NewDependencies(cfg *config.Config) (*Dependencies, error) {
	ctx := context.Background()

	db, err := initPostgres(ctx, cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to init postgres: %w", err)
	}
	log.Println("✅ PostgreSQL connected")

	redisClient, err := initRedis(ctx, cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("failed to init redis: %w", err)
	}
	log.Println("✅ Redis connected")

	jwtManager := jwt.NewJWTManager(cfg.JWT.Secret)

	resumeRepo := postgres.NewResumeRepository(db)
	resumeCache := redis.NewResumeCache(redisClient)

	resumeService := service.NewResumeService(resumeRepo, resumeCache)

	return &Dependencies{
		DB:            db,
		RedisClient:   redisClient,
		JWTManager:    jwtManager,
		ResumeService: resumeService,
	}, nil
}

func initPostgres(ctx context.Context, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}

	poolCfg.MaxConns = 25
	poolCfg.MinConns = 5

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

func initRedis(ctx context.Context, cfg config.RedisConfig) (*redisClient.Client, error) {
	client := redisClient.NewClient(&redisClient.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return client, nil
}
