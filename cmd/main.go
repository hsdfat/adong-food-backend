package main

import (
	"adong-be/server"
	"adong-be/store"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	var err error
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=adong password=adong123 dbname=adongfood port=5432 sslmode=disable"
	}

	store.DB.GormClient, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")
	s := server.SetupRouter() 
	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "18080"
	}

	log.Printf("Server starting on port %s", port)
	if err := s.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
