package api

import (
	"database/sql"
	"fmt"
	"log"
	"sphere/internal/core"
	"sphere/internal/modules/dex"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func StartServer(bc *core.Blockchain, port string, db *sql.DB) error {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://www.thesphere.online, http://localhost:3000",
		AllowMethods:     "GET,POST,OPTIONS",
		AllowHeaders:     "Content-Type",
		AllowCredentials: true,
	}))

	dexModule := dex.New(db, 5)

	// Blockchain routes
	app.Get("/blocks", func(c *fiber.Ctx) error {
		return c.JSON(bc.Chain)
	})

	app.Post("/tx", func(c *fiber.Ctx) error {
		tx := new(core.Transaction)
		if err := c.BodyParser(tx); err != nil {
			log.Printf("❌ Failed to parse tx: %v", err)
			return c.Status(400).SendString(err.Error())
		}
		block := bc.AddBlock([]core.Transaction{*tx}, "validator-1")
		return c.JSON(block)
	})

	// DEX routes
	registerDexRoutes(app, db, dexModule)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("🌐 Listening on %s ...", addr)
	return app.Listen(addr)
}
