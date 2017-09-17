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
	fuzz_helper.CoverTab[22588]++
	x, y := stack.pop(), stack.pop()
	stack.push(math.U256(x.Add(x, y)))

	evm.interpreter.intPool.put(y)

	return nil, nil
}

func opSub(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[44810]++
	x, y := stack.pop(), stack.pop()
	stack.push(math.U256(x.Sub(x, y)))

	evm.interpreter.intPool.put(y)

	return nil, nil
}

func opMul(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[5262]++
	x, y := stack.pop(), stack.pop()
	stack.push(math.U256(x.Mul(x, y)))

	evm.interpreter.intPool.put(y)

	return nil, nil
}

func opDiv(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[17878]++
	x, y := stack.pop(), stack.pop()
	if y.Sign() != 0 {
		fuzz_helper.CoverTab[39040]++
		stack.push(math.U256(x.Div(x, y)))
	} else {
		fuzz_helper.CoverTab[2095]++
		stack.push(new(big.Int))
	}
	fuzz_helper.CoverTab[45021]++

	evm.interpreter.intPool.put(y)

	return nil, nil
}

func opSdiv(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[21668]++
	x, y := math.S256(stack.pop()), math.S256(stack.pop())
	if y.Sign() == 0 {
		fuzz_helper.CoverTab[16619]++
		stack.push(new(big.Int))
		return nil, nil
	} else {
		fuzz_helper.CoverTab[12692]++
		n := new(big.Int)
		if evm.interpreter.intPool.get().Mul(x, y).Sign() < 0 {
			fuzz_helper.CoverTab[6577]++
			n.SetInt64(-1)
		} else {
			fuzz_helper.CoverTab[17393]++
			n.SetInt64(1)
		}
		fuzz_helper.CoverTab[42483]++

		res := x.Div(x.Abs(x), y.Abs(y))
		res.Mul(res, n)

		stack.push(math.U256(res))
	}
	fuzz_helper.CoverTab[45213]++
	evm.interpreter.intPool.put(y)
	return nil, nil
}

func opMod(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[64174]++
	x, y := stack.pop(), stack.pop()
	if y.Sign() == 0 {
		fuzz_helper.CoverTab[35657]++
		stack.push(new(big.Int))
	} else {
		fuzz_helper.CoverTab[30358]++
		stack.push(math.U256(x.Mod(x, y)))
	}
	fuzz_helper.CoverTab[38740]++
	evm.interpreter.intPool.put(y)
	return nil, nil
}

func opSmod(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[23294]++
	x, y := math.S256(stack.pop()), math.S256(stack.pop())

	if y.Sign() == 0 {
		fuzz_helper.CoverTab[11162]++
		stack.push(new(big.Int))
	} else {
		fuzz_helper.CoverTab[49217]++
		n := new(big.Int)
		if x.Sign() < 0 {
			fuzz_helper.CoverTab[64074]++
			n.SetInt64(-1)
		} else {
			fuzz_helper.CoverTab[28614]++
			n.SetInt64(1)
		}
		fuzz_helper.CoverTab[34511]++

		res := x.Mod(x.Abs(x), y.Abs(y))
		res.Mul(res, n)

		stack.push(math.U256(res))
	}
	fuzz_helper.CoverTab[61639]++
	evm.interpreter.intPool.put(y)
	return nil, nil
}

func opExp(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[39226]++
	base, exponent := stack.pop(), stack.pop()
	stack.push(math.Exp(base, exponent))

	evm.interpreter.intPool.put(base, exponent)

	return nil, nil
}

func opSignExtend(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[2297]++
	back := stack.pop()
	if back.Cmp(big.NewInt(31)) < 0 {
		fuzz_helper.CoverTab[52877]++
		bit := uint(back.Uint64()*8 + 7)
		num := stack.pop()
		mask := back.Lsh(common.Big1, bit)
		mask.Sub(mask, common.Big1)
		if num.Bit(int(bit)) > 0 {
			fuzz_helper.CoverTab[33340]++
			num.Or(num, mask.Not(mask))
		} else {
			fuzz_helper.CoverTab[15638]++
			num.And(num, mask)
		}
		fuzz_helper.CoverTab[778]++

		stack.push(math.U256(num))
	} else {
		fuzz_helper.CoverTab[45869]++
	}
	fuzz_helper.CoverTab[40870]++

	evm.interpreter.intPool.put(back)
	return nil, nil
}

func opNot(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[23368]++
	x := stack.pop()
	stack.push(math.U256(x.Not(x)))
	return nil, nil
}

func opLt(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[12901]++
	x, y := stack.pop(), stack.pop()
	if x.Cmp(y) < 0 {
		fuzz_helper.CoverTab[42993]++
		stack.push(evm.interpreter.intPool.get().SetUint64(1))
	} else {
		fuzz_helper.CoverTab[30301]++
		stack.push(new(big.Int))
	}
	fuzz_helper.CoverTab[12499]++

	evm.interpreter.intPool.put(x, y)
	return nil, nil
}

func opGt(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[45210]++
	x, y := stack.pop(), stack.pop()
	if x.Cmp(y) > 0 {
		fuzz_helper.CoverTab[3566]++
		stack.push(evm.interpreter.intPool.get().SetUint64(1))
	} else {
		fuzz_helper.CoverTab[47636]++
		stack.push(new(big.Int))
	}
	fuzz_helper.CoverTab[264]++

	evm.interpreter.intPool.put(x, y)
	return nil, nil
}

func opSlt(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[8730]++
	x, y := math.S256(stack.pop()), math.S256(stack.pop())
	if x.Cmp(math.S256(y)) < 0 {
		fuzz_helper.CoverTab[63931]++
		stack.push(evm.interpreter.intPool.get().SetUint64(1))
	} else {
		fuzz_helper.CoverTab[19009]++
		stack.push(new(big.Int))
	}
	fuzz_helper.CoverTab[20539]++

	evm.interpreter.intPool.put(x, y)
	return nil, nil
}

func opSgt(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[64748]++
	x, y := math.S256(stack.pop()), math.S256(stack.pop())
	if x.Cmp(y) > 0 {
		fuzz_helper.CoverTab[18500]++
		stack.push(evm.interpreter.intPool.get().SetUint64(1))
	} else {
		fuzz_helper.CoverTab[52152]++
		stack.push(new(big.Int))
	}
	fuzz_helper.CoverTab[50446]++

	evm.interpreter.intPool.put(x, y)
	return nil, nil
}

func opEq(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[17111]++
	x, y := stack.pop(), stack.pop()
	if x.Cmp(y) == 0 {
		fuzz_helper.CoverTab[55848]++
		stack.push(evm.interpreter.intPool.get().SetUint64(1))
	} else {
		fuzz_helper.CoverTab[50755]++
		stack.push(new(big.Int))
	}
	fuzz_helper.CoverTab[9670]++

	evm.interpreter.intPool.put(x, y)
	return nil, nil
}

func opIszero(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[912]++
	x := stack.pop()
	if x.Sign() > 0 {
		fuzz_helper.CoverTab[15513]++
		stack.push(new(big.Int))
	} else {
		fuzz_helper.CoverTab[17300]++
		stack.push(evm.interpreter.intPool.get().SetUint64(1))
	}
	fuzz_helper.CoverTab[64631]++

	evm.interpreter.intPool.put(x)
	return nil, nil
}

func opAnd(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[16403]++
	x, y := stack.pop(), stack.pop()
	stack.push(x.And(x, y))

	evm.interpreter.intPool.put(y)
	return nil, nil
}

func opOr(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[40937]++
	x, y := stack.pop(), stack.pop()
	stack.push(x.Or(x, y))

	evm.interpreter.intPool.put(y)
	return nil, nil
}

func opXor(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[33825]++
	x, y := stack.pop(), stack.pop()
	stack.push(x.Xor(x, y))

	evm.interpreter.intPool.put(y)
	return nil, nil
}

func opByte(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[7237]++
	th, val := stack.pop(), stack.peek()
	if th.Cmp(common.Big32) < 0 {
		fuzz_helper.CoverTab[52715]++
		b := math.Byte(val, 32, int(th.Int64()))
		val.SetUint64(uint64(b))
	} else {
		fuzz_helper.CoverTab[11389]++
		val.SetUint64(0)
	}
	fuzz_helper.CoverTab[23248]++
	evm.interpreter.intPool.put(th)
	return nil, nil
}

func opAddmod(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[60629]++
	x, y, z := stack.pop(), stack.pop(), stack.pop()
	if z.Cmp(bigZero) > 0 {
		fuzz_helper.CoverTab[5383]++
		add := x.Add(x, y)
		add.Mod(add, z)
		stack.push(math.U256(add))
	} else {
		fuzz_helper.CoverTab[52957]++
		stack.push(new(big.Int))
	}
	fuzz_helper.CoverTab[23245]++

	evm.interpreter.intPool.put(y, z)
	return nil, nil
}

func opMulmod(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[6211]++
	x, y, z := stack.pop(), stack.pop(), stack.pop()
	if z.Cmp(bigZero) > 0 {
		fuzz_helper.CoverTab[15785]++
		mul := x.Mul(x, y)
		mul.Mod(mul, z)
		stack.push(math.U256(mul))
	} else {
		fuzz_helper.CoverTab[9735]++
		stack.push(new(big.Int))
	}
	fuzz_helper.CoverTab[49245]++

	evm.interpreter.intPool.put(y, z)
	return nil, nil
}

func opSha3(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[45823]++
	offset, size := stack.pop(), stack.pop()
	data := memory.Get(offset.Int64(), size.Int64())
	hash := crypto.Keccak256(data)

	if evm.vmConfig.EnablePreimageRecording {
		fuzz_helper.CoverTab[24978]++
		evm.StateDB.AddPreimage(common.BytesToHash(hash), data)
	} else {
		fuzz_helper.CoverTab[61755]++
	}
	fuzz_helper.CoverTab[48647]++

	stack.push(new(big.Int).SetBytes(hash))

	evm.interpreter.intPool.put(offset, size)
	return nil, nil
}

func opAddress(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[19607]++
	stack.push(contract.Address().Big())
	return nil, nil
}

func opBalance(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[28743]++
	addr := common.BigToAddress(stack.pop())
	balance := evm.StateDB.GetBalance(addr)

	stack.push(new(big.Int).Set(balance))
	return nil, nil
}

func opOrigin(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[8832]++
	stack.push(evm.Origin.Big())
	return nil, nil
}

func opCaller(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[40052]++
	stack.push(contract.Caller().Big())
	return nil, nil
}

func opCallValue(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[63449]++
	stack.push(evm.interpreter.intPool.get().Set(contract.value))
	return nil, nil
}

func opCallDataLoad(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[47485]++
	stack.push(new(big.Int).SetBytes(getDataBig(contract.Input, stack.pop(), big32)))
	return nil, nil
}

func opCallDataSize(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[60075]++
	stack.push(evm.interpreter.intPool.get().SetInt64(int64(len(contract.Input))))
	return nil, nil
}

func opCallDataCopy(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[21817]++
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
	fuzz_helper.CoverTab[57682]++
	stack.push(evm.interpreter.intPool.get().SetUint64(uint64(len(evm.interpreter.returnData))))
	return nil, nil
}

func opReturnDataCopy(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[45496]++
	var (
		memOffset  = stack.pop()
		dataOffset = stack.pop()
		length     = stack.pop()
	)
	defer evm.interpreter.intPool.put(memOffset, dataOffset, length)

	end := new(big.Int).Add(dataOffset, length)
	if end.BitLen() > 64 || uint64(len(evm.interpreter.returnData)) < end.Uint64() {
		fuzz_helper.CoverTab[22210]++
		return nil, errReturnDataOutOfBounds
	} else {
		fuzz_helper.CoverTab[4417]++
	}
	fuzz_helper.CoverTab[3661]++
	memory.Set(memOffset.Uint64(), length.Uint64(), evm.interpreter.returnData[dataOffset.Uint64():end.Uint64()])

	return nil, nil
}

func opExtCodeSize(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[5093]++
	a := stack.pop()

	addr := common.BigToAddress(a)
	a.SetInt64(int64(evm.StateDB.GetCodeSize(addr)))
	stack.push(a)

	return nil, nil
}

func opCodeSize(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[56274]++
	l := evm.interpreter.intPool.get().SetInt64(int64(len(contract.Code)))
	stack.push(l)
	return nil, nil
}

func opCodeCopy(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[1404]++
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
	fuzz_helper.CoverTab[34815]++
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
	fuzz_helper.CoverTab[58334]++
	stack.push(evm.interpreter.intPool.get().Set(evm.GasPrice))
	return nil, nil
}

func opBlockhash(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[46899]++
	num := stack.pop()

	n := evm.interpreter.intPool.get().Sub(evm.BlockNumber, common.Big257)
	if num.Cmp(n) > 0 && num.Cmp(evm.BlockNumber) < 0 {
		fuzz_helper.CoverTab[10840]++
		stack.push(evm.GetHash(num.Uint64()).Big())
	} else {
		fuzz_helper.CoverTab[43066]++
		stack.push(new(big.Int))
	}
	fuzz_helper.CoverTab[40738]++

	evm.interpreter.intPool.put(num, n)
	return nil, nil
}

func opCoinbase(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[14816]++
	stack.push(evm.Coinbase.Big())
	return nil, nil
}

func opTimestamp(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[12109]++
	stack.push(math.U256(new(big.Int).Set(evm.Time)))
	return nil, nil
}

func opNumber(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[29966]++
	stack.push(math.U256(new(big.Int).Set(evm.BlockNumber)))
	return nil, nil
}

func opDifficulty(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[27153]++
	stack.push(math.U256(new(big.Int).Set(evm.Difficulty)))
	return nil, nil
}

func opGasLimit(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[51386]++
	stack.push(math.U256(new(big.Int).Set(evm.GasLimit)))
	return nil, nil
}

func opPop(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[4676]++
	evm.interpreter.intPool.put(stack.pop())
	return nil, nil
}

func opMload(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[61379]++
	offset := stack.pop()
	val := new(big.Int).SetBytes(memory.Get(offset.Int64(), 32))
	stack.push(val)

	evm.interpreter.intPool.put(offset)
	return nil, nil
}

func opMstore(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[16873]++

	mStart, val := stack.pop(), stack.pop()
	memory.Set(mStart.Uint64(), 32, math.PaddedBigBytes(val, 32))

	evm.interpreter.intPool.put(mStart, val)
	return nil, nil
}

func opMstore8(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[20099]++
	off, val := stack.pop().Int64(), stack.pop().Int64()
	memory.store[off] = byte(val & 0xff)

	return nil, nil
}

func opSload(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[61712]++
	loc := common.BigToHash(stack.pop())
	val := evm.StateDB.GetState(contract.Address(), loc).Big()
	stack.push(val)
	return nil, nil
}

func opSstore(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[12762]++
	loc := common.BigToHash(stack.pop())
	val := stack.pop()
	evm.StateDB.SetState(contract.Address(), loc, common.BigToHash(val))

	evm.interpreter.intPool.put(val)
	return nil, nil
}

func opJump(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[17951]++
	pos := stack.pop()
	if !contract.jumpdests.has(contract.CodeHash, contract.Code, pos) {
		fuzz_helper.CoverTab[64454]++
		nop := contract.GetOp(pos.Uint64())
		return nil, fmt.Errorf("invalid jump destination (%v) %v", nop, pos)
	} else {
		fuzz_helper.CoverTab[7160]++
	}
	fuzz_helper.CoverTab[53729]++
	*pc = pos.Uint64()

	evm.interpreter.intPool.put(pos)
	return nil, nil
}

func opJumpi(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[7189]++
	pos, cond := stack.pop(), stack.pop()
	if cond.Sign() != 0 {
		fuzz_helper.CoverTab[55434]++
		if !contract.jumpdests.has(contract.CodeHash, contract.Code, pos) {
			fuzz_helper.CoverTab[51703]++
			nop := contract.GetOp(pos.Uint64())
			return nil, fmt.Errorf("invalid jump destination (%v) %v", nop, pos)
		} else {
			fuzz_helper.CoverTab[39545]++
		}
		fuzz_helper.CoverTab[35217]++
		*pc = pos.Uint64()
	} else {
		fuzz_helper.CoverTab[19417]++
		*pc++
	}
	fuzz_helper.CoverTab[54641]++

	evm.interpreter.intPool.put(pos, cond)
	return nil, nil
}

func opJumpdest(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[47979]++
	return nil, nil
}

func opPc(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[22809]++
	stack.push(evm.interpreter.intPool.get().SetUint64(*pc))
	return nil, nil
}

func opMsize(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[20851]++
	stack.push(evm.interpreter.intPool.get().SetInt64(int64(memory.Len())))
	return nil, nil
}

func opGas(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[29885]++
	stack.push(evm.interpreter.intPool.get().SetUint64(contract.Gas))
	return nil, nil
}

func opCreate(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[17574]++
	var (
		value        = stack.pop()
		offset, size = stack.pop(), stack.pop()
		input        = memory.Get(offset.Int64(), size.Int64())
		gas          = contract.Gas
	)
	if evm.ChainConfig().IsEIP150(evm.BlockNumber) {
		fuzz_helper.CoverTab[56379]++
		gas -= gas / 64
	} else {
		fuzz_helper.CoverTab[7887]++
	}
	fuzz_helper.CoverTab[29431]++

	contract.UseGas(gas)
	res, addr, returnGas, suberr := evm.Create(contract, input, gas, value)

	if evm.ChainConfig().IsHomestead(evm.BlockNumber) && suberr == ErrCodeStoreOutOfGas {
		fuzz_helper.CoverTab[33884]++
		stack.push(new(big.Int))
	} else {
		fuzz_helper.CoverTab[49673]++
		if suberr != nil && suberr != ErrCodeStoreOutOfGas {
			fuzz_helper.CoverTab[4765]++
			stack.push(new(big.Int))
		} else {
			fuzz_helper.CoverTab[30464]++
			stack.push(addr.Big())
		}
	}
	fuzz_helper.CoverTab[31201]++
	contract.Gas += returnGas
	evm.interpreter.intPool.put(value, offset, size)

	if suberr == errExecutionReverted {
		fuzz_helper.CoverTab[62973]++
		return res, nil
	} else {
		fuzz_helper.CoverTab[21243]++
	}
	fuzz_helper.CoverTab[25315]++
	return nil, nil
}

func opCall(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[29150]++
	gas := stack.pop().Uint64()

	addr, value := stack.pop(), stack.pop()
	value = math.U256(value)

	inOffset, inSize := stack.pop(), stack.pop()

	retOffset, retSize := stack.pop(), stack.pop()

	address := common.BigToAddress(addr)

	args := memory.Get(inOffset.Int64(), inSize.Int64())

	if value.Sign() != 0 {
		fuzz_helper.CoverTab[28688]++
		gas += params.CallStipend
	} else {
		fuzz_helper.CoverTab[16762]++
	}
	fuzz_helper.CoverTab[26509]++
	ret, returnGas, err := evm.Call(contract, address, args, gas, value)
	if err != nil {
		fuzz_helper.CoverTab[60446]++
		stack.push(new(big.Int))
	} else {
		fuzz_helper.CoverTab[58133]++
		stack.push(big.NewInt(1))
	}
	fuzz_helper.CoverTab[62902]++
	if err == nil || err == errExecutionReverted {
		fuzz_helper.CoverTab[20790]++
		memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
	} else {
		fuzz_helper.CoverTab[63383]++
	}
	fuzz_helper.CoverTab[2802]++
	contract.Gas += returnGas

	evm.interpreter.intPool.put(addr, value, inOffset, inSize, retOffset, retSize)
	return ret, nil
}

func opCallCode(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[58404]++
	gas := stack.pop().Uint64()

	addr, value := stack.pop(), stack.pop()
	value = math.U256(value)

	inOffset, inSize := stack.pop(), stack.pop()

	retOffset, retSize := stack.pop(), stack.pop()

	address := common.BigToAddress(addr)

	args := memory.Get(inOffset.Int64(), inSize.Int64())

	if value.Sign() != 0 {
		fuzz_helper.CoverTab[55641]++
		gas += params.CallStipend
	} else {
		fuzz_helper.CoverTab[11885]++
	}
	fuzz_helper.CoverTab[54412]++

	ret, returnGas, err := evm.CallCode(contract, address, args, gas, value)
	if err != nil {
		fuzz_helper.CoverTab[60640]++
		stack.push(new(big.Int))
	} else {
		fuzz_helper.CoverTab[32340]++
		stack.push(big.NewInt(1))
	}
	fuzz_helper.CoverTab[62727]++
	if err == nil || err == errExecutionReverted {
		fuzz_helper.CoverTab[49446]++
		memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
	} else {
		fuzz_helper.CoverTab[55378]++
	}
	fuzz_helper.CoverTab[29728]++
	contract.Gas += returnGas

	evm.interpreter.intPool.put(addr, value, inOffset, inSize, retOffset, retSize)
	return ret, nil
}

func opDelegateCall(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[45742]++
	gas, to, inOffset, inSize, outOffset, outSize := stack.pop().Uint64(), stack.pop(), stack.pop(), stack.pop(), stack.pop(), stack.pop()

	toAddr := common.BigToAddress(to)
	args := memory.Get(inOffset.Int64(), inSize.Int64())

	ret, returnGas, err := evm.DelegateCall(contract, toAddr, args, gas)
	if err != nil {
		fuzz_helper.CoverTab[1596]++
		stack.push(new(big.Int))
	} else {
		fuzz_helper.CoverTab[48751]++
		stack.push(big.NewInt(1))
	}
	fuzz_helper.CoverTab[58797]++
	if err == nil || err == errExecutionReverted {
		fuzz_helper.CoverTab[27607]++
		memory.Set(outOffset.Uint64(), outSize.Uint64(), ret)
	} else {
		fuzz_helper.CoverTab[60004]++
	}
	fuzz_helper.CoverTab[45308]++
	contract.Gas += returnGas

	evm.interpreter.intPool.put(to, inOffset, inSize, outOffset, outSize)
	return ret, nil
}

func opStaticCall(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[13056]++

	gas := stack.pop().Uint64()

	addr := stack.pop()

	inOffset, inSize := stack.pop(), stack.pop()

	retOffset, retSize := stack.pop(), stack.pop()

	address := common.BigToAddress(addr)

	args := memory.Get(inOffset.Int64(), inSize.Int64())

	ret, returnGas, err := evm.StaticCall(contract, address, args, gas)
	if err != nil {
		fuzz_helper.CoverTab[47390]++
		stack.push(new(big.Int))
	} else {
		fuzz_helper.CoverTab[47458]++
		stack.push(big.NewInt(1))
	}
	fuzz_helper.CoverTab[12285]++
	if err == nil || err == errExecutionReverted {
		fuzz_helper.CoverTab[10614]++
		memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
	} else {
		fuzz_helper.CoverTab[9340]++
	}
	fuzz_helper.CoverTab[10163]++
	contract.Gas += returnGas

	evm.interpreter.intPool.put(addr, inOffset, inSize, retOffset, retSize)
	return ret, nil
}

func opReturn(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[25744]++
	offset, size := stack.pop(), stack.pop()
	ret := memory.GetPtr(offset.Int64(), size.Int64())

	evm.interpreter.intPool.put(offset, size)
	return ret, nil
}

func opRevert(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[65311]++
	offset, size := stack.pop(), stack.pop()
	ret := memory.GetPtr(offset.Int64(), size.Int64())

	evm.interpreter.intPool.put(offset, size)
	return ret, nil
}

func opStop(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[19813]++
	return nil, nil
}

func opSuicide(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	fuzz_helper.CoverTab[24779]++
	balance := evm.StateDB.GetBalance(contract.Address())
	evm.StateDB.AddBalance(common.BigToAddress(stack.pop()), balance)

	evm.StateDB.Suicide(contract.Address())
	return nil, nil
}

// make log instruction function
func makeLog(size int) executionFunc {
	fuzz_helper.CoverTab[179]++
	return func(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
		fuzz_helper.CoverTab[4752]++
		topics := make([]common.Hash, size)
		mStart, mSize := stack.pop(), stack.pop()
		for i := 0; i < size; i++ {
			fuzz_helper.CoverTab[60128]++
			topics[i] = common.BigToHash(stack.pop())
		}
		fuzz_helper.CoverTab[58259]++

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
	fuzz_helper.CoverTab[62145]++
	return func(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
		fuzz_helper.CoverTab[12227]++
		codeLen := len(contract.Code)

		startMin := codeLen
		if int(*pc+1) < startMin {
			fuzz_helper.CoverTab[30484]++
			startMin = int(*pc + 1)
		} else {
			fuzz_helper.CoverTab[40730]++
		}
		fuzz_helper.CoverTab[10795]++

		endMin := codeLen
		if startMin+pushByteSize < endMin {
			fuzz_helper.CoverTab[60454]++
			endMin = startMin + pushByteSize
		} else {
			fuzz_helper.CoverTab[38671]++
		}
		fuzz_helper.CoverTab[53895]++

		integer := evm.interpreter.intPool.get()
		stack.push(integer.SetBytes(common.RightPadBytes(contract.Code[startMin:endMin], pushByteSize)))

		*pc += size
		return nil, nil
	}
}

// make push instruction function
func makeDup(size int64) executionFunc {
	fuzz_helper.CoverTab[55665]++
	return func(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
		fuzz_helper.CoverTab[57927]++
		stack.dup(evm.interpreter.intPool, int(size))
		return nil, nil
	}
}

// make swap instruction function
func makeSwap(size int64) executionFunc {
	fuzz_helper.CoverTab[65356]++

	size += 1
	return func(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
		fuzz_helper.CoverTab[8275]++
		stack.swap(int(size))
		return nil, nil
	}
}

var _ = fuzz_helper.CoverTab
