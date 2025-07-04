package database

import (
	"log"
	"movie-api/internal/model"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() *gorm.DB {
	dsn := os.Getenv("POSTGRES")
	if dsn == "" {
		log.Fatal("POSTGRES environment variable is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Println("Connected to the database successfully")
	log.Println("Migrating database schema...")

	err = db.AutoMigrate(&model.Movie{}, &model.User{})
	if err != nil {
		log.Fatalf("failed to migrate database schema: %v", err)
	}
	log.Println("Database schema migrated successfully")

	return db
}
