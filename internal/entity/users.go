package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key;"`
	TenantID   *uuid.UUID `gorm:"type:uuid;index"`
	LocationID *uuid.UUID `gorm:"type:uuid;index"`

	Email         *string `gorm:"type:varchar(255);unique"`
	Password      *string `gorm:"type:varchar(255)"`
	Name          string  `gorm:"type:varchar(255);not null"`
	Role          string  `gorm:"type:varchar(255)" json:"role"`
	AccountStatus string  `gorm:"type:varchar(50);default:'unclaimed'"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

type UserDetail struct {
	ID     uuid.UUID `gorm:"type:uuid;primary_key;"`
	UserID uuid.UUID `gorm:"type:uuid;unique;not null"`
	User   User      `gorm:"foreignKey:UserID"`

	FullName    string    `gorm:"type:varchar(255);not null"`
	Gender      string    `gorm:"type:varchar(10)"`
	DateOfBirth time.Time `gorm:"type:date"`

	BloodType   *string `gorm:"type:varchar(2)"`
	Rhesus      *string `gorm:"type:varchar(8)"`
	PhoneNumber string  `gorm:"type:varchar(20)" json:"phone_number"`

	Address       string `gorm:"type:text"`
	IsActiveDonor bool   `gorm:"default:true"`

	Latitude  *float64 `gorm:"type:decimal(10,8)"`
	Longitude *float64 `gorm:"type:decimal(11,8)"`

	Weight    float64 `gorm:"type:decimal(5,2)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (ud *UserDetail) BeforeCreate(tx *gorm.DB) (err error) {
	ud.ID = uuid.New()
	return
}
