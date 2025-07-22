package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Event struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	EventName   string    `gorm:"type:varchar(255);not null" json:"event_name"`
	Description string    `gorm:"type:text" json:"description"`
	StartDate   time.Time `gorm:"type:date" json:"start_date"`
	EndDate     time.Time `gorm:"type:date" json:"end_date"`
	LocationID  uuid.UUID `gorm:"type:uuid;not null" json:"location_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *Event) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
