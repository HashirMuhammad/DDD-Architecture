package main

import (
	"ddd-user-service/internal/application/service"
	"ddd-user-service/internal/infrastructure/config"
	"ddd-user-service/internal/infrastructure/repository"
	"ddd-user-service/internal/interfaces/http/handler"
	"ddd-user-service/internal/interfaces/http/router"
	"log"
	"os"
)

func main() {
	log.Println("Attempting to connect to MongoDB...")

	var userHandler *handler.UserHandler

	mongoConfig := config.NewMongoConfig()
	db, err := mongoConfig.Connect()
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v", err)
		log.Println("Falling back to in-memory repository")
		userRepo := repository.NewMemoryUserRepository()
		userService := service.NewUserService(userRepo)
		userHandler = handler.NewUserHandler(userService)
	} else {
		log.Println("✅ Connected to MongoDB successfully - using persistent storage!")
		userRepo := repository.NewMongoUserRepository(db)
		userService := service.NewUserService(userRepo)
		userHandler = handler.NewUserHandler(userService)
	}

	r := router.SetupRouter(userHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
