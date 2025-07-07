package amm

import "errors"

type Pool struct {
	TokenA   string
	TokenB   string
	ReserveA uint64
	ReserveB uint64
}

// Swap using constant product formula: x * y = k
func GetAmountOut(amountIn, reserveIn, reserveOut uint64) (uint64, error) {
	if amountIn == 0 || reserveIn == 0 || reserveOut == 0 {
		return 0, errors.New("invalid reserves or input")
	}
	amountInWithFee := amountIn * 997 // 0.3% fee
	numerator := amountInWithFee * reserveOut
	denominator := reserveIn*1000 + amountInWithFee
	return numerator / denominator, nil
}
