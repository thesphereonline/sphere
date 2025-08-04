package entities

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"
)

type Transaction struct {
	ID        string    `json:"id"`
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
	Signature string    `json:"signature"`
}

// Hash returns a unique hash of the transaction (excluding signature)
func (tx *Transaction) Hash() string {
	record := tx.Sender + tx.Receiver + fmt.Sprintf("%f", tx.Amount) + tx.Timestamp.String()
	hash := sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}

// SignTransaction signs the transaction using sender's private key
func (tx *Transaction) SignTransaction(privKey *ecdsa.PrivateKey) error {
	hash := sha256.Sum256([]byte(tx.Hash()))
	r, s, err := ecdsa.Sign(rand.Reader, privKey, hash[:])
	if err != nil {
		return err
	}
	signature := append(r.Bytes(), s.Bytes()...)
	tx.Signature = hex.EncodeToString(signature)
	return nil
}

// VerifySignature verifies the transaction's signature
func (tx *Transaction) VerifySignature(pubKey *ecdsa.PublicKey) bool {
	if tx.Signature == "" {
		return false
	}

	sigBytes, err := hex.DecodeString(tx.Signature)
	if err != nil || len(sigBytes) != 64 {
		return false
	}

	r := new(big.Int).SetBytes(sigBytes[:32])
	s := new(big.Int).SetBytes(sigBytes[32:])
	hash := sha256.Sum256([]byte(tx.Hash()))

	return ecdsa.Verify(pubKey, hash[:], r, s)
}
