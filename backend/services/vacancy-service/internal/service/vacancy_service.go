package service

import (
	"context"

	"vacancy-service/internal/domain/entity"
	domainErrors "vacancy-service/internal/domain/errors"
	reqModel "vacancy-service/internal/domain/models/request"
	resModel "vacancy-service/internal/domain/models/response"
	"vacancy-service/internal/repository"
)

type vacancyService struct {
	repo repository.VacancyRepository
}

func NewVacancyService(repo repository.VacancyRepository) VacancyService {
	return &vacancyService{repo: repo}
}

func (s *vacancyService) Create(ctx context.Context, employerID int64, req reqModel.CreateVacancyRequest) (*resModel.VacancyResponse, error) {
	v := &entity.Vacancy{
		EmployerID:  employerID,
		Title:       req.Title,
		Salary:      req.Salary,
		Description: req.Description,
		Hours:       req.Hours,
		Format:      entity.VacancyFormat(req.Format),
		Status:      entity.VacancyStatusInactive,
	}
	if err := s.repo.Create(ctx, v); err != nil {
		return nil, err
	}
	return toVacancyResponse(v), nil
}

func (s *vacancyService) GetByID(ctx context.Context, id int64) (*resModel.VacancyResponse, error) {
	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toVacancyResponse(v), nil
}

func (s *vacancyService) GetByEmployer(ctx context.Context, employerID int64, limit, offset int) (*resModel.VacancyListResponse, error) {
	items, total, err := s.repo.GetByEmployer(ctx, employerID, limit, offset)
	if err != nil {
		return nil, err
	}
	resp := &resModel.VacancyListResponse{
		Items: make([]resModel.VacancyResponse, 0, len(items)),
		Total: total,
	}
	for i := range items {
		resp.Items = append(resp.Items, *toVacancyResponse(&items[i]))
	}
	return resp, nil
}

func (s *vacancyService) GetFeed(ctx context.Context, q reqModel.VacancyFeedQuery) (*resModel.VacancyListResponse, error) {
	var format *string
	if q.Format != "" {
		format = &q.Format
	}
	items, total, err := s.repo.GetFeed(ctx, format, q.Limit, q.Offset)
	if err != nil {
		return nil, err
	}
	resp := &resModel.VacancyListResponse{
		Items: make([]resModel.VacancyResponse, 0, len(items)),
		Total: total,
	}
	for i := range items {
		resp.Items = append(resp.Items, *toVacancyResponse(&items[i]))
	}
	return resp, nil
}

func (s *vacancyService) Update(ctx context.Context, employerID, id int64, req reqModel.UpdateVacancyRequest) (*resModel.VacancyResponse, error) {
	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if v.EmployerID != employerID {
		return nil, domainErrors.ErrForbidden
	}
	if req.Title != nil {
		v.Title = *req.Title
	}
	if req.Salary != nil {
		v.Salary = *req.Salary
	}
	if req.Description != nil {
		v.Description = *req.Description
	}
	if req.Hours != nil {
		v.Hours = *req.Hours
	}
	if req.Format != nil {
		v.Format = entity.VacancyFormat(*req.Format)
	}
	if err := s.repo.Update(ctx, v); err != nil {
		return nil, err
	}
	return toVacancyResponse(v), nil
}

func (s *vacancyService) Delete(ctx context.Context, employerID, id int64) error {
	return s.repo.Delete(ctx, id, employerID)
}

func (s *vacancyService) Publish(ctx context.Context, employerID, id int64) error {
	return s.repo.UpdateStatus(ctx, id, employerID, entity.VacancyStatusActive)
}

func toVacancyResponse(v *entity.Vacancy) *resModel.VacancyResponse {
	return &resModel.VacancyResponse{
		ID:          v.ID,
		EmployerID:  v.EmployerID,
		Title:       v.Title,
		Salary:      v.Salary,
		Description: v.Description,
		Hours:       v.Hours,
		Format:      string(v.Format),
		Status:      string(v.Status),
		CreatedAt:   v.CreatedAt,
		UpdatedAt:   v.UpdatedAt,
	}
}
