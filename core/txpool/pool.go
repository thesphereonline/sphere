package txpool

import (
	"sync"

	"github.com/thesphereonline/sphere/core/types"
)

type TxPool struct {
	Txs []types.Transaction
	mu  sync.Mutex
}

func NewPool() *TxPool {
	return &TxPool{
		Txs: make([]types.Transaction, 0),
	}
}

func (p *TxPool) Add(tx types.Transaction) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Txs = append(p.Txs, tx)
}

func (p *TxPool) Flush() []types.Transaction {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := p.Txs
	p.Txs = []types.Transaction{}
	return out
}
