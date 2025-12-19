package pagination

type PaginationParams struct {
	Limit uint `json:"limit"`
	Page  uint `json:"page"`
}

type PaginationResult[T any] struct {
	Data       []T  `json:"data"`
	Page       uint `json:"page"`
	PageCount  uint `json:"pageCount"`
	TotalPages uint `json:"totalPages"`
	TotalCount uint `json:"totalCount"`
}
