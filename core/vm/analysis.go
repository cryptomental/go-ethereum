package vm

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// destinations stores one map per contract (keyed by hash of code).
// The maps contain an entry for each location of a JUMPDEST
// instruction.
type destinations map[common.Hash][]byte

// has checks whether code has a JUMPDEST at dest.
func (d destinations) has(codehash common.Hash, code []byte, dest *big.Int) bool {
	fuzz_helper.AddCoverage(22588)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()

	udest := dest.Uint64()
	if dest.BitLen() >= 63 || udest >= uint64(len(code)) {
		fuzz_helper.AddCoverage(17878)
		return false
	} else {
		fuzz_helper.AddCoverage(45021)
	}
	fuzz_helper.AddCoverage(44810)

	m, analysed := d[codehash]
	if !analysed {
		fuzz_helper.AddCoverage(39040)
		m = jumpdests(code)
		d[codehash] = m
	} else {
		fuzz_helper.AddCoverage(2095)
	}
	fuzz_helper.AddCoverage(5262)
	return (m[udest/8] & (1 << (udest % 8))) != 0
}

// jumpdests creates a map that contains an entry for each
// PC location that is a JUMPDEST instruction.
func jumpdests(code []byte) []byte {
	fuzz_helper.AddCoverage(21668)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	m := make([]byte, len(code)/8+1)
	for pc := uint64(0); pc < uint64(len(code)); pc++ {
		fuzz_helper.AddCoverage(16619)
		op := OpCode(code[pc])
		if op == JUMPDEST {
			fuzz_helper.AddCoverage(12692)
			m[pc/8] |= 1 << (pc % 8)
		} else {
			fuzz_helper.AddCoverage(42483)
			if op >= PUSH1 && op <= PUSH32 {
				fuzz_helper.AddCoverage(6577)
				a := uint64(op) - uint64(PUSH1) + 1
				pc += a
			} else {
				fuzz_helper.AddCoverage(17393)
			}
		}
	}
	fuzz_helper.AddCoverage(45213)
	return m
}

var _ = fuzz_helper.AddCoverage
