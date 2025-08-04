package entities

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
	Address    string
}

// NewWallet creates a new wallet with fresh keys
func NewWallet() (*Wallet, error) {
	privKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	pubKey := append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...)
	address := generateAddress(pubKey)
	return &Wallet{PrivateKey: privKey, PublicKey: pubKey, Address: address}, nil
}

// generateAddress creates an address by hashing the public key
func generateAddress(pubKey []byte) string {
	hash := sha256.Sum256(pubKey)
	return hex.EncodeToString(hash[:20]) // Take first 20 bytes as address (like Ethereum)
}

// Sign signs data with the wallet private key
func (w *Wallet) Sign(data []byte) (r, s *big.Int, err error) {
	hash := sha256.Sum256(data)
	return ecdsa.Sign(rand.Reader, w.PrivateKey, hash[:])
}

// Verify verifies a signature for data with the public key
func (w *Wallet) Verify(data []byte, r, s *big.Int) bool {
	hash := sha256.Sum256(data)
	return ecdsa.Verify(&w.PrivateKey.PublicKey, hash[:], r, s)
}

// PublicKeyHex returns the public key as a hex string
func (w *Wallet) PublicKeyHex() string {
	return hex.EncodeToString(w.PublicKey)
}
