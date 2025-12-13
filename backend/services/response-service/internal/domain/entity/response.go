package entity

import "time"

type ResponseStatus string

const (
	ResponseStatusReview    ResponseStatus = "review"
	ResponseStatusApproved  ResponseStatus = "approved"
	ResponseStatusDiscarded ResponseStatus = "discarded"
	ResponseStatusInvited   ResponseStatus = "invited"
)

type Response struct {
	ID        int64          `db:"id"`
	VacancyID int64          `db:"vacancy_id"`
	StudentID int64          `db:"student_id"`
	Status    ResponseStatus `db:"status"`
	Message   string         `db:"message"`
	CreatedAt time.Time      `db:"created_dt"`
}
