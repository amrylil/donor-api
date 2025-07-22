package dto

import (
	"time"

	"github.com/google/uuid"
)

type StockRequest struct {
	BloodType   string    `json:"blood_type" binding:"required,oneof=A B AB O"`
	Rhesus      string    `json:"rhesus" binding:"required,oneof=+ -"`
	BagQuantity int       `json:"bag_quantity" binding:"required,gte=0"`
	LocationID  uuid.UUID `json:"location_id" binding:"required"`
}

type UpdateQuantityRequest struct {
	BagQuantity int `json:"bag_quantity" binding:"required,gte=0"`
}

type StockResponse struct {
	ID          string    `json:"id"`
	BloodType   string    `json:"blood_type"`
	Rhesus      string    `json:"rhesus"`
	BagQuantity int       `json:"bag_quantity"`
	LocationID  string    `json:"location_id"`
	UpdatedAt   time.Time `json:"updated_at"`
}
