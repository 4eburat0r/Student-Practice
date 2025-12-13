package service

import (
	"context"
	"strings"

	"backend/services/user-service/internal/domain/entity"
	"backend/services/user-service/internal/repository"
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
	if err := s.studentRepo.UpdateByUserID(ctx, st); err != nil {
		// Если профиль ещё не создан, создаём его
		if strings.Contains(err.Error(), "not found") {
			_, createErr := s.studentRepo.Create(ctx, st)
			return createErr
		}
		return err
	}
	return nil
}

func (s *StudentService) DeleteAccount(ctx context.Context, userID int64) error {
	if err := s.studentRepo.DeleteByUserID(ctx, userID); err != nil {
		return err
	}
	return s.userRepo.Delete(ctx, userID)
}
