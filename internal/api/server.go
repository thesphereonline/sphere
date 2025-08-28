package api

import (
	"sphere/internal/core"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func StartServer(bc *core.Blockchain) {
	app := fiber.New()

	// Enable CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://www.thesphere.online", // your frontend
		AllowMethods:     "GET,POST,OPTIONS",
		AllowHeaders:     "Content-Type",
		AllowCredentials: true,
	}))

	// Get all blocks
	app.Get("/blocks", func(c *fiber.Ctx) error {
		return c.JSON(bc.Chain)
	})

	// Submit transaction (creates a block)
	app.Post("/tx", func(c *fiber.Ctx) error {
		tx := new(core.Transaction)
		if err := c.BodyParser(tx); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		block := bc.AddBlock([]core.Transaction{*tx}, "validator-1")
		return c.JSON(block)
	})

	// Start server
	app.Listen(":8080")
}
