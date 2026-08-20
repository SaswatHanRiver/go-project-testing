package config

import (
	"fmt"
	"log/slog"
	"os"

	"go-project-testing/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	// AutoMigrate both models - like ddl-auto=update for both entities
	err = db.AutoMigrate(&models.Product{}, &models.User{})
	if err != nil {
		slog.Error("Failed to migrate database", "error", err)
		os.Exit(1)
	}

	DB = db
	slog.Info("Database connected", "host", os.Getenv("DB_HOST"), "db", os.Getenv("DB_NAME"))
}
