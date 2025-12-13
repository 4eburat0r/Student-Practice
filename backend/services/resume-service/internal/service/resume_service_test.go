package service_test

import (
	"context"
	"errors"
	"testing"

	"resume-service/internal/domain/entity"
	domainErrors "resume-service/internal/domain/errors"
	"resume-service/internal/domain/models/request"
	"resume-service/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockResumeRepository struct {
	mock.Mock
}

func (m *MockResumeRepository) Create(ctx context.Context, resume *entity.Resume) error {
	args := m.Called(ctx, resume)
	if args.Get(0) == nil {
		resume.ID = 1
	}
	return args.Error(0)
}

func (m *MockResumeRepository) GetByID(ctx context.Context, id int) (*entity.Resume, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Resume), args.Error(1)
}

func (m *MockResumeRepository) GetByStudentID(ctx context.Context, studentID int) ([]*entity.Resume, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Resume), args.Error(1)
}

func (m *MockResumeRepository) GetFeed(ctx context.Context, limit, offset int) ([]*entity.Resume, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Resume), args.Error(1)
}

func (m *MockResumeRepository) Update(ctx context.Context, resume *entity.Resume) error {
	args := m.Called(ctx, resume)
	return args.Error(0)
}

func (m *MockResumeRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockResumeRepository) Publish(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockResumeCache struct {
	mock.Mock
}

func (m *MockResumeCache) Get(ctx context.Context, key string) (*entity.Resume, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Resume), args.Error(1)
}

func (m *MockResumeCache) Set(ctx context.Context, key string, resume *entity.Resume) error {
	args := m.Called(ctx, key, resume)
	return args.Error(0)
}

func (m *MockResumeCache) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockResumeCache) DeletePattern(ctx context.Context, pattern string) error {
	args := m.Called(ctx, pattern)
	return args.Error(0)
}

func TestResumeService_Create(t *testing.T) {
	mockRepo := new(MockResumeRepository)
	mockCache := new(MockResumeCache)
	svc := service.NewResumeService(mockRepo, mockCache)

	ctx := context.Background()
	studentID := 1
	req := &request.CreateResumeRequest{
		Title:   "Go Backend Developer",
		Summary: "Experienced in microservices and PostgreSQL",
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.Resume")).Return(nil)
	mockCache.On("DeletePattern", ctx, mock.Anything).Return(nil)

	resume, err := svc.Create(ctx, studentID, req)

	assert.NoError(t, err)
	assert.NotNil(t, resume)
	assert.Equal(t, req.Title, resume.Title)
	assert.Equal(t, studentID, resume.StudentID)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestResumeService_GetByID_Success(t *testing.T) {
	mockRepo := new(MockResumeRepository)
	mockCache := new(MockResumeCache)
	svc := service.NewResumeService(mockRepo, mockCache)

	ctx := context.Background()
	resumeID := 1

	expectedResume := &entity.Resume{
		ID:        resumeID,
		StudentID: 1,
		Title:     "Go Developer",
		Summary:   "5 years experience",
		Status:    entity.ResumeStatusActive,
	}

	mockCache.On("Get", ctx, "resume:1").Return(nil, errors.New("cache miss"))
	mockRepo.On("GetByID", ctx, resumeID).Return(expectedResume, nil)
	mockCache.On("Set", ctx, "resume:1", expectedResume).Return(nil)

	resume, err := svc.GetByID(ctx, resumeID)

	assert.NoError(t, err)
	assert.NotNil(t, resume)
	assert.Equal(t, expectedResume.Title, resume.Title)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestResumeService_Update_Forbidden(t *testing.T) {
	mockRepo := new(MockResumeRepository)
	mockCache := new(MockResumeCache)
	svc := service.NewResumeService(mockRepo, mockCache)

	ctx := context.Background()
	resumeID := 1
	studentID := 2

	existingResume := &entity.Resume{
		ID:        resumeID,
		StudentID: 1,
		Title:     "Some Title",
		Summary:   "Some Summary",
	}

	req := &request.UpdateResumeRequest{
		Title: "Hacked Title",
	}

	mockRepo.On("GetByID", ctx, resumeID).Return(existingResume, nil)

	resume, err := svc.Update(ctx, resumeID, studentID, req)

	assert.Error(t, err)
	assert.Nil(t, resume)
	assert.Equal(t, domainErrors.ErrForbidden, err)
	mockRepo.AssertExpectations(t)
}
