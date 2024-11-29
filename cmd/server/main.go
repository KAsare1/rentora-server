package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"rentora-go/internal/handler"
	"rentora-go/internal/model"
	"rentora-go/internal/repository"
	"rentora-go/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	jwtSecretStr := (os.Getenv("JWT_SECRET"))

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" || jwtSecretStr == "" {
		log.Fatal("Required environment variables are missing")
	}

	jwtSecret := []byte(jwtSecretStr)

	// Initialize database
	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Auto-migrate
	db.AutoMigrate(&model.User{}, 
		&model.Car{},
	&model.Booking{})

	// Initialize dependencies
	userRepo := repository.NewUserRepository(db)
	carRepo := repository.NewCarRepository(db)
	bookingRepo := repository.NewBookingRepository(db)

	authService := service.NewAuthService(userRepo, []byte(jwtSecret))
	carService := service.NewCarService(carRepo)
	bookingService := service.NewBookingService(bookingRepo)

	authHandler := handler.NewAuthHandler(authService)
	carHandler := handler.NewCarHandler(carService)
	bookingHandler := handler.NewBookingHandler(bookingService)

	// Set up routes
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	


    handler.RegisterUserRoutes(r, authHandler)
    handler.RegisterCarRoutes(r, carHandler)
	handler.RegisterBookingRoutes(r, bookingHandler)


	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		sqlDB, err := db.DB() 
		if err != nil {
			http.Error(w, "Failed to retrieve database instance", http.StatusInternalServerError)
			return
		}
	
		// Ping the database to check its health
		if err := sqlDB.Ping(); err != nil {
			http.Error(w, "Database unreachable", http.StatusServiceUnavailable)
			return
		}
	
		w.Write([]byte("OK"))
	})

	// Start server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Println("Server is running on port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped")
}
