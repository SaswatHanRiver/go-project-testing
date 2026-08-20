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

	// ── Public routes (no token needed) ──────────────────────────────────────
	// Like permitAll() in Spring Security's SecurityConfig
	auth := router.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	// ── Protected routes (JWT required) ──────────────────────────────────────
	// Like .anyRequest().authenticated() in Spring Security
	// Every request to /api/* goes through AuthMiddleware first
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware()) // attach middleware to this group
	{
		products := api.Group("/products")
		{
			products.GET("", productController.GetAllProducts)
			products.GET("/:id", productController.GetProductByID)
			products.POST("", productController.CreateProduct)
			products.PUT("/:id", productController.UpdateProduct)
			products.DELETE("/:id", productController.DeleteProduct)
		}

		// /auth/me is protected - requires valid token
		api.GET("/auth/me", authController.Me)
	}

	// Swagger UI - public
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
