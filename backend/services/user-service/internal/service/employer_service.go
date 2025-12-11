package service

import (
	"context"

	"backend/services/user-service/internal/repository"
	"backend/services/user-service/internal/domain/entity"
)

type EmployerService struct {
	employerRepo repository.EmployerRepository
	userRepo     repository.UserRepository
}

func NewEmployerService(er repository.EmployerRepository, ur repository.UserRepository) *EmployerService {
	return &EmployerService{employerRepo: er, userRepo: ur}
}

func (s *EmployerService) GetProfile(ctx context.Context, userID int64) (*entity.Employer, error) {
	return s.employerRepo.GetByUserID(ctx, userID)
}

func (s *EmployerService) UpdateProfile(ctx context.Context, e *entity.Employer) error {
	return s.employerRepo.UpdateByUserID(ctx, e)
}

func (s *EmployerService) DeleteAccount(ctx context.Context, userID int64) error {
	if err := s.employerRepo.DeleteByUserID(ctx, userID); err != nil {
		return err
	}
	return s.userRepo.Delete(ctx, userID)
}
