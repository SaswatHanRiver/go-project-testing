package models

import (
	"time"

	"gorm.io/gorm"
)

// Product - equivalent to @Entity in Spring Boot
// CategoryID is a foreign key - like @ManyToOne @JoinColumn(name="category_id") in JPA
// *uint means it is a pointer = nullable (category is optional for a product)
type Product struct {
	ID          uint           `json:"id"          gorm:"primaryKey;autoIncrement"`
	Name        string         `json:"name"        gorm:"not null"`
	Description string         `json:"description"`
	Price       float64        `json:"price"       gorm:"not null"`
	Stock       int            `json:"stock"       gorm:"default:0"`
	CategoryID  *uint          `json:"category_id"`                              // FK - nullable
	Category    *Category      `json:"category,omitempty" gorm:"foreignKey:CategoryID"` // preloaded relation
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-"           gorm:"index"`
}

// CreateProductRequest - DTO for create/update
type CreateProductRequest struct {
	Name        string  `json:"name"        binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price"       binding:"required,gt=0"`
	Stock       int     `json:"stock"`
	CategoryID  *uint   `json:"category_id"` // optional - pointer so it can be null
}
