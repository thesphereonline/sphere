package core

import (
	"crypto/ed25519"
	"encoding/hex"
	"time"
)

// Transaction represents a transaction in the Sphere blockchain
type Transaction struct {
	ID        string `json:"id"`
	Inputs    []UTXO `json:"inputs"`
	Outputs   []UTXO `json:"outputs"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
	PublicKey string `json:"public_key"`
}

// UTXO (Unspent Transaction Output) Model
type UTXO struct {
	TxID    string  `json:"tx_id"`
	Index   int     `json:"index"`
	Address string  `json:"address"`
	Amount  float64 `json:"amount"`
}

// SignTransaction signs a transaction using the user's private key
func (tx *Transaction) SignTransaction(privateKey ed25519.PrivateKey) {
	data := tx.ID + tx.PublicKey
	signature := ed25519.Sign(privateKey, []byte(data))
	tx.Signature = hex.EncodeToString(signature)
}

// VerifyTransaction verifies if a transaction is valid
func (tx *Transaction) VerifyTransaction() bool {
	data := tx.ID + tx.PublicKey
	signature, err := hex.DecodeString(tx.Signature)
	if err != nil {
		return false
	}
	return ed25519.Verify([]byte(tx.PublicKey), []byte(data), signature)
}

// NewTransaction creates a new transaction
func NewTransaction(inputs []UTXO, outputs []UTXO, publicKey string) Transaction {
	tx := Transaction{
		ID:        hex.EncodeToString([]byte(time.Now().String())),
		Inputs:    inputs,
		Outputs:   outputs,
		Timestamp: time.Now().Unix(),
		PublicKey: publicKey,
	}
	return tx
}
