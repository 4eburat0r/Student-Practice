package request

type VacancyFeedQuery struct {
	Limit  int    `form:"limit,default=20" binding:"gte=1,lte=100"`
	Offset int    `form:"offset,default=0" binding:"gte=0"`
	Format string `form:"format" binding:"omitempty,oneof=onsite remote hybrid"`
}
