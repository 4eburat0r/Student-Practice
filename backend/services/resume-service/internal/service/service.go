package service

import (
	"context"
	"resume-service/internal/domain/models/request"
	"resume-service/internal/domain/models/response"
)

type Service interface {
	Create(ctx context.Context, studentID int, req *request.CreateResumeRequest) (*response.ResumeResponse, error)
	GetByID(ctx context.Context, id int) (*response.ResumeResponse, error)
	GetMyResumes(ctx context.Context, studentID int) (*response.ResumeListResponse, error)
	GetFeed(ctx context.Context, limit, offset int) (*response.ResumeListResponse, error)
	Update(ctx context.Context, id, studentID int, req *request.UpdateResumeRequest) (*response.ResumeResponse, error)
	Delete(ctx context.Context, id, studentID int) error
	Publish(ctx context.Context, id, studentID int) error
}
