package service

import (
    "auth-service/internal/domain/models/request"
    "auth-service/internal/domain/models/response"
    "auth-service/pkg/jwt"
    "context"
)

type Auth interface {
    Register(ctx context.Context, req *request.RegisterRequest) error
    Login(ctx context.Context, req *request.LoginRequest) (*response.AuthResponse, error)
    RefreshToken(ctx context.Context, req *request.RefreshTokenRequest) (*response.AuthResponse, error)
    ValidateToken(ctx context.Context, token string) (*jwt.Claims, error)
}
