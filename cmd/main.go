package main

import (
	"bufio"
	"fmt"
	"os"
	"sphere/api"
	"sphere/core"
	"sphere/network"
	"strings"
)

func main() {
	fmt.Println("Starting Sphere Blockchain Node...")

	// Initialize blockchain
	blockchain := core.NewBlockchain()

	// Start RPC server and get assigned port
	rpcAddress, err := network.StartRPCServer(blockchain)
	if err != nil {
		fmt.Println("Failed to start RPC server:", err)
		return
	}

	fmt.Printf("Node started at %s\n", rpcAddress)

	// Ask for peer node to connect to
	fmt.Print("Enter peer node address (or press Enter to skip): ")
	reader := bufio.NewReader(os.Stdin)
	peerAddress, _ := reader.ReadString('\n')
	peerAddress = strings.TrimSpace(peerAddress)

	// If user provided a peer, attempt to sync
	if peerAddress != "" {
		fmt.Printf("Attempting to sync with peer at %s...\n", peerAddress)
		syncedBlockchain, err := network.SyncBlockchain(peerAddress)
		if err == nil {
			blockchain.Chain = syncedBlockchain
			fmt.Println("Blockchain synced successfully!")
		} else {
			fmt.Println("Failed to sync blockchain:", err)
		}

		// Send a test transaction
		tx := core.NewTransaction(
			[]core.UTXO{{Address: "Alice", Amount: 10}},
			[]core.UTXO{{Address: "Bob", Amount: 10}},
			"AlicePublicKey",
		)

		err = network.SendTransaction(peerAddress, tx)
		if err != nil {
			fmt.Println("Transaction failed:", err)
		}
	}

	// Start REST API
	go api.StartServer()

	// Keep the node running
	select {}
}
