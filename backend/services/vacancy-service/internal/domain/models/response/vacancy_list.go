package response

type VacancyListResponse struct {
	Items []VacancyResponse `json:"items"`
	Total int               `json:"total"`
}
