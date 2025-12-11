package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	ServerPort string
	DBURL      string
}

func Load() *Config {
	port := getEnv("SERVER_PORT", "8083") // например, 8083 для vacancy-service
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	return &Config{
		ServerPort: port,
		DBURL:      dbURL,
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}
