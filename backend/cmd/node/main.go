package main

import (
	"fmt"
	"time"

	"github.com/thesphereonline/sphere/internal/api"
	"github.com/thesphereonline/sphere/internal/core/block"
	"github.com/thesphereonline/sphere/internal/core/chain"
	"github.com/thesphereonline/sphere/internal/core/consensus"
	"github.com/thesphereonline/sphere/internal/core/mempool"
	"github.com/thesphereonline/sphere/internal/core/state"
	"github.com/thesphereonline/sphere/internal/db"
	"github.com/thesphereonline/sphere/internal/p2p"
)

func main() {
	// Connect to Railway-hosted PostgreSQL
	dbConn := "postgres://postgres:dSCrFhYHcYCQFFubeRWGUcZUHjCyTOCu@crossover.proxy.rlwy.net:43753/railway?sslmode=require"
	pg, err := db.NewPostgres(dbConn)
	if err != nil {
		panic(fmt.Errorf("❌ failed to connect to PostgreSQL: %w", err))
	}

	// Load blockchain from DB or create genesis block
	blocks, err := pg.LoadBlocks()
	if err != nil {
		panic(fmt.Errorf("❌ failed to load blockchain from DB: %w", err))
	}

	var bc *chain.Blockchain
	if len(blocks) == 0 {
		fmt.Println("🔧 No blocks found. Creating genesis block.")
		bc = chain.NewBlockchain()
		if err := pg.SaveBlock(bc.LastBlock()); err != nil {
			panic(fmt.Errorf("❌ failed to save genesis block: %w", err))
		}
	} else {
		bc = &chain.Blockchain{Blocks: blocks}
		fmt.Println("✅ Blockchain loaded. Current height:", bc.Height())
	}

	// Core components
	mem := mempool.NewMempool()
	st := state.NewState(pg)

	// Start REST API
	apiServer := api.NewServer(mem, st)
	go apiServer.Start("8080")

	// Start P2P server
	p2pServer := p2p.NewP2PServer("9000")
	p2pServer.Chain = bc
	p2pServer.Database = pg
	_ = p2pServer.Start()

	// Manually connect to other peer (optional)
	_ = p2pServer.ConnectToPeer("localhost:9001") // Replace with real peer

	// Block producer (runs every 10 seconds)
	go func() {
		for {
			time.Sleep(10 * time.Second)

			txs := mem.All()
			if len(txs) == 0 {
				continue
			}

			validator := consensus.PickValidator()
			prev := bc.LastBlock().Hash
			newBlock := block.NewBlock(uint64(bc.Height()), prev, validator, txs)

			fmt.Println("🧱 New block produced by", validator, ":", newBlock.Hash)

			// Save to DB
			if err := pg.SaveBlock(newBlock); err != nil {
				fmt.Println("❌ Failed to save block:", err)
				continue
			}

			// Add to in-memory blockchain
			bc.AddBlock(newBlock)

			// TODO: Clear mempool txs once included in block
		}
	}()

	fmt.Println("🚀 Sphere Protocol node is live.")
	select {} // Keep process alive
}
