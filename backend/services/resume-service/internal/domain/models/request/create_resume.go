package request

type CreateResumeRequest struct {
	Title   string `json:"title" binding:"required,min=3,max=255"`
	Summary string `json:"summary" binding:"required,min=10"`
}
