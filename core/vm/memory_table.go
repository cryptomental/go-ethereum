package vm

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/math"
)

func memorySha3(stack *Stack) *big.Int {
	fuzz_helper.AddCoverage(47070)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return calcMemSize(stack.Back(0), stack.Back(1))
}

func memoryCallDataCopy(stack *Stack) *big.Int {
	fuzz_helper.AddCoverage(681)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return calcMemSize(stack.Back(0), stack.Back(2))
}

func memoryReturnDataCopy(stack *Stack) *big.Int {
	fuzz_helper.AddCoverage(50658)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return calcMemSize(stack.Back(0), stack.Back(2))
}

func memoryCodeCopy(stack *Stack) *big.Int {
	fuzz_helper.AddCoverage(40310)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return calcMemSize(stack.Back(0), stack.Back(2))
}

func memoryExtCodeCopy(stack *Stack) *big.Int {
	fuzz_helper.AddCoverage(28441)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return calcMemSize(stack.Back(1), stack.Back(3))
}

func memoryMLoad(stack *Stack) *big.Int {
	fuzz_helper.AddCoverage(60192)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return calcMemSize(stack.Back(0), big.NewInt(32))
}

func memoryMStore8(stack *Stack) *big.Int {
	fuzz_helper.AddCoverage(36930)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return calcMemSize(stack.Back(0), big.NewInt(1))
}

func memoryMStore(stack *Stack) *big.Int {
	fuzz_helper.AddCoverage(48894)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return calcMemSize(stack.Back(0), big.NewInt(32))
}

func memoryCreate(stack *Stack) *big.Int {
	fuzz_helper.AddCoverage(48452)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return calcMemSize(stack.Back(1), stack.Back(2))
}

func memoryCall(stack *Stack) *big.Int {
	fuzz_helper.AddCoverage(24027)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := calcMemSize(stack.Back(5), stack.Back(6))
	y := calcMemSize(stack.Back(3), stack.Back(4))

	return math.BigMax(x, y)
}

func memoryCallCode(stack *Stack) *big.Int {
	fuzz_helper.AddCoverage(34031)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := calcMemSize(stack.Back(5), stack.Back(6))
	y := calcMemSize(stack.Back(3), stack.Back(4))

	return math.BigMax(x, y)
}
func memoryDelegateCall(stack *Stack) *big.Int {
	fuzz_helper.AddCoverage(38081)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := calcMemSize(stack.Back(4), stack.Back(5))
	y := calcMemSize(stack.Back(2), stack.Back(3))

	return math.BigMax(x, y)
}

func memoryStaticCall(stack *Stack) *big.Int {
	fuzz_helper.AddCoverage(62138)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := calcMemSize(stack.Back(4), stack.Back(5))
	y := calcMemSize(stack.Back(2), stack.Back(3))

	return math.BigMax(x, y)
}

func memoryReturn(stack *Stack) *big.Int {
	fuzz_helper.AddCoverage(58303)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return calcMemSize(stack.Back(0), stack.Back(1))
}

func memoryRevert(stack *Stack) *big.Int {
	fuzz_helper.AddCoverage(62967)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return calcMemSize(stack.Back(0), stack.Back(1))
}

func memoryLog(stack *Stack) *big.Int {
	fuzz_helper.AddCoverage(118)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	mSize, mStart := stack.Back(1), stack.Back(0)
	return calcMemSize(mStart, mSize)
}

var _ = fuzz_helper.AddCoverage
