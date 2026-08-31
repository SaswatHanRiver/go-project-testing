package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"go-project-testing/models"
	"go-project-testing/services"
	"go-project-testing/utils"
	"go-project-testing/worker"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	service *services.ProductService
	worker  *worker.JobWorker
}

func NewProductController(w *worker.JobWorker) *ProductController {
	return &ProductController{
		service: services.NewProductService(),
		worker:  w,
	}
}

func requestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

// GetAllProducts godoc
// @Summary      Get all products
// @Description  Returns a list of all products
// @Tags         products
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.ApiResponse
// @Failure      500  {object}  utils.ApiResponse
// @Router       /api/products [get]
func (c *ProductController) GetAllProducts(ctx *gin.Context) {
	_, cancel := requestContext()
	defer cancel()

	products, err := c.service.GetAllProducts()
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to fetch products")
		return
	}
	utils.Success(ctx, http.StatusOK, "Products fetched successfully", products)
}

// GetProductByID godoc
// @Summary      Get product by ID
// @Description  Returns a single product by its ID
// @Tags         products
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  utils.ApiResponse
// @Failure      404  {object}  utils.ApiResponse
// @Router       /api/products/{id} [get]
func (c *ProductController) GetProductByID(ctx *gin.Context) {
	_, cancel := requestContext()
	defer cancel()

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "Invalid ID format")
		return
	}
	product, err := c.service.GetProductByID(uint(id))
	if err != nil {
		utils.Error(ctx, http.StatusNotFound, "Product not found")
		return
	}
	utils.Success(ctx, http.StatusOK, "Product fetched successfully", product)
}

// CreateProduct godoc
// @Summary      Create a new product
// @Description  Creates a new product with the given data
// @Tags         products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        product  body      models.CreateProductRequest  true  "Product data"
// @Success      201      {object}  utils.ApiResponse
// @Failure      400      {object}  utils.ApiResponse
// @Router       /api/products [post]
func (c *ProductController) CreateProduct(ctx *gin.Context) {
	_, cancel := requestContext()
	defer cancel()

	var req models.CreateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	product, err := c.service.CreateProduct(&req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to create product")
		return
	}

	c.worker.Submit(worker.Job{
		Type:      worker.JobProductCreated,
		ProductID: product.ID,
		Details:   "Product created: " + product.Name,
	})

	utils.Success(ctx, http.StatusCreated, "Product created successfully", product)
}

// UpdateProduct godoc
// @Summary      Update a product
// @Description  Updates an existing product by ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                          true  "Product ID"
// @Param        product  body      models.CreateProductRequest  true  "Updated product data"
// @Success      200      {object}  utils.ApiResponse
// @Failure      404      {object}  utils.ApiResponse
// @Router       /api/products/{id} [put]
func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	_, cancel := requestContext()
	defer cancel()

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "Invalid ID format")
		return
	}
	var req models.CreateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	product, err := c.service.UpdateProduct(uint(id), &req)
	if err != nil {
		utils.Error(ctx, http.StatusNotFound, "Product not found")
		return
	}

	c.worker.Submit(worker.Job{
		Type:      worker.JobProductUpdated,
		ProductID: product.ID,
		Details:   "Product updated: " + product.Name,
	})

	utils.Success(ctx, http.StatusOK, "Product updated successfully", product)
}

// DeleteProduct godoc
// @Summary      Delete a product
// @Description  Soft-deletes a product by ID
// @Tags         products
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  utils.ApiResponse
// @Failure      404  {object}  utils.ApiResponse
// @Router       /api/products/{id} [delete]
func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	_, cancel := requestContext()
	defer cancel()

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "Invalid ID format")
		return
	}
	if err := c.service.DeleteProduct(uint(id)); err != nil {
		utils.Error(ctx, http.StatusNotFound, "Product not found")
		return
	}

	c.worker.Submit(worker.Job{
		Type:      worker.JobProductDeleted,
		ProductID: uint(id),
		Details:   "Product deleted",
	})

	utils.Success(ctx, http.StatusOK, "Product deleted successfully", nil)
}
