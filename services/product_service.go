package services

import (
	"go-project-testing/models"
	"go-project-testing/repositories"
)

// ProductService - equivalent to @Service in Spring Boot
type ProductService struct {
	repo repositories.ProductRepositoryInterface // interface, not concrete struct
}

// NewProductService - used in production (real DB repository)
func NewProductService() *ProductService {
	return &ProductService{repo: &repositories.ProductRepository{}}
}

// NewProductServiceWithRepo - used in tests (mock repository injected)
// Like @MockBean in Spring Boot tests
func NewProductServiceWithRepo(repo repositories.ProductRepositoryInterface) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	return s.repo.FindAll()
}

func (s *ProductService) GetProductByID(id uint) (models.Product, error) {
	return s.repo.FindByID(id)
}

func (s *ProductService) CreateProduct(req *models.CreateProductRequest) (models.Product, error) {
	product := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}
	err := s.repo.Create(&product)
	return product, err
}

func (s *ProductService) UpdateProduct(id uint, req *models.CreateProductRequest) (models.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return product, err
	}
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock
	err = s.repo.Update(&product)
	return product, err
}

func (s *ProductService) DeleteProduct(id uint) error {
	return s.repo.Delete(id)
}
