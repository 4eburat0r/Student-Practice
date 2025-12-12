package config

import "os"

type Config struct {
	Port           string
	DatabaseURL    string
	AuthIntrospect string
}

func LoadFromEnv() *Config {
	return &Config{
		Port:           getenv("PORT", "8081"),
		DatabaseURL:    getenv("DATABASE_URL", "postgres://postgres:postgres@db:5432/userdb?sslmode=disable"),
		AuthIntrospect: getenv("AUTH_INTROSPECT_URL", "http://auth-service:8080/introspect"),
	}
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
