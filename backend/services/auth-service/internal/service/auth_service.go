package service

import (
	"auth-service/internal/domain/entity"
	"auth-service/internal/domain/errors"
	"auth-service/internal/domain/models/request"
	"auth-service/internal/domain/models/response"
	"auth-service/internal/repository"
	syncsvc "auth-service/internal/service/sync"
	"auth-service/pkg/hash"
	"auth-service/pkg/jwt"
	"context"
	"fmt"
)

type AuthService struct {
	userRepo     repository.UserRepository
	studentRepo  repository.StudentRepository
	employerRepo repository.EmployerRepository
	jwtManager   *jwt.JWTManager
	hasher       hash.Hasher
	userSync     *syncsvc.UserSyncClient
}

func NewAuthService(
	userRepo repository.UserRepository,
	studentRepo repository.StudentRepository,
	employerRepo repository.EmployerRepository,
	jwtManager *jwt.JWTManager,
	hasher hash.Hasher,
	userSync *syncsvc.UserSyncClient,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		studentRepo:  studentRepo,
		employerRepo: employerRepo,
		jwtManager:   jwtManager,
		hasher:       hasher,
		userSync:     userSync,
	}
}

func (s *AuthService) Register(ctx context.Context, req *request.RegisterRequest) error {
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return errors.ErrUserAlreadyExists
	}

	passwordHash, err := s.hasher.Hash(req.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := &entity.User{
		Email:        req.Email,
		PasswordHash: passwordHash,
		Role:         req.Role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// best-effort: синхронизируем в user-service
	if s.userSync != nil {
		if err := s.userSync.CreateUser(ctx, user.ID, req.Email, req.Password, req.Role); err != nil {
			// не падаем, но логируем для дальнейшего анализа
			fmt.Printf("user-service sync failed: %v\n", err)
		}
	}

	if user.IsStudent() {
		_, err = s.studentRepo.Create(ctx, user.ID)
		if err != nil {
			return fmt.Errorf("failed to create student profile: %w", err)
		}
	} else if user.IsEmployer() {
		_, err = s.employerRepo.Create(ctx, user.ID)
		if err != nil {
			return fmt.Errorf("failed to create employer profile: %w", err)
		}
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, req *request.LoginRequest) (*response.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return nil, errors.ErrInvalidCredentials
	}

	if err := s.hasher.Verify(user.PasswordHash, req.Password); err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	accessToken, expiresAt, err := s.jwtManager.GenerateAccessToken(
		user.ID,
		user.Email,
		user.Role,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &response.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User: response.UserInfo{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *request.RefreshTokenRequest) (*response.AuthResponse, error) {
	claims, err := s.jwtManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil || user == nil {
		return nil, errors.ErrUserNotFound
	}

	accessToken, expiresAt, err := s.jwtManager.GenerateAccessToken(
		user.ID,
		user.Email,
		user.Role,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &response.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User: response.UserInfo{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*jwt.Claims, error) {
	claims, err := s.jwtManager.ValidateAccessToken(token)
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	return claims, nil
}
