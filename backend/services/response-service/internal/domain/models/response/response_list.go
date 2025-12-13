package response

type ResponseListResponse struct {
	Items []ResponseResponse `json:"items"`
	Total int                `json:"total"`
}
