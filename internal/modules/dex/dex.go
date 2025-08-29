package dex

import (
	"database/sql"
	"errors"
	"fmt"
)

// Pool represents a liquidity pool
type Pool struct {
	ID       int     `json:"id"`
	TokenA   string  `json:"token_a"`
	TokenB   string  `json:"token_b"`
	ReserveA float64 `json:"reserve_a"`
	ReserveB float64 `json:"reserve_b"`
	FeeBps   int     `json:"fee_bps"`
}

type Module struct {
	db  *sql.DB
	fee int // protocol fee in basis points (e.g. 5 = 0.05%)
}

// New creates a new DEX module instance
func New(db *sql.DB, feeBps int) *Module {
	m := &Module{db: db, fee: feeBps}
	m.ensureSchema()
	return m
}

// ensureSchema creates the pools table if not exists
func (m *Module) ensureSchema() {
	query := `
	CREATE TABLE IF NOT EXISTS pools (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		token_a TEXT NOT NULL,
		token_b TEXT NOT NULL,
		reserve_a REAL NOT NULL,
		reserve_b REAL NOT NULL,
		fee_bps INTEGER NOT NULL
	);`
	_, _ = m.db.Exec(query)
}

// ListPools returns all pools
func (m *Module) ListPools() ([]Pool, error) {
	rows, err := m.db.Query("SELECT id, token_a, token_b, reserve_a, reserve_b, fee_bps FROM pools")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pools []Pool
	for rows.Next() {
		var p Pool
		if err := rows.Scan(&p.ID, &p.TokenA, &p.TokenB, &p.ReserveA, &p.ReserveB, &p.FeeBps); err != nil {
			return nil, err
		}
		pools = append(pools, p)
	}
	return pools, nil
}

// AddPool inserts a new pool
func (m *Module) AddPool(tokenA, tokenB string, reserveA, reserveB float64) (*Pool, error) {
	query := `INSERT INTO pools (token_a, token_b, reserve_a, reserve_b, fee_bps) VALUES (?, ?, ?, ?, ?)`
	res, err := m.db.Exec(query, tokenA, tokenB, reserveA, reserveB, m.fee)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &Pool{
		ID:       int(id),
		TokenA:   tokenA,
		TokenB:   tokenB,
		ReserveA: reserveA,
		ReserveB: reserveB,
		FeeBps:   m.fee,
	}, nil
}

// Swap performs a token swap (constant product AMM: x*y=k)
func (m *Module) Swap(poolID int, fromToken string, amountIn, minOut float64, trader string) (float64, error) {
	// Load pool
	var pool Pool
	row := m.db.QueryRow("SELECT id, token_a, token_b, reserve_a, reserve_b, fee_bps FROM pools WHERE id=?", poolID)
	if err := row.Scan(&pool.ID, &pool.TokenA, &pool.TokenB, &pool.ReserveA, &pool.ReserveB, &pool.FeeBps); err != nil {
		return 0, errors.New("pool not found")
	}

	// Apply fee
	amountInWithFee := amountIn * (1.0 - float64(pool.FeeBps)/10000.0)

	var amountOut float64
	if fromToken == pool.TokenA {
		// Swap A → B
		amountOut = (pool.ReserveB * amountInWithFee) / (pool.ReserveA + amountInWithFee)
		if amountOut < minOut {
			return 0, errors.New("slippage: insufficient output")
		}
		pool.ReserveA += amountIn
		pool.ReserveB -= amountOut
	} else if fromToken == pool.TokenB {
		// Swap B → A
		amountOut = (pool.ReserveA * amountInWithFee) / (pool.ReserveB + amountInWithFee)
		if amountOut < minOut {
			return 0, errors.New("slippage: insufficient output")
		}
		pool.ReserveB += amountIn
		pool.ReserveA -= amountOut
	} else {
		return 0, fmt.Errorf("invalid token: %s", fromToken)
	}

	// Update DB
	_, err := m.db.Exec(
		"UPDATE pools SET reserve_a=?, reserve_b=? WHERE id=?",
		pool.ReserveA, pool.ReserveB, pool.ID,
	)
	if err != nil {
		return 0, err
	}

	return amountOut, nil
}
