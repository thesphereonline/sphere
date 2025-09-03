package api

import (
	"database/sql"
	"fmt"
	"log"
	"sphere/internal/core"
	"sphere/internal/db"
	"sphere/internal/modules/dex"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// StartServer starts the HTTP API. It wires a mempool and persists blocks to DB when mined.
func StartServer(bc *core.Blockchain, port string, dbConn *sql.DB) error {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://www.thesphere.online, http://localhost:3000",
		AllowMethods:     "GET,POST,OPTIONS",
		AllowHeaders:     "Content-Type",
		AllowCredentials: true,
	}))

	dexModule := dex.New(dbConn, 30) // fee default 30 bps

	// Create mempool and attach to runtime
	mempool := core.NewMempool(0)

	// Blockchain routes
	app.Get("/blocks", func(c *fiber.Ctx) error {
		// return in-memory chain (fast)
		return c.JSON(bc.Chain)
	})

	// Get latest persisted blocks (from DB)
	app.Get("/blocks/persisted", func(c *fiber.Ctx) error {
		rows, err := dbConn.Query("SELECT id, height, hash, prev_hash, timestamp, validator, created_at FROM blocks ORDER BY height DESC LIMIT 50")
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		defer rows.Close()
		type Row struct {
			ID        int    `json:"id"`
			Height    int    `json:"height"`
			Hash      string `json:"hash"`
			PrevHash  string `json:"prev_hash"`
			Timestamp int64  `json:"timestamp"`
			Validator string `json:"validator"`
		}
		var out []Row
		for rows.Next() {
			var r Row
			var createdAt interface{}
			if err := rows.Scan(&r.ID, &r.Height, &r.Hash, &r.PrevHash, &r.Timestamp, &r.Validator, &createdAt); err != nil {
				return c.Status(500).SendString(err.Error())
			}
			out = append(out, r)
		}
		return c.JSON(out)
	})

	// Submit tx -> goes into mempool
	app.Post("/tx", func(c *fiber.Ctx) error {
		tx := new(core.Transaction)
		if err := c.BodyParser(tx); err != nil {
			log.Printf("❌ Failed to parse tx: %v", err)
			return c.Status(400).SendString(err.Error())
		}
		if err := mempool.AddTx(*tx); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON(fiber.Map{"status": "queued"})
	})

	// Mine / flush mempool into a block (dev endpoint)
	app.Post("/mine", func(c *fiber.Ctx) error {
		pending := mempool.Flush()
		if len(pending) == 0 {
			return c.Status(400).SendString("no pending txs")
		}
		block := bc.AddBlock(pending, "local-miner-1")
		// persist block & transactions
		if _, err := db.SaveBlock(dbConn, block); err != nil {
			log.Printf("❌ failed to save block: %v", err)
			// note: don't rollback in-memory chain — in prod you'd coordinate differently
			return c.Status(500).SendString(fmt.Sprintf("save block: %v", err))
		}
		return c.JSON(block)
	})

	// Validator endpoints (basic)
	app.Post("/validators/register", func(c *fiber.Ctx) error {
		var body struct {
			Address string `json:"address"`
			Stake   string `json:"stake"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		_, err := dbConn.Exec("INSERT INTO validators (address, stake, active, created_at) VALUES ($1, $2, TRUE, now()) ON CONFLICT (address) DO UPDATE SET stake = validators.stake + EXCLUDED.stake", body.Address, body.Stake)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON(fiber.Map{"status": "registered"})
	})

	app.Post("/validators/delegate", func(c *fiber.Ctx) error {
		var body struct {
			Delegator string `json:"delegator"`
			Validator string `json:"validator"`
			Amount    string `json:"amount"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		_, err := dbConn.Exec("INSERT INTO delegations (delegator, validator, amount, created_at) VALUES ($1, $2, $3, now()) ON CONFLICT (delegator, validator) DO UPDATE SET amount = (delegations.amount::numeric + EXCLUDED.amount::numeric)", body.Delegator, body.Validator, body.Amount)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON(fiber.Map{"status": "delegated"})
	})

	app.Get("/validators", func(c *fiber.Ctx) error {
		rows, err := dbConn.Query("SELECT address, stake, commission_bps, active, created_at FROM validators ORDER BY stake::numeric DESC LIMIT 50")
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		defer rows.Close()
		type V struct {
			Address       string `json:"address"`
			Stake         string `json:"stake"`
			CommissionBps int    `json:"commission_bps"`
			Active        bool   `json:"active"`
		}
		out := []V{}
		for rows.Next() {
			var v V
			var createdAt interface{}
			if err := rows.Scan(&v.Address, &v.Stake, &v.CommissionBps, &v.Active, &createdAt); err != nil {
				return c.Status(500).SendString(err.Error())
			}
			out = append(out, v)
		}
		return c.JSON(out)
	})

	// DEX routes
	registerDexRoutes(app, dbConn, dexModule)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("🌐 Listening on %s ...", addr)
	return app.Listen(addr)
}
