package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"go-project-testing/models"
	"go-project-testing/services"
	"go-project-testing/worker"

	"github.com/gin-gonic/gin"
)

// ProductController - equivalent to @RestController in Spring Boot
type ProductController struct {
	service *services.ProductService
	worker  *worker.JobWorker // injected - like @Autowired in Spring Boot
}

func NewProductController(w *worker.JobWorker) *ProductController {
	return &ProductController{
		service: services.NewProductService(),
		worker:  w,
	}
}

// requestContext creates a context with a 5-second timeout for each request
// Like Spring Boot's @Transactional timeout or RestTemplate timeout
// If the DB takes more than 5s, the request is cancelled automatically
func requestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

// GetAllProducts godoc
// @Summary      Get all products
// @Description  Returns a list of all products
// @Tags         products
// @Produce      json
// @Success      200  {array}   models.Product
// @Failure      500  {object}  map[string]string
// @Router       /api/products [get]
func (c *ProductController) GetAllProducts(ctx *gin.Context) {
	_, cancel := requestContext()
	defer cancel() // always release the context when done

	products, err := c.service.GetAllProducts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, products)
}

// GetProductByID godoc
// @Summary      Get product by ID
// @Description  Returns a single product by its ID
// @Tags         products
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  models.Product
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/products/{id} [get]
func (c *ProductController) GetProductByID(ctx *gin.Context) {
	_, cancel := requestContext()
	defer cancel()

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	product, err := c.service.GetProductByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	ctx.JSON(http.StatusOK, product)
}

// CreateProduct godoc
// @Summary      Create a new product
// @Description  Creates a new product with the given data
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body      models.CreateProductRequest  true  "Product data"
// @Success      201      {object}  models.Product
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/products [post]
func (c *ProductController) CreateProduct(ctx *gin.Context) {
	_, cancel := requestContext()
	defer cancel()

	var req models.CreateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product, err := c.service.CreateProduct(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Submit background job AFTER responding - fire and forget via channel
	// Like calling an @Async method in Spring Boot after saving to DB
	c.worker.Submit(worker.Job{
		Type:      worker.JobProductCreated,
		ProductID: product.ID,
		Details:   "Product created: " + product.Name,
	})

	ctx.JSON(http.StatusCreated, product)
}

// UpdateProduct godoc
// @Summary      Update a product
// @Description  Updates an existing product by ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id       path      int                          true  "Product ID"
// @Param        product  body      models.CreateProductRequest  true  "Updated product data"
// @Success      200      {object}  models.Product
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Router       /api/products/{id} [put]
func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	_, cancel := requestContext()
	defer cancel()

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req models.CreateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product, err := c.service.UpdateProduct(uint(id), &req)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.worker.Submit(worker.Job{
		Type:      worker.JobProductUpdated,
		ProductID: product.ID,
		Details:   "Product updated: " + product.Name,
	})

	ctx.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary      Delete a product
// @Description  Soft-deletes a product by ID
// @Tags         products
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/products/{id} [delete]
func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	_, cancel := requestContext()
	defer cancel()

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	if err := c.service.DeleteProduct(uint(id)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.worker.Submit(worker.Job{
		Type:      worker.JobProductDeleted,
		ProductID: uint(id),
		Details:   "Product deleted",
	})

	ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
