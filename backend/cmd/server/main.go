package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/thesphereonline/sphere/config"
	"github.com/thesphereonline/sphere/internal/infra/db"
	"github.com/thesphereonline/sphere/internal/infra/logger"

	"github.com/thesphereonline/sphere/cmd/server/api" // import your API router package
	"github.com/thesphereonline/sphere/entities"
	"github.com/thesphereonline/sphere/node"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()
	log := logger.New()

	// Init DB
	pg, err := db.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer pg.Close()

	conn, err := pg.Pool.Acquire(context.Background())
	if err != nil {
		log.Fatalf("failed to acquire db connection: %v", err)
	}
	defer conn.Release()

	if err := db.RunMigrations(conn.Conn()); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Create genesis block (simplified for example)
	genesis := entities.Block{
		Index:     0,
		Timestamp: time.Now(),
		PrevHash:  "0",
		Hash:      "",
	}
	genesis.Hash = genesis.ComputeHash()

	// Initialize Global Node
	node.GlobalNode = node.NewNode(genesis)

	// Start node loop (mining etc) in background
	go node.GlobalNode.Start()
	defer node.GlobalNode.Stop()

	// Initialize HTTP server with API router
	httpServer := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: api.NewRouter(),
	}

	// Graceful shutdown setup
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("server started on port %s", cfg.Server.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down...")

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctxTimeout); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
}
