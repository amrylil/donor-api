package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Stock struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	BloodType   string    `gorm:"type:varchar(2);not null;uniqueIndex:idx_stock_location" json:"blood_type"`
	Rhesus      string    `gorm:"type:varchar(1);not null;uniqueIndex:idx_stock_location" json:"rhesus"`
	BagQuantity int       `gorm:"not null" json:"bag_quantity"`
	LocationID  uuid.UUID `gorm:"type:uuid;index" json:"location_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *Stock) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
