package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"backend/services/user-service/internal/domain/entity"
	"backend/services/user-service/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, email, password, role string, idOverride *int64) (int64, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	return s.repo.Create(ctx, email, string(hash), role, idOverride)
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*entity.User, string, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
