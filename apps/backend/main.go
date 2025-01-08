package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/thesphereonline/sphere/backend/database"
	"github.com/thesphereonline/sphere/backend/routes"
	"github.com/thesphereonline/sphere/backend/services"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to Redis and PostgreSQL
	database.ConnectDB()
	database.ConnectRedis()

	// Start Redis pub/sub
	go services.SubscribeToNotifications()

	// Initialize Fiber app
	app := fiber.New()

	// Middlewares
	app.Use(logger.New())  // Logs requests
	app.Use(recover.New()) // Recovers from panics
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Setup routes
	routes.SetupRoutes(app)

	// Start the server
	port := os.Getenv("PORT")
	log.Fatal(app.Listen(":" + port))
}
