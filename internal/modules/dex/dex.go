package dex

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strconv"
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

func (m *Module) AddLiquidity(ctx context.Context, id int, owner string, a string, b string) (map[string]any, error) {
	// parse amounts (frontend may send strings)
	amtA, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid amountA: %w", err)
	}
	amtB, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid amountB: %w", err)
	}
	if amtA <= 0 || amtB <= 0 {
		return nil, errors.New("amounts must be > 0")
	}

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Lock pool row
	var pool Pool
	var totalLPStr string
	row := tx.QueryRowContext(ctx, "SELECT id, token_a, token_b, reserve_a, reserve_b, fee_bps, total_lp FROM pools WHERE id=$1 FOR UPDATE", id)
	var reserveAStr, reserveBStr string
	if err := row.Scan(&pool.ID, &pool.TokenA, &pool.TokenB, &reserveAStr, &reserveBStr, &pool.FeeBps, &totalLPStr); err != nil {
		tx.Rollback()
		return nil, errors.New("pool not found")
	}

	// parse reserves & total_lp (migrations use TEXT to allow big numeric strings)
	pool.ReserveA, _ = strconv.ParseFloat(reserveAStr, 64)
	pool.ReserveB, _ = strconv.ParseFloat(reserveBStr, 64)
	var totalLP float64
	if totalLPStr != "" {
		totalLP, _ = strconv.ParseFloat(totalLPStr, 64)
	}

	var mintedLP float64
	if totalLP == 0 {
		// initial provider gets sqrt(a*b)
		mintedLP = math.Sqrt(amtA * amtB)
		totalLP = mintedLP
	} else {
		// must preserve ratio; compute LP based on smallest proportional contribution
		lpFromA := (amtA * totalLP) / pool.ReserveA
		lpFromB := (amtB * totalLP) / pool.ReserveB
		mintedLP = math.Min(lpFromA, lpFromB)
		totalLP += mintedLP
	}

	// update reserves
	pool.ReserveA += amtA
	pool.ReserveB += amtB

	// write pool updates
	_, err = tx.ExecContext(ctx, "UPDATE pools SET reserve_a=$1, reserve_b=$2, total_lp=$3 WHERE id=$4",
		fmt.Sprintf("%f", pool.ReserveA), fmt.Sprintf("%f", pool.ReserveB), fmt.Sprintf("%f", totalLP), pool.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// upsert lp_positions
	// try update first
	res, err := tx.ExecContext(ctx, "UPDATE lp_positions SET lp_amount = (lp_amount::numeric + $1::numeric) WHERE pool_id=$2 AND owner=$3", fmt.Sprintf("%f", mintedLP), pool.ID, owner)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	rowsAff, _ := res.RowsAffected()
	if rowsAff == 0 {
		// insert
		_, err = tx.ExecContext(ctx, "INSERT INTO lp_positions (pool_id, owner, lp_amount) VALUES ($1, $2, $3)", pool.ID, owner, fmt.Sprintf("%f", mintedLP))
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return map[string]any{
		"pool_id":   pool.ID,
		"minted_lp": mintedLP,
		"total_lp":  totalLP,
		"reserve_a": pool.ReserveA,
		"reserve_b": pool.ReserveB,
		"owner":     owner,
	}, nil
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
	query := `INSERT INTO pools (token_a, token_b, reserve_a, reserve_b, fee_bps, total_lp) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	// initial total_lp = sqrt(reserveA*reserveB)
	initialLP := math.Sqrt(reserveA * reserveB)
	var id int
	if err := m.db.QueryRow(query, tokenA, tokenB, fmt.Sprintf("%f", reserveA), fmt.Sprintf("%f", reserveB), m.fee, fmt.Sprintf("%f", initialLP)).Scan(&id); err != nil {
		return nil, err
	}
	return &Pool{ID: id, TokenA: tokenA, TokenB: tokenB, ReserveA: reserveA, ReserveB: reserveB, FeeBps: m.fee}, nil
}

// Swap performs a token swap (constant product AMM: x*y=k)
func (m *Module) Swap(poolID int, fromToken string, amountIn, minOut float64, trader string) (float64, error) {
	// Load pool
	var pool Pool
	row := m.db.QueryRow("SELECT id, token_a, token_b, reserve_a, reserve_b, fee_bps FROM pools WHERE id=$1", poolID)
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
		"UPDATE pools SET reserve_a=$1, reserve_b=$2 WHERE id=$3",
		fmt.Sprintf("%f", pool.ReserveA), fmt.Sprintf("%f", pool.ReserveB), pool.ID,
	)
	if err != nil {
		return 0, err
	}

	return amountOut, nil
}
