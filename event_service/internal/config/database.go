package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/mountainman199231/event_service/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DatabaseConnect(logger *slog.Logger) (*gorm.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbName == "" {
		return nil, fmt.Errorf("database environment variables not set")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	logger.Info("connected to database",
		"host", dbHost,
		"port", dbPort,
		"dbname", dbName,
	)

	if err := db.AutoMigrate(
		&models.MatchEvent{},
		&models.Commentary{},
		&models.Comment{},
		&models.Reaction{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
