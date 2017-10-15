package vm

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

var (
	bigZero                  = new(big.Int)
	errWriteProtection       = errors.New("evm: write protection")
	errReturnDataOutOfBounds = errors.New("evm: return data out of bounds")
	errExecutionReverted     = errors.New("evm: execution reverted")
	errMaxCodeSizeExceeded   = errors.New("evm: max code size exceeded")
)

func opAdd(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(47662)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x, y := stack.pop(), stack.pop()
	stack.push(math.U256(x.Add(x, y)))

	evm.interpreter.intPool.put(y)

	return nil, nil
}

func opSub(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(61477)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x, y := stack.pop(), stack.pop()
	stack.push(math.U256(x.Sub(x, y)))

	evm.interpreter.intPool.put(y)

	return nil, nil
}

func opMul(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(31033)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x, y := stack.pop(), stack.pop()
	stack.push(math.U256(x.Mul(x, y)))

	evm.interpreter.intPool.put(y)

	return nil, nil
}

func opDiv(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(34735)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x, y := stack.pop(), stack.pop()
	if y.Sign() != 0 {
		fuzz_helper.AddCoverage(7335)
		stack.push(math.U256(x.Div(x, y)))
	} else {
		fuzz_helper.AddCoverage(26237)
		stack.push(new(big.Int))
	}
	fuzz_helper.AddCoverage(50912)

	evm.interpreter.intPool.put(y)

	return nil, nil
}

func opSdiv(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(27802)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x, y := math.S256(stack.pop()), math.S256(stack.pop())
	if y.Sign() == 0 {
		fuzz_helper.AddCoverage(8166)
		stack.push(new(big.Int))
		return nil, nil
	} else {
		fuzz_helper.AddCoverage(47787)
		n := new(big.Int)
		if evm.interpreter.intPool.get().Mul(x, y).Sign() < 0 {
			fuzz_helper.AddCoverage(60019)
			n.SetInt64(-1)
		} else {
			fuzz_helper.AddCoverage(28532)
			n.SetInt64(1)
		}
		fuzz_helper.AddCoverage(29361)

		res := x.Div(x.Abs(x), y.Abs(y))
		res.Mul(res, n)

		stack.push(math.U256(res))
	}
	fuzz_helper.AddCoverage(16892)
	evm.interpreter.intPool.put(y)
	return nil, nil
}

func opMod(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(23159)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x, y := stack.pop(), stack.pop()
	if y.Sign() == 0 {
		fuzz_helper.AddCoverage(40384)
		stack.push(new(big.Int))
	} else {
		fuzz_helper.AddCoverage(19679)
		stack.push(math.U256(x.Mod(x, y)))
	}
	fuzz_helper.AddCoverage(16491)
	evm.interpreter.intPool.put(y)
	return nil, nil
}

func opSmod(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(43575)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x, y := math.S256(stack.pop()), math.S256(stack.pop())

	if y.Sign() == 0 {
		fuzz_helper.AddCoverage(57474)
		stack.push(new(big.Int))
	} else {
		fuzz_helper.AddCoverage(57224)
		n := new(big.Int)
		if x.Sign() < 0 {
			fuzz_helper.AddCoverage(53237)
			n.SetInt64(-1)
		} else {
			fuzz_helper.AddCoverage(54351)
			n.SetInt64(1)
		}
		fuzz_helper.AddCoverage(63341)

		res := x.Mod(x.Abs(x), y.Abs(y))
		res.Mul(res, n)

		stack.push(math.U256(res))
	}
	fuzz_helper.AddCoverage(16162)
	evm.interpreter.intPool.put(y)
	return nil, nil
}

func opExp(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(43457)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	base, exponent := stack.pop(), stack.pop()
	stack.push(math.Exp(base, exponent))

	evm.interpreter.intPool.put(base, exponent)

	return nil, nil
}

func opSignExtend(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(33005)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	back := stack.pop()
	if back.Cmp(big.NewInt(31)) < 0 {
		fuzz_helper.AddCoverage(49025)
		bit := uint(back.Uint64()*8 + 7)
		num := stack.pop()
		mask := back.Lsh(common.Big1, bit)
		mask.Sub(mask, common.Big1)
		if num.Bit(int(bit)) > 0 {
			fuzz_helper.AddCoverage(57001)
			num.Or(num, mask.Not(mask))
		} else {
			fuzz_helper.AddCoverage(33639)
			num.And(num, mask)
		}
		fuzz_helper.AddCoverage(20123)

		stack.push(math.U256(num))
	} else {
		fuzz_helper.AddCoverage(6504)
	}
	fuzz_helper.AddCoverage(36120)

	evm.interpreter.intPool.put(back)
	return nil, nil
}

func opNot(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(63668)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := stack.pop()
	stack.push(math.U256(x.Not(x)))
	return nil, nil
}

func opLt(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(58296)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x, y := stack.pop(), stack.pop()
	if x.Cmp(y) < 0 {
		fuzz_helper.AddCoverage(5059)
		stack.push(evm.interpreter.intPool.get().SetUint64(1))
	} else {
		fuzz_helper.AddCoverage(2794)
		stack.push(new(big.Int))
	}
	fuzz_helper.AddCoverage(64985)

	evm.interpreter.intPool.put(x, y)
	return nil, nil
}

func opGt(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(64290)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x, y := stack.pop(), stack.pop()
	if x.Cmp(y) > 0 {
		fuzz_helper.AddCoverage(4032)
		stack.push(evm.interpreter.intPool.get().SetUint64(1))
	} else {
		fuzz_helper.AddCoverage(58678)
		stack.push(new(big.Int))
	}
	fuzz_helper.AddCoverage(23582)

	evm.interpreter.intPool.put(x, y)
	return nil, nil
}

func opSlt(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(52438)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x, y := math.S256(stack.pop()), math.S256(stack.pop())
	if x.Cmp(math.S256(y)) < 0 {
		fuzz_helper.AddCoverage(2334)
		stack.push(evm.interpreter.intPool.get().SetUint64(1))
	} else {
		fuzz_helper.AddCoverage(53963)
		stack.push(new(big.Int))
	}
	fuzz_helper.AddCoverage(6473)

	evm.interpreter.intPool.put(x, y)
	return nil, nil
}

func opSgt(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(65501)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x, y := math.S256(stack.pop()), math.S256(stack.pop())
	if x.Cmp(y) > 0 {
		fuzz_helper.AddCoverage(25468)
		stack.push(evm.interpreter.intPool.get().SetUint64(1))
	} else {
		fuzz_helper.AddCoverage(23110)
		stack.push(new(big.Int))
	}
	fuzz_helper.AddCoverage(9291)

	evm.interpreter.intPool.put(x, y)
	return nil, nil
}

func opEq(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(27667)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x, y := stack.pop(), stack.pop()
	if x.Cmp(y) == 0 {
		fuzz_helper.AddCoverage(31403)
		stack.push(evm.interpreter.intPool.get().SetUint64(1))
	} else {
		fuzz_helper.AddCoverage(21562)
		stack.push(new(big.Int))
	}
	fuzz_helper.AddCoverage(50961)

	evm.interpreter.intPool.put(x, y)
	return nil, nil
}

func opIszero(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(42411)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x := stack.pop()
	if x.Sign() > 0 {
		fuzz_helper.AddCoverage(37741)
		stack.push(new(big.Int))
	} else {
		fuzz_helper.AddCoverage(11155)
		stack.push(evm.interpreter.intPool.get().SetUint64(1))
	}
	fuzz_helper.AddCoverage(61278)

	evm.interpreter.intPool.put(x)
	return nil, nil
}

func opAnd(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(2147)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x, y := stack.pop(), stack.pop()
	stack.push(x.And(x, y))

	evm.interpreter.intPool.put(y)
	return nil, nil
}

func opOr(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(478)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x, y := stack.pop(), stack.pop()
	stack.push(x.Or(x, y))

	evm.interpreter.intPool.put(y)
	return nil, nil
}

func opXor(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(48634)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x, y := stack.pop(), stack.pop()
	stack.push(x.Xor(x, y))

	evm.interpreter.intPool.put(y)
	return nil, nil
}

func opByte(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(49535)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	th, val := stack.pop(), stack.peek()
	if th.Cmp(common.Big32) < 0 {
		fuzz_helper.AddCoverage(49456)
		b := math.Byte(val, 32, int(th.Int64()))
		val.SetUint64(uint64(b))
	} else {
		fuzz_helper.AddCoverage(45291)
		val.SetUint64(0)
	}
	fuzz_helper.AddCoverage(413)
	evm.interpreter.intPool.put(th)
	return nil, nil
}

func opAddmod(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(30119)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x, y, z := stack.pop(), stack.pop(), stack.pop()
	if z.Cmp(bigZero) > 0 {
		fuzz_helper.AddCoverage(14176)
		add := x.Add(x, y)
		add.Mod(add, z)
		stack.push(math.U256(add))
	} else {
		fuzz_helper.AddCoverage(51212)
		stack.push(new(big.Int))
	}
	fuzz_helper.AddCoverage(3321)

	evm.interpreter.intPool.put(y, z)
	return nil, nil
}

func opMulmod(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(52714)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x, y, z := stack.pop(), stack.pop(), stack.pop()
	if z.Cmp(bigZero) > 0 {
		fuzz_helper.AddCoverage(38571)
		mul := x.Mul(x, y)
		mul.Mod(mul, z)
		stack.push(math.U256(mul))
	} else {
		fuzz_helper.AddCoverage(52556)
		stack.push(new(big.Int))
	}
	fuzz_helper.AddCoverage(41907)

	evm.interpreter.intPool.put(y, z)
	return nil, nil
}

func opSha3(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(41975)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	offset, size := stack.pop(), stack.pop()
	data := memory.Get(offset.Int64(), size.Int64())
	hash := crypto.Keccak256(data)

	if evm.vmConfig.EnablePreimageRecording {
		fuzz_helper.AddCoverage(37531)
		evm.StateDB.AddPreimage(common.BytesToHash(hash), data)
	} else {
		fuzz_helper.AddCoverage(33863)
	}
	fuzz_helper.AddCoverage(27403)

	stack.push(new(big.Int).SetBytes(hash))

	evm.interpreter.intPool.put(offset, size)
	return nil, nil
}

func opAddress(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(24366)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	stack.push(contract.Address().Big())
	return nil, nil
}

func opBalance(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(14693)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	addr := common.BigToAddress(stack.pop())
	balance := evm.StateDB.GetBalance(addr)

	stack.push(new(big.Int).Set(balance))
	return nil, nil
}

func opOrigin(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(6696)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	stack.push(evm.Origin.Big())
	return nil, nil
}

func opCaller(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(708)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	stack.push(contract.Caller().Big())
	return nil, nil
}

func opCallValue(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(17273)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	stack.push(evm.interpreter.intPool.get().Set(contract.value))
	return nil, nil
}

func opCallDataLoad(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(12629)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	stack.push(new(big.Int).SetBytes(getDataBig(contract.Input, stack.pop(), big32)))
	return nil, nil
}

func opCallDataSize(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(51895)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	stack.push(evm.interpreter.intPool.get().SetInt64(int64(len(contract.Input))))
	return nil, nil
}

func opCallDataCopy(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(32875)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	var (
		memOffset  = stack.pop()
		dataOffset = stack.pop()
		length     = stack.pop()
	)
	memory.Set(memOffset.Uint64(), length.Uint64(), getDataBig(contract.Input, dataOffset, length))

	evm.interpreter.intPool.put(memOffset, dataOffset, length)
	return nil, nil
}

func opReturnDataSize(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(2352)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	stack.push(evm.interpreter.intPool.get().SetUint64(uint64(len(evm.interpreter.returnData))))
	return nil, nil
}

func opReturnDataCopy(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(6089)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	var (
		memOffset  = stack.pop()
		dataOffset = stack.pop()
		length     = stack.pop()
	)
	defer evm.interpreter.intPool.put(memOffset, dataOffset, length)

	end := new(big.Int).Add(dataOffset, length)
	if end.BitLen() > 64 || uint64(len(evm.interpreter.returnData)) < end.Uint64() {
		fuzz_helper.AddCoverage(46184)
		return nil, errReturnDataOutOfBounds
	} else {
		fuzz_helper.AddCoverage(5156)
	}
	fuzz_helper.AddCoverage(25407)
	memory.Set(memOffset.Uint64(), length.Uint64(), evm.interpreter.returnData[dataOffset.Uint64():end.Uint64()])

	return nil, nil
}

func opExtCodeSize(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(17378)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	a := stack.pop()

	addr := common.BigToAddress(a)
	a.SetInt64(int64(evm.StateDB.GetCodeSize(addr)))
	stack.push(a)

	return nil, nil
}

func opCodeSize(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(9796)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	l := evm.interpreter.intPool.get().SetInt64(int64(len(contract.Code)))
	stack.push(l)
	return nil, nil
}

func opCodeCopy(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(39315)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	var (
		memOffset  = stack.pop()
		codeOffset = stack.pop()
		length     = stack.pop()
	)
	codeCopy := getDataBig(contract.Code, codeOffset, length)
	memory.Set(memOffset.Uint64(), length.Uint64(), codeCopy)

	evm.interpreter.intPool.put(memOffset, codeOffset, length)
	return nil, nil
}

func opExtCodeCopy(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(44778)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	var (
		addr       = common.BigToAddress(stack.pop())
		memOffset  = stack.pop()
		codeOffset = stack.pop()
		length     = stack.pop()
	)
	codeCopy := getDataBig(evm.StateDB.GetCode(addr), codeOffset, length)
	memory.Set(memOffset.Uint64(), length.Uint64(), codeCopy)

	evm.interpreter.intPool.put(memOffset, codeOffset, length)
	return nil, nil
}

func opGasprice(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(45806)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	stack.push(evm.interpreter.intPool.get().Set(evm.GasPrice))
	return nil, nil
}

func opBlockhash(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(3985)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	num := stack.pop()

	n := evm.interpreter.intPool.get().Sub(evm.BlockNumber, common.Big257)
	if num.Cmp(n) > 0 && num.Cmp(evm.BlockNumber) < 0 {
		fuzz_helper.AddCoverage(9778)
		stack.push(evm.GetHash(num.Uint64()).Big())
	} else {
		fuzz_helper.AddCoverage(41307)
		stack.push(new(big.Int))
	}
	fuzz_helper.AddCoverage(15060)

	evm.interpreter.intPool.put(num, n)
	return nil, nil
}

func opCoinbase(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(3028)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	stack.push(evm.Coinbase.Big())
	return nil, nil
}

func opTimestamp(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(11144)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	stack.push(math.U256(new(big.Int).Set(evm.Time)))
	return nil, nil
}

func opNumber(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(48281)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	stack.push(math.U256(new(big.Int).Set(evm.BlockNumber)))
	return nil, nil
}

func opDifficulty(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(4326)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	stack.push(math.U256(new(big.Int).Set(evm.Difficulty)))
	return nil, nil
}

func opGasLimit(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(24314)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	stack.push(math.U256(new(big.Int).Set(evm.GasLimit)))
	return nil, nil
}

func opPop(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(46928)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	evm.interpreter.intPool.put(stack.pop())
	return nil, nil
}

func opMload(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(21553)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	offset := stack.pop()
	val := new(big.Int).SetBytes(memory.Get(offset.Int64(), 32))
	stack.push(val)

	evm.interpreter.intPool.put(offset)
	return nil, nil
}

func opMstore(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(42579)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()

	mStart, val := stack.pop(), stack.pop()
	memory.Set(mStart.Uint64(), 32, math.PaddedBigBytes(val, 32))

	evm.interpreter.intPool.put(mStart, val)
	return nil, nil
}

func opMstore8(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(31252)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	off, val := stack.pop().Int64(), stack.pop().Int64()
	memory.store[off] = byte(val & 0xff)

	return nil, nil
}

func opSload(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(40934)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	loc := common.BigToHash(stack.pop())
	val := evm.StateDB.GetState(contract.Address(), loc).Big()
	stack.push(val)
	return nil, nil
}

func opSstore(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(10798)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	loc := common.BigToHash(stack.pop())
	val := stack.pop()
	evm.StateDB.SetState(contract.Address(), loc, common.BigToHash(val))

	evm.interpreter.intPool.put(val)
	return nil, nil
}

func opJump(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(6834)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
    /* GUIDO */ return nil, fmt.Errorf("invalid jump destination")
	pos := stack.pop()
	if !contract.jumpdests.has(contract.CodeHash, contract.Code, pos) {
		fuzz_helper.AddCoverage(2401)
		nop := contract.GetOp(pos.Uint64())
		return nil, fmt.Errorf("invalid jump destination (%v) %v", nop, pos)
	} else {
		fuzz_helper.AddCoverage(28464)
	}
	fuzz_helper.AddCoverage(39786)
	*pc = pos.Uint64()

	evm.interpreter.intPool.put(pos)
	return nil, nil
}

func opJumpi(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(41154)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
    /* GUIDO */ return nil, fmt.Errorf("invalid jump destination")
	pos, cond := stack.pop(), stack.pop()
	if cond.Sign() != 0 {
		fuzz_helper.AddCoverage(48780)
		if !contract.jumpdests.has(contract.CodeHash, contract.Code, pos) {
			fuzz_helper.AddCoverage(47782)
			nop := contract.GetOp(pos.Uint64())
			return nil, fmt.Errorf("invalid jump destination (%v) %v", nop, pos)
		} else {
			fuzz_helper.AddCoverage(11015)
		}
		fuzz_helper.AddCoverage(64489)
		*pc = pos.Uint64()
	} else {
		fuzz_helper.AddCoverage(43113)
		*pc++
	}
	fuzz_helper.AddCoverage(59174)

	evm.interpreter.intPool.put(pos, cond)
	return nil, nil
}

func opJumpdest(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(30435)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return nil, nil
}

func opPc(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(44124)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	stack.push(evm.interpreter.intPool.get().SetUint64(*pc))
	return nil, nil
}

func opMsize(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(14220)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	stack.push(evm.interpreter.intPool.get().SetInt64(int64(memory.Len())))
	return nil, nil
}

func opGas(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(17845)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	stack.push(evm.interpreter.intPool.get().SetUint64(contract.Gas))
	return nil, nil
}

func opCreate(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(7703)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	var (
		value        = stack.pop()
		offset, size = stack.pop(), stack.pop()
		input        = memory.Get(offset.Int64(), size.Int64())
		gas          = contract.Gas
	)
	if evm.ChainConfig().IsEIP150(evm.BlockNumber) {
		fuzz_helper.AddCoverage(49477)
		gas -= gas / 64
	} else {
		fuzz_helper.AddCoverage(39615)
	}
	fuzz_helper.AddCoverage(51344)

	contract.UseGas(gas)
	res, addr, returnGas, suberr := evm.Create(contract, input, gas, value)

	if evm.ChainConfig().IsHomestead(evm.BlockNumber) && suberr == ErrCodeStoreOutOfGas {
		fuzz_helper.AddCoverage(30441)
		stack.push(new(big.Int))
	} else {
		fuzz_helper.AddCoverage(2510)
		if suberr != nil && suberr != ErrCodeStoreOutOfGas {
			fuzz_helper.AddCoverage(61940)
			stack.push(new(big.Int))
		} else {
			fuzz_helper.AddCoverage(14597)
			stack.push(addr.Big())
		}
	}
	fuzz_helper.AddCoverage(51806)
	contract.Gas += returnGas
	evm.interpreter.intPool.put(value, offset, size)

	if suberr == errExecutionReverted {
		fuzz_helper.AddCoverage(32844)
		return res, nil
	} else {
		fuzz_helper.AddCoverage(53872)
	}
	fuzz_helper.AddCoverage(45980)
	return nil, nil
}

func opCall(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(36572)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	gas := stack.pop().Uint64()

	addr, value := stack.pop(), stack.pop()
	value = math.U256(value)

	inOffset, inSize := stack.pop(), stack.pop()

	retOffset, retSize := stack.pop(), stack.pop()

	address := common.BigToAddress(addr)

	args := memory.Get(inOffset.Int64(), inSize.Int64())

	if value.Sign() != 0 {
		fuzz_helper.AddCoverage(40206)
		gas += params.CallStipend
	} else {
		fuzz_helper.AddCoverage(12610)
	}
	fuzz_helper.AddCoverage(3587)
	ret, returnGas, err := evm.Call(contract, address, args, gas, value)
	if err != nil {
		fuzz_helper.AddCoverage(7269)
		stack.push(new(big.Int))
	} else {
		fuzz_helper.AddCoverage(35205)
		stack.push(big.NewInt(1))
	}
	fuzz_helper.AddCoverage(62840)
	if err == nil || err == errExecutionReverted {
		fuzz_helper.AddCoverage(38538)
		memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
	} else {
		fuzz_helper.AddCoverage(35676)
	}
	fuzz_helper.AddCoverage(20665)
	contract.Gas += returnGas

	evm.interpreter.intPool.put(addr, value, inOffset, inSize, retOffset, retSize)
	return ret, nil
}

func opCallCode(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(2961)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	gas := stack.pop().Uint64()

	addr, value := stack.pop(), stack.pop()
	value = math.U256(value)

	inOffset, inSize := stack.pop(), stack.pop()

	retOffset, retSize := stack.pop(), stack.pop()

	address := common.BigToAddress(addr)

	args := memory.Get(inOffset.Int64(), inSize.Int64())

	if value.Sign() != 0 {
		fuzz_helper.AddCoverage(26142)
		gas += params.CallStipend
	} else {
		fuzz_helper.AddCoverage(44665)
	}
	fuzz_helper.AddCoverage(30773)

	ret, returnGas, err := evm.CallCode(contract, address, args, gas, value)
	if err != nil {
		fuzz_helper.AddCoverage(63536)
		stack.push(new(big.Int))
	} else {
		fuzz_helper.AddCoverage(30800)
		stack.push(big.NewInt(1))
	}
	fuzz_helper.AddCoverage(20366)
	if err == nil || err == errExecutionReverted {
		fuzz_helper.AddCoverage(32682)
		memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
	} else {
		fuzz_helper.AddCoverage(62554)
	}
	fuzz_helper.AddCoverage(44276)
	contract.Gas += returnGas

	evm.interpreter.intPool.put(addr, value, inOffset, inSize, retOffset, retSize)
	return ret, nil
}

func opDelegateCall(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(58244)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	gas, to, inOffset, inSize, outOffset, outSize := stack.pop().Uint64(), stack.pop(), stack.pop(), stack.pop(), stack.pop(), stack.pop()

	toAddr := common.BigToAddress(to)
	args := memory.Get(inOffset.Int64(), inSize.Int64())

	ret, returnGas, err := evm.DelegateCall(contract, toAddr, args, gas)
	if err != nil {
		fuzz_helper.AddCoverage(19116)
		stack.push(new(big.Int))
	} else {
		fuzz_helper.AddCoverage(40623)
		stack.push(big.NewInt(1))
	}
	fuzz_helper.AddCoverage(7855)
	if err == nil || err == errExecutionReverted {
		fuzz_helper.AddCoverage(10558)
		memory.Set(outOffset.Uint64(), outSize.Uint64(), ret)
	} else {
		fuzz_helper.AddCoverage(55739)
	}
	fuzz_helper.AddCoverage(2627)
	contract.Gas += returnGas

	evm.interpreter.intPool.put(to, inOffset, inSize, outOffset, outSize)
	return ret, nil
}

func opStaticCall(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(44486)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()

	gas := stack.pop().Uint64()

	addr := stack.pop()

	inOffset, inSize := stack.pop(), stack.pop()

	retOffset, retSize := stack.pop(), stack.pop()

	address := common.BigToAddress(addr)

	args := memory.Get(inOffset.Int64(), inSize.Int64())

	ret, returnGas, err := evm.StaticCall(contract, address, args, gas)
	if err != nil {
		fuzz_helper.AddCoverage(38903)
		stack.push(new(big.Int))
	} else {
		fuzz_helper.AddCoverage(44209)
		stack.push(big.NewInt(1))
	}
	fuzz_helper.AddCoverage(50458)
	if err == nil || err == errExecutionReverted {
		fuzz_helper.AddCoverage(36418)
		memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
	} else {
		fuzz_helper.AddCoverage(9823)
	}
	fuzz_helper.AddCoverage(53727)
	contract.Gas += returnGas

	evm.interpreter.intPool.put(addr, inOffset, inSize, retOffset, retSize)
	return ret, nil
}

func opReturn(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(22432)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	offset, size := stack.pop(), stack.pop()
	ret := memory.GetPtr(offset.Int64(), size.Int64())

	evm.interpreter.intPool.put(offset, size)
	return ret, nil
}

func opRevert(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(16062)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	offset, size := stack.pop(), stack.pop()
	ret := memory.GetPtr(offset.Int64(), size.Int64())

	evm.interpreter.intPool.put(offset, size)
	return ret, nil
}

func opStop(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(47431)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return nil, nil
}

func opSuicide(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.AddCoverage(64337)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	balance := evm.StateDB.GetBalance(contract.Address())
	evm.StateDB.AddBalance(common.BigToAddress(stack.pop()), balance)

	evm.StateDB.Suicide(contract.Address())
	return nil, nil
}

// make log instruction function
func makeLog(size int) executionFunc {
	fuzz_helper.AddCoverage(868)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return func(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
		fuzz_helper.AddCoverage(3074)
		topics := make([]common.Hash, size)
		mStart, mSize := stack.pop(), stack.pop()
		for i := 0; i < size; i++ {
			fuzz_helper.AddCoverage(28786)
			topics[i] = common.BigToHash(stack.pop())
		}
		fuzz_helper.AddCoverage(63763)

		d := memory.Get(mStart.Int64(), mSize.Int64())
		evm.StateDB.AddLog(&types.Log{
			Address: contract.Address(),
			Topics:  topics,
			Data:    d,

			BlockNumber: evm.BlockNumber.Uint64(),
		})

		evm.interpreter.intPool.put(mStart, mSize)
		return nil, nil
	}
}

// make push instruction function
func makePush(size uint64, pushByteSize int) executionFunc {
	fuzz_helper.AddCoverage(23847)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return func(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
		fuzz_helper.AddCoverage(48969)
		codeLen := len(contract.Code)

		startMin := codeLen
		if int(*pc+1) < startMin {
			fuzz_helper.AddCoverage(7912)
			startMin = int(*pc + 1)
		} else {
			fuzz_helper.AddCoverage(61323)
		}
		fuzz_helper.AddCoverage(19394)

		endMin := codeLen
		if startMin+pushByteSize < endMin {
			fuzz_helper.AddCoverage(1310)
			endMin = startMin + pushByteSize
		} else {
			fuzz_helper.AddCoverage(57110)
		}
		fuzz_helper.AddCoverage(26860)

		integer := evm.interpreter.intPool.get()
		stack.push(integer.SetBytes(common.RightPadBytes(contract.Code[startMin:endMin], pushByteSize)))

		*pc += size
		return nil, nil
	}
}

// make push instruction function
func makeDup(size int64) executionFunc {
	fuzz_helper.AddCoverage(35135)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return func(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
		fuzz_helper.AddCoverage(18036)
		stack.dup(evm.interpreter.intPool, int(size))
		return nil, nil
	}
}

// make swap instruction function
func makeSwap(size int64) executionFunc {
	fuzz_helper.AddCoverage(2544)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()

	size += 1
	return func(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
		fuzz_helper.AddCoverage(60942)
		stack.swap(int(size))
		return nil, nil
	}
}

var _ = fuzz_helper.AddCoverage
