package transaction

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Transaction struct {
	ID        string
	From      string
	To        string
	Amount    float64
	Timestamp int64
	Signature []byte
}

type TxInput struct {
	TxID      string
	OutIndex  int
	Signature string
}

type TxOutput struct {
	Value      float64
	PubKeyHash string
}

func NewTransaction(from, to string, amount float64) *Transaction {
	tx := &Transaction{
		From:      from,
		To:        to,
		Amount:    amount,
		Timestamp: time.Now().Unix(),
	}
	tx.ID = tx.calculateHash()
	return tx
}

func (tx *Transaction) calculateHash() string {
	record := tx.From + tx.To + string(tx.Timestamp) + string(int64(tx.Amount))
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}
