package api

import (
	"database/sql"
	"strconv"

	"sphere/internal/modules/dex"

	"github.com/gofiber/fiber/v2"
)

func registerDexRoutes(app *fiber.App, db *sql.DB, dexModule *dex.Module) {
	// List all pools
	app.Get("/dex/pools", func(c *fiber.Ctx) error {
		pools, err := dexModule.ListPools()
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON(pools)
	})

	// Create a new pool
	app.Post("/dex/pools", func(c *fiber.Ctx) error {
		var body struct {
			TokenA   string  `json:"tokenA"`
			TokenB   string  `json:"tokenB"`
			ReserveA float64 `json:"reserveA"`
			ReserveB float64 `json:"reserveB"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		pool, err := dexModule.AddPool(body.TokenA, body.TokenB, body.ReserveA, body.ReserveB)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.Status(201).JSON(pool)
	})

	// Swap tokens
	app.Post("/dex/pools/:id/swap", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(400).SendString("invalid pool id")
		}

		var body struct {
			FromToken string  `json:"fromToken"`
			AmountIn  float64 `json:"amountIn"`
			MinOut    float64 `json:"minOut"`
			Trader    string  `json:"trader"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		amountOut, err := dexModule.Swap(id, body.FromToken, body.AmountIn, body.MinOut, body.Trader)
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}
		return c.JSON(fiber.Map{"amountOut": amountOut})
	})
}
