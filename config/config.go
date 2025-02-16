package config

import (
	"path/filepath"
)

type Config struct {
	DataDir string
	Port    string
}

func NewConfig(dataDir, port string) *Config {
	return &Config{
		DataDir: dataDir,
		Port:    port,
	}
}

func (c *Config) GetBlockchainPath() string {
	return filepath.Join(c.DataDir, "blockchain.db")
}

func (c *Config) GetWalletsPath() string {
	return filepath.Join(c.DataDir, "wallets")
}
