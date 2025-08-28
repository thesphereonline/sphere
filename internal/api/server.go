// internal/api/server.go
package api

import (
	"sphere/internal/core"

	"github.com/gofiber/fiber/v2"
)

func StartServer(bc *core.Blockchain) {
	app := fiber.New()

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

	app.Listen(":8080")
}
