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
func main() {
	// Structured JSON logging - like Logback in Spring Boot
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	// Load .env - like application.properties
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using environment variables")
	}

	// Connect DB
	config.ConnectDatabase()

	// context.WithCancel gives us a way to shut down all goroutines cleanly
	// Like Spring Boot's ApplicationContext shutdown hook
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create and start the background job worker
	// bufferSize=100 means up to 100 jobs can queue before Submit() blocks
	jobWorker := worker.NewJobWorker(100)
	jobWorker.Start(ctx) // launches goroutine internally

	// Gin router
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

	// Graceful shutdown - listen for Ctrl+C or kill signal
	// Like Spring Boot's graceful shutdown (server.shutdown=graceful)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		slog.Info("Shutdown signal received")
		cancel() // cancels ctx → stops the job worker goroutine
	}()

	if err := router.Run(":" + port); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
