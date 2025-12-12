package repository

import (
	"context"

	"response-service/internal/domain/entity"
)

type ResponseRepository interface {
	Create(ctx context.Context, r *entity.Response) error
	GetByID(ctx context.Context, id int64, studentID int64) (*entity.Response, error)
	ListByStudent(ctx context.Context, studentID int64, limit, offset int) ([]entity.Response, int, error)
	Delete(ctx context.Context, id int64, studentID int64) error
}
