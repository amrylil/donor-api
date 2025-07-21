package dto

// PaginatedResponse adalah struktur DTO generik untuk respons pagination.
type PaginatedResponse[T any] struct {
	Data       []T   `json:"data"`
	TotalItems int64 `json:"total_items"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
}
