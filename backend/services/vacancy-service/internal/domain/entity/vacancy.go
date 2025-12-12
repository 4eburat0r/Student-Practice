package entity

import "time"

type VacancyStatus string
type VacancyFormat string

const (
	VacancyStatusActive   VacancyStatus = "active"
	VacancyStatusInactive VacancyStatus = "inactive"

	VacancyFormatOnsite  VacancyFormat = "onsite"
	VacancyFormatRemote  VacancyFormat = "remote"
	VacancyFormatHybrid  VacancyFormat = "hybrid"
)

type Vacancy struct {
	ID         int64          `db:"id"`
	EmployerID int64          `db:"employer_id"`
	Title      string         `db:"title"`
	Salary     float64        `db:"salary"`
	Description string        `db:"description"`
	Hours      int            `db:"hours"`
	Format     VacancyFormat  `db:"format"`
	Status     VacancyStatus  `db:"status"`
	CreatedAt  time.Time      `db:"created_dt"`
	UpdatedAt  time.Time      `db:"updated_dt"`
}
