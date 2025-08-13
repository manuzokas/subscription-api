package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/manuzokas/subscription-api/internal/adapters/database"
	"github.com/manuzokas/subscription-api/internal/adapters/web"
	"github.com/manuzokas/subscription-api/internal/core/auth"
	"github.com/manuzokas/subscription-api/internal/core/subscription"
)

func main() {
	log.Println("Starting Subscription API...")

	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables from OS")
	}

	databaseUrl := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	apiPort := os.Getenv("API_PORT")

	if databaseUrl == "" || jwtSecret == "" || apiPort == "" {
		log.Fatal("Error: DATABASE_URL, JWT_SECRET, and API_PORT must be set")
	}

	pool, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()
	log.Println("Successfully connected to the database!")

	subRepo := database.NewPostgresRepository(pool)
	userRepo := database.NewPostgresUserRepository(pool)

	subService := subscription.NewService(subRepo)
	authService := auth.NewAuthService(userRepo)

	subHandler := web.NewSubscriptionHandler(subService)
	authHandler := web.NewAuthHandler(authService, jwtSecret)

	router := web.SetupRouter(subHandler, authHandler, jwtSecret)

	port := fmt.Sprintf(":%s", apiPort)
	log.Printf("Server is running on port %s", port)

	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
