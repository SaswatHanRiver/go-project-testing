package repositories

import (
	"go-project-testing/config"
	"go-project-testing/models"
)

// ProductRepositoryInterface - the interface used for mocking in tests
// In Spring Boot: Mockito.mock(ProductRepository.class) mocks the whole interface
// In Go: we define the interface explicitly, then create a mock struct that implements it
type ProductRepositoryInterface interface {
	FindAll() ([]models.Product, error)
	FindByID(id uint) (models.Product, error)
	Create(product *models.Product) error
	Update(product *models.Product) error
	Delete(id uint) error
}

// ProductRepository - the real implementation that hits PostgreSQL
type ProductRepository struct{}

func (r *ProductRepository) FindAll() ([]models.Product, error) {
	var products []models.Product
	result := config.DB.Find(&products)
	return products, result.Error
}

func (r *ProductRepository) FindByID(id uint) (models.Product, error) {
	var product models.Product
	result := config.DB.First(&product, id)
	return product, result.Error
}

func (r *ProductRepository) Create(product *models.Product) error {
	result := config.DB.Create(product)
	return result.Error
}

func (r *ProductRepository) Update(product *models.Product) error {
	result := config.DB.Save(product)
	return result.Error
}

func (r *ProductRepository) Delete(id uint) error {
	result := config.DB.Delete(&models.Product{}, id)
	return result.Error
}
