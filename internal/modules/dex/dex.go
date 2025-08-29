package dex

import (
	"context"
	"database/sql"
	"errors"
	"math/big"
)

type Dex struct {
	db *sql.DB
	// Optional: protocol fee basis points
	ProtocolFeeBps int64
}

func New(db *sql.DB, protocolFeeBps int64) *Dex {
	return &Dex{db: db, ProtocolFeeBps: protocolFeeBps}
}

type Pool struct {
	ID       int
	TokenA   string
	TokenB   string
	ReserveA *big.Int
	ReserveB *big.Int
	TotalLP  *big.Int
	FeeBps   int64
}

func parseBig(s string) *big.Int {
	i := new(big.Int)
	i.SetString(s, 10)
	return i
}

func (d *Dex) GetPool(ctx context.Context, poolID int) (*Pool, error) {
	var resA, resB, totalLP string
	var tokenA, tokenB string
	var feeBps int64
	row := d.db.QueryRowContext(ctx, `SELECT token_a, token_b, reserve_a, reserve_b, total_lp, fee_bps FROM pools WHERE id=$1`, poolID)
	if err := row.Scan(&tokenA, &tokenB, &resA, &resB, &totalLP, &feeBps); err != nil {
		return nil, err
	}
	return &Pool{
		ID:       poolID,
		TokenA:   tokenA,
		TokenB:   tokenB,
		ReserveA: parseBig(resA),
		ReserveB: parseBig(resB),
		TotalLP:  parseBig(totalLP),
		FeeBps:   feeBps,
	}, nil
}

// AddLiquidity: owner must have transferred tokens in design. For demo we skip balance transfers.
func (d *Dex) AddLiquidity(ctx context.Context, poolID int, owner string, amountAStr, amountBStr string) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// lock pool row
	var resA, resB, totalLP string
	if err := tx.QueryRowContext(ctx, `SELECT reserve_a, reserve_b, total_lp FROM pools WHERE id=$1 FOR UPDATE`, poolID).
		Scan(&resA, &resB, &totalLP); err != nil {
		return err
	}

	reserveA := parseBig(resA)
	reserveB := parseBig(resB)
	total := parseBig(totalLP)
	amountA := parseBig(amountAStr)
	amountB := parseBig(amountBStr)

	var lpMint *big.Int
	zero := big.NewInt(0)
	if total.Cmp(zero) == 0 {
		// initial liquidity -> mint sqrt(amountA * amountB)
		prod := new(big.Int).Mul(amountA, amountB)
		lpMint = integerSqrt(prod)
	} else {
		// lpMint = min(amountA * total / reserveA, amountB * total / reserveB)
		lpA := new(big.Int).Div(new(big.Int).Mul(amountA, total), reserveA)
		lpB := new(big.Int).Div(new(big.Int).Mul(amountB, total), reserveB)
		if lpA.Cmp(lpB) < 0 {
			lpMint = lpA
		} else {
			lpMint = lpB
		}
	}

	reserveA.Add(reserveA, amountA)
	reserveB.Add(reserveB, amountB)
	total.Add(total, lpMint)

	if _, err := tx.ExecContext(ctx, `UPDATE pools SET reserve_a=$1, reserve_b=$2, total_lp=$3 WHERE id=$4`,
		reserveA.String(), reserveB.String(), total.String(), poolID); err != nil {
		return err
	}

	// upsert lp position
	_, err = tx.ExecContext(ctx, `
      INSERT INTO lp_positions (pool_id, owner, lp_amount) VALUES ($1,$2,$3)
      ON CONFLICT (pool_id, owner) DO UPDATE SET lp_amount = (lp_positions.lp_amount::numeric + $3::numeric)::text`, poolID, owner, lpMint.String())
	if err != nil {
		return err
	}

	return nil
}

func (d *Dex) Swap(ctx context.Context, poolID int, fromToken string, amountInStr string, minOutStr string, trader string) (string, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var tokenA, tokenB, resAStr, resBStr string
	var feeBps int64
	if err := tx.QueryRowContext(ctx, `SELECT token_a, token_b, reserve_a, reserve_b, fee_bps FROM pools WHERE id=$1 FOR UPDATE`, poolID).
		Scan(&tokenA, &tokenB, &resAStr, &resBStr, &feeBps); err != nil {
		return "", err
	}

	reserveA := parseBig(resAStr)
	reserveB := parseBig(resBStr)
	amountIn := parseBig(amountInStr)
	minOut := parseBig(minOutStr)

	var reserveIn, reserveOut *big.Int
	if fromToken == tokenA {
		reserveIn = reserveA
		reserveOut = reserveB
	} else if fromToken == tokenB {
		reserveIn = reserveB
		reserveOut = reserveA
	} else {
		return "", errors.New("invalid token")
	}

	// amountInWithFee = amountIn * (10000 - feeBps) / 10000
	numerator := new(big.Int).Mul(amountIn, big.NewInt(int64(10000-feeBps)))
	denom := new(big.Int).Add(new(big.Int).Mul(reserveIn, big.NewInt(10000)), numerator)
	amountOut := new(big.Int).Div(new(big.Int).Mul(numerator, reserveOut), denom)

	if amountOut.Cmp(minOut) < 0 {
		return "", errors.New("insufficient output amount")
	}

	// Update reserves
	reserveIn.Add(reserveIn, amountIn)
	reserveOut.Sub(reserveOut, amountOut)

	if fromToken == tokenA {
		if _, err := tx.ExecContext(ctx, `UPDATE pools SET reserve_a=$1, reserve_b=$2 WHERE id=$3`, reserveIn.String(), reserveOut.String(), poolID); err != nil {
			return "", err
		}
	} else {
		if _, err := tx.ExecContext(ctx, `UPDATE pools SET reserve_a=$1, reserve_b=$2 WHERE id=$3`, reserveOut.String(), reserveIn.String(), poolID); err != nil {
			return "", err
		}
	}

	// TODO: handle balances transfer, protocol fee accounting
	return amountOut.String(), nil
}

func integerSqrt(a *big.Int) *big.Int {
	if a.Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(0)
	}
	z := new(big.Int).Rsh(a, uint(a.BitLen()/2))
	big.NewInt(1)
	for {
		y := new(big.Int).Div(new(big.Int).Add(new(big.Int).Div(a, z), z), big.NewInt(2))
		if y.Cmp(z) >= 0 {
			return z
		}
		z = y
	}
}
