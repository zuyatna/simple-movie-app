package main

import (
	"log"
	"movie-api/pkg/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default environment variables")
	}

	gormDB := database.InitDB()
	if gormDB == nil {
		log.Fatal("Failed to initialize database connection")
	}
	defer func() {
		if sqlDB, err := gormDB.DB(); err != nil {
			log.Printf("Failed to get database instance: %v", err)
		} else if err := sqlDB.Close(); err != nil {
			log.Printf("Failed to close database connection: %v", err)
		} else {
			log.Println("Database connection closed successfully")
		}
	}()
	log.Println("Database initialized successfully")

	redisClient := database.InitRedis()
	if redisClient == nil {
		log.Fatal("Failed to initialize Redis connection")
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("Failed to close Redis connection: %v", err)
		} else {
			log.Println("Redis connection closed successfully")
		}
	}()
	log.Println("Redis initialized successfully")

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the Movie API",
		})
	})

	log.Println("Starting server on :8081")
	router.Run(":8081")
}
