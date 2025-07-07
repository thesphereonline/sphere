package mempool

import (
	"sync"

	"github.com/thesphereonline/sphere/internal/core/transaction"
)

type Mempool struct {
	TxMap map[string]transaction.Transaction
	Lock  sync.Mutex
}

func NewMempool() *Mempool {
	return &Mempool{TxMap: make(map[string]transaction.Transaction)}
}

func (m *Mempool) Add(tx transaction.Transaction) bool {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	hash := tx.Hash()
	if _, exists := m.TxMap[hash]; exists {
		return false
	}
	m.TxMap[hash] = tx
	return true
}

func (m *Mempool) All() []transaction.Transaction {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	txs := make([]transaction.Transaction, 0, len(m.TxMap))
	for _, tx := range m.TxMap {
		txs = append(txs, tx)
	}
	return txs
}
