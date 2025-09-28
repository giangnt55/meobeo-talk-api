package utils

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func NewPagination(page, limit int) *Pagination {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return &Pagination{
		Page:  page,
		Limit: limit,
	}
}
