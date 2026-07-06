package dto

type PaginateQuery struct {
	Page  *int64 `form:"page" binding:"omitempty,min=1"`
	Limit *int64 `form:"limit" binding:"omitempty,min=1,max=100"`
}

func (p *PaginateQuery) SetPage(page int64) {
	p.Page = &page
}
func (p *PaginateQuery) GetPage() int64 {
	if p.Page == nil {
		return 0
	}
	return *p.Page
}
func (p *PaginateQuery) GetLimit() int64 {
	if p.Limit == nil {
		return 0
	}
	return *p.Limit
}
