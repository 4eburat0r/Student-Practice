package config

import (
    "fmt"
    "os"
    "strconv"
    "time"

    "github.com/joho/godotenv"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    JWT      JWTConfig
}

type ServerConfig struct {
    Host string
    Port string
    Env  string
}

type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
}

type JWTConfig struct {
    Secret     string
    AccessTTL  time.Duration
    RefreshTTL time.Duration
}

func Load() (*Config, error) {
    _ = godotenv.Load()

    accessTTL, _ := strconv.Atoi(getEnv("JWT_ACCESS_TTL_MINUTES", "15"))
    refreshTTL, _ := strconv.Atoi(getEnv("JWT_REFRESH_TTL_HOURS", "168"))

    cfg := &Config{
        Server: ServerConfig{
            Host: getEnv("SERVER_HOST", "0.0.0.0"),
            Port: getEnv("SERVER_PORT", "8080"),
            Env:  getEnv("ENV", "development"),
        },
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnv("DB_PORT", "5432"),
            User:     getEnv("DB_USER", "postgres"),
            Password: getEnv("DB_PASSWORD", "postgres"),
            DBName:   getEnv("DB_NAME", "student_practice"),
            SSLMode:  getEnv("DB_SSLMODE", "disable"),
        },
        JWT: JWTConfig{
            Secret:     getEnv("JWT_SECRET", ""),
            AccessTTL:  time.Duration(accessTTL) * time.Minute,
            RefreshTTL: time.Duration(refreshTTL) * time.Hour,
        },
    }

    if err := cfg.validate(); err != nil {
        return nil, err
    }

    return cfg, nil
}

func (c *Config) validate() error {
    if c.JWT.Secret == "" {
        return fmt.Errorf("JWT_SECRET must be set")
    }

    return nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
