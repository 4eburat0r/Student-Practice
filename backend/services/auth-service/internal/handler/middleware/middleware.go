package middleware

import (
    "auth-service/pkg/jwt"
    "log"
)

type Middleware struct {
    jwtManager *jwt.JWTManager
    logger     *log.Logger
}

func NewMiddleware(jwtManager *jwt.JWTManager, logger *log.Logger) *Middleware {
    return &Middleware{
        jwtManager: jwtManager,
        logger:     logger,
    }
}
