// internal/core/blockchain.go
package core

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"sync"
	"time"
)

type Transaction struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
	Fee    float64 `json:"fee"`
	Data   string  `json:"data"`
	Sig    string  `json:"sig"`
}

type Block struct {
	Height       int           `json:"height"`
	Hash         string        `json:"hash"`
	PrevHash     string        `json:"prevHash"`
	Timestamp    int64         `json:"timestamp"`
	Validator    string        `json:"validator"`
	Transactions []Transaction `json:"transactions"`
}

type Blockchain struct {
	Chain     []Block
	PendingTx []Transaction
	mu        sync.Mutex
}

func NewBlockchain() *Blockchain {
	genesis := Block{
		Height:    0,
		Hash:      "genesis",
		PrevHash:  "",
		Timestamp: time.Now().Unix(),
		Validator: "genesis-validator",
	}
	return &Blockchain{Chain: []Block{genesis}, PendingTx: []Transaction{}}
}

// AddTx adds a transaction to the mempool
func (bc *Blockchain) AddTx(tx Transaction) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.PendingTx = append(bc.PendingTx, tx)
}

// MinePending creates a new block if there are pending txs
func (bc *Blockchain) MinePending(validator string) *Block {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if len(bc.PendingTx) == 0 {
		return nil
	}

	prev := bc.Chain[len(bc.Chain)-1]
	block := Block{
		Height:       prev.Height + 1,
		PrevHash:     prev.Hash,
		Timestamp:    time.Now().Unix(),
		Validator:    validator,
		Transactions: bc.PendingTx,
	}
	block.Hash = bc.computeHash(block)

	bc.Chain = append(bc.Chain, block)
	bc.PendingTx = []Transaction{} // clear mempool

	return &block
}

func (bc *Blockchain) computeHash(b Block) string {
	data := strconv.Itoa(b.Height) + b.PrevHash + strconv.FormatInt(b.Timestamp, 10) + b.Validator
	for _, tx := range b.Transactions {
		data += tx.From +
			tx.To +
			strconv.FormatFloat(tx.Amount, 'f', -1, 64) +
			strconv.FormatFloat(tx.Fee, 'f', -1, 64) +
			tx.Data +
			tx.Sig
	}
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func (bc *Blockchain) GetLatestBlock() *Block {
	return &bc.Chain[len(bc.Chain)-1]
}
