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

// StartServer runs the API server
func StartServer(bc *core.Blockchain, port string, db *sql.DB) error {
	app := fiber.New()

	// Enable CORS for frontend
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://www.thesphere.online, http://localhost:3000",
		AllowMethods:     "GET,POST,OPTIONS",
		AllowHeaders:     "Content-Type",
		AllowCredentials: true,
	}))

	// --- Initialize modules ---
	dexModule := dex.New(db, 5) // protocolFeeBps

	// --- Blockchain routes ---
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
		log.Printf("✅ New block added: %+v", block)
		return c.JSON(block)
	})

	// --- Health check ---
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	// --- DEX routes ---
	app.Get("/dex/pools", func(c *fiber.Ctx) error {
		pools, err := dexModule.ListPools()
		if err != nil {
			log.Printf("❌ ListPools error: %v", err)
			return c.Status(500).SendString("failed to fetch pools")
		}
		return c.JSON(pools)
	})

	app.Post("/dex/pools", func(c *fiber.Ctx) error {
		var req struct {
			TokenA   string  `json:"token_a"`
			TokenB   string  `json:"token_b"`
			ReserveA float64 `json:"reserve_a"`
			ReserveB float64 `json:"reserve_b"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		pool, err := dexModule.AddPool(req.TokenA, req.TokenB, req.ReserveA, req.ReserveB)
		if err != nil {
			log.Printf("❌ AddPool error: %v", err)
			return c.Status(500).SendString("failed to add pool")
		}
		return c.JSON(pool)
	})

	app.Post("/dex/pools/:id/swap", func(c *fiber.Ctx) error {
		var req struct {
			FromToken string  `json:"fromToken"`
			AmountIn  float64 `json:"amountIn"`
			MinOut    float64 `json:"minOut"`
			Trader    string  `json:"trader"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		poolID, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(400).SendString("invalid pool id")
		}

		amountOut, err := dexModule.Swap(poolID, req.FromToken, req.AmountIn, req.MinOut, req.Trader)
		if err != nil {
			log.Printf("❌ Swap error: %v", err)
			return c.Status(400).SendString(err.Error())
		}

		return c.JSON(fiber.Map{
			"amountOut": amountOut,
		})
	})

	// Start listening
	addr := fmt.Sprintf(":%s", port)
	log.Printf("🌐 Listening on %s ...", addr)
	return app.Listen(addr)
}
