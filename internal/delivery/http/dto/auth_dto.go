package dto

// RegisterRequest adalah DTO untuk data yang masuk saat registrasi.
type RegisterRequest struct {
	Name       string  `json:"name" binding:"required"`
	Email      string  `json:"email" binding:"required,email"`
	Password   string  `json:"password" binding:"required,min=8"`
	LocationID *string `json:"location_id,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
