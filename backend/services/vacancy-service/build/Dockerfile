# build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o vacancy-service ./cmd/api

# run stage
FROM alpine:3.19

WORKDIR /app

ENV GIN_MODE=release

COPY --from=builder /app/vacancy-service /app/vacancy-service
COPY ./migrations /app/migrations

EXPOSE 8083

CMD ["/app/vacancy-service"]
