package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Location struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	TenantID     uuid.UUID `gorm:"type:uuid;not null" json:"tenant_id"`
	LocationName string    `gorm:"type:varchar(255);not null" json:"location_name"`
	Slug         string    `gorm:"uniqueIndex;not null" `
	Address      string    `gorm:"type:text;not null" json:"address"`
	City         string    `gorm:"type:varchar(100);not null" json:"city"`
	Latitude     *float64  `gorm:"type:decimal(10,8)" json:"latitude"`
	Longitude    *float64  `gorm:"type:decimal(11,8)" json:"longitude"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Tenant Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Users  []User `gorm:"foreignKey:LocationID" json:"users,omitempty"`
}

func (p *Location) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
