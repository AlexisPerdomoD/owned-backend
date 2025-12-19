package dto

type PaginationParams struct {
	Limit uint `json:"limit"`
	Page  uint `json:"page"`
}
