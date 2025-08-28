package main

import (
	"log"
	"os"
	"sphere/internal/api"
	"sphere/internal/core"
)

func main() {
	log.Println("🚀 Starting The Sphere backend...")

	// Create a new blockchain
	bc := core.NewBlockchain()

	// Get port from Railway (default 8080 locally)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println("⚠️  PORT not set, defaulting to 8080")
	}

	// Start API server
	if err := api.StartServer(bc, port); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}
