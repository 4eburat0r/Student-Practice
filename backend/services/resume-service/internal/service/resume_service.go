package service

import (
	"context"
	"fmt"
	"log"

	"resume-service/internal/domain/entity"
	"resume-service/internal/domain/errors"
	"resume-service/internal/domain/models/request"
	"resume-service/internal/domain/models/response"
	"resume-service/internal/repository"
)

type ResumeService struct {
	repo  repository.ResumeRepository
	cache repository.ResumeCache
}

func NewResumeService(repo repository.ResumeRepository, cache repository.ResumeCache) *ResumeService {
	return &ResumeService{
		repo:  repo,
		cache: cache,
	}
}

func (s *ResumeService) Create(ctx context.Context, studentID int, req *request.CreateResumeRequest) (*response.ResumeResponse, error) {
	resume := &entity.Resume{
		StudentID: studentID,
		Title:     req.Title,
		Summary:   req.Summary,
		Status:    entity.ResumeStatusInactive,
	}

	if err := s.repo.Create(ctx, resume); err != nil {
		return nil, fmt.Errorf("failed to create resume: %w", err)
	}

	_ = s.cache.DeletePattern(ctx, fmt.Sprintf("resume:student:%d:*", studentID))

	return response.ToResumeResponse(resume), nil
}

func (s *ResumeService) GetByID(ctx context.Context, id int) (*response.ResumeResponse, error) {
	cacheKey := fmt.Sprintf("resume:%d", id)
	cachedResume, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		log.Printf("Cache error: %v", err)
	}
	if cachedResume != nil {
		return response.ToResumeResponse(cachedResume), nil
	}

	resume, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.cache.Set(ctx, cacheKey, resume); err != nil {
		log.Printf("Failed to cache resume: %v", err)
	}

	return response.ToResumeResponse(resume), nil
}

func (s *ResumeService) GetMyResumes(ctx context.Context, studentID int) (*response.ResumeListResponse, error) {
	resumes, err := s.repo.GetByStudentID(ctx, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student resumes: %w", err)
	}

	resumeResponses := make([]response.ResumeResponse, 0, len(resumes))
	for _, resume := range resumes {
		resumeResponses = append(resumeResponses, *response.ToResumeResponse(resume))
	}

	return &response.ResumeListResponse{
		Resumes: resumeResponses,
		Total:   len(resumeResponses),
	}, nil
}

func (s *ResumeService) GetFeed(ctx context.Context, limit, offset int) (*response.ResumeListResponse, error) {
	resumes, err := s.repo.GetFeed(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get resume feed: %w", err)
	}

	resumeResponses := make([]response.ResumeResponse, 0, len(resumes))
	for _, resume := range resumes {
		resumeResponses = append(resumeResponses, *response.ToResumeResponse(resume))
	}

	return &response.ResumeListResponse{
		Resumes: resumeResponses,
		Total:   len(resumeResponses),
	}, nil
}

func (s *ResumeService) Update(ctx context.Context, id, studentID int, req *request.UpdateResumeRequest) (*response.ResumeResponse, error) {
	resume, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !resume.IsOwnedBy(studentID) {
		return nil, errors.ErrForbidden
	}

	if req.Title != "" {
		resume.Title = req.Title
	}
	if req.Summary != "" {
		resume.Summary = req.Summary
	}

	if err := s.repo.Update(ctx, resume); err != nil {
		return nil, fmt.Errorf("failed to update resume: %w", err)
	}

	cacheKey := fmt.Sprintf("resume:%d", id)
	_ = s.cache.Delete(ctx, cacheKey)
	_ = s.cache.DeletePattern(ctx, fmt.Sprintf("resume:student:%d:*", studentID))

	return response.ToResumeResponse(resume), nil
}

func (s *ResumeService) Delete(ctx context.Context, id, studentID int) error {
	resume, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !resume.IsOwnedBy(studentID) {
		return errors.ErrForbidden
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete resume: %w", err)
	}

	cacheKey := fmt.Sprintf("resume:%d", id)
	_ = s.cache.Delete(ctx, cacheKey)
	_ = s.cache.DeletePattern(ctx, fmt.Sprintf("resume:student:%d:*", studentID))

	return nil
}

func (s *ResumeService) Publish(ctx context.Context, id, studentID int) error {
	resume, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !resume.IsOwnedBy(studentID) {
		return errors.ErrForbidden
	}

	if err := s.repo.Publish(ctx, id); err != nil {
		return fmt.Errorf("failed to publish resume: %w", err)
	}

	cacheKey := fmt.Sprintf("resume:%d", id)
	_ = s.cache.Delete(ctx, cacheKey)

	return nil
}
