package main

import (
	"log"
	"net/http"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"rentora-go/internal/handler"
	"rentora-go/internal/model"
	"rentora-go/internal/repository"
	"rentora-go/internal/service"
)

func main() {
	// Initialize the database
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Auto-migrate the User model
	db.AutoMigrate(&model.User{})

	// Dependency injection
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, []byte(os.Getenv("JWT_SECRET")))
	authHandler := handler.NewAuthHandler(authService)

	// Set up routes
	http.HandleFunc("/login", authHandler.Login)
	http.HandleFunc("/refresh", authHandler.Refresh)

	// Start the server
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
