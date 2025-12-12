FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o response-service ./cmd/api

FROM alpine:3.19

WORKDIR /app

ENV GIN_MODE=release

COPY --from=builder /app/response-service /app/response-service
COPY ./migrations /app/migrations

EXPOSE 8084

CMD ["/app/response-service"]
