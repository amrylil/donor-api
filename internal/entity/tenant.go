package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tenant struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" `
	Name      string    `gorm:"not null" `
	Slug      string    `gorm:"uniqueIndex;not null" `
	CreatedAt time.Time
	UpdatedAt time.Time

	Locations []Location `gorm:"foreignKey:TenantID"`
	Users     []User     `gorm:"foreignKey:TenantID"`
}

func (p *Tenant) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
