package dto

import "time"

// DTO untuk request body (Create & Update)
type LocationRequest struct {
	Title string `json:"title" binding:"required"`
}

// DTO untuk response (data aman untuk publik)
type LocationResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
