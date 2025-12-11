package entity

import "time"

type Employer struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	CompanyName string    `json:"company_name"`
	Description string    `json:"description,omitempty"`
	PhotoPath   string    `json:"photo_path,omitempty"`
	WebsiteURL  string    `json:"website_url,omitempty"`
	Requisites  string    `json:"requisites,omitempty"`
	CreatedAt   time.Time `json:"created_dt"`
	UpdatedAt   time.Time `json:"updated_dt"`
}
