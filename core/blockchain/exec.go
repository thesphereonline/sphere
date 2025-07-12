package blockchain

import (
	"github.com/thesphereonline/sphere/core/types"
	"github.com/thesphereonline/sphere/core/vm"
)

type Executor struct {
	State *types.State
}

func NewExecutor() *Executor {
	return &Executor{
		State: &types.State{
			Balances: make(map[string]uint64),
			Tokens:   make(map[string]map[string]uint64),
			NFTs:     make(map[string]string),
			Metadata: make(map[string]string),
		},
	}
}

func (e *Executor) ApplyTx(tx types.Transaction) error {
	interpreter := vm.NewVM(tx.From, tx.GasLimit, e.State)
	return interpreter.Execute(tx.Data)
}
