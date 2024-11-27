package main

import (
	"log"
	"net/http"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/joho/godotenv"

	
	"rentora-go/internal/handler"
	"rentora-go/internal/model"
	"rentora-go/internal/repository"
	"rentora-go/internal/service"
)

func main() {
	// Get database credentials from environment variables
	err := godotenv.Load()
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Construct DSN (Data Source Name)
	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"

	// Initialize the database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
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
