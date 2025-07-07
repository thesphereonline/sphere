package block

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/thesphereonline/sphere/internal/core/transaction"
)

type Block struct {
	Index     uint64                    // block height
	Timestamp int64                     // unix time
	PrevHash  string                    // parent block hash
	Hash      string                    // current block hash
	Validator string                    // validator address (public key hash)
	Txs       []transaction.Transaction // transactions in block
	Signature string                    // validator's signature over block hash
}

// HashBlock calculates SHA256 hash of block contents
func (b *Block) HashBlock() string {
	raw := string(b.Index) + b.PrevHash + b.Validator + string(b.Timestamp)
	for _, tx := range b.Txs {
		raw += tx.Hash()
	}
	hash := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hash[:])
}

// NewBlock constructs and hashes a new block
func NewBlock(index uint64, prev string, validator string, txs []transaction.Transaction) *Block {
	b := &Block{
		Index:     index,
		Timestamp: time.Now().Unix(),
		PrevHash:  prev,
		Validator: validator,
		Txs:       txs,
	}
	b.Hash = b.HashBlock()
	return b
}
