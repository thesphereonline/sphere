package core

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Block struct {
	Index        int           `json:"index"`
	Timestamp    int64         `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	PrevHash     string        `json:"prev_hash"`
	Hash         string        `json:"hash"`
	Validator    string        `json:"validator"`
}

func (b *Block) CalculateHash() string {
	record := string(b.Index) + string(b.Timestamp) + b.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	return hex.EncodeToString(h.Sum(nil))
}

func NewBlock(prevBlock Block, transactions []Transaction, validator string) Block {
	block := Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
		PrevHash:     prevBlock.Hash,
		Validator:    validator,
	}
	block.Hash = block.CalculateHash()
	return block
}
