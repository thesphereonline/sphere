package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
)

type Wallet struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
	Address    string // derived from pubkey
}

// NewWallet generates a new wallet with ed25519 keys
func NewWallet() *Wallet {
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	return &Wallet{
		PrivateKey: priv,
		PublicKey:  pub,
		Address:    hex.EncodeToString(pub[:]),
	}
}

// Sign signs a message using the wallet's private key
func (w *Wallet) Sign(msg []byte) string {
	sig := ed25519.Sign(w.PrivateKey, msg)
	return hex.EncodeToString(sig)
}

// VerifySignature validates a signature
func VerifySignature(pubHex, msgHex, sigHex string) bool {
	pubBytes, err := hex.DecodeString(pubHex)
	if err != nil {
		return false
	}
	msgBytes, err := hex.DecodeString(msgHex)
	if err != nil {
		return false
	}
	sigBytes, err := hex.DecodeString(sigHex)
	if err != nil {
		return false
	}
	return ed25519.Verify(pubBytes, msgBytes, sigBytes)
}
