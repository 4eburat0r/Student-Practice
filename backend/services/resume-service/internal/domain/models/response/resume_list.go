package response

type ResumeListResponse struct {
	Resumes []ResumeResponse `json:"resumes"`
	Total   int              `json:"total"`
}
