package service

import (
	"context"

	"response-service/internal/domain/entity"
	reqModel "response-service/internal/domain/models/request"
	resModel "response-service/internal/domain/models/response"
	"response-service/internal/repository"
)

type responseService struct {
	repo repository.ResponseRepository
}

func NewResponseService(repo repository.ResponseRepository) ResponseService {
	return &responseService{repo: repo}
}

func (s *responseService) Create(ctx context.Context, studentID int64, req reqModel.CreateResponseRequest) (*resModel.ResponseResponse, error) {
	resp := &entity.Response{
		VacancyID: req.VacancyID,
		StudentID: studentID,
		Status: entity.ResponseStatusReview,
		Message:   req.Message,
	}
	if err := s.repo.Create(ctx, resp); err != nil {
		return nil, err
	}
	return toResponse(resp), nil
}

func (s *responseService) GetByID(ctx context.Context, studentID, id int64) (*resModel.ResponseResponse, error) {
	resp, err := s.repo.GetByID(ctx, id, studentID)
	if err != nil {
		return nil, err
	}
	return toResponse(resp), nil
}

func (s *responseService) ListByStudent(ctx context.Context, studentID int64, q reqModel.ResponseListQuery) (*resModel.ResponseListResponse, error) {
	items, total, err := s.repo.ListByStudent(ctx, studentID, q.Limit, q.Offset)
	if err != nil {
		return nil, err
	}
	out := &resModel.ResponseListResponse{
		Items: make([]resModel.ResponseResponse, 0, len(items)),
		Total: total,
	}
	for i := range items {
		out.Items = append(out.Items, *toResponse(&items[i]))
	}
	return out, nil
}

func (s *responseService) Delete(ctx context.Context, studentID, id int64) error {
	return s.repo.Delete(ctx, id, studentID)
}

func toResponse(r *entity.Response) *resModel.ResponseResponse {
	return &resModel.ResponseResponse{
		ID:        r.ID,
		VacancyID: r.VacancyID,
		StudentID: r.StudentID,
		Status:    string(r.Status),
		Message:   r.Message,
		CreatedAt: r.CreatedAt,
	}
}
