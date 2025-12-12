package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort int
	Env        string

	AuthServiceURL     string
	UserServiceURL     string
	ResumeServiceURL   string
	VacancyServiceURL  string
	ResponseServiceURL string

	JWTSecret string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	port, err := strconv.Atoi(getEnv("SERVER_PORT", "8000"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
	}

	return &Config{
		ServerPort: port,
		Env:        getEnv("ENV", "development"),

		AuthServiceURL:     getEnv("AUTH_SERVICE_URL", "http://auth-service:8080"),
		UserServiceURL:     getEnv("USER_SERVICE_URL", "http://user-service:8081"),
		ResumeServiceURL:   getEnv("RESUME_SERVICE_URL", "http://resume-service:8082"),
		VacancyServiceURL:  getEnv("VACANCY_SERVICE_URL", "http://vacancy-service:8083"),
		ResponseServiceURL: getEnv("RESPONSE_SERVICE_URL", "http://response-service:8084"),

		JWTSecret: getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
