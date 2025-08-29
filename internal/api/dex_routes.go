package api

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"sphere/internal/modules/dex"

	"github.com/gofiber/fiber/v2"
)

func registerDexRoutes(app *fiber.App, db *sql.DB, dexModule *dex.Dex) {
	// list pools
	app.Get("/dex/pools", func(c *fiber.Ctx) error {
		rows, err := db.QueryContext(context.Background(), `SELECT id, token_a, token_b, reserve_a, reserve_b, total_lp FROM pools ORDER BY id`)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		defer rows.Close()
		var out []any
		for rows.Next() {
			var id int
			var a, b, ra, rb, tl string
			if err := rows.Scan(&id, &a, &b, &ra, &rb, &tl); err != nil {
				return err
			}
			out = append(out, fiber.Map{"id": id, "token_a": a, "token_b": b, "reserve_a": ra, "reserve_b": rb, "total_lp": tl})
		}
		return c.JSON(out)
	})

	// create pool
	app.Post("/dex/pools", func(c *fiber.Ctx) error {
		var body struct {
			TokenA string `json:"tokenA"`
			TokenB string `json:"tokenB"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		if _, err := db.ExecContext(context.Background(), `INSERT INTO pools (token_a, token_b) VALUES ($1,$2)`, body.TokenA, body.TokenB); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.SendStatus(201)
	})

	// add liquidity
	app.Post("/dex/pools/:id/add", func(c *fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		var body struct {
			Owner   string `json:"owner"`
			AmountA string `json:"amountA"`
			AmountB string `json:"amountB"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		if err := dexModule.AddLiquidity(context.Background(), id, body.Owner, body.AmountA, body.AmountB); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.SendStatus(200)
	})

	// swap
	app.Post("/dex/pools/:id/swap", func(c *fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		var body struct {
			FromToken string `json:"fromToken"`
			AmountIn  string `json:"amountIn"`
			MinOut    string `json:"minOut"`
			Trader    string `json:"trader"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		out, err := dexModule.Swap(context.Background(), id, body.FromToken, body.AmountIn, body.MinOut, body.Trader)
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}
		return c.JSON(fiber.Map{"amountOut": out})
	})
}
