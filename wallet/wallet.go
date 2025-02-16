package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Address    string
}

func NewWallet() *Wallet {
	private, public := generateKeyPair()
	wallet := Wallet{
		PrivateKey: private,
		PublicKey:  public,
		Address:    generateAddress(public),
	}
	return &wallet
}

func generateKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	return private, &private.PublicKey
}

func generateAddress(pubKey *ecdsa.PublicKey) string {
	pubKeyBytes := elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)
	pubKeyHash := sha256.Sum256(pubKeyBytes)
	return hex.EncodeToString(pubKeyHash[:])
}

func (w *Wallet) Sign(data []byte) []byte {
	signature, err := ecdsa.SignASN1(rand.Reader, w.PrivateKey, data)
	if err != nil {
		log.Panic(err)
	}
	return signature
}

func LoadWalletFromFile(address string) *Wallet {
	data, err := os.ReadFile(fmt.Sprintf("wallet_%s.dat", address))
	if err != nil {
		log.Printf("Failed to load wallet: %v", err)
		return NewWallet()
	}

	var wallet Wallet
	if err := json.Unmarshal(data, &wallet); err != nil {
		log.Printf("Failed to unmarshal wallet: %v", err)
		return NewWallet()
	}

	return &wallet
}

func (w *Wallet) SaveToFile() error {
	data, err := json.Marshal(w)
	if err != nil {
		return err
	}
	return os.WriteFile(fmt.Sprintf("wallet_%s.dat", w.Address), data, 0644)
}
