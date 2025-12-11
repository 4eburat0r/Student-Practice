package repository

import (
	"context"

	"backend/services/user-service/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, email, passwordHash, role string) (int64, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, string, error) // returns user and passwordHash
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	Delete(ctx context.Context, id int64) error
}

type StudentRepository interface {
	Create(ctx context.Context, s *entity.Student) (int64, error)
	GetByUserID(ctx context.Context, userID int64) (*entity.Student, error)
	UpdateByUserID(ctx context.Context, s *entity.Student) error
	DeleteByUserID(ctx context.Context, userID int64) error
}

type EmployerRepository interface {
	Create(ctx context.Context, e *entity.Employer) (int64, error)
	GetByUserID(ctx context.Context, userID int64) (*entity.Employer, error)
	UpdateByUserID(ctx context.Context, e *entity.Employer) error
	DeleteByUserID(ctx context.Context, userID int64) error
}
