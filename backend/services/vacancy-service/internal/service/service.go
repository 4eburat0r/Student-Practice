package service

import (
	"context"

	reqModel "vacancy-service/internal/domain/models/request"
	resModel "vacancy-service/internal/domain/models/response"
)

type VacancyService interface {
	Create(ctx context.Context, employerID int64, req reqModel.CreateVacancyRequest) (*resModel.VacancyResponse, error)
	GetByID(ctx context.Context, id int64) (*resModel.VacancyResponse, error)
	GetByEmployer(ctx context.Context, employerID int64, limit, offset int) (*resModel.VacancyListResponse, error)
	GetFeed(ctx context.Context, q reqModel.VacancyFeedQuery) (*resModel.VacancyListResponse, error)
	Update(ctx context.Context, employerID, id int64, req reqModel.UpdateVacancyRequest) (*resModel.VacancyResponse, error)
	Delete(ctx context.Context, employerID, id int64) error
	Publish(ctx context.Context, employerID, id int64) error
}
