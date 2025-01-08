package network

import (
	"fmt"
	"net/rpc"
	"sphere/core"
	"time"
)

// SyncBlockchain requests blockchain data from a peer
func SyncBlockchain(peerAddress string) ([]core.Block, error) {
	var blockchain []core.Block
	retryCount := 3

	for i := 0; i < retryCount; i++ {
		client, err := rpc.Dial("tcp", peerAddress)
		if err != nil {
			fmt.Printf("Failed to connect to peer (attempt %d/%d): %v\n", i+1, retryCount, err)
			time.Sleep(2 * time.Second) // Wait before retrying
			continue
		}
		defer client.Close()

		err = client.Call("BlockchainService.SyncBlockchain", struct{}{}, &blockchain)
		if err != nil {
			fmt.Println("Error syncing blockchain:", err)
			return nil, err
		}

		fmt.Println("Blockchain synced successfully.")
		return blockchain, nil
	}

	return nil, fmt.Errorf("could not connect to peer after %d attempts", retryCount)
}

// SendTransaction sends a transaction to a peer
func SendTransaction(peerAddress string, tx core.Transaction) error {
	client, err := rpc.Dial("tcp", peerAddress)
	if err != nil {
		fmt.Println("Failed to connect to peer:", err)
		return err
	}
	defer client.Close()

	var reply bool
	err = client.Call("BlockchainService.SendTransaction", tx, &reply)
	if err != nil {
		fmt.Println("Error sending transaction:", err)
		return err
	}

	fmt.Println("Transaction sent successfully.")
	return nil
}
