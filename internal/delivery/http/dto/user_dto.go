package dto

import "time"

type UserRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role"`
}
type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type ProfileResponse struct {
	User    UserResponse        `json:"user"`
	Details *UserDetailResponse `json:"details,omitempty"`
}

type UserDetailRequest struct {
	FullName      string    `json:"full_name" binding:"required"`
	NIK           string    `json:"nik" binding:"required,len=16"`
	Gender        string    `json:"gender" binding:"required,oneof=L P"`
	DateOfBirth   time.Time `json:"date_of_birth" binding:"required"`
	BloodType     string    `json:"blood_type" binding:"required,oneof=A B AB O"`
	Rhesus        string    `json:"rhesus" binding:"required,oneof=+ -"`
	PhoneNumber   string    `json:"phone_number" binding:"required,min=10,max=15"`
	Address       string    `json:"address" binding:"required"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	Weight        float64   `json:"weight"`
	IsActiveDonor bool      `json:"is_active_donor" binding:"required"`
}

type UserDetailResponse struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	FullName      string    `json:"full_name"`
	NIK           string    `json:"nik"`
	Gender        string    `json:"gender"`
	DateOfBirth   time.Time `json:"date_of_birth"`
	BloodType     string    `json:"blood_type"`
	Rhesus        string    `json:"rhesus"`
	PhoneNumber   string    `json:"phone_number"`
	Address       string    `json:"address"`
	IsActiveDonor bool      `json:"is_active_donor"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	Weight        float64   `json:"weight"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
