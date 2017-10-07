package vm

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

func TestByteOp(t *testing.T) {
	fuzz_helper.AddCoverage(21402)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	var (
		env   = NewEVM(Context{}, nil, params.TestChainConfig, Config{EnableJit: false, ForceJit: false})
		stack = newstack()
	)
	tests := []struct {
		v        string
		th       uint64
		expected *big.Int
	}{
		{"ABCDEF0908070605040302010000000000000000000000000000000000000000", 0, big.NewInt(0xAB)},
		{"ABCDEF0908070605040302010000000000000000000000000000000000000000", 1, big.NewInt(0xCD)},
		{"00CDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff", 0, big.NewInt(0x00)},
		{"00CDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff", 1, big.NewInt(0xCD)},
		{"0000000000000000000000000000000000000000000000000000000000102030", 31, big.NewInt(0x30)},
		{"0000000000000000000000000000000000000000000000000000000000102030", 30, big.NewInt(0x20)},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 32, big.NewInt(0x0)},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 0xFFFFFFFFFFFFFFFF, big.NewInt(0x0)},
	}
	pc := uint64(0)
	for _, test := range tests {
		fuzz_helper.AddCoverage(42439)
		val := new(big.Int).SetBytes(common.Hex2Bytes(test.v))
		th := new(big.Int).SetUint64(test.th)
		stack.push(val)
		stack.push(th)
		opByte(&pc, env, nil, nil, stack)
		actual := stack.pop()
		if actual.Cmp(test.expected) != 0 {
			fuzz_helper.AddCoverage(9505)
			t.Fatalf("Expected  [%v] %v:th byte to be %v, was %v.", test.v, test.th, test.expected, actual)
		} else {
			fuzz_helper.AddCoverage(50669)
		}
	}
}

func opBenchmark(bench *testing.B, op func(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error), args ...string) {
	fuzz_helper.AddCoverage(8818)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	var (
		env   = NewEVM(Context{}, nil, params.TestChainConfig, Config{EnableJit: false, ForceJit: false})
		stack = newstack()
	)

	byteArgs := make([][]byte, len(args))
	for i, arg := range args {
		fuzz_helper.AddCoverage(7658)
		byteArgs[i] = common.Hex2Bytes(arg)
	}
	fuzz_helper.AddCoverage(7644)
	pc := uint64(0)
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		fuzz_helper.AddCoverage(28012)
		for _, arg := range byteArgs {
			fuzz_helper.AddCoverage(57284)
			a := new(big.Int).SetBytes(arg)
			stack.push(a)
		}
		fuzz_helper.AddCoverage(5559)
		op(&pc, env, nil, nil, stack)
		stack.pop()
	}
}

func BenchmarkOpAdd64(b *testing.B) {
	fuzz_helper.AddCoverage(59440)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := "ffffffff"
	y := "fd37f3e2bba2c4f"

	opBenchmark(b, opAdd, x, y)
}

func BenchmarkOpAdd128(b *testing.B) {
	fuzz_helper.AddCoverage(8312)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := "ffffffffffffffff"
	y := "f5470b43c6549b016288e9a65629687"

	opBenchmark(b, opAdd, x, y)
}

func BenchmarkOpAdd256(b *testing.B) {
	fuzz_helper.AddCoverage(30333)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := "0802431afcbce1fc194c9eaa417b2fb67dc75a95db0bc7ec6b1c8af11df6a1da9"
	y := "a1f5aac137876480252e5dcac62c354ec0d42b76b0642b6181ed099849ea1d57"

	opBenchmark(b, opAdd, x, y)
}

func BenchmarkOpSub64(b *testing.B) {
	fuzz_helper.AddCoverage(61307)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := "51022b6317003a9d"
	y := "a20456c62e00753a"

	opBenchmark(b, opSub, x, y)
}

func BenchmarkOpSub128(b *testing.B) {
	fuzz_helper.AddCoverage(4167)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := "4dde30faaacdc14d00327aac314e915d"
	y := "9bbc61f5559b829a0064f558629d22ba"

	opBenchmark(b, opSub, x, y)
}

func BenchmarkOpSub256(b *testing.B) {
	fuzz_helper.AddCoverage(1602)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := "4bfcd8bb2ac462735b48a17580690283980aa2d679f091c64364594df113ea37"
	y := "97f9b1765588c4e6b69142eb00d20507301545acf3e1238c86c8b29be227d46e"

	opBenchmark(b, opSub, x, y)
}

func BenchmarkOpMul(b *testing.B) {
	fuzz_helper.AddCoverage(56128)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opMul, x, y)
}

func BenchmarkOpDiv256(b *testing.B) {
	fuzz_helper.AddCoverage(28653)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := "ff3f9014f20db29ae04af2c2d265de17"
	y := "fe7fb0d1f59dfe9492ffbf73683fd1e870eec79504c60144cc7f5fc2bad1e611"
	opBenchmark(b, opDiv, x, y)
}

func BenchmarkOpDiv128(b *testing.B) {
	fuzz_helper.AddCoverage(41118)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := "fdedc7f10142ff97"
	y := "fbdfda0e2ce356173d1993d5f70a2b11"
	opBenchmark(b, opDiv, x, y)
}

func BenchmarkOpDiv64(b *testing.B) {
	fuzz_helper.AddCoverage(9151)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := "fcb34eb3"
	y := "f97180878e839129"
	opBenchmark(b, opDiv, x, y)
}

func BenchmarkOpSdiv(b *testing.B) {
	fuzz_helper.AddCoverage(59255)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := "ff3f9014f20db29ae04af2c2d265de17"
	y := "fe7fb0d1f59dfe9492ffbf73683fd1e870eec79504c60144cc7f5fc2bad1e611"

	opBenchmark(b, opSdiv, x, y)
}

func BenchmarkOpMod(b *testing.B) {
	fuzz_helper.AddCoverage(4209)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opMod, x, y)
}

func BenchmarkOpSmod(b *testing.B) {
	fuzz_helper.AddCoverage(16091)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opSmod, x, y)
}

func BenchmarkOpExp(b *testing.B) {
	fuzz_helper.AddCoverage(7163)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opExp, x, y)
}

func BenchmarkOpSignExtend(b *testing.B) {
	fuzz_helper.AddCoverage(33640)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opSignExtend, x, y)
}

func BenchmarkOpLt(b *testing.B) {
	fuzz_helper.AddCoverage(48920)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opLt, x, y)
}

func BenchmarkOpGt(b *testing.B) {
	fuzz_helper.AddCoverage(39497)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opGt, x, y)
}

func BenchmarkOpSlt(b *testing.B) {
	fuzz_helper.AddCoverage(2405)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opSlt, x, y)
}

func BenchmarkOpSgt(b *testing.B) {
	fuzz_helper.AddCoverage(1448)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opSgt, x, y)
}

func BenchmarkOpEq(b *testing.B) {
	fuzz_helper.AddCoverage(31247)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opEq, x, y)
}

func BenchmarkOpAnd(b *testing.B) {
	fuzz_helper.AddCoverage(10799)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opAnd, x, y)
}

func BenchmarkOpOr(b *testing.B) {
	fuzz_helper.AddCoverage(1094)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opOr, x, y)
}

func BenchmarkOpXor(b *testing.B) {
	fuzz_helper.AddCoverage(10390)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opXor, x, y)
}

func BenchmarkOpByte(b *testing.B) {
	fuzz_helper.AddCoverage(43817)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opByte, x, y)
}

func BenchmarkOpAddmod(b *testing.B) {
	fuzz_helper.AddCoverage(42883)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	z := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opAddmod, x, y, z)
}

func BenchmarkOpMulmod(b *testing.B) {
	fuzz_helper.AddCoverage(35848)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	z := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opMulmod, x, y, z)
}

var _ = fuzz_helper.AddCoverage
