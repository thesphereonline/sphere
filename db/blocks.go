package db

import (
	"log"

	"github.com/thesphereonline/sphere/core/types"
)

func InsertBlock(block types.Block) error {
	_, err := DB.Exec(`
		INSERT INTO blocks (index, timestamp, prev_hash, hash, validator, state_root)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, block.Index, block.Timestamp, block.PrevHash, block.Hash, block.Validator, block.StateRoot)

	if err != nil {
		log.Printf("InsertBlock error: %v", err)
	}
	return err
}

func GetBlockByIndex(index uint64) (types.Block, error) {
	var b types.Block
	row := DB.QueryRow(`SELECT index, timestamp, prev_hash, hash, validator, state_root FROM blocks WHERE index = $1`, index)
	err := row.Scan(&b.Index, &b.Timestamp, &b.PrevHash, &b.Hash, &b.Validator, &b.StateRoot)
	if err != nil {
		return types.Block{}, err
	}
	return b, nil
}

func ListBlocks() ([]types.Block, error) {
	rows, err := DB.Query(`SELECT index, timestamp, prev_hash, hash, validator, state_root FROM blocks ORDER BY index DESC LIMIT 50`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []types.Block
	for rows.Next() {
		var b types.Block
		if err := rows.Scan(&b.Index, &b.Timestamp, &b.PrevHash, &b.Hash, &b.Validator, &b.StateRoot); err != nil {
			return nil, err
		}
		blocks = append(blocks, b)
	}
	return blocks, nil
}

func GetTransactionByHash(hash string) (types.Transaction, error) {
	var tx types.Transaction
	row := DB.QueryRow(`SELECT hash, block_index, from_addr, to_addr, nonce, gas_limit, data, signature FROM transactions WHERE hash = $1`, hash)
	err := row.Scan(&tx.Hash, &tx.BlockIndex, &tx.From, &tx.To, &tx.Nonce, &tx.GasLimit, &tx.Data, &tx.Signature)
	if err != nil {
		return types.Transaction{}, err
	}
	return tx, nil
}

func GetUserPortfolio(address string) (types.Portfolio, error) {
	// Define Portfolio struct as needed.
	// Example: get balances + NFTs + LP shares
	// Implement query logic here.

	// For now, return empty portfolio or error
	return types.Portfolio{}, nil
}
