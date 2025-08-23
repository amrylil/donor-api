package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Donation struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;" `
	Name         string     `gorm:"type:varchar(100)" `
	UserID       *uuid.UUID `gorm:"type:uuid;not null" `
	LocationID   uuid.UUID  `gorm:"type:uuid;index" `
	EventID      *uuid.UUID `gorm:"type:uuid" `
	DonationDate time.Time  `gorm:"type:date" `
	Status       string     `gorm:"type:varchar(50);default:'pending'" `
	CreatedAt    time.Time  ``
	UpdatedAt    time.Time  ``

	// Definisi relasi (opsional, untuk preloading)
	// User         User       `gorm:"foreignKey:UserID"`
	// Location     Location   `gorm:"foreignKey:LocationID"`
	// Event        Event      `gorm:"foreignKey:EventID"`
}

func (p *Donation) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
