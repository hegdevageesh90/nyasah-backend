package main

import (
	"log"
	"nyasah/api"
	"nyasah/config"
	"nyasah/database"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize and start the server
	server := api.NewServer(cfg, db)
	server.Start()
}