package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/thesphereonline/sphere/api/rest"
	"github.com/thesphereonline/sphere/db"
)

func main() {
	// Load .env if present
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize DB connection
	db.Init()

	// Start HTTP server with our router
	http.Handle("/", rest.Router())

	log.Println("🚀 API Server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
