package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;" json:"id"`
	Name      string         `gorm:"type:varchar(255)" json:"name"`
	Role      string         `gorm:"type:varchar(255)" json:"role"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex" json:"email"`
	Password  string         `gorm:"type:varchar(255)" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

type UserDetail struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid;unique" json:"user_id"`
	User          User      `gorm:"foreignKey:UserID" json:"-"`
	FullName      string    `gorm:"type:varchar(255)" json:"full_name"`
	NIK           string    `gorm:"type:varchar(16);unique" json:"nik"`
	Gender        string    `gorm:"type:varchar(1)" json:"gender"`
	DateOfBirth   time.Time `gorm:"type:date" json:"date_of_birth"`
	BloodType     string    `gorm:"type:varchar(2)" json:"blood_type"`
	Rhesus        string    `gorm:"type:varchar(1)" json:"rhesus"`
	PhoneNumber   string    `gorm:"type:varchar(20)" json:"phone_number"`
	Address       string    `gorm:"type:text" json:"address"`
	IsActiveDonor bool      `gorm:"default:true" json:"is_active_donor"`
	Latitude      float64   `gorm:"type:decimal(9,6)" `
	Longitude     float64   `gorm:"type:decimal(9,6)" `
	Weight        float64   `gorm:"type:decimal(5,2)" `

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ud *UserDetail) BeforeCreate(tx *gorm.DB) (err error) {
	ud.ID = uuid.New()
	return
}
