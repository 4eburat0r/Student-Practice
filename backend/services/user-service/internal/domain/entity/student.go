package entity

import "time"

type Student struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	MiddleName string    `json:"middle_name,omitempty"`
	PhotoPath  string    `json:"photo_path,omitempty"`
	Facility   string    `json:"facility,omitempty"`
	Course     int       `json:"course,omitempty"`
	Description string   `json:"description,omitempty"`
	CreatedAt  time.Time `json:"created_dt"`
	UpdatedAt  time.Time `json:"updated_dt"`
}
