package response

import (
	"resume-service/internal/domain/entity"
	"time"
)

type ResumeResponse struct {
	ID        int                 `json:"id"`
	StudentID int                 `json:"student_id"`
	Title     string              `json:"title"`
	Summary   string              `json:"summary"`
	Status    entity.ResumeStatus `json:"status"`
	CreatedDt time.Time           `json:"created_dt"`
	UpdatedDt time.Time           `json:"updated_dt"`
}

func ToResumeResponse(resume *entity.Resume) *ResumeResponse {
	return &ResumeResponse{
		ID:        resume.ID,
		StudentID: resume.StudentID,
		Title:     resume.Title,
		Summary:   resume.Summary,
		Status:    resume.Status,
		CreatedDt: resume.CreatedAt,
		UpdatedDt: resume.UpdatedAt,
	}
}
