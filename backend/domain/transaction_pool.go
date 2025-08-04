package domain

import (
	"sync"

	"github.com/thesphereonline/sphere/entities"
)

type TransactionPool struct {
	mu           sync.RWMutex
	transactions map[string]entities.Transaction
}

func NewTransactionPool() *TransactionPool {
	return &TransactionPool{
		transactions: make(map[string]entities.Transaction),
	}
}

func (tp *TransactionPool) Add(tx entities.Transaction) {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	tp.transactions[tx.ID] = tx
}

func (tp *TransactionPool) Remove(txID string) {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	delete(tp.transactions, txID)
}

func (tp *TransactionPool) GetAll() []entities.Transaction {
	tp.mu.RLock()
	defer tp.mu.RUnlock()
	txs := make([]entities.Transaction, 0, len(tp.transactions))
	for _, tx := range tp.transactions {
		txs = append(txs, tx)
	}
	return txs
}
