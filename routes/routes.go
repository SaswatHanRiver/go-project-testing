package routes

import (
	"go-project-testing/controllers"
	"go-project-testing/worker"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes - pass the worker into the controller (manual dependency injection)
// In Spring Boot this is done automatically via @Autowired
// In Go we pass dependencies explicitly - no magic, fully traceable
func SetupRoutes(router *gin.Engine, w *worker.JobWorker) {
	productController := controllers.NewProductController(w)

	api := router.Group("/api")
	{
		products := api.Group("/products")
		{
			products.GET("", productController.GetAllProducts)
			products.GET("/:id", productController.GetProductByID)
			products.POST("", productController.CreateProduct)
			products.PUT("/:id", productController.UpdateProduct)
			products.DELETE("/:id", productController.DeleteProduct)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
