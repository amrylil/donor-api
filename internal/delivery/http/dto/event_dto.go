package dto

import (
	"time"

	"github.com/google/uuid"
)

type EventRequest struct {
	EventName   string    `json:"event_name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
	LocationID  uuid.UUID `json:"location_id" binding:"required"`
}

type EventResponse struct {
	ID          string    `json:"id"`
	EventName   string    `json:"event_name"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	LocationID  string    `json:"location_id"`
	CreatedAt   time.Time `json:"created_at"`
}
