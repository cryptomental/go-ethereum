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
	fuzz_helper.AddCoverage(2940)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()

	if newMemSize == 0 {
		fuzz_helper.AddCoverage(453)
		return 0, nil
	} else {
		fuzz_helper.AddCoverage(1360)
	}
	fuzz_helper.AddCoverage(13966)

	if newMemSize > 0xffffffffe0 {
		fuzz_helper.AddCoverage(323)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(31721)
	}
	fuzz_helper.AddCoverage(32830)

	newMemSizeWords := toWordSize(newMemSize)
	newMemSize = newMemSizeWords * 32

	if newMemSize > uint64(mem.Len()) {
		fuzz_helper.AddCoverage(46704)
		square := newMemSizeWords * newMemSizeWords
		linCoef := newMemSizeWords * params.MemoryGas
		quadCoef := square / params.QuadCoeffDiv
		newTotalFee := linCoef + quadCoef

		fee := newTotalFee - mem.lastGasCost
		mem.lastGasCost = newTotalFee

		return fee, nil
	} else {
		fuzz_helper.AddCoverage(45053)
	}
	fuzz_helper.AddCoverage(46491)
	return 0, nil
}

func constGasFunc(gas uint64) gasFunc {
	fuzz_helper.AddCoverage(7700)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return func(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
		fuzz_helper.AddCoverage(8681)
		return gas, nil
	}
}

func gasCallDataCopy(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(16069)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.AddCoverage(51199)
		return 0, err
	} else {
		fuzz_helper.AddCoverage(30834)
	}
	fuzz_helper.AddCoverage(3519)

	var overflow bool
	if gas, overflow = math.SafeAdd(gas, GasFastestStep); overflow {
		fuzz_helper.AddCoverage(17128)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(5569)
	}
	fuzz_helper.AddCoverage(2374)

	words, overflow := bigUint64(stack.Back(2))
	if overflow {
		fuzz_helper.AddCoverage(3611)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(40490)
	}
	fuzz_helper.AddCoverage(38762)

	if words, overflow = math.SafeMul(toWordSize(words), params.CopyGas); overflow {
		fuzz_helper.AddCoverage(42753)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(39187)
	}
	fuzz_helper.AddCoverage(30627)

	if gas, overflow = math.SafeAdd(gas, words); overflow {
		fuzz_helper.AddCoverage(8037)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(5640)
	}
	fuzz_helper.AddCoverage(11606)
	return gas, nil
}

func gasReturnDataCopy(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(52665)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.AddCoverage(6261)
		return 0, err
	} else {
		fuzz_helper.AddCoverage(22680)
	}
	fuzz_helper.AddCoverage(36951)

	var overflow bool
	if gas, overflow = math.SafeAdd(gas, GasFastestStep); overflow {
		fuzz_helper.AddCoverage(9542)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(25302)
	}
	fuzz_helper.AddCoverage(20517)

	words, overflow := bigUint64(stack.Back(2))
	if overflow {
		fuzz_helper.AddCoverage(14997)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(17944)
	}
	fuzz_helper.AddCoverage(64951)

	if words, overflow = math.SafeMul(toWordSize(words), params.CopyGas); overflow {
		fuzz_helper.AddCoverage(32968)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(62502)
	}
	fuzz_helper.AddCoverage(6339)

	if gas, overflow = math.SafeAdd(gas, words); overflow {
		fuzz_helper.AddCoverage(32280)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(196)
	}
	fuzz_helper.AddCoverage(55742)
	return gas, nil
}

func gasSStore(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(40041)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	var (
		y, x = stack.Back(1), stack.Back(0)
		val  = evm.StateDB.GetState(contract.Address(), common.BigToHash(x))
	)

	if common.EmptyHash(val) && !common.EmptyHash(common.BigToHash(y)) {
		fuzz_helper.AddCoverage(24269)

		return params.SstoreSetGas, nil
	} else {
		fuzz_helper.AddCoverage(15925)
		if !common.EmptyHash(val) && common.EmptyHash(common.BigToHash(y)) {
			fuzz_helper.AddCoverage(4784)
			evm.StateDB.AddRefund(new(big.Int).SetUint64(params.SstoreRefundGas))

			return params.SstoreClearGas, nil
		} else {
			fuzz_helper.AddCoverage(57679)

			return params.SstoreResetGas, nil
		}
	}
}

func makeGasLog(n uint64) gasFunc {
	fuzz_helper.AddCoverage(26801)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return func(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
		fuzz_helper.AddCoverage(30782)
		requestedSize, overflow := bigUint64(stack.Back(1))
		if overflow {
			fuzz_helper.AddCoverage(52069)
			return 0, errGasUintOverflow
		} else {
			fuzz_helper.AddCoverage(64301)
		}
		fuzz_helper.AddCoverage(32228)

		gas, err := memoryGasCost(mem, memorySize)
		if err != nil {
			fuzz_helper.AddCoverage(17230)
			return 0, err
		} else {
			fuzz_helper.AddCoverage(36487)
		}
		fuzz_helper.AddCoverage(17680)

		if gas, overflow = math.SafeAdd(gas, params.LogGas); overflow {
			fuzz_helper.AddCoverage(42553)
			return 0, errGasUintOverflow
		} else {
			fuzz_helper.AddCoverage(52272)
		}
		fuzz_helper.AddCoverage(4655)
		if gas, overflow = math.SafeAdd(gas, n*params.LogTopicGas); overflow {
			fuzz_helper.AddCoverage(23495)
			return 0, errGasUintOverflow
		} else {
			fuzz_helper.AddCoverage(48542)
		}
		fuzz_helper.AddCoverage(6483)

		var memorySizeGas uint64
		if memorySizeGas, overflow = math.SafeMul(requestedSize, params.LogDataGas); overflow {
			fuzz_helper.AddCoverage(38511)
			return 0, errGasUintOverflow
		} else {
			fuzz_helper.AddCoverage(42614)
		}
		fuzz_helper.AddCoverage(27367)
		if gas, overflow = math.SafeAdd(gas, memorySizeGas); overflow {
			fuzz_helper.AddCoverage(42350)
			return 0, errGasUintOverflow
		} else {
			fuzz_helper.AddCoverage(54091)
		}
		fuzz_helper.AddCoverage(9538)
		return gas, nil
	}
}

func gasSha3(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(2452)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	var overflow bool
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.AddCoverage(35467)
		return 0, err
	} else {
		fuzz_helper.AddCoverage(39609)
	}
	fuzz_helper.AddCoverage(9768)

	if gas, overflow = math.SafeAdd(gas, params.Sha3Gas); overflow {
		fuzz_helper.AddCoverage(438)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(58861)
	}
	fuzz_helper.AddCoverage(57173)

	wordGas, overflow := bigUint64(stack.Back(1))
	if overflow {
		fuzz_helper.AddCoverage(58613)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(26710)
	}
	fuzz_helper.AddCoverage(48102)
	if wordGas, overflow = math.SafeMul(toWordSize(wordGas), params.Sha3WordGas); overflow {
		fuzz_helper.AddCoverage(51224)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(23908)
	}
	fuzz_helper.AddCoverage(56520)
	if gas, overflow = math.SafeAdd(gas, wordGas); overflow {
		fuzz_helper.AddCoverage(63440)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(64416)
	}
	fuzz_helper.AddCoverage(1729)
	return gas, nil
}

func gasCodeCopy(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(30881)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.AddCoverage(25220)
		return 0, err
	} else {
		fuzz_helper.AddCoverage(12290)
	}
	fuzz_helper.AddCoverage(63770)

	var overflow bool
	if gas, overflow = math.SafeAdd(gas, GasFastestStep); overflow {
		fuzz_helper.AddCoverage(25388)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(64747)
	}
	fuzz_helper.AddCoverage(40174)

	wordGas, overflow := bigUint64(stack.Back(2))
	if overflow {
		fuzz_helper.AddCoverage(62246)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(17910)
	}
	fuzz_helper.AddCoverage(54922)
	if wordGas, overflow = math.SafeMul(toWordSize(wordGas), params.CopyGas); overflow {
		fuzz_helper.AddCoverage(17467)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(4556)
	}
	fuzz_helper.AddCoverage(62100)
	if gas, overflow = math.SafeAdd(gas, wordGas); overflow {
		fuzz_helper.AddCoverage(24085)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(20769)
	}
	fuzz_helper.AddCoverage(21472)
	return gas, nil
}

func gasExtCodeCopy(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(21440)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.AddCoverage(4684)
		return 0, err
	} else {
		fuzz_helper.AddCoverage(41053)
	}
	fuzz_helper.AddCoverage(42282)

	var overflow bool
	if gas, overflow = math.SafeAdd(gas, gt.ExtcodeCopy); overflow {
		fuzz_helper.AddCoverage(19631)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(39133)
	}
	fuzz_helper.AddCoverage(43291)

	wordGas, overflow := bigUint64(stack.Back(3))
	if overflow {
		fuzz_helper.AddCoverage(47366)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(39547)
	}
	fuzz_helper.AddCoverage(4838)

	if wordGas, overflow = math.SafeMul(toWordSize(wordGas), params.CopyGas); overflow {
		fuzz_helper.AddCoverage(12115)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(3712)
	}
	fuzz_helper.AddCoverage(4804)

	if gas, overflow = math.SafeAdd(gas, wordGas); overflow {
		fuzz_helper.AddCoverage(62467)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(35608)
	}
	fuzz_helper.AddCoverage(33734)
	return gas, nil
}

func gasMLoad(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(29564)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	var overflow bool
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.AddCoverage(59732)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(60130)
	}
	fuzz_helper.AddCoverage(19470)
	if gas, overflow = math.SafeAdd(gas, GasFastestStep); overflow {
		fuzz_helper.AddCoverage(59032)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(3792)
	}
	fuzz_helper.AddCoverage(32533)
	return gas, nil
}

func gasMStore8(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(22517)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	var overflow bool
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.AddCoverage(62341)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(19573)
	}
	fuzz_helper.AddCoverage(59155)
	if gas, overflow = math.SafeAdd(gas, GasFastestStep); overflow {
		fuzz_helper.AddCoverage(3257)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(56741)
	}
	fuzz_helper.AddCoverage(46943)
	return gas, nil
}

func gasMStore(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(9951)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	var overflow bool
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.AddCoverage(41895)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(40834)
	}
	fuzz_helper.AddCoverage(40752)
	if gas, overflow = math.SafeAdd(gas, GasFastestStep); overflow {
		fuzz_helper.AddCoverage(34173)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(22765)
	}
	fuzz_helper.AddCoverage(6173)
	return gas, nil
}

func gasCreate(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(53607)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	var overflow bool
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.AddCoverage(55228)
		return 0, err
	} else {
		fuzz_helper.AddCoverage(26880)
	}
	fuzz_helper.AddCoverage(24881)
	if gas, overflow = math.SafeAdd(gas, params.CreateGas); overflow {
		fuzz_helper.AddCoverage(6142)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(55310)
	}
	fuzz_helper.AddCoverage(61261)
	return gas, nil
}

func gasBalance(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(32371)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return gt.Balance, nil
}

func gasExtCodeSize(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(39979)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return gt.ExtcodeSize, nil
}

func gasSLoad(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(22558)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return gt.SLoad, nil
}

func gasExp(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(24532)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	expByteLen := uint64((stack.data[stack.len()-2].BitLen() + 7) / 8)

	var (
		gas      = expByteLen * gt.ExpByte // no overflow check required. Max is 256 * ExpByte gas
		overflow bool
	)
	if gas, overflow = math.SafeAdd(gas, GasSlowStep); overflow {
		fuzz_helper.AddCoverage(56146)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(23537)
	}
	fuzz_helper.AddCoverage(64097)
	return gas, nil
}

func gasCall(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(59533)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	var (
		gas            = gt.Calls
		transfersValue = stack.Back(2).Sign() != 0
		address        = common.BigToAddress(stack.Back(1))
		eip158         = evm.ChainConfig().IsEIP158(evm.BlockNumber)
	)
	if eip158 {
		fuzz_helper.AddCoverage(11807)
		if transfersValue && evm.StateDB.Empty(address) {
			fuzz_helper.AddCoverage(42583)
			gas += params.CallNewAccountGas
		} else {
			fuzz_helper.AddCoverage(60449)
		}
	} else {
		fuzz_helper.AddCoverage(245)
		if !evm.StateDB.Exist(address) {
			fuzz_helper.AddCoverage(27634)
			gas += params.CallNewAccountGas
		} else {
			fuzz_helper.AddCoverage(59341)
		}
	}
	fuzz_helper.AddCoverage(4373)
	if transfersValue {
		fuzz_helper.AddCoverage(11294)
		gas += params.CallValueTransferGas
	} else {
		fuzz_helper.AddCoverage(49855)
	}
	fuzz_helper.AddCoverage(22157)
	memoryGas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.AddCoverage(64281)
		return 0, err
	} else {
		fuzz_helper.AddCoverage(64045)
	}
	fuzz_helper.AddCoverage(2937)
	var overflow bool
	if gas, overflow = math.SafeAdd(gas, memoryGas); overflow {
		fuzz_helper.AddCoverage(48676)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(60356)
	}
	fuzz_helper.AddCoverage(20807)

	cg, err := callGas(gt, contract.Gas, gas, stack.Back(0))
	if err != nil {
		fuzz_helper.AddCoverage(12397)
		return 0, err
	} else {
		fuzz_helper.AddCoverage(2437)
	}
	fuzz_helper.AddCoverage(15626)

	stack.data[stack.len()-1] = new(big.Int).SetUint64(cg)

	if gas, overflow = math.SafeAdd(gas, cg); overflow {
		fuzz_helper.AddCoverage(30100)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(7103)
	}
	fuzz_helper.AddCoverage(51516)
	return gas, nil
}

func gasCallCode(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(8732)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	gas := gt.Calls
	if stack.Back(2).Sign() != 0 {
		fuzz_helper.AddCoverage(4173)
		gas += params.CallValueTransferGas
	} else {
		fuzz_helper.AddCoverage(22384)
	}
	fuzz_helper.AddCoverage(61046)
	memoryGas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.AddCoverage(46319)
		return 0, err
	} else {
		fuzz_helper.AddCoverage(46581)
	}
	fuzz_helper.AddCoverage(22403)
	var overflow bool
	if gas, overflow = math.SafeAdd(gas, memoryGas); overflow {
		fuzz_helper.AddCoverage(28074)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(38451)
	}
	fuzz_helper.AddCoverage(4708)

	cg, err := callGas(gt, contract.Gas, gas, stack.Back(0))
	if err != nil {
		fuzz_helper.AddCoverage(53791)
		return 0, err
	} else {
		fuzz_helper.AddCoverage(15063)
	}
	fuzz_helper.AddCoverage(17692)

	stack.data[stack.len()-1] = new(big.Int).SetUint64(cg)

	if gas, overflow = math.SafeAdd(gas, cg); overflow {
		fuzz_helper.AddCoverage(36589)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(13747)
	}
	fuzz_helper.AddCoverage(4531)
	return gas, nil
}

func gasReturn(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(48379)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return memoryGasCost(mem, memorySize)
}

func gasRevert(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(7959)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return memoryGasCost(mem, memorySize)
}

func gasSuicide(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(20422)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	var gas uint64

	if evm.ChainConfig().IsEIP150(evm.BlockNumber) {
		fuzz_helper.AddCoverage(34874)
		gas = gt.Suicide
		var (
			address = common.BigToAddress(stack.Back(0))
			eip158  = evm.ChainConfig().IsEIP158(evm.BlockNumber)
		)

		if eip158 {
			fuzz_helper.AddCoverage(61284)

			if evm.StateDB.Empty(address) && evm.StateDB.GetBalance(contract.Address()).Sign() != 0 {
				fuzz_helper.AddCoverage(43732)
				gas += gt.CreateBySuicide
			} else {
				fuzz_helper.AddCoverage(20700)
			}
		} else {
			fuzz_helper.AddCoverage(35894)
			if !evm.StateDB.Exist(address) {
				fuzz_helper.AddCoverage(40322)
				gas += gt.CreateBySuicide
			} else {
				fuzz_helper.AddCoverage(61113)
			}
		}
	} else {
		fuzz_helper.AddCoverage(39411)
	}
	fuzz_helper.AddCoverage(49817)

	if !evm.StateDB.HasSuicided(contract.Address()) {
		fuzz_helper.AddCoverage(53301)
		evm.StateDB.AddRefund(new(big.Int).SetUint64(params.SuicideRefundGas))
	} else {
		fuzz_helper.AddCoverage(47061)
	}
	fuzz_helper.AddCoverage(55187)
	return gas, nil
}

func gasDelegateCall(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(46144)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.AddCoverage(50479)
		return 0, err
	} else {
		fuzz_helper.AddCoverage(1653)
	}
	fuzz_helper.AddCoverage(61708)
	var overflow bool
	if gas, overflow = math.SafeAdd(gas, gt.Calls); overflow {
		fuzz_helper.AddCoverage(36485)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(16021)
	}
	fuzz_helper.AddCoverage(63720)

	cg, err := callGas(gt, contract.Gas, gas, stack.Back(0))
	if err != nil {
		fuzz_helper.AddCoverage(44401)
		return 0, err
	} else {
		fuzz_helper.AddCoverage(6222)
	}
	fuzz_helper.AddCoverage(14112)

	stack.data[stack.len()-1] = new(big.Int).SetUint64(cg)

	if gas, overflow = math.SafeAdd(gas, cg); overflow {
		fuzz_helper.AddCoverage(44470)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(26441)
	}
	fuzz_helper.AddCoverage(37712)
	return gas, nil
}

func gasStaticCall(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(11728)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.AddCoverage(9524)
		return 0, err
	} else {
		fuzz_helper.AddCoverage(23690)
	}
	fuzz_helper.AddCoverage(63606)
	var overflow bool
	if gas, overflow = math.SafeAdd(gas, gt.Calls); overflow {
		fuzz_helper.AddCoverage(6117)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(22638)
	}
	fuzz_helper.AddCoverage(62902)

	cg, err := callGas(gt, contract.Gas, gas, stack.Back(0))
	if err != nil {
		fuzz_helper.AddCoverage(43179)
		return 0, err
	} else {
		fuzz_helper.AddCoverage(23276)
	}
	fuzz_helper.AddCoverage(27169)

	stack.data[stack.len()-1] = new(big.Int).SetUint64(cg)

	if gas, overflow = math.SafeAdd(gas, cg); overflow {
		fuzz_helper.AddCoverage(29679)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(32965)
	}
	fuzz_helper.AddCoverage(10317)
	return gas, nil
}

func gasPush(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(12275)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return GasFastestStep, nil
}

func gasSwap(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(57347)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return GasFastestStep, nil
}

func gasDup(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(55413)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return GasFastestStep, nil
}

var _ = fuzz_helper.AddCoverage
