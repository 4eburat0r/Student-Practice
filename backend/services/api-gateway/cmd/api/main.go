package main

import (
	"log"

	"api-gateway/internal/config"
	"api-gateway/internal/app"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := app.Run(cfg); err != nil {
		log.Fatalf("Failed to run application: %v", err)
	}
}
