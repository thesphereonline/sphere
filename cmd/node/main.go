package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"your-username/blockchain/blockchain"
	"your-username/blockchain/config"
	"your-username/blockchain/network"
	"your-username/blockchain/wallet"
)

func main() {
	// Parse command line flags
	port := flag.String("port", "8080", "Port to listen on")
	peers := flag.String("peers", "", "Comma-separated list of peer addresses")
	dataDir := flag.String("datadir", filepath.Join(".", "data"), "Directory to store blockchain data")
	minerAddress := flag.String("miner", "", "Miner wallet address")
	flag.Parse()

	// Create data directory with Windows-compatible path
	if err := os.MkdirAll(filepath.Clean(*dataDir), 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Initialize configuration
	cfg := config.NewConfig(*dataDir, *port)

	// Initialize or load wallet
	var minerWallet *wallet.Wallet
	if *minerAddress == "" {
		minerWallet = wallet.NewWallet()
		if err := minerWallet.SaveToFile(); err != nil {
			log.Printf("Warning: Failed to save wallet: %v", err)
		}
		log.Printf("Created new miner wallet: %s", minerWallet.Address)
	} else {
		var err error
		minerWallet = wallet.LoadWalletFromFile(*minerAddress)
		if minerWallet == nil {
			log.Fatalf("Failed to load wallet: %v", err)
		}
	}

	// Initialize blockchain
	bc := blockchain.CreateBlockchain(cfg.GetBlockchainPath())
	if bc == nil {
		log.Fatal("Failed to initialize blockchain")
	}

	// Initialize network server
	server := network.NewServer(bc)
	server.MinerAddress = minerWallet.Address

	// Connect to initial peers
	if *peers != "" {
		for _, peer := range strings.Split(*peers, ",") {
			peer = strings.TrimSpace(peer)
			if peer != "" {
				server.AddPeer(peer)
			}
		}
	}

	// Start server
	errChan := make(chan error, 1)
	go func() {
		log.Printf("Starting node on port %s", *port)
		errChan <- server.Start(*port)
	}()

	// Handle shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		log.Fatalf("Server error: %v", err)
	case sig := <-sigChan:
		log.Printf("Received signal %v, shutting down...", sig)
		server.Stop()
		time.Sleep(time.Second) // Give time for cleanup
	}
}
