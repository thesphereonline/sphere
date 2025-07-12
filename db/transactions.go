package db

import (
	"log"

	"github.com/thesphereonline/sphere/core/types"
)

func InsertTransaction(tx types.Transaction, blockIndex uint64) error {
	_, err := DB.Exec(`
		INSERT INTO transactions (hash, block_index, from_addr, to_addr, nonce, gas_limit, data, signature)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, tx.Hash, blockIndex, tx.From, tx.To, tx.Nonce, tx.GasLimit, tx.Data, tx.Signature)

	if err != nil {
		log.Printf("InsertTransaction error: %v", err)
	}
	return err
}
