package config

import (
	"fmt"
	"log"

	"realworld-api/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func InitDB() error {
	cfg := GetConfig()

	dsn := cfg.Database.GetDSN()

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established")

	// Auto migrate models
	if err := db.AutoMigrate(
		&models.User{},
		&models.Follow{},
		&models.Tag{},
		&models.Article{},
		&models.Comment{},
		&models.Favorite{},
	); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed")

	return nil
}

func GetDB() *gorm.DB {
	return db
}
