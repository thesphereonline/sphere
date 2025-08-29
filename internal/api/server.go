package api

import (
	"database/sql"
	"fmt"
	"log"
	"sphere/internal/core"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// StartServer runs the API server
func StartServer(bc *core.Blockchain, port string, db *sql.DB) error {
	app := fiber.New()

	// Enable CORS for frontend
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://www.thesphere.online", // your frontend
		AllowMethods:     "GET,POST,OPTIONS",
		AllowHeaders:     "Content-Type",
		AllowCredentials: true,
	}))

	// Blockchain routes
	app.Get("/blocks", func(c *fiber.Ctx) error {
		return c.JSON(bc.Chain)
	})

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	app.Post("/tx", func(c *fiber.Ctx) error {
		tx := new(core.Transaction)
		if err := c.BodyParser(tx); err != nil {
			log.Printf("❌ Failed to parse tx: %v", err)
			return c.Status(400).SendString(err.Error())
		}
		block := bc.AddBlock([]core.Transaction{*tx}, "validator-1")
		log.Printf("✅ New block added: %+v", block)
		return c.JSON(block)
	})

	// Start listening
	addr := fmt.Sprintf(":%s", port)
	log.Printf("🌐 Listening on %s ...", addr)
	return app.Listen(addr)
}
