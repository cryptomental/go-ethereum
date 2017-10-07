package vm

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

// calculates the memory size required for a step
func calcMemSize(off, l *big.Int) *big.Int {
	fuzz_helper.AddCoverage(64174)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	if l.Sign() == 0 {
		fuzz_helper.AddCoverage(35657)
		return common.Big0
	} else {
		fuzz_helper.AddCoverage(30358)
	}
	fuzz_helper.AddCoverage(38740)

	return new(big.Int).Add(off, l)
}

// getData returns a slice from the data based on the start and size and pads
// up to size with zero's. This function is overflow safe.
func getData(data []byte, start uint64, size uint64) []byte {
	fuzz_helper.AddCoverage(23294)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	length := uint64(len(data))
	if start > length {
		fuzz_helper.AddCoverage(49217)
		start = length
	} else {
		fuzz_helper.AddCoverage(34511)
	}
	fuzz_helper.AddCoverage(61639)
	end := start + size
	if end > length {
		fuzz_helper.AddCoverage(64074)
		end = length
	} else {
		fuzz_helper.AddCoverage(28614)
	}
	fuzz_helper.AddCoverage(11162)
	return common.RightPadBytes(data[start:end], int(size))
}

// getDataBig returns a slice from the data based on the start and size and pads
// up to size with zero's. This function is overflow safe.
func getDataBig(data []byte, start *big.Int, size *big.Int) []byte {
	fuzz_helper.AddCoverage(39226)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	dlen := big.NewInt(int64(len(data)))

	s := math.BigMin(start, dlen)
	e := math.BigMin(new(big.Int).Add(s, size), dlen)
	return common.RightPadBytes(data[s.Uint64():e.Uint64()], int(size.Uint64()))
}

// bigUint64 returns the integer casted to a uint64 and returns whether it
// overflowed in the process.
func bigUint64(v *big.Int) (uint64, bool) {
	fuzz_helper.AddCoverage(2297)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return v.Uint64(), v.BitLen() > 64
}

// toWordSize returns the ceiled word size required for memory expansion.
func toWordSize(size uint64) uint64 {
	fuzz_helper.AddCoverage(40870)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	if size > math.MaxUint64-31 {
		fuzz_helper.AddCoverage(778)
		return math.MaxUint64/32 + 1
	} else {
		fuzz_helper.AddCoverage(33340)
	}
	fuzz_helper.AddCoverage(52877)

	return (size + 31) / 32
}

func allZero(b []byte) bool {
	fuzz_helper.AddCoverage(15638)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	for _, byte := range b {
		fuzz_helper.AddCoverage(23368)
		if byte != 0 {
			fuzz_helper.AddCoverage(12901)
			return false
		} else {
			fuzz_helper.AddCoverage(12499)
		}
	}
	fuzz_helper.AddCoverage(45869)
	return true
}

var _ = fuzz_helper.AddCoverage
