package services

import (
	"errors"

	"go-project-testing/models"
	"go-project-testing/repositories"
)

// CategoryService - equivalent to @Service in Spring Boot
type CategoryService struct {
	repo repositories.CategoryRepository
}

func NewCategoryService() *CategoryService {
	return &CategoryService{repo: repositories.CategoryRepository{}}
}

func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	return s.repo.FindAll()
}

func (s *CategoryService) GetCategoryByID(id uint) (models.Category, error) {
	return s.repo.FindByID(id)
}

func (s *CategoryService) GetCategoryWithProducts(id uint) (models.Category, error) {
	return s.repo.FindByIDWithProducts(id)
}

func (s *CategoryService) CreateCategory(req *models.CategoryRequest) (models.Category, error) {
	if s.repo.ExistsByName(req.Name) {
		return models.Category{}, errors.New("category name already exists")
	}

	category := models.Category{
		Name:        req.Name,
		Description: req.Description,
	}
	err := s.repo.Create(&category)
	return category, err
}

func (s *CategoryService) UpdateCategory(id uint, req *models.CategoryRequest) (models.Category, error) {
	category, err := s.repo.FindByID(id)
	if err != nil {
		return category, errors.New("category not found")
	}
	category.Name = req.Name
	category.Description = req.Description
	err = s.repo.Update(&category)
	return category, err
}

func (s *CategoryService) DeleteCategory(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("category not found")
	}
	return s.repo.Delete(id)
}
