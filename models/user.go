package models

import (
	"time"

	"gorm.io/gorm"
)

// User entity
type User struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Password  string         `json:"-"` // hidden from all responses
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// RegisterRequest DTO
type RegisterRequest struct {
	Name     string `json:"name"     binding:"required"       example:"Saswat Kumar Sahoo"`
	Email    string `json:"email"    binding:"required,email" example:"saswat.sahoo@hanriver.in"`
	Password string `json:"password" binding:"required,min=6" example:"Saswat@123"`
}

// LoginRequest DTO
type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email" example:"saswat.sahoo@hanriver.in"`
	Password string `json:"password" binding:"required"       example:"Saswat@123"`
}

// AuthResponse - returned after login/register
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
