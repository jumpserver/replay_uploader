package model

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	IsValid  bool   `json:"is_valid"`
	IsActive bool   `json:"is_active"`
	OTPLevel int    `json:"otp_level"`
}

