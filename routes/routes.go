package routes

import (
	"go-project-testing/controllers"
	"go-project-testing/middleware"
	"go-project-testing/worker"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine, w *worker.JobWorker) {
	authController := controllers.NewAuthController()
	productController := controllers.NewProductController(w)
	categoryController := controllers.NewCategoryController()

	// ── Public routes ─────────────────────────────────────────────────────────
	auth := router.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	// ── Protected routes (JWT required) ───────────────────────────────────────
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// Auth
		api.GET("/auth/me", authController.Me)

		// Products
		products := api.Group("/products")
		{
			products.GET("", productController.GetAllProducts)
			products.GET("/:id", productController.GetProductByID)
			products.POST("", productController.CreateProduct)
			products.PUT("/:id", productController.UpdateProduct)
			products.DELETE("/:id", productController.DeleteProduct)
		}

		// Categories
		categories := api.Group("/categories")
		{
			categories.GET("", categoryController.GetAllCategories)
			categories.GET("/:id", categoryController.GetCategoryByID)
			categories.GET("/:id/products", categoryController.GetCategoryProducts) // nested route
			categories.POST("", categoryController.CreateCategory)
			categories.PUT("/:id", categoryController.UpdateCategory)
			categories.DELETE("/:id", categoryController.DeleteCategory)
		}
	}

	// Swagger UI - public
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
