package request

type PaginationRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	SortBy   string `form:"sort_by" binding:"omitempty"`
	Order    string `form:"order" binding:"omitempty,oneof=asc desc"`
}

func (r *PaginationRequest) Normalize() {
	if r.Page < 1 {
		r.Page = 1
	}
	if r.PageSize < 1 || r.PageSize > 100 {
		r.PageSize = 20
	}
	if r.Order != "asc" && r.Order != "desc" {
		r.Order = "desc"
	}
	if r.SortBy == "" {
		r.SortBy = "created_at"
	}
}

func (r *PaginationRequest) GetOffset() int {
	return (r.Page - 1) * r.PageSize
}
