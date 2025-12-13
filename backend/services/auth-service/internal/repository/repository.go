package repository

import (
    "auth-service/internal/domain/entity"
    "context"
)

type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    GetByEmail(ctx context.Context, email string) (*entity.User, error)
    GetByID(ctx context.Context, id int) (*entity.User, error)
    Update(ctx context.Context, user *entity.User) error
}

type StudentRepository interface {
    Create(ctx context.Context, userID int) (*entity.Student, error)
    GetByID(ctx context.Context, id int) (*entity.Student, error)
    GetByUserID(ctx context.Context, userID int) (*entity.Student, error)
    Update(ctx context.Context, student *entity.Student) error
    Delete(ctx context.Context, id int) error
}

type EmployerRepository interface {
    Create(ctx context.Context, userID int) (*entity.Employer, error)
    GetByID(ctx context.Context, id int) (*entity.Employer, error)
    GetByUserID(ctx context.Context, userID int) (*entity.Employer, error)
    Update(ctx context.Context, employer *entity.Employer) error
    Delete(ctx context.Context, id int) error
}
