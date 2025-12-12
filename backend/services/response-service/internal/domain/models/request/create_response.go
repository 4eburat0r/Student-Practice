package request

type CreateResponseRequest struct {
	VacancyID int64  `json:"vacancy_id" binding:"required,gt=0"`
	Message   string `json:"message" binding:"omitempty,max=2000"`
}
