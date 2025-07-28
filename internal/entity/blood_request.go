package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BloodRequest struct {
	ID          uuid.UUID
	LocationID  uuid.UUID
	BloodType   string
	Quantity    int
	Status      string
	Description string
	CreatedBy   uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *BloodRequest) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
