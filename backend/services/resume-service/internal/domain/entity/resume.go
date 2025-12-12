package entity

import "time"

type ResumeStatus string

const (
	ResumeStatusActive   ResumeStatus = "active"
	ResumeStatusInactive ResumeStatus = "inactive"
)

type Resume struct {
	ID        int          `json:"id" db:"id"`
	StudentID int          `json:"student_id" db:"student_id"`
	Title     string       `json:"title" db:"title"`
	Summary   string       `json:"summary" db:"summary"`
	Status    ResumeStatus `json:"status" db:"status"`
	CreatedAt time.Time    `json:"created_at" db:"created_dt"`
	UpdatedAt time.Time    `json:"updated_at" db:"updated_dt"`
}

func (r *Resume) IsActive() bool {
	return r.Status == ResumeStatusActive
}

func (r *Resume) Publish() {
	r.Status = ResumeStatusActive
}

func (r *Resume) Unpublish() {
	r.Status = ResumeStatusInactive
}

func (r *Resume) IsOwnedBy(studentID int) bool {
	return r.StudentID == studentID
}
