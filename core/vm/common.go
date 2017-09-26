package vm

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

// calculates the memory size required for a step
func calcMemSize(off, l *big.Int) *big.Int {
	fuzz_helper.AddCoverage(22588)
	if l.Sign() == 0 {
		fuzz_helper.AddCoverage(5262)
		return common.Big0
	} else {
		fuzz_helper.AddCoverage(17878)
	}
	fuzz_helper.AddCoverage(44810)

	return new(big.Int).Add(off, l)
}

// getData returns a slice from the data based on the start and size and pads
// up to size with zero's. This function is overflow safe.
func getData(data []byte, start uint64, size uint64) []byte {
	fuzz_helper.AddCoverage(45021)
	length := uint64(len(data))
	if start > length {
		fuzz_helper.AddCoverage(21668)
		start = length
	} else {
		fuzz_helper.AddCoverage(45213)
	}
	fuzz_helper.AddCoverage(39040)
	end := start + size
	if end > length {
		fuzz_helper.AddCoverage(16619)
		end = length
	} else {
		fuzz_helper.AddCoverage(12692)
	}
	fuzz_helper.AddCoverage(2095)
	return common.RightPadBytes(data[start:end], int(size))
}

// getDataBig returns a slice from the data based on the start and size and pads
// up to size with zero's. This function is overflow safe.
func getDataBig(data []byte, start *big.Int, size *big.Int) []byte {
	fuzz_helper.AddCoverage(42483)
	dlen := big.NewInt(int64(len(data)))

	s := math.BigMin(start, dlen)
	e := math.BigMin(new(big.Int).Add(s, size), dlen)
	return common.RightPadBytes(data[s.Uint64():e.Uint64()], int(size.Uint64()))
}

// bigUint64 returns the integer casted to a uint64 and returns whether it
// overflowed in the process.
func bigUint64(v *big.Int) (uint64, bool) {
	fuzz_helper.AddCoverage(6577)
	return v.Uint64(), v.BitLen() > 64
}

// toWordSize returns the ceiled word size required for memory expansion.
func toWordSize(size uint64) uint64 {
	fuzz_helper.AddCoverage(17393)
	if size > math.MaxUint64-31 {
		fuzz_helper.AddCoverage(38740)
		return math.MaxUint64/32 + 1
	} else {
		fuzz_helper.AddCoverage(35657)
	}
	fuzz_helper.AddCoverage(64174)

	return (size + 31) / 32
}

func allZero(b []byte) bool {
	fuzz_helper.AddCoverage(30358)
	for _, byte := range b {
		fuzz_helper.AddCoverage(61639)
		if byte != 0 {
			fuzz_helper.AddCoverage(11162)
			return false
		} else {
			fuzz_helper.AddCoverage(49217)
		}
	}
	fuzz_helper.AddCoverage(23294)
	return true
}

var _ = fuzz_helper.AddCoverage
