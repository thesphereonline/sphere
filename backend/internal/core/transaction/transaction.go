package transaction

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"math/big"
)

type Transaction struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Amount    uint64 `json:"amount"`
	Nonce     uint64 `json:"nonce"`
	Signature string `json:"signature"`
	Timestamp int64  `json:"timestamp"`
}

func (tx *Transaction) Hash() string {
	txCopy := *tx
	txCopy.Signature = ""
	bytes, _ := json.Marshal(txCopy)
	hash := sha256.Sum256(bytes)
	return hex.EncodeToString(hash[:])
}

func (tx *Transaction) VerifySignature(pubKey *ecdsa.PublicKey) bool {
	hash := sha256.Sum256([]byte(tx.From + tx.To + string(tx.Amount) + string(tx.Nonce)))
	r := new(big.Int)
	s := new(big.Int)
	sigBytes, _ := hex.DecodeString(tx.Signature)
	r.SetBytes(sigBytes[:len(sigBytes)/2])
	s.SetBytes(sigBytes[len(sigBytes)/2:])
	return ecdsa.Verify(pubKey, hash[:], r, s)
}
