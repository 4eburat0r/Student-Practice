package request

type UpdateVacancyRequest struct {
	Title       *string  `json:"title" binding:"omitempty,min=3,max=255"`
	Salary      *float64 `json:"salary" binding:"omitempty,gte=0"`
	Description *string  `json:"description" binding:"omitempty,min=10"`
	Hours       *int     `json:"hours" binding:"omitempty,gte=2,lte=8"`
	Format      *string  `json:"format" binding:"omitempty,oneof=onsite remote hybrid"`
}
