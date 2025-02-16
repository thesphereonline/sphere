package storage

import (
	"encoding/json"
	"os"
	"sync"
	"your-username/blockchain/block"
)

type Storage struct {
	DbFile string
	mu     sync.RWMutex
}

func NewStorage(dbFile string) *Storage {
	return &Storage{
		DbFile: dbFile,
	}
}

func (s *Storage) SaveChain(blocks []*block.Block) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.Marshal(blocks)
	if err != nil {
		return err
	}

	return os.WriteFile(s.DbFile, data, 0644)
}

func (s *Storage) LoadChain() ([]*block.Block, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var blocks []*block.Block

	data, err := os.ReadFile(s.DbFile)
	if err != nil {
		if os.IsNotExist(err) {
			return blocks, nil
		}
		return nil, err
	}

	err = json.Unmarshal(data, &blocks)
	return blocks, err
}

func (s *Storage) SaveWallet(address, walletData string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return os.WriteFile("wallet_"+address+".dat", []byte(walletData), 0644)
}

func (s *Storage) LoadWallet(address string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile("wallet_" + address + ".dat")
	if err != nil {
		return "", err
	}

	return string(data), nil
}
