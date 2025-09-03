package core

import (
	"sync"
)

// Mempool is a very small in-memory transaction pool.
type Mempool struct {
	mu  sync.Mutex
	txs []Transaction
	max int
}

// NewMempool creates a mempool with optional max size (0 = unlimited)
func NewMempool(max int) *Mempool {
	return &Mempool{txs: make([]Transaction, 0), max: max}
}

// AddTx adds a tx to the mempool
func (m *Mempool) AddTx(t Transaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.max > 0 && len(m.txs) >= m.max {
		// drop oldest (simple behavior) to make room
		m.txs = m.txs[1:]
	}
	m.txs = append(m.txs, t)
	return nil
}

// Flush returns all pending txs and clears the mempool (for mining)
func (m *Mempool) Flush() []Transaction {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := append([]Transaction(nil), m.txs...)
	m.txs = m.txs[:0]
	return out
}

// Len returns number of pending txs
func (m *Mempool) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.txs)
}
