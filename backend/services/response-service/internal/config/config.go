package config

import (
	"log"
	"os"
)

type Config struct {
	ServerPort string
	DBURL      string
}

func Load() *Config {
	port := getEnv("SERVER_PORT", "8084")
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
