package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	_ "go-project-testing/docs"

	"go-project-testing/config"
	"go-project-testing/routes"
	"go-project-testing/worker"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title           Go Project Testing API
// @version         1.0
// @description     A CRUD REST API built with Gin + GORM - Go equivalent of Spring Boot
// @host            localhost:8080
// @BasePath        /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter: Bearer {your_token}
func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using environment variables")
	}

	config.ConnectDatabase()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobWorker := worker.NewJobWorker(100)
	jobWorker.Start(ctx)

	router := gin.Default()
	routes.SetupRoutes(router, jobWorker)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("Server started",
		"port", port,
		"swagger", "http://localhost:"+port+"/swagger/index.html",
		"api", "http://localhost:"+port+"/api/products",
	)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		slog.Info("Shutdown signal received")
		cancel()
	}()

	if err := router.Run(":" + port); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
