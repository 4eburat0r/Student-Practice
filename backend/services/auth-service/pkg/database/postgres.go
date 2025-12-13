package database

import (
    "auth-service/internal/config"
    "context"
    "fmt"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(ctx context.Context, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
    )

    poolConfig, err := pgxpool.ParseConfig(dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to parse pool config: %w", err)
    }

    poolConfig.MaxConns = 25
    poolConfig.MinConns = 5
    poolConfig.MaxConnLifetime = 1 * time.Hour
    poolConfig.MaxConnIdleTime = 15 * time.Minute
    poolConfig.HealthCheckPeriod = 1 * time.Minute

    pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to create connection pool: %w", err)
    }

    if err := pool.Ping(ctx); err != nil {
        pool.Close()
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    return pool, nil
}
