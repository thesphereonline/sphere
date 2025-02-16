package main

import (
	"fmt"
	"log"
	"time"
	"your-username/blockchain/blockchain"
	"your-username/blockchain/network"
	"your-username/blockchain/transaction"
	"your-username/blockchain/wallet"
)

func main() {
	// Create blockchain
	fmt.Println("Creating blockchain...")
	bc := blockchain.CreateBlockchain("blockchain.db")
	if bc == nil {
		log.Fatal("Failed to create blockchain")
	}

	// Create two wallets
	fmt.Println("\nCreating wallets...")
	wallet1 := wallet.NewWallet()
	wallet2 := wallet.NewWallet()

	fmt.Printf("Wallet 1 address: %s\n", wallet1.Address)
	fmt.Printf("Wallet 2 address: %s\n", wallet2.Address)

	// Start network server
	fmt.Println("\nStarting network server...")
	server := network.NewServer(bc)
	go func() {
		err := server.Start("8080")
		if err != nil {
			log.Printf("Server error: %v", err)
		}
	}()
	fmt.Println("Server started on port 8080")

	// Wait for server to start
	time.Sleep(time.Second)

	// Create and process transactions
	fmt.Println("\nCreating transactions...")

	// Transaction 1: Send 50 coins
	tx1 := transaction.NewTransaction(wallet1.Address, wallet2.Address, 50)
	if tx1 == nil {
		log.Fatal("Failed to create transaction")
	}

	signature1 := wallet1.Sign([]byte(tx1.ID))
	tx1.Signature = signature1

	err := bc.AddTransaction(tx1)
	if err != nil {
		log.Printf("Transaction error: %v", err)
	}

	fmt.Printf("Transaction created: %s -> %s, Amount: 50\n",
		truncateString(wallet1.Address, 10),
		truncateString(wallet2.Address, 10))

	// Mine a block
	fmt.Println("\nMining first block...")
	block1 := bc.AddBlock(bc.PendingTxs)
	fmt.Printf("Block mined! Hash: %s\n", block1.Hash)

	// Check balances
	fmt.Println("\nChecking balances after first transaction...")
	fmt.Printf("Wallet 1 balance: %.2f\n", bc.GetBalance(wallet1.Address))
	fmt.Printf("Wallet 2 balance: %.2f\n", bc.GetBalance(wallet2.Address))

	// Keep program running
	select {}
}

func truncateString(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
