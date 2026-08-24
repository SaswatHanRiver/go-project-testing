package controllers

import (
	"net/http"
	"strconv"

	"go-project-testing/models"
	"go-project-testing/services"

	"github.com/gin-gonic/gin"
)

// CategoryController - equivalent to @RestController in Spring Boot
type CategoryController struct {
	service *services.CategoryService
}

func NewCategoryController() *CategoryController {
	return &CategoryController{service: services.NewCategoryService()}
}

// GetAllCategories godoc
// @Summary      Get all categories
// @Description  Returns a list of all categories
// @Tags         categories
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.Category
// @Failure      500  {object}  map[string]string
// @Router       /api/categories [get]
func (c *CategoryController) GetAllCategories(ctx *gin.Context) {
	categories, err := c.service.GetAllCategories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, categories)
}

// GetCategoryByID godoc
// @Summary      Get category by ID
// @Description  Returns a single category with its products
// @Tags         categories
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  models.Category
// @Failure      404  {object}  map[string]string
// @Router       /api/categories/{id} [get]
func (c *CategoryController) GetCategoryByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	category, err := c.service.GetCategoryByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	ctx.JSON(http.StatusOK, category)
}

// GetCategoryProducts godoc
// @Summary      Get products by category
// @Description  Returns all products belonging to a specific category
// @Tags         categories
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Category ID"
// @Success      200  {array}   models.Product
// @Failure      404  {object}  map[string]string
// @Router       /api/categories/{id}/products [get]
func (c *CategoryController) GetCategoryProducts(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	category, err := c.service.GetCategoryWithProducts(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	ctx.JSON(http.StatusOK, category.Products)
}

// CreateCategory godoc
// @Summary      Create a new category
// @Description  Creates a new product category
// @Tags         categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        category  body      models.CategoryRequest  true  "Category data"
// @Success      201       {object}  models.Category
// @Failure      400       {object}  map[string]string
// @Failure      409       {object}  map[string]string
// @Router       /api/categories [post]
func (c *CategoryController) CreateCategory(ctx *gin.Context) {
	var req models.CategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	category, err := c.service.CreateCategory(&req)
	if err != nil {
		if err.Error() == "category name already exists" {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, category)
}

// UpdateCategory godoc
// @Summary      Update a category
// @Description  Updates an existing category by ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id        path      int                     true  "Category ID"
// @Param        category  body      models.CategoryRequest  true  "Updated category data"
// @Success      200       {object}  models.Category
// @Failure      400       {object}  map[string]string
// @Failure      404       {object}  map[string]string
// @Router       /api/categories/{id} [put]
func (c *CategoryController) UpdateCategory(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req models.CategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	category, err := c.service.UpdateCategory(uint(id), &req)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, category)
}

// DeleteCategory godoc
// @Summary      Delete a category
// @Description  Soft-deletes a category by ID
// @Tags         categories
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/categories/{id} [delete]
func (c *CategoryController) DeleteCategory(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	if err := c.service.DeleteCategory(uint(id)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
