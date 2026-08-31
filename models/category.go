package models

import (
	"time"

	"gorm.io/gorm"
)

// Category - equivalent to @Entity Category in Spring Boot
// Relationship: One Category has many Products
// Like @OneToMany(mappedBy = "category") in Spring Boot JPA
type Category struct {
	ID          uint           `json:"id"          gorm:"primaryKey;autoIncrement"`
	Name        string         `json:"name"        gorm:"uniqueIndex;not null"`
	Description string         `json:"description"`
	Products    []Product      `json:"products,omitempty" gorm:"foreignKey:CategoryID"` // omitempty = hide if empty
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-"           gorm:"index"`
}

// CategoryRequest - DTO for create and update
type CategoryRequest struct {
	Name        string `json:"name"        binding:"required"`
	Description string `json:"description"`
}
