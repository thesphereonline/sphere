package network

import (
	"fmt"
	"net"
	"net/rpc"
	"sphere/core"
)

// BlockchainService provides RPC methods for blockchain operations
type BlockchainService struct {
	Blockchain *core.Blockchain
}

// SyncBlockchain returns the latest blockchain state to the requesting node
func (b *BlockchainService) SyncBlockchain(_ struct{}, reply *[]core.Block) error {
	*reply = b.Blockchain.Chain
	fmt.Println("Blockchain sync requested. Sending blockchain data.")
	return nil
}

// SendTransaction receives a transaction from another node and adds it to the mempool
func (b *BlockchainService) SendTransaction(tx core.Transaction, reply *bool) error {
	fmt.Printf("Received transaction from %s to %s for amount: %.2f\n", tx.Inputs[0].Address, tx.Outputs[0].Address, tx.Outputs[0].Amount)
	*reply = true
	return nil
}

// StartRPCServer initializes and starts the RPC server on a dynamic port
func StartRPCServer(blockchain *core.Blockchain) (string, error) {
	blockchainService := &BlockchainService{Blockchain: blockchain}
	rpc.Register(blockchainService)

	listener, err := net.Listen("tcp", ":0") // Use ":0" to dynamically assign an available port
	if err != nil {
		fmt.Println("Failed to start RPC server:", err)
		return "", err
	}

	addr := listener.Addr().String() // Get the dynamically assigned port
	fmt.Println("RPC server started on", addr)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Connection error:", err)
				continue
			}
			go rpc.ServeConn(conn)
		}
	}()
	return addr, nil
}
