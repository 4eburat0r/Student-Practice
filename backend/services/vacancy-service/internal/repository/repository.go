package repository

import (
	"context"

	"vacancy-service/internal/domain/entity"
)

type VacancyRepository interface {
	Create(ctx context.Context, v *entity.Vacancy) error
	GetByID(ctx context.Context, id int64) (*entity.Vacancy, error)
	GetByEmployer(ctx context.Context, employerID int64, limit, offset int) ([]entity.Vacancy, int, error)
	GetFeed(ctx context.Context, format *string, limit, offset int) ([]entity.Vacancy, int, error)
	Update(ctx context.Context, v *entity.Vacancy) error
	Delete(ctx context.Context, id, employerID int64) error
	UpdateStatus(ctx context.Context, id, employerID int64, status entity.VacancyStatus) error
}
