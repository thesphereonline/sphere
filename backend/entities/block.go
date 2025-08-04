package entities

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Block struct {
	Index        uint64        `json:"index"`
	Timestamp    time.Time     `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	PrevHash     string        `json:"prev_hash"`
	Hash         string        `json:"hash"`
	Nonce        uint64        `json:"nonce"` // For PoW or consensus
}

// ComputeHash calculates the block's hash based on its content.
func (b *Block) ComputeHash() string {
	record := string(b.Index) + b.Timestamp.String() + b.PrevHash + string(b.Nonce)
	for _, tx := range b.Transactions {
		record += tx.Hash()
	}
	hash := sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}
