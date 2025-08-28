// internal/consensus/dpos.go
package consensus

import (
	"math/rand"
	"time"
)

type Validator struct {
	Address string
	Stake   uint64
}

type DPoS struct {
	Validators []Validator
}

func NewDPoS(validators []Validator) *DPoS {
	return &DPoS{Validators: validators}
}

func (d *DPoS) SelectValidator() Validator {
	totalStake := uint64(0)
	for _, v := range d.Validators {
		totalStake += v.Stake
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano())).Uint64() % totalStake
	cum := uint64(0)
	for _, v := range d.Validators {
		cum += v.Stake
		if r < cum {
			return v
		}
	}
	return d.Validators[0]
}
