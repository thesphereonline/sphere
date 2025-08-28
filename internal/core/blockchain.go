// internal/core/blockchain.go
package core

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Transaction struct {
	From   string `json:"From"`
	To     string `json:"To"`
	Amount int    `json:"Amount"`
	Fee    int    `json:"Fee"`
	Data   string `json:"Data"`
	Sig    string `json:"Sig"` // <-- was []byte or base64, now plain string
}

type Block struct {
	Height       uint64
	PrevHash     string
	Timestamp    int64
	Transactions []Transaction
	Validator    string
	Hash         string
}

type Blockchain struct {
	Chain []*Block
}

func NewBlockchain() *Blockchain {
	genesis := &Block{
		Height:    0,
		PrevHash:  "",
		Timestamp: time.Now().Unix(),
	}
	genesis.Hash = hashBlock(genesis)
	return &Blockchain{Chain: []*Block{genesis}}
}

func (bc *Blockchain) AddBlock(txs []Transaction, validator string) *Block {
	prev := bc.Chain[len(bc.Chain)-1]
	block := &Block{
		Height:       prev.Height + 1,
		PrevHash:     prev.Hash,
		Timestamp:    time.Now().Unix(),
		Transactions: txs,
		Validator:    validator,
	}
	block.Hash = hashBlock(block)
	bc.Chain = append(bc.Chain, block)
	return block
}

func hashBlock(b *Block) string {
	h := sha256.New()
	h.Write([]byte(string(b.Height)))
	h.Write([]byte(b.PrevHash))
	h.Write([]byte(string(b.Timestamp)))
	for _, tx := range b.Transactions {
		h.Write([]byte(tx.From + tx.To))
	}
	return hex.EncodeToString(h.Sum(nil))
}
