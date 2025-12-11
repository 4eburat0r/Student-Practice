package entity

import "time"

type Employer struct {
    ID          int       `json:"id" db:"id"`
    UserID      int       `json:"user_id" db:"user_id"`
    CompanyName *string   `json:"company_name,omitempty" db:"company_name"`
    Description *string   `json:"description,omitempty" db:"description"`
    PhotoPath   *string   `json:"photo_path,omitempty" db:"photo_path"`
    WebsiteURL  *string   `json:"website_url,omitempty" db:"website_url"`
    Requisites  *string   `json:"requisites,omitempty" db:"requisites"`
    CreatedDt   time.Time `json:"created_at" db:"created_dt"`
    UpdatedDt   time.Time `json:"updated_at" db:"updated_dt"`
}

func (e *Employer) IsProfileComplete() bool {
    return e.CompanyName != nil && 
           e.Description != nil && 
           e.WebsiteURL != nil && 
           e.Requisites != nil
}

func (e *Employer) HasLogo() bool {
    return e.PhotoPath != nil
}
