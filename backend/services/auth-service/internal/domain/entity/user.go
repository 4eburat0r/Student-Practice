package entity

import "time"

type User struct {
    ID           int       `json:"id" db:"id"`
    Email        string    `json:"email" db:"email"`
    PasswordHash string    `json:"-" db:"password_hash"`
    Role         string    `json:"role" db:"role"`
    CreatedDt    time.Time `json:"created_at" db:"created_dt"`
    UpdatedDt    time.Time `json:"updated_at" db:"updated_dt"`
}

const (
    RoleStudent  = "student"
    RoleEmployer = "employer"
    RoleAdmin    = "admin"
)

func (u *User) IsStudent() bool {
    return u.Role == RoleStudent
}

func (u *User) IsEmployer() bool {
    return u.Role == RoleEmployer
}

func (u *User) IsAdmin() bool {
    return u.Role == RoleAdmin
}
