package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/thesphereonline/sphere/core/types"
	"github.com/thesphereonline/sphere/db"
)

type Blockchain struct {
	Chain    []types.Block
	Executor *Executor // manages global state
}

func NewBlockchain() *Blockchain {
	genesis := types.Block{
		Index:     0,
		Timestamp: time.Now().Unix(),
		PrevHash:  "0",
		Validator: "genesis",
		Txs:       []types.Transaction{},
		StateRoot: "genesis-root", // placeholder
	}
	genesis.Hash = computeHash(genesis)

	return &Blockchain{
		Chain:    []types.Block{genesis},
		Executor: NewExecutor(), // holds global state
	}
}

func (bc *Blockchain) AddBlock(txs []types.Transaction, validator string) types.Block {
	prev := bc.Chain[len(bc.Chain)-1]

	validTxs := make([]types.Transaction, 0)

	// Process each transaction through VM
	for _, tx := range txs {
		err := bc.Executor.ApplyTx(tx)
		if err != nil {
			log.Printf("❌ Invalid tx from %s: %v", tx.From, err)
			continue // skip invalid txs
		}
		validTxs = append(validTxs, tx)
	}

	// Create new block with valid txs only
	newBlock := types.Block{
		Index:     prev.Index + 1,
		Timestamp: time.Now().Unix(),
		PrevHash:  prev.Hash,
		Txs:       validTxs,
		Validator: validator,
		StateRoot: computeStateRoot(bc.Executor.State), // placeholder
	}

	newBlock.Hash = computeHash(newBlock)
	bc.Chain = append(bc.Chain, newBlock)

	log.Printf("✅ Mined Block #%d with %d txs", newBlock.Index, len(validTxs))

	err := db.InsertBlock(newBlock)
	if err != nil {
		log.Printf("DB insert block failed: %v", err)
	}

	for _, tx := range validTxs {
		err := db.InsertTransaction(tx, newBlock.Index)
		if err != nil {
			log.Printf("DB insert tx failed: %v", err)
		}
	}

	return newBlock
}

func computeHash(b types.Block) string {
	data := fmt.Sprintf("%d%s%s%d", b.Index, b.PrevHash, b.Validator, b.Timestamp)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func computeStateRoot(state *types.State) string {
	data, err := json.Marshal(state)
	if err != nil {
		return ""
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
