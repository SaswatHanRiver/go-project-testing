package repositories

import (
	"go-project-testing/config"
	"go-project-testing/models"
)

// CategoryRepository - like CategoryRepository extends JpaRepository<Category, Long>
type CategoryRepository struct{}

// FindAll - returns all categories without products
func (r *CategoryRepository) FindAll() ([]models.Category, error) {
	var categories []models.Category
	result := config.DB.Find(&categories)
	return categories, result.Error
}

// FindByID - returns a category WITH its products preloaded
// Like @EntityGraph or JOIN FETCH in Spring Boot JPA
// Preload("Products") = automatically loads the related products
func (r *CategoryRepository) FindByID(id uint) (models.Category, error) {
	var category models.Category
	result := config.DB.Preload("Products").First(&category, id)
	return category, result.Error
}

// FindByIDWithProducts - explicitly load products for a category
func (r *CategoryRepository) FindByIDWithProducts(id uint) (models.Category, error) {
	var category models.Category
	result := config.DB.Preload("Products").First(&category, id)
	return category, result.Error
}

func (r *CategoryRepository) Create(category *models.Category) error {
	result := config.DB.Create(category)
	return result.Error
}

func (r *CategoryRepository) Update(category *models.Category) error {
	result := config.DB.Save(category)
	return result.Error
}

func (r *CategoryRepository) Delete(id uint) error {
	result := config.DB.Delete(&models.Category{}, id)
	return result.Error
}

func (r *CategoryRepository) ExistsByName(name string) bool {
	var count int64
	config.DB.Model(&models.Category{}).Where("name = ?", name).Count(&count)
	return count > 0
}
