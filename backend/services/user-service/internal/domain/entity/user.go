package entity

import "time"

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	EmailVerified bool     `json:"email_verified"`
	CreatedAt    time.Time `json:"created_dt"`
	UpdatedAt    time.Time `json:"updated_dt"`
}
