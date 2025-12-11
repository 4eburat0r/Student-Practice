APP_NAME=vacancy-service

.PHONY: run migrate-up migrate-down build docker-build

run:
	go run ./cmd/api

migrate-up:
	./scripts/migrate-up.sh

migrate-down:
	./scripts/migrate-down.sh

build:
	go build -o bin/$(APP_NAME) ./cmd/api

docker-build:
	docker build -t $(APP_NAME) -f build/Dockerfile .
