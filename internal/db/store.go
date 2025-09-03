package db

import (
	"database/sql"
	"fmt"
	"sphere/internal/core"
	"time"
)

// SaveBlock persists a block and its transactions. Returns the inserted block row id.
func SaveBlock(conn *sql.DB, b *core.Block) (int, error) {
	tx, err := conn.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		// if commit not done, rollback
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	// Insert block
	var blockID int
	err = tx.QueryRow(`
		INSERT INTO blocks (height, hash, prev_hash, timestamp, validator, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, b.Height, b.Hash, b.PrevHash, b.Timestamp, b.Validator, time.Now()).Scan(&blockID)
	if err != nil {
		return 0, err
	}

	// Insert transactions
	stmt, err := tx.Prepare(`
		INSERT INTO transactions (block_id, "from", "to", amount, fee, data, sig, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	for _, txv := range b.Transactions {
		_, err := stmt.Exec(blockID, txv.From, txv.To, fmt.Sprintf("%f", txv.Amount), fmt.Sprintf("%f", txv.Fee), txv.Data, txv.Sig, time.Now())
		if err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	// set tx to nil to avoid deferred rollback
	tx = nil
	return blockID, nil
}
