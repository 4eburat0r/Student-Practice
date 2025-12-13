package app

import (
	"context"
	"fmt"

	"auth-service/internal/config"
	"auth-service/internal/repository/postgres"
	"auth-service/internal/service"
	syncsvc "auth-service/internal/service/sync"
	"auth-service/pkg/database"
	"auth-service/pkg/hash"
	"auth-service/pkg/jwt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependencies struct {
	DB          *pgxpool.Pool
	JWTManager  *jwt.JWTManager
	AuthService service.AuthService
}

func NewDependencies(cfg *config.Config) (*Dependencies, error) {
	ctx := context.Background()
	db, err := database.NewPostgresPool(ctx, cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to init database: %w", err)
	}

	userRepo := postgres.NewUserRepository(db)
	studentRepo := postgres.NewStudentRepository(db)
	employerRepo := postgres.NewEmployerRepository(db)

	jwtManager := jwt.NewJWTManager(cfg.JWT.Secret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)

	hasher := hash.NewBcryptHasher()

	// sync клиент в user-service
	userSync := syncsvc.NewUserSyncClient(cfg.Sync.UserServiceURL)

	authService := service.NewAuthService(userRepo, studentRepo, employerRepo, jwtManager, hasher, userSync)

	return &Dependencies{
		DB:          db,
		JWTManager:  jwtManager,
		AuthService: *authService,
	}, nil
}
