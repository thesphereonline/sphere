package main

import (
	"database/sql"
	"log"
	"os"
	"sphere/internal/api"
	"sphere/internal/core"
	"sphere/internal/db"

	_ "github.com/lib/pq"
)

func main() {
	log.Println("🚀 Starting The Sphere backend...")

	// Database
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalf("DATABASE_URL is required")
	}

	sqlDB, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("db ping: %v", err)
	}
	log.Println("✅ Connected to Postgres")

	// Apply migrations
	if err := db.ApplyMigrations(sqlDB, "./migrations"); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}
	log.Println("✅ Migrations applied")

	// Blockchain
	bc := core.NewBlockchain()

	// Start HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println("⚠️  PORT not set, defaulting to 8080")
	}

	if err := api.StartServer(bc, port, sqlDB); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}
