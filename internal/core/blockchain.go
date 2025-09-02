// internal/core/blockchain.go
package core

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

// Transaction represents a blockchain transaction
type Transaction struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
	Fee    float64 `json:"fee"`
	Data   string  `json:"data"`
	Sig    string  `json:"sig"`
}

// Block represents a blockchain block
type Block struct {
	Height       int           `json:"height"`
	Hash         string        `json:"hash"`
	PrevHash     string        `json:"prevHash"`
	Timestamp    int64         `json:"timestamp"`
	Validator    string        `json:"validator"`
	Transactions []Transaction `json:"transactions"`
}

// Blockchain is a chain of blocks
type Blockchain struct {
	Chain []Block
}

// NewBlockchain initializes the blockchain with a genesis block
func NewBlockchain() *Blockchain {
	genesis := Block{
		Height:    0,
		Hash:      "genesis",
		PrevHash:  "",
		Timestamp: time.Now().Unix(),
		Validator: "genesis-validator",
	}
	return &Blockchain{Chain: []Block{genesis}}
}

// AddBlock adds a new block with transactions to the chain
func (bc *Blockchain) AddBlock(txs []Transaction, validator string) *Block {
	prev := bc.Chain[len(bc.Chain)-1]
	block := Block{
		Height:       prev.Height + 1,
		PrevHash:     prev.Hash,
		Timestamp:    time.Now().Unix(),
		Validator:    validator,
		Transactions: txs,
	}
	block.Hash = bc.computeHash(block)
	bc.Chain = append(bc.Chain, block)
	return &block
}

// computeHash generates a SHA256 hash of the block
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

// GetLatestBlock returns the last block in the chain
func (bc *Blockchain) GetLatestBlock() *Block {
	return &bc.Chain[len(bc.Chain)-1]
}
