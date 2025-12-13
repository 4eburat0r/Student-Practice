package response

import "time"

type VacancyResponse struct {
	ID          int64     `json:"id"`
	EmployerID  int64     `json:"employer_id"`
	Title       string    `json:"title"`
	Salary      float64   `json:"salary"`
	Description string    `json:"description"`
	Hours       int       `json:"hours"`
	Format      string    `json:"format"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
