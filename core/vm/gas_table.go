package vm

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/params"
)

// memoryGasCosts calculates the quadratic gas for memory expansion. It does so
// only for the memory region that is expanded, not the total memory.
func memoryGasCost(mem *Memory, newMemSize uint64) (uint64, error) {
	fuzz_helper.CoverTab[22588]++

	if newMemSize == 0 {
		fuzz_helper.CoverTab[45021]++
		return 0, nil
	} else {
		fuzz_helper.CoverTab[39040]++
	}
	fuzz_helper.CoverTab[44810]++

	if newMemSize > 0xffffffffe0 {
		fuzz_helper.CoverTab[2095]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[21668]++
	}
	fuzz_helper.CoverTab[5262]++

	newMemSizeWords := toWordSize(newMemSize)
	newMemSize = newMemSizeWords * 32

	if newMemSize > uint64(mem.Len()) {
		fuzz_helper.CoverTab[45213]++
		square := newMemSizeWords * newMemSizeWords
		linCoef := newMemSizeWords * params.MemoryGas
		quadCoef := square / params.QuadCoeffDiv
		newTotalFee := linCoef + quadCoef

		fee := newTotalFee - mem.lastGasCost
		mem.lastGasCost = newTotalFee

		return fee, nil
	} else {
		fuzz_helper.CoverTab[16619]++
	}
	fuzz_helper.CoverTab[17878]++
	return 0, nil
}

func constGasFunc(gas uint64) gasFunc {
	fuzz_helper.CoverTab[12692]++
	return func(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
		fuzz_helper.CoverTab[42483]++
		return gas, nil
	}
}

func gasCallDataCopy(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.CoverTab[6577]++
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.CoverTab[23294]++
		return 0, err
	} else {
		fuzz_helper.CoverTab[61639]++
	}
	fuzz_helper.CoverTab[17393]++

	var overflow bool
	if gas, overflow = math.SafeAdd(gas, GasFastestStep); overflow {
		fuzz_helper.CoverTab[11162]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[49217]++
	}
	fuzz_helper.CoverTab[64174]++

	words, overflow := bigUint64(stack.Back(2))
	if overflow {
		fuzz_helper.CoverTab[34511]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[64074]++
	}
	fuzz_helper.CoverTab[38740]++

	if words, overflow = math.SafeMul(toWordSize(words), params.CopyGas); overflow {
		fuzz_helper.CoverTab[28614]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[39226]++
	}
	fuzz_helper.CoverTab[35657]++

	if gas, overflow = math.SafeAdd(gas, words); overflow {
		fuzz_helper.CoverTab[2297]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[40870]++
	}
	fuzz_helper.CoverTab[30358]++
	return gas, nil
}

func gasReturnDataCopy(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.CoverTab[52877]++
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.CoverTab[12901]++
		return 0, err
	} else {
		fuzz_helper.CoverTab[12499]++
	}
	fuzz_helper.CoverTab[778]++

	var overflow bool
	if gas, overflow = math.SafeAdd(gas, GasFastestStep); overflow {
		fuzz_helper.CoverTab[42993]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[30301]++
	}
	fuzz_helper.CoverTab[33340]++

	words, overflow := bigUint64(stack.Back(2))
	if overflow {
		fuzz_helper.CoverTab[45210]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[264]++
	}
	fuzz_helper.CoverTab[15638]++

	if words, overflow = math.SafeMul(toWordSize(words), params.CopyGas); overflow {
		fuzz_helper.CoverTab[3566]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[47636]++
	}
	fuzz_helper.CoverTab[45869]++

	if gas, overflow = math.SafeAdd(gas, words); overflow {
		fuzz_helper.CoverTab[8730]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[20539]++
	}
	fuzz_helper.CoverTab[23368]++
	return gas, nil
}

func gasSStore(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.CoverTab[63931]++
	var (
		y, x = stack.Back(1), stack.Back(0)
		val  = evm.StateDB.GetState(contract.Address(), common.BigToHash(x))
	)

	if common.EmptyHash(val) && !common.EmptyHash(common.BigToHash(y)) {
		fuzz_helper.CoverTab[19009]++

		return params.SstoreSetGas, nil
	} else {
		fuzz_helper.CoverTab[64748]++
		if !common.EmptyHash(val) && common.EmptyHash(common.BigToHash(y)) {
			fuzz_helper.CoverTab[50446]++
			evm.StateDB.AddRefund(new(big.Int).SetUint64(params.SstoreRefundGas))

			return params.SstoreClearGas, nil
		} else {
			fuzz_helper.CoverTab[18500]++

			return params.SstoreResetGas, nil
		}
	}
}

func makeGasLog(n uint64) gasFunc {
	fuzz_helper.CoverTab[52152]++
	return func(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
		fuzz_helper.CoverTab[17111]++
		requestedSize, overflow := bigUint64(stack.Back(1))
		if overflow {
			fuzz_helper.CoverTab[17300]++
			return 0, errGasUintOverflow
		} else {
			fuzz_helper.CoverTab[16403]++
		}
		fuzz_helper.CoverTab[9670]++

		gas, err := memoryGasCost(mem, memorySize)
		if err != nil {
			fuzz_helper.CoverTab[40937]++
			return 0, err
		} else {
			fuzz_helper.CoverTab[33825]++
		}
		fuzz_helper.CoverTab[55848]++

		if gas, overflow = math.SafeAdd(gas, params.LogGas); overflow {
			fuzz_helper.CoverTab[7237]++
			return 0, errGasUintOverflow
		} else {
			fuzz_helper.CoverTab[23248]++
		}
		fuzz_helper.CoverTab[50755]++
		if gas, overflow = math.SafeAdd(gas, n*params.LogTopicGas); overflow {
			fuzz_helper.CoverTab[52715]++
			return 0, errGasUintOverflow
		} else {
			fuzz_helper.CoverTab[11389]++
		}
		fuzz_helper.CoverTab[912]++

		var memorySizeGas uint64
		if memorySizeGas, overflow = math.SafeMul(requestedSize, params.LogDataGas); overflow {
			fuzz_helper.CoverTab[60629]++
			return 0, errGasUintOverflow
		} else {
			fuzz_helper.CoverTab[23245]++
		}
		fuzz_helper.CoverTab[64631]++
		if gas, overflow = math.SafeAdd(gas, memorySizeGas); overflow {
			fuzz_helper.CoverTab[5383]++
			return 0, errGasUintOverflow
		} else {
			fuzz_helper.CoverTab[52957]++
		}
		fuzz_helper.CoverTab[15513]++
		return gas, nil
	}
}

func gasSha3(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.CoverTab[6211]++
	var overflow bool
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.CoverTab[24978]++
		return 0, err
	} else {
		fuzz_helper.CoverTab[61755]++
	}
	fuzz_helper.CoverTab[49245]++

	if gas, overflow = math.SafeAdd(gas, params.Sha3Gas); overflow {
		fuzz_helper.CoverTab[19607]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[28743]++
	}
	fuzz_helper.CoverTab[15785]++

	wordGas, overflow := bigUint64(stack.Back(1))
	if overflow {
		fuzz_helper.CoverTab[8832]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[40052]++
	}
	fuzz_helper.CoverTab[9735]++
	if wordGas, overflow = math.SafeMul(toWordSize(wordGas), params.Sha3WordGas); overflow {
		fuzz_helper.CoverTab[63449]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[47485]++
	}
	fuzz_helper.CoverTab[45823]++
	if gas, overflow = math.SafeAdd(gas, wordGas); overflow {
		fuzz_helper.CoverTab[60075]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[21817]++
	}
	fuzz_helper.CoverTab[48647]++
	return gas, nil
}

func gasCodeCopy(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.CoverTab[57682]++
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.CoverTab[56274]++
		return 0, err
	} else {
		fuzz_helper.CoverTab[1404]++
	}
	fuzz_helper.CoverTab[45496]++

	var overflow bool
	if gas, overflow = math.SafeAdd(gas, GasFastestStep); overflow {
		fuzz_helper.CoverTab[34815]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[58334]++
	}
	fuzz_helper.CoverTab[3661]++

	wordGas, overflow := bigUint64(stack.Back(2))
	if overflow {
		fuzz_helper.CoverTab[46899]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[40738]++
	}
	fuzz_helper.CoverTab[22210]++
	if wordGas, overflow = math.SafeMul(toWordSize(wordGas), params.CopyGas); overflow {
		fuzz_helper.CoverTab[10840]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[43066]++
	}
	fuzz_helper.CoverTab[4417]++
	if gas, overflow = math.SafeAdd(gas, wordGas); overflow {
		fuzz_helper.CoverTab[14816]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[12109]++
	}
	fuzz_helper.CoverTab[5093]++
	return gas, nil
}

func gasExtCodeCopy(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.CoverTab[29966]++
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.CoverTab[20099]++
		return 0, err
	} else {
		fuzz_helper.CoverTab[61712]++
	}
	fuzz_helper.CoverTab[27153]++

	var overflow bool
	if gas, overflow = math.SafeAdd(gas, gt.ExtcodeCopy); overflow {
		fuzz_helper.CoverTab[12762]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[17951]++
	}
	fuzz_helper.CoverTab[51386]++

	wordGas, overflow := bigUint64(stack.Back(3))
	if overflow {
		fuzz_helper.CoverTab[53729]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[64454]++
	}
	fuzz_helper.CoverTab[4676]++

	if wordGas, overflow = math.SafeMul(toWordSize(wordGas), params.CopyGas); overflow {
		fuzz_helper.CoverTab[7160]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[7189]++
	}
	fuzz_helper.CoverTab[61379]++

	if gas, overflow = math.SafeAdd(gas, wordGas); overflow {
		fuzz_helper.CoverTab[54641]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[55434]++
	}
	fuzz_helper.CoverTab[16873]++
	return gas, nil
}

func gasMLoad(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.CoverTab[35217]++
	var overflow bool
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.CoverTab[19417]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[47979]++
	}
	fuzz_helper.CoverTab[51703]++
	if gas, overflow = math.SafeAdd(gas, GasFastestStep); overflow {
		fuzz_helper.CoverTab[22809]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[20851]++
	}
	fuzz_helper.CoverTab[39545]++
	return gas, nil
}

func gasMStore8(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.CoverTab[29885]++
	var overflow bool
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.CoverTab[31201]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[25315]++
	}
	fuzz_helper.CoverTab[17574]++
	if gas, overflow = math.SafeAdd(gas, GasFastestStep); overflow {
		fuzz_helper.CoverTab[56379]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[7887]++
	}
	fuzz_helper.CoverTab[29431]++
	return gas, nil
}

func gasMStore(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.CoverTab[33884]++
	var overflow bool
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.CoverTab[30464]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[62973]++
	}
	fuzz_helper.CoverTab[49673]++
	if gas, overflow = math.SafeAdd(gas, GasFastestStep); overflow {
		fuzz_helper.CoverTab[21243]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[29150]++
	}
	fuzz_helper.CoverTab[4765]++
	return gas, nil
}

func gasCreate(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.CoverTab[26509]++
	var overflow bool
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.CoverTab[28688]++
		return 0, err
	} else {
		fuzz_helper.CoverTab[16762]++
	}
	fuzz_helper.CoverTab[62902]++
	if gas, overflow = math.SafeAdd(gas, params.CreateGas); overflow {
		fuzz_helper.CoverTab[60446]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[58133]++
	}
	fuzz_helper.CoverTab[2802]++
	return gas, nil
}

func gasBalance(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.CoverTab[20790]++
	return gt.Balance, nil
}

func gasExtCodeSize(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.CoverTab[63383]++
	return gt.ExtcodeSize, nil
}

func gasSLoad(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.CoverTab[58404]++
	return gt.SLoad, nil
}

func gasExp(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.CoverTab[54412]++
	expByteLen := uint64((stack.data[stack.len()-2].BitLen() + 7) / 8)

	var (
		gas      = expByteLen * gt.ExpByte // no overflow check required. Max is 256 * ExpByte gas
		overflow bool
	)
	if gas, overflow = math.SafeAdd(gas, GasSlowStep); overflow {
		fuzz_helper.CoverTab[29728]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[55641]++
	}
	fuzz_helper.CoverTab[62727]++
	return gas, nil
}

func gasCall(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.CoverTab[11885]++
	var (
		gas            = gt.Calls
		transfersValue = stack.Back(2).Sign() != 0
		address        = common.BigToAddress(stack.Back(1))
		eip158         = evm.ChainConfig().IsEIP158(evm.BlockNumber)
	)
	if eip158 {
		fuzz_helper.CoverTab[45308]++
		if transfersValue && evm.StateDB.Empty(address) {
			fuzz_helper.CoverTab[1596]++
			gas += params.CallNewAccountGas
		} else {
			fuzz_helper.CoverTab[48751]++
		}
	} else {
		fuzz_helper.CoverTab[27607]++
		if !evm.StateDB.Exist(address) {
			fuzz_helper.CoverTab[60004]++
			gas += params.CallNewAccountGas
		} else {
			fuzz_helper.CoverTab[13056]++
		}
	}
	fuzz_helper.CoverTab[60640]++
	if transfersValue {
		fuzz_helper.CoverTab[12285]++
		gas += params.CallValueTransferGas
	} else {
		fuzz_helper.CoverTab[10163]++
	}
	fuzz_helper.CoverTab[32340]++
	memoryGas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.CoverTab[47390]++
		return 0, err
	} else {
		fuzz_helper.CoverTab[47458]++
	}
	fuzz_helper.CoverTab[49446]++
	var overflow bool
	if gas, overflow = math.SafeAdd(gas, memoryGas); overflow {
		fuzz_helper.CoverTab[10614]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[9340]++
	}
	fuzz_helper.CoverTab[55378]++

	cg, err := callGas(gt, contract.Gas, gas, stack.Back(0))
	if err != nil {
		fuzz_helper.CoverTab[25744]++
		return 0, err
	} else {
		fuzz_helper.CoverTab[65311]++
	}
	fuzz_helper.CoverTab[45742]++

	stack.data[stack.len()-1] = new(big.Int).SetUint64(cg)

	if gas, overflow = math.SafeAdd(gas, cg); overflow {
		fuzz_helper.CoverTab[19813]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[24779]++
	}
	fuzz_helper.CoverTab[58797]++
	return gas, nil
}

func gasCallCode(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.CoverTab[179]++
	gas := gt.Calls
	if stack.Back(2).Sign() != 0 {
		fuzz_helper.CoverTab[10795]++
		gas += params.CallValueTransferGas
	} else {
		fuzz_helper.CoverTab[53895]++
	}
	fuzz_helper.CoverTab[4752]++
	memoryGas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.CoverTab[30484]++
		return 0, err
	} else {
		fuzz_helper.CoverTab[40730]++
	}
	fuzz_helper.CoverTab[58259]++
	var overflow bool
	if gas, overflow = math.SafeAdd(gas, memoryGas); overflow {
		fuzz_helper.CoverTab[60454]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[38671]++
	}
	fuzz_helper.CoverTab[60128]++

	cg, err := callGas(gt, contract.Gas, gas, stack.Back(0))
	if err != nil {
		fuzz_helper.CoverTab[55665]++
		return 0, err
	} else {
		fuzz_helper.CoverTab[57927]++
	}
	fuzz_helper.CoverTab[62145]++

	stack.data[stack.len()-1] = new(big.Int).SetUint64(cg)

	if gas, overflow = math.SafeAdd(gas, cg); overflow {
		fuzz_helper.CoverTab[65356]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[8275]++
	}
	fuzz_helper.CoverTab[12227]++
	return gas, nil
}

func gasReturn(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.CoverTab[41896]++
	return memoryGasCost(mem, memorySize)
}

func gasRevert(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.CoverTab[5111]++
	return memoryGasCost(mem, memorySize)
}

func gasSuicide(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.CoverTab[2960]++
	var gas uint64

	if evm.ChainConfig().IsEIP150(evm.BlockNumber) {
		fuzz_helper.CoverTab[25547]++
		gas = gt.Suicide
		var (
			address = common.BigToAddress(stack.Back(0))
			eip158  = evm.ChainConfig().IsEIP158(evm.BlockNumber)
		)

		if eip158 {
			fuzz_helper.CoverTab[6705]++

			if evm.StateDB.Empty(address) && evm.StateDB.GetBalance(contract.Address()).Sign() != 0 {
				fuzz_helper.CoverTab[12502]++
				gas += gt.CreateBySuicide
			} else {
				fuzz_helper.CoverTab[5209]++
			}
		} else {
			fuzz_helper.CoverTab[3613]++
			if !evm.StateDB.Exist(address) {
				fuzz_helper.CoverTab[9728]++
				gas += gt.CreateBySuicide
			} else {
				fuzz_helper.CoverTab[20787]++
			}
		}
	} else {
		fuzz_helper.CoverTab[14592]++
	}
	fuzz_helper.CoverTab[35072]++

	if !evm.StateDB.HasSuicided(contract.Address()) {
		fuzz_helper.CoverTab[1274]++
		evm.StateDB.AddRefund(new(big.Int).SetUint64(params.SuicideRefundGas))
	} else {
		fuzz_helper.CoverTab[56143]++
	}
	fuzz_helper.CoverTab[55654]++
	return gas, nil
}

func gasDelegateCall(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.CoverTab[17283]++
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.CoverTab[20593]++
		return 0, err
	} else {
		fuzz_helper.CoverTab[26119]++
	}
	fuzz_helper.CoverTab[11455]++
	var overflow bool
	if gas, overflow = math.SafeAdd(gas, gt.Calls); overflow {
		fuzz_helper.CoverTab[38633]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[27520]++
	}
	fuzz_helper.CoverTab[38950]++

	cg, err := callGas(gt, contract.Gas, gas, stack.Back(0))
	if err != nil {
		fuzz_helper.CoverTab[54694]++
		return 0, err
	} else {
		fuzz_helper.CoverTab[7779]++
	}
	fuzz_helper.CoverTab[65052]++

	stack.data[stack.len()-1] = new(big.Int).SetUint64(cg)

	if gas, overflow = math.SafeAdd(gas, cg); overflow {
		fuzz_helper.CoverTab[50594]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[53346]++
	}
	fuzz_helper.CoverTab[35391]++
	return gas, nil
}

func gasStaticCall(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.CoverTab[38099]++
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.CoverTab[55029]++
		return 0, err
	} else {
		fuzz_helper.CoverTab[24853]++
	}
	fuzz_helper.CoverTab[28290]++
	var overflow bool
	if gas, overflow = math.SafeAdd(gas, gt.Calls); overflow {
		fuzz_helper.CoverTab[8]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[39569]++
	}
	fuzz_helper.CoverTab[46467]++

	cg, err := callGas(gt, contract.Gas, gas, stack.Back(0))
	if err != nil {
		fuzz_helper.CoverTab[14681]++
		return 0, err
	} else {
		fuzz_helper.CoverTab[32228]++
	}
	fuzz_helper.CoverTab[2463]++

	stack.data[stack.len()-1] = new(big.Int).SetUint64(cg)

	if gas, overflow = math.SafeAdd(gas, cg); overflow {
		fuzz_helper.CoverTab[42372]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[15925]++
	}
	fuzz_helper.CoverTab[41764]++
	return gas, nil
}

func gasPush(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.CoverTab[61802]++
	return GasFastestStep, nil
}

func gasSwap(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.CoverTab[63150]++
	return GasFastestStep, nil
}

func gasDup(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.CoverTab[26030]++
	return GasFastestStep, nil
}

var _ = fuzz_helper.CoverTab
