package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"wallet-service/internal/api"
	"wallet-service/internal/repository"
	"wallet-service/internal/server"
	"wallet-service/internal/service"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load("config.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	serverPort := os.Getenv("SERVER_PORT")

	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"wallet-db",
		"5432",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	var db *sqlx.DB
	var err error

	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		db, err = sqlx.Connect("postgres", conn)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxAttempts, err)
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database after %d attempts: %v", maxAttempts, err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	walletRepo := repository.NewWallet(db)
	walletService := service.NewWalletService(walletRepo)
	walletHandler := api.NewHandler(walletService)

	log.Printf("Starting server on port: %s", serverPort)
	if err := server.Run(serverPort, walletHandler); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
