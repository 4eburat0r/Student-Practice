package service

import (
	"context"

	reqModel "response-service/internal/domain/models/request"
	resModel "response-service/internal/domain/models/response"
)

type ResponseService interface {
	Create(ctx context.Context, studentID int64, req reqModel.CreateResponseRequest) (*resModel.ResponseResponse, error)
	GetByID(ctx context.Context, studentID, id int64) (*resModel.ResponseResponse, error)
	ListByStudent(ctx context.Context, studentID int64, q reqModel.ResponseListQuery) (*resModel.ResponseListResponse, error)
	Delete(ctx context.Context, studentID, id int64) error
}
