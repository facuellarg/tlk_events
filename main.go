package main

import (
	"challenge/repository"
	"challenge/service"
	"context"
	"log"
	"os"
)

func main() {
	ctx := context.Background()

	// Get database path from environment variable
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./events.db"
	}

	// Create database connection
	db, err := repository.NewDatabase(ctx, dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create table if it doesn't exist
	if err := db.CreateTable(ctx); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Create and start server
	server := service.NewServer(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := server.Start(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
