package repository

import (
	"context"
	"resume-service/internal/domain/entity"
)

type ResumeRepository interface {
	Create(ctx context.Context, resume *entity.Resume) error
	GetByID(ctx context.Context, id int) (*entity.Resume, error)
	GetByStudentID(ctx context.Context, studentID int) ([]*entity.Resume, error)
	GetFeed(ctx context.Context, limit, offset int) ([]*entity.Resume, error)
	Update(ctx context.Context, resume *entity.Resume) error
	Delete(ctx context.Context, id int) error
	Publish(ctx context.Context, id int) error
}

type ResumeCache interface {
	Get(ctx context.Context, key string) (*entity.Resume, error)
	Set(ctx context.Context, key string, resume *entity.Resume) error
	Delete(ctx context.Context, key string) error
	DeletePattern(ctx context.Context, pattern string) error
}
