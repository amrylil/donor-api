package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Location struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;" json:"id"`
	Title     string         `gorm:"type:varchar(255)" json:"title"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *Location) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
