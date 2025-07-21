package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateDonationRequest struct {
	LocationID   uuid.UUID  `json:"location_id" binding:"required"`
	EventID      *uuid.UUID `json:"event_id"` // Opsional, bisa tidak diisi
	DonationDate time.Time  `json:"donation_date" binding:"required"`
	Status       string     `json:"status" binding:"required,oneof=selesai batal pending"`
}

type UpdateDonationRequest struct {
	Status string `json:"status" binding:"required,oneof=selesai batal pending"`
}

type DonationResponse struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	LocationID   string    `json:"location_id"`
	EventID      *string   `json:"event_id,omitempty"` // omitempty agar hilang jika null
	DonationDate time.Time `json:"donation_date"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}
