package entity

import "time"

type Student struct {
    ID          int        `json:"id" db:"id"`
    UserID      int        `json:"user_id" db:"user_id"`
    FirstName   *string    `json:"first_name,omitempty" db:"first_name"`
    LastName    *string    `json:"last_name,omitempty" db:"last_name"`
    MiddleName  *string    `json:"middle_name,omitempty" db:"middle_name"`
    BirthDate   *time.Time `json:"birth_date,omitempty" db:"birth_date"`
    PhotoPath   *string    `json:"photo_path,omitempty" db:"photo_path"`
    Facility    *string    `json:"facility,omitempty" db:"facility"`
    Speciality  *string    `json:"speciality,omitempty" db:"speciality"`
    Course      *int       `json:"course,omitempty" db:"course"`
    Description *string    `json:"description,omitempty" db:"description"`
    CreatedDt   time.Time  `json:"created_at" db:"created_dt"`
    UpdatedDt   time.Time  `json:"updated_at" db:"updated_dt"`
}

func (s *Student) HasFullName() bool {
    return s.FirstName != nil && s.LastName != nil && s.MiddleName != nil
}

func (s *Student) GetFullName() string {
    if !s.HasFullName() {
        return ""
    }
    return *s.LastName + " " + *s.FirstName + " " + *s.MiddleName
}

func (s *Student) IsProfileComplete() bool {
    return s.HasFullName() && 
           s.Facility != nil && 
           s.Course != nil && 
           s.Description != nil
}

func (s *Student) HasPhoto() bool {
    return s.PhotoPath != nil
}
