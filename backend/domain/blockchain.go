package domain

import (
	"errors"

	"github.com/thesphereonline/sphere/entities"
)

type Blockchain struct {
	chain []entities.Block
}

func (bc *Blockchain) Chain() []entities.Block {
	return bc.chain
}

func NewBlockchain(genesis entities.Block) *Blockchain {
	return &Blockchain{
		chain: []entities.Block{genesis},
	}
}

func (bc *Blockchain) GetLastBlock() entities.Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) AddBlock(block entities.Block) error {
	if block.PrevHash != bc.GetLastBlock().Hash {
		return errors.New("invalid previous hash")
	}
	if block.Hash != block.ComputeHash() {
		return errors.New("invalid block hash")
	}
	bc.chain = append(bc.chain, block)
	return nil
}
