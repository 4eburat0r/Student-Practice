package request

type UpdateResumeRequest struct {
	Title   string `json:"title" binding:"omitempty,min=3,max=255"`
	Summary string `json:"summary" binding:"omitempty,min=10"`
}
