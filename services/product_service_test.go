package services

// In Spring Boot: @ExtendWith(MockitoExtension.class) + @MockBean
// In Go: no framework needed - we create a mock struct that implements the interface

import (
	"errors"
	"testing"

	"go-project-testing/models"
)

// ── Mock Repository ──────────────────────────────────────────────────────────
// MockProductRepository - equivalent to Mockito.mock(ProductRepository.class)
// It implements ProductRepositoryInterface with controlled, fake responses
type MockProductRepository struct {
	products []models.Product
	err      error // set this to simulate DB errors
}

func (m *MockProductRepository) FindAll() ([]models.Product, error) {
	return m.products, m.err
}

func (m *MockProductRepository) FindByID(id uint) (models.Product, error) {
	for _, p := range m.products {
		if p.ID == id {
			return p, m.err
		}
	}
	return models.Product{}, errors.New("record not found")
}

func (m *MockProductRepository) Create(product *models.Product) error {
	product.ID = uint(len(m.products) + 1) // simulate auto-increment
	m.products = append(m.products, *product)
	return m.err
}

func (m *MockProductRepository) Update(product *models.Product) error {
	for i, p := range m.products {
		if p.ID == product.ID {
			m.products[i] = *product
			return m.err
		}
	}
	return errors.New("record not found")
}

func (m *MockProductRepository) Delete(id uint) error {
	return m.err
}

// ── Tests ────────────────────────────────────────────────────────────────────

// TestGetAllProducts - like @Test void getAllProducts_shouldReturnList()
func TestGetAllProducts(t *testing.T) {
	mock := &MockProductRepository{
		products: []models.Product{
			{ID: 1, Name: "Laptop", Price: 999.99, Stock: 10},
			{ID: 2, Name: "Phone", Price: 499.99, Stock: 20},
		},
	}
	svc := NewProductServiceWithRepo(mock)

	result, err := svc.GetAllProducts()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 products, got %d", len(result))
	}
}

// TestGetAllProducts_DBError - simulate DB failure
func TestGetAllProducts_DBError(t *testing.T) {
	mock := &MockProductRepository{err: errors.New("DB connection lost")}
	svc := NewProductServiceWithRepo(mock)

	_, err := svc.GetAllProducts()

	if err == nil {
		t.Fatal("expected an error but got nil")
	}
}

// TestCreateProduct - verify data is mapped correctly from request to model
func TestCreateProduct(t *testing.T) {
	mock := &MockProductRepository{}
	svc := NewProductServiceWithRepo(mock)

	req := &models.CreateProductRequest{
		Name:        "Keyboard",
		Description: "Mechanical",
		Price:       79.99,
		Stock:       50,
	}

	product, err := svc.CreateProduct(req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if product.Name != "Keyboard" {
		t.Errorf("expected name 'Keyboard', got '%s'", product.Name)
	}
	if product.Price != 79.99 {
		t.Errorf("expected price 79.99, got %f", product.Price)
	}
	if product.Stock != 50 {
		t.Errorf("expected stock 50, got %d", product.Stock)
	}
}

// TestGetProductByID - verify correct product returned
func TestGetProductByID(t *testing.T) {
	mock := &MockProductRepository{
		products: []models.Product{
			{ID: 1, Name: "Monitor", Price: 299.99},
		},
	}
	svc := NewProductServiceWithRepo(mock)

	product, err := svc.GetProductByID(1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if product.Name != "Monitor" {
		t.Errorf("expected 'Monitor', got '%s'", product.Name)
	}
}

// TestGetProductByID_NotFound - verify error when product doesn't exist
func TestGetProductByID_NotFound(t *testing.T) {
	mock := &MockProductRepository{products: []models.Product{}}
	svc := NewProductServiceWithRepo(mock)

	_, err := svc.GetProductByID(99)

	if err == nil {
		t.Fatal("expected 'not found' error, got nil")
	}
}

// TestUpdateProduct - verify fields are updated correctly
func TestUpdateProduct(t *testing.T) {
	mock := &MockProductRepository{
		products: []models.Product{
			{ID: 1, Name: "Old Name", Price: 10.00, Stock: 5},
		},
	}
	svc := NewProductServiceWithRepo(mock)

	req := &models.CreateProductRequest{
		Name:  "New Name",
		Price: 20.00,
		Stock: 10,
	}

	updated, err := svc.UpdateProduct(1, req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updated.Name != "New Name" {
		t.Errorf("expected name 'New Name', got '%s'", updated.Name)
	}
	if updated.Price != 20.00 {
		t.Errorf("expected price 20.00, got %f", updated.Price)
	}
}

// TestDeleteProduct - verify delete passes through without error
func TestDeleteProduct(t *testing.T) {
	mock := &MockProductRepository{}
	svc := NewProductServiceWithRepo(mock)

	err := svc.DeleteProduct(1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
