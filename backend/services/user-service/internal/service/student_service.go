package service

import (
	"context"

	"backend/services/user-service/internal/repository"
	"backend/services/user-service/internal/domain/entity"
)

type StudentService struct {
	studentRepo repository.StudentRepository
	userRepo    repository.UserRepository
}

func NewStudentService(sr repository.StudentRepository, ur repository.UserRepository) *StudentService {
	return &StudentService{studentRepo: sr, userRepo: ur}
}

func (s *StudentService) GetProfile(ctx context.Context, userID int64) (*entity.Student, error) {
	return s.studentRepo.GetByUserID(ctx, userID)
}

func (s *StudentService) UpdateProfile(ctx context.Context, st *entity.Student) error {
	return s.studentRepo.UpdateByUserID(ctx, st)
}

func (s *StudentService) DeleteAccount(ctx context.Context, userID int64) error {
	if err := s.studentRepo.DeleteByUserID(ctx, userID); err != nil {
		return err
	}
	return s.userRepo.Delete(ctx, userID)
}
