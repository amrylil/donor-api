package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BloodRequest struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	LocationID  uuid.UUID `gorm:"type:uuid;index" json:"location_id"`
	BloodType   string    `gorm:"type:varchar(2);not null" json:"blood_type"`       // A, B, AB, O
	Quantity    int       `gorm:"not null" json:"quantity"`                         // in bags
	Status      string    `gorm:"type:varchar(50);default:'pending'" json:"status"` // pending, fulfilled, cancelled
	Description string    `gorm:"type:text" json:"description"`
	CreatedBy   uuid.UUID `gorm:"type:uuid;not null" json:"created_by"` // UserID of the requester
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *BloodRequest) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
