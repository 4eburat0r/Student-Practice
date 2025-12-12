package request

type CreateVacancyRequest struct {
	Title       string  `json:"title" binding:"required,min=3,max=255"`
	Salary      float64 `json:"salary" binding:"gte=0"`
	Description string  `json:"description" binding:"required,min=10"`
	Hours       int     `json:"hours" binding:"required,gte=2,lte=8"`
	Format      string  `json:"format" binding:"required,oneof=onsite remote hybrid"`
}
