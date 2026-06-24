package dto

type PaginateQuery struct {
	Page  int64 `form:"page" binding:"omitempty,min=1"`
	Limit int64 `form:"limit" binding:"omitempty,min=1,max=100"`
}
