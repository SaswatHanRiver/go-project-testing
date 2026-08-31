package models

import (
	"time"

	"gorm.io/gorm"
)

// User - equivalent to @Entity + UserDetails in Spring Boot
// json:"-" on Password means it is NEVER included in JSON responses
// like @JsonIgnore in Jackson
type User struct {
	ID        uint           `json:"id"         gorm:"primaryKey;autoIncrement"`
	Name      string         `json:"name"       gorm:"not null"`
	Email     string         `json:"email"      gorm:"uniqueIndex;not null"`
	Password  string         `json:"-"          gorm:"not null"` // hidden from all responses
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"          gorm:"index"`
}

// RegisterRequest - DTO for registration (like a Spring Boot @RequestBody DTO)
type RegisterRequest struct {
	Name     string `json:"name"     binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest - DTO for login
type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse - what we return after login/register (token + user info)
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
