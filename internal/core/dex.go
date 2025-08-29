package core

import "sync"

// Pool represents a simple liquidity pool
type Pool struct {
	ID       int     `json:"id"`
	TokenA   string  `json:"token_a"`
	TokenB   string  `json:"token_b"`
	ReserveA float64 `json:"reserve_a"`
	ReserveB float64 `json:"reserve_b"`
}

type DEX struct {
	mu     sync.Mutex
	pools  []Pool
	nextID int
}

func NewDEX() *DEX {
	return &DEX{pools: []Pool{}, nextID: 1}
}

// AddPool adds a new pool
func (d *DEX) AddPool(tokenA, tokenB string, reserveA, reserveB float64) Pool {
	d.mu.Lock()
	defer d.mu.Unlock()
	p := Pool{
		ID:       d.nextID,
		TokenA:   tokenA,
		TokenB:   tokenB,
		ReserveA: reserveA,
		ReserveB: reserveB,
	}
	d.pools = append(d.pools, p)
	d.nextID++
	return p
}

// ListPools returns all pools
func (d *DEX) ListPools() []Pool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return append([]Pool(nil), d.pools...) // return copy
}
