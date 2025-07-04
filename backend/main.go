package main

import (
	"context"
	"log"
	"movie-api/internal/handler"
	"movie-api/internal/middleware"
	"movie-api/internal/repository"
	"movie-api/internal/usecase"
	"movie-api/pkg/database"

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

	userRepo := repository.NewUserRepository(gormDB)
	authUsecase := usecase.NewAuthUsecase(userRepo)
	authHandler := handler.NewAuthHandler(authUsecase)

	movieRepo := repository.NewMovieRepository(gormDB)
	movieUsecase := usecase.NewMovieUsecase(movieRepo, redisClient, context.Background())
	movieHandler := handler.NewMovieHandler(movieUsecase)

	router := gin.Default()

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		api.GET("/movies", movieHandler.FindAllMovie)
		api.GET("/movies/:id", movieHandler.FindByMovieID)
		api.GET("/movies/:id/poster", movieHandler.GetMoviePoster)

		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware())
		admin.Use(middleware.AdminMiddleware())
		{
			admin.POST("/movies", movieHandler.CreateMovie)
			// TODO: Implement UpdateMovie and DeleteMovie handlers
		}
	}

	log.Println("Starting server on port :8081")
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
