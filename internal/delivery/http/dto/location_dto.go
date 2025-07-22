package dto

import "time"

type LocationRequest struct {
	LocationName string   `json:"location_name" binding:"required"`
	Address      string   `json:"address" binding:"required"`
	City         string   `json:"city" binding:"required"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
}

type LocationResponse struct {
	ID           string    `json:"id"`
	LocationName string    `json:"location_name"`
	Address      string    `json:"address"`
	City         string    `json:"city"`
	Latitude     *float64  `json:"latitude,omitempty"`
	Longitude    *float64  `json:"longitude,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
