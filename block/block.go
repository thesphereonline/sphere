package block

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"
	"your-username/blockchain/transaction"
)

type Block struct {
	Timestamp     int64
	Transactions  []*transaction.Transaction
	PrevBlockHash string
	Hash          string
	Nonce         int
	Difficulty    int
}

// CreateBlock creates a new block with the given data and previous block's hash
func CreateBlock(transactions []*transaction.Transaction, prevBlockHash string, difficulty int) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
		Difficulty:    difficulty,
	}

	block.Mine()
	return block
}

// CalculateHash calculates the hash of the block
func (b *Block) CalculateHash() string {
	record := string(b.Timestamp) + b.PrevBlockHash + string(b.Nonce)
	for _, tx := range b.Transactions {
		record += tx.ID
	}
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// Mine performs proof-of-work
func (b *Block) Mine() {
	target := strings.Repeat("0", b.Difficulty)

	for {
		b.Hash = b.CalculateHash()
		if b.Hash[:b.Difficulty] == target {
			return
		}
		b.Nonce++
	}
}
