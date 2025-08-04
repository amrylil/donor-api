package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Donation struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;" json:"id"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	LocationID   uuid.UUID  `gorm:"type:uuid;index" json:"location_id"`
	EventID      *uuid.UUID `gorm:"type:uuid" json:"event_id"`
	DonationDate time.Time  `gorm:"type:date" json:"donation_date"`
	Status       string     `gorm:"type:varchar(50);default:'pending'" json:"status"` // pending, selesai, batal
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	// Definisi relasi (opsional, untuk preloading)
	// User         User       `gorm:"foreignKey:UserID"`
	// Location     Location   `gorm:"foreignKey:LocationID"`
	// Event        Event      `gorm:"foreignKey:EventID"`
}

func (p *Donation) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
