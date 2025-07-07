package chain

import (
	"sync"

	"github.com/thesphereonline/sphere/internal/core/block"
)

type Blockchain struct {
	Blocks []*block.Block
	Lock   sync.Mutex
}

func NewBlockchain() *Blockchain {
	genesis := block.NewBlock(0, "", "genesis", nil)
	return &Blockchain{Blocks: []*block.Block{genesis}}
}

func (bc *Blockchain) AddBlock(b *block.Block) {
	bc.Lock.Lock()
	defer bc.Lock.Unlock()
	bc.Blocks = append(bc.Blocks, b)
}

func (bc *Blockchain) LastBlock() *block.Block {
	bc.Lock.Lock()
	defer bc.Lock.Unlock()
	return bc.Blocks[len(bc.Blocks)-1]
}

func (bc *Blockchain) Height() int {
	bc.Lock.Lock()
	defer bc.Lock.Unlock()
	return len(bc.Blocks)
}

func (bc *Blockchain) Contains(hash string) bool {
	bc.Lock.Lock()
	defer bc.Lock.Unlock()

	for _, b := range bc.Blocks {
		if b.Hash == hash {
			return true
		}
	}
	return false
}
