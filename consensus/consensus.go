package consensus

import (
	"errors"
	"your-username/blockchain/block"
)

type Consensus struct {
	MinDifficulty int
	MaxDifficulty int
}

func NewConsensus() *Consensus {
	return &Consensus{
		MinDifficulty: 4,
		MaxDifficulty: 6,
	}
}

func (c *Consensus) ValidateBlock(block *block.Block, prevBlock *block.Block) error {
	// Validate block hash
	if block.CalculateHash() != block.Hash {
		return errors.New("invalid block hash")
	}

	// Validate previous hash
	if block.PrevBlockHash != prevBlock.Hash {
		return errors.New("invalid previous block hash")
	}

	// Validate timestamp
	if block.Timestamp <= prevBlock.Timestamp {
		return errors.New("invalid block timestamp")
	}

	// Validate proof of work
	target := string(make([]byte, block.Difficulty))
	for i := 0; i < block.Difficulty; i++ {
		target += "0"
	}

	if block.Hash[:block.Difficulty] != target {
		return errors.New("invalid proof of work")
	}

	return nil
}

func (c *Consensus) AdjustDifficulty(lastBlock *block.Block, blockTime int64) int {
	// Adjust difficulty based on block time
	// Target block time: 10 seconds
	if blockTime < 5 { // Too fast
		return min(lastBlock.Difficulty+1, c.MaxDifficulty)
	} else if blockTime > 15 { // Too slow
		return max(lastBlock.Difficulty-1, c.MinDifficulty)
	}
	return lastBlock.Difficulty
}
