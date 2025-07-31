package dto

type APIResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

type SuccessWrapper struct {
	APIResponse[any]
}

type ErrorWrapper struct {
	APIResponse[any]
}
