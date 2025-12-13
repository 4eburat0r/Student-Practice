services/response-service/
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── app/
│   │   ├── app.go
│   │   └── dependencies.go (можно не заводить отдельно)
│   ├── config/
│   │   └── config.go
│   ├── domain/
│   │   ├── entity/
│   │   │   └── response.go
│   │   ├── errors/
│   │   │   └── errors.go        # можно переиспользовать общую схему как в vacancy
│   │   └── models/
│   │       ├── request/
│   │       │   ├── create_response.go
│   │       │   └── list_query.go
│   │       └── response/
│   │           ├── response.go
│   │           └── response_list.go
│   ├── repository/
│   │   ├── repository.go
│   │   └── postgres/
│   │       └── response_repository.go
│   ├── service/
│   │   ├── service.go
│   │   └── response_service.go
│   ├── handler/
│   │   ├── middleware/
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   ├── logger.go
│   │   │   └── recovery.go
│   │   ├── response.go
│   │   ├── handler.go
│   │   └── routes.go (по желанию, можно всё в handler.go)
│   └── server/
│       └── server.go
│
├── pkg/
│   ├── logger/
│   │   └── logger.go
│   ├── validator/
│   │   └── validator.go
│   └── response/
│       ├── success.go
│       └── error.go
│
├── migrations/
│   ├── 001_create_responses_table.up.sql
│   └── 001_create_responses_table.down.sql
│
├── build/
│   ├── Dockerfile
│   └── docker-compose.yml
│
├── configs/
│   └── .env.example
├── scripts/
│   ├── migrate-up.sh
│   └── migrate-up.sh
├── .env
├── Makefile
├── go.mod
└── go.sum
