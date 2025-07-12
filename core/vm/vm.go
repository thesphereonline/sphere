package vm

import (
	"encoding/binary"
	"fmt"

	"github.com/thesphereonline/sphere/core/types"
)

type VM struct {
	Caller string
	Gas    uint64
	State  *types.State
}

func NewVM(caller string, gas uint64, state *types.State) *VM {
	return &VM{
		Caller: caller,
		Gas:    gas,
		State:  state,
	}
}

func (vm *VM) Execute(bytecode []byte) error {
	pc := 0
	for pc < len(bytecode) {
		op := bytecode[pc]
		pc++

		switch op {
		case OP_NOP:
			continue

		case OP_MINT:
			addr := readAddr(bytecode[pc : pc+20])
			pc += 20
			amt := binary.BigEndian.Uint64(bytecode[pc : pc+8])
			pc += 8
			vm.State.Balances[addr] += amt

		case OP_TRANSFER:
			to := readAddr(bytecode[pc : pc+20])
			pc += 20
			amt := binary.BigEndian.Uint64(bytecode[pc : pc+8])
			pc += 8

			if vm.State.Balances[vm.Caller] < amt {
				return fmt.Errorf("insufficient balance")
			}
			vm.State.Balances[vm.Caller] -= amt
			vm.State.Balances[to] += amt

		case OP_MINT_NFT:
			tokenId := readString(bytecode[pc : pc+32])
			pc += 32
			vm.State.NFTs[tokenId] = vm.Caller

		case OP_STORE_META:
			tokenId := readString(bytecode[pc : pc+32])
			pc += 32
			uriLen := int(bytecode[pc])
			pc++
			uri := string(bytecode[pc : pc+uriLen])
			pc += uriLen
			vm.State.Metadata[tokenId] = uri

		default:
			return fmt.Errorf("invalid opcode: %x", op)
		}
	}
	return nil
}

// helpers

func readAddr(b []byte) string {
	return fmt.Sprintf("0x%x", b)
}

func readString(b []byte) string {
	return string(b)
}
