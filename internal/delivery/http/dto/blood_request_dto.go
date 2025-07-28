package dto

import "time"

type BloodRequestRequest struct {
	LocationID  string `json:"location_id" binding:"required"`
	BloodType   string `json:"blood_type" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type BloodRequestResponse struct {
	ID          string    `json:"id"`
	LocationID  string    `json:"location_id"`
	BloodType   string    `json:"blood_type"`
	Quantity    int       `json:"quantity"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateBloodRequestStatusDTO struct {
	Status string `json:"status" binding:"required"`
}
