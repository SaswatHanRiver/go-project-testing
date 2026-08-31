package controllers

import (
	"net/http"
	"strconv"

	"go-project-testing/models"
	"go-project-testing/services"
	"go-project-testing/utils"

	"github.com/gin-gonic/gin"
)

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
// @Success      200  {object}  utils.ApiResponse
// @Failure      500  {object}  utils.ApiResponse
// @Router       /api/categories [get]
func (c *CategoryController) GetAllCategories(ctx *gin.Context) {
	categories, err := c.service.GetAllCategories()
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to fetch categories")
		return
	}
	utils.Success(ctx, http.StatusOK, "Categories fetched successfully", categories)
}

// GetCategoryByID godoc
// @Summary      Get category by ID
// @Description  Returns a single category with its products
// @Tags         categories
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  utils.ApiResponse
// @Failure      404  {object}  utils.ApiResponse
// @Router       /api/categories/{id} [get]
func (c *CategoryController) GetCategoryByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "Invalid ID format")
		return
	}
	category, err := c.service.GetCategoryByID(uint(id))
	if err != nil {
		utils.Error(ctx, http.StatusNotFound, "Category not found")
		return
	}
	utils.Success(ctx, http.StatusOK, "Category fetched successfully", category)
}

// GetCategoryProducts godoc
// @Summary      Get products by category
// @Description  Returns all products belonging to a specific category
// @Tags         categories
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  utils.ApiResponse
// @Failure      404  {object}  utils.ApiResponse
// @Router       /api/categories/{id}/products [get]
func (c *CategoryController) GetCategoryProducts(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "Invalid ID format")
		return
	}
	category, err := c.service.GetCategoryWithProducts(uint(id))
	if err != nil {
		utils.Error(ctx, http.StatusNotFound, "Category not found")
		return
	}
	utils.Success(ctx, http.StatusOK, "Products fetched successfully", category.Products)
}

// CreateCategory godoc
// @Summary      Create a new category
// @Description  Creates a new product category
// @Tags         categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        category  body      models.CategoryRequest  true  "Category data"
// @Success      201       {object}  utils.ApiResponse
// @Failure      400       {object}  utils.ApiResponse
// @Router       /api/categories [post]
func (c *CategoryController) CreateCategory(ctx *gin.Context) {
	var req models.CategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	category, err := c.service.CreateCategory(&req)
	if err != nil {
		if err.Error() == "category name already exists" {
			utils.Error(ctx, http.StatusConflict, err.Error())
			return
		}
		utils.Error(ctx, http.StatusInternalServerError, "Failed to create category")
		return
	}
	utils.Success(ctx, http.StatusCreated, "Category created successfully", category)
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
// @Success      200       {object}  utils.ApiResponse
// @Failure      404       {object}  utils.ApiResponse
// @Router       /api/categories/{id} [put]
func (c *CategoryController) UpdateCategory(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "Invalid ID format")
		return
	}
	var req models.CategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	category, err := c.service.UpdateCategory(uint(id), &req)
	if err != nil {
		utils.Error(ctx, http.StatusNotFound, err.Error())
		return
	}
	utils.Success(ctx, http.StatusOK, "Category updated successfully", category)
}

// DeleteCategory godoc
// @Summary      Delete a category
// @Description  Soft-deletes a category by ID
// @Tags         categories
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  utils.ApiResponse
// @Failure      404  {object}  utils.ApiResponse
// @Router       /api/categories/{id} [delete]
func (c *CategoryController) DeleteCategory(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "Invalid ID format")
		return
	}
	if err := c.service.DeleteCategory(uint(id)); err != nil {
		utils.Error(ctx, http.StatusNotFound, err.Error())
		return
	}
	utils.Success(ctx, http.StatusOK, "Category deleted successfully", nil)
}
