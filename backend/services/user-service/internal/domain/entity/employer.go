package entity

type Employer struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	CompanyName string `json:"company_name"`
	Description string `json:"description"`
	PhotoPath   string `json:"photo_path"`
	WebsiteURL  string `json:"website_url"`
	Requisites  string `json:"requisites"`
}
