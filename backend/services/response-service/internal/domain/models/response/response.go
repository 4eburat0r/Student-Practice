package response

import "time"

type ResponseResponse struct {
	ID        int64     `json:"id"`
	VacancyID int64     `json:"vacancy_id"`
	StudentID int64     `json:"student_id"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
