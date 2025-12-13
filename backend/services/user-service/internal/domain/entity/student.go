package entity

type Student struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	MiddleName  string `json:"middle_name"`
	BirthDate   string `json:"birth_date"` // YYYY-MM-DD
	PhotoPath   string `json:"photo_path"`
	Facility    string `json:"facility"`
	Speciality  string `json:"speciality"`
	Course      int    `json:"course"`
	Description string `json:"description"`
}
