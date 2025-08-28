package api

import (
	"log"
	"os"
	"sphere/internal/core"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func StartServer(bc *core.Blockchain) {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://www.thesphere.online",
		AllowMethods:     "GET,POST,OPTIONS",
		AllowHeaders:     "Content-Type",
		AllowCredentials: true,
	}))

	app.Get("/blocks", func(c *fiber.Ctx) error {
		return c.JSON(bc.Chain)
	})

	app.Post("/tx", func(c *fiber.Ctx) error {
		tx := new(core.Transaction)
		if err := c.BodyParser(tx); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		block := bc.AddBlock([]core.Transaction{*tx}, "validator-1")
		return c.JSON(block)
	})

	// Use Railway PORT environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s\n", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
