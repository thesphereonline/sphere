package node

import (
	"log"
	"time"

	"github.com/thesphereonline/sphere/domain"
	"github.com/thesphereonline/sphere/entities"
)

var GlobalNode *Node

type Node struct {
	Blockchain      *domain.Blockchain
	TransactionPool *domain.TransactionPool
	quit            chan struct{}
}

func NewNode(genesis entities.Block) *Node {
	return &Node{
		Blockchain:      domain.NewBlockchain(genesis),
		TransactionPool: domain.NewTransactionPool(),
		quit:            make(chan struct{}),
	}
}

func (n *Node) Start() {
	ticker := time.NewTicker(10 * time.Second) // block time interval
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			n.mineBlock()
		case <-n.quit:
			log.Println("Node shutting down")
			return
		}
	}
}

func (n *Node) Stop() {
	close(n.quit)
}

func (n *Node) mineBlock() {
	txs := n.TransactionPool.GetAll()
	if len(txs) == 0 {
		log.Println("No transactions to include in block")
		return
	}

	prevBlock := n.Blockchain.GetLastBlock()
	newBlock := entities.Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now(),
		Transactions: txs,
		PrevHash:     prevBlock.Hash,
		Nonce:        0, // placeholder for consensus nonce
	}

	newBlock.Hash = newBlock.ComputeHash()

	err := n.Blockchain.AddBlock(newBlock)
	if err != nil {
		log.Printf("Failed to add block: %v", err)
		return
	}

	// Clear included transactions
	for _, tx := range txs {
		n.TransactionPool.Remove(tx.ID)
	}

	log.Printf("Mined new block #%d with %d transactions", newBlock.Index, len(txs))
}
