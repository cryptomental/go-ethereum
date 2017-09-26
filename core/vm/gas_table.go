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
	fuzz_helper.AddCoverage(22588)

	if newMemSize == 0 {
		fuzz_helper.AddCoverage(45021)
		return 0, nil
	} else {
		fuzz_helper.AddCoverage(39040)
	}
	fuzz_helper.AddCoverage(44810)

	if newMemSize > 0xffffffffe0 {
		fuzz_helper.AddCoverage(2095)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(21668)
	}
	fuzz_helper.AddCoverage(5262)

	newMemSizeWords := toWordSize(newMemSize)
	newMemSize = newMemSizeWords * 32

	if newMemSize > uint64(mem.Len()) {
		fuzz_helper.AddCoverage(45213)
		square := newMemSizeWords * newMemSizeWords
		linCoef := newMemSizeWords * params.MemoryGas
		quadCoef := square / params.QuadCoeffDiv
		newTotalFee := linCoef + quadCoef

		fee := newTotalFee - mem.lastGasCost
		mem.lastGasCost = newTotalFee

		return fee, nil
	} else {
		fuzz_helper.AddCoverage(16619)
	}
	fuzz_helper.AddCoverage(17878)
	return 0, nil
}

func constGasFunc(gas uint64) gasFunc {
	fuzz_helper.AddCoverage(12692)
	return func(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
		fuzz_helper.AddCoverage(42483)
		return gas, nil
	}
}

func gasCallDataCopy(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(6577)
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.AddCoverage(23294)
		return 0, err
	} else {
		fuzz_helper.AddCoverage(61639)
	}
	fuzz_helper.AddCoverage(17393)

	var overflow bool
	if gas, overflow = math.SafeAdd(gas, GasFastestStep); overflow {
		fuzz_helper.AddCoverage(11162)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(49217)
	}
	fuzz_helper.AddCoverage(64174)

	words, overflow := bigUint64(stack.Back(2))
	if overflow {
		fuzz_helper.AddCoverage(34511)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(64074)
	}
	fuzz_helper.AddCoverage(38740)

	if words, overflow = math.SafeMul(toWordSize(words), params.CopyGas); overflow {
		fuzz_helper.AddCoverage(28614)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(39226)
	}
	fuzz_helper.AddCoverage(35657)

	if gas, overflow = math.SafeAdd(gas, words); overflow {
		fuzz_helper.AddCoverage(2297)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(40870)
	}
	fuzz_helper.AddCoverage(30358)
	return gas, nil
}

func gasReturnDataCopy(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(52877)
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.AddCoverage(12901)
		return 0, err
	} else {
		fuzz_helper.AddCoverage(12499)
	}
	fuzz_helper.AddCoverage(778)

	var overflow bool
	if gas, overflow = math.SafeAdd(gas, GasFastestStep); overflow {
		fuzz_helper.AddCoverage(42993)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(30301)
	}
	fuzz_helper.AddCoverage(33340)

	words, overflow := bigUint64(stack.Back(2))
	if overflow {
		fuzz_helper.AddCoverage(45210)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(264)
	}
	fuzz_helper.AddCoverage(15638)

	if words, overflow = math.SafeMul(toWordSize(words), params.CopyGas); overflow {
		fuzz_helper.AddCoverage(3566)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(47636)
	}
	fuzz_helper.AddCoverage(45869)

	if gas, overflow = math.SafeAdd(gas, words); overflow {
		fuzz_helper.AddCoverage(8730)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(20539)
	}
	fuzz_helper.AddCoverage(23368)
	return gas, nil
}

func gasSStore(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(63931)
	var (
		y, x = stack.Back(1), stack.Back(0)
		val  = evm.StateDB.GetState(contract.Address(), common.BigToHash(x))
	)

	if common.EmptyHash(val) && !common.EmptyHash(common.BigToHash(y)) {
		fuzz_helper.AddCoverage(19009)

		return params.SstoreSetGas, nil
	} else {
		fuzz_helper.AddCoverage(64748)
		if !common.EmptyHash(val) && common.EmptyHash(common.BigToHash(y)) {
			fuzz_helper.AddCoverage(50446)
			evm.StateDB.AddRefund(new(big.Int).SetUint64(params.SstoreRefundGas))

			return params.SstoreClearGas, nil
		} else {
			fuzz_helper.AddCoverage(18500)

			return params.SstoreResetGas, nil
		}
	}
}

func makeGasLog(n uint64) gasFunc {
	fuzz_helper.AddCoverage(52152)
	return func(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
		fuzz_helper.AddCoverage(17111)
		requestedSize, overflow := bigUint64(stack.Back(1))
		if overflow {
			fuzz_helper.AddCoverage(17300)
			return 0, errGasUintOverflow
		} else {
			fuzz_helper.AddCoverage(16403)
		}
		fuzz_helper.AddCoverage(9670)

		gas, err := memoryGasCost(mem, memorySize)
		if err != nil {
			fuzz_helper.AddCoverage(40937)
			return 0, err
		} else {
			fuzz_helper.AddCoverage(33825)
		}
		fuzz_helper.AddCoverage(55848)

		if gas, overflow = math.SafeAdd(gas, params.LogGas); overflow {
			fuzz_helper.AddCoverage(7237)
			return 0, errGasUintOverflow
		} else {
			fuzz_helper.AddCoverage(23248)
		}
		fuzz_helper.AddCoverage(50755)
		if gas, overflow = math.SafeAdd(gas, n*params.LogTopicGas); overflow {
			fuzz_helper.AddCoverage(52715)
			return 0, errGasUintOverflow
		} else {
			fuzz_helper.AddCoverage(11389)
		}
		fuzz_helper.AddCoverage(912)

		var memorySizeGas uint64
		if memorySizeGas, overflow = math.SafeMul(requestedSize, params.LogDataGas); overflow {
			fuzz_helper.AddCoverage(60629)
			return 0, errGasUintOverflow
		} else {
			fuzz_helper.AddCoverage(23245)
		}
		fuzz_helper.AddCoverage(64631)
		if gas, overflow = math.SafeAdd(gas, memorySizeGas); overflow {
			fuzz_helper.AddCoverage(5383)
			return 0, errGasUintOverflow
		} else {
			fuzz_helper.AddCoverage(52957)
		}
		fuzz_helper.AddCoverage(15513)
		return gas, nil
	}
}

func gasSha3(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(6211)
	var overflow bool
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.AddCoverage(24978)
		return 0, err
	} else {
		fuzz_helper.AddCoverage(61755)
	}
	fuzz_helper.AddCoverage(49245)

	if gas, overflow = math.SafeAdd(gas, params.Sha3Gas); overflow {
		fuzz_helper.AddCoverage(19607)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(28743)
	}
	fuzz_helper.AddCoverage(15785)

	wordGas, overflow := bigUint64(stack.Back(1))
	if overflow {
		fuzz_helper.AddCoverage(8832)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(40052)
	}
	fuzz_helper.AddCoverage(9735)
	if wordGas, overflow = math.SafeMul(toWordSize(wordGas), params.Sha3WordGas); overflow {
		fuzz_helper.AddCoverage(63449)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(47485)
	}
	fuzz_helper.AddCoverage(45823)
	if gas, overflow = math.SafeAdd(gas, wordGas); overflow {
		fuzz_helper.AddCoverage(60075)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(21817)
	}
	fuzz_helper.AddCoverage(48647)
	return gas, nil
}

func gasCodeCopy(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(57682)
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.AddCoverage(56274)
		return 0, err
	} else {
		fuzz_helper.AddCoverage(1404)
	}
	fuzz_helper.AddCoverage(45496)

	var overflow bool
	if gas, overflow = math.SafeAdd(gas, GasFastestStep); overflow {
		fuzz_helper.AddCoverage(34815)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(58334)
	}
	fuzz_helper.AddCoverage(3661)

	wordGas, overflow := bigUint64(stack.Back(2))
	if overflow {
		fuzz_helper.AddCoverage(46899)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(40738)
	}
	fuzz_helper.AddCoverage(22210)
	if wordGas, overflow = math.SafeMul(toWordSize(wordGas), params.CopyGas); overflow {
		fuzz_helper.AddCoverage(10840)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(43066)
	}
	fuzz_helper.AddCoverage(4417)
	if gas, overflow = math.SafeAdd(gas, wordGas); overflow {
		fuzz_helper.AddCoverage(14816)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(12109)
	}
	fuzz_helper.AddCoverage(5093)
	return gas, nil
}

func gasExtCodeCopy(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(29966)
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.AddCoverage(20099)
		return 0, err
	} else {
		fuzz_helper.AddCoverage(61712)
	}
	fuzz_helper.AddCoverage(27153)

	var overflow bool
	if gas, overflow = math.SafeAdd(gas, gt.ExtcodeCopy); overflow {
		fuzz_helper.AddCoverage(12762)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(17951)
	}
	fuzz_helper.AddCoverage(51386)

	wordGas, overflow := bigUint64(stack.Back(3))
	if overflow {
		fuzz_helper.AddCoverage(53729)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(64454)
	}
	fuzz_helper.AddCoverage(4676)

	if wordGas, overflow = math.SafeMul(toWordSize(wordGas), params.CopyGas); overflow {
		fuzz_helper.AddCoverage(7160)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(7189)
	}
	fuzz_helper.AddCoverage(61379)

	if gas, overflow = math.SafeAdd(gas, wordGas); overflow {
		fuzz_helper.AddCoverage(54641)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(55434)
	}
	fuzz_helper.AddCoverage(16873)
	return gas, nil
}

func gasMLoad(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(35217)
	var overflow bool
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.AddCoverage(19417)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(47979)
	}
	fuzz_helper.AddCoverage(51703)
	if gas, overflow = math.SafeAdd(gas, GasFastestStep); overflow {
		fuzz_helper.AddCoverage(22809)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(20851)
	}
	fuzz_helper.AddCoverage(39545)
	return gas, nil
}

func gasMStore8(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(29885)
	var overflow bool
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.AddCoverage(31201)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(25315)
	}
	fuzz_helper.AddCoverage(17574)
	if gas, overflow = math.SafeAdd(gas, GasFastestStep); overflow {
		fuzz_helper.AddCoverage(56379)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(7887)
	}
	fuzz_helper.AddCoverage(29431)
	return gas, nil
}

func gasMStore(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(33884)
	var overflow bool
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.AddCoverage(30464)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(62973)
	}
	fuzz_helper.AddCoverage(49673)
	if gas, overflow = math.SafeAdd(gas, GasFastestStep); overflow {
		fuzz_helper.AddCoverage(21243)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(29150)
	}
	fuzz_helper.AddCoverage(4765)
	return gas, nil
}

func gasCreate(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(26509)
	var overflow bool
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.AddCoverage(28688)
		return 0, err
	} else {
		fuzz_helper.AddCoverage(16762)
	}
	fuzz_helper.AddCoverage(62902)
	if gas, overflow = math.SafeAdd(gas, params.CreateGas); overflow {
		fuzz_helper.AddCoverage(60446)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(58133)
	}
	fuzz_helper.AddCoverage(2802)
	return gas, nil
}

func gasBalance(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(20790)
	return gt.Balance, nil
}

func gasExtCodeSize(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(63383)
	return gt.ExtcodeSize, nil
}

func gasSLoad(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(58404)
	return gt.SLoad, nil
}

func gasExp(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(54412)
	expByteLen := uint64((stack.data[stack.len()-2].BitLen() + 7) / 8)

	var (
		gas      = expByteLen * gt.ExpByte // no overflow check required. Max is 256 * ExpByte gas
		overflow bool
	)
	if gas, overflow = math.SafeAdd(gas, GasSlowStep); overflow {
		fuzz_helper.AddCoverage(29728)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(55641)
	}
	fuzz_helper.AddCoverage(62727)
	return gas, nil
}

func gasCall(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(11885)
	var (
		gas            = gt.Calls
		transfersValue = stack.Back(2).Sign() != 0
		address        = common.BigToAddress(stack.Back(1))
		eip158         = evm.ChainConfig().IsEIP158(evm.BlockNumber)
	)
	if eip158 {
		fuzz_helper.AddCoverage(45308)
		if transfersValue && evm.StateDB.Empty(address) {
			fuzz_helper.AddCoverage(1596)
			gas += params.CallNewAccountGas
		} else {
			fuzz_helper.AddCoverage(48751)
		}
	} else {
		fuzz_helper.AddCoverage(27607)
		if !evm.StateDB.Exist(address) {
			fuzz_helper.AddCoverage(60004)
			gas += params.CallNewAccountGas
		} else {
			fuzz_helper.AddCoverage(13056)
		}
	}
	fuzz_helper.AddCoverage(60640)
	if transfersValue {
		fuzz_helper.AddCoverage(12285)
		gas += params.CallValueTransferGas
	} else {
		fuzz_helper.AddCoverage(10163)
	}
	fuzz_helper.AddCoverage(32340)
	memoryGas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.AddCoverage(47390)
		return 0, err
	} else {
		fuzz_helper.AddCoverage(47458)
	}
	fuzz_helper.AddCoverage(49446)
	var overflow bool
	if gas, overflow = math.SafeAdd(gas, memoryGas); overflow {
		fuzz_helper.AddCoverage(10614)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(9340)
	}
	fuzz_helper.AddCoverage(55378)

	cg, err := callGas(gt, contract.Gas, gas, stack.Back(0))
	if err != nil {
		fuzz_helper.AddCoverage(25744)
		return 0, err
	} else {
		fuzz_helper.AddCoverage(65311)
	}
	fuzz_helper.AddCoverage(45742)

	stack.data[stack.len()-1] = new(big.Int).SetUint64(cg)

	if gas, overflow = math.SafeAdd(gas, cg); overflow {
		fuzz_helper.AddCoverage(19813)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(24779)
	}
	fuzz_helper.AddCoverage(58797)
	return gas, nil
}

func gasCallCode(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(179)
	gas := gt.Calls
	if stack.Back(2).Sign() != 0 {
		fuzz_helper.AddCoverage(10795)
		gas += params.CallValueTransferGas
	} else {
		fuzz_helper.AddCoverage(53895)
	}
	fuzz_helper.AddCoverage(4752)
	memoryGas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.AddCoverage(30484)
		return 0, err
	} else {
		fuzz_helper.AddCoverage(40730)
	}
	fuzz_helper.AddCoverage(58259)
	var overflow bool
	if gas, overflow = math.SafeAdd(gas, memoryGas); overflow {
		fuzz_helper.AddCoverage(60454)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(38671)
	}
	fuzz_helper.AddCoverage(60128)

	cg, err := callGas(gt, contract.Gas, gas, stack.Back(0))
	if err != nil {
		fuzz_helper.AddCoverage(55665)
		return 0, err
	} else {
		fuzz_helper.AddCoverage(57927)
	}
	fuzz_helper.AddCoverage(62145)

	stack.data[stack.len()-1] = new(big.Int).SetUint64(cg)

	if gas, overflow = math.SafeAdd(gas, cg); overflow {
		fuzz_helper.AddCoverage(65356)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(8275)
	}
	fuzz_helper.AddCoverage(12227)
	return gas, nil
}

func gasReturn(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(41896)
	return memoryGasCost(mem, memorySize)
}

func gasRevert(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(5111)
	return memoryGasCost(mem, memorySize)
}

func gasSuicide(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(2960)
	var gas uint64

	if evm.ChainConfig().IsEIP150(evm.BlockNumber) {
		fuzz_helper.AddCoverage(25547)
		gas = gt.Suicide
		var (
			address = common.BigToAddress(stack.Back(0))
			eip158  = evm.ChainConfig().IsEIP158(evm.BlockNumber)
		)

		if eip158 {
			fuzz_helper.AddCoverage(6705)

			if evm.StateDB.Empty(address) && evm.StateDB.GetBalance(contract.Address()).Sign() != 0 {
				fuzz_helper.AddCoverage(12502)
				gas += gt.CreateBySuicide
			} else {
				fuzz_helper.AddCoverage(5209)
			}
		} else {
			fuzz_helper.AddCoverage(3613)
			if !evm.StateDB.Exist(address) {
				fuzz_helper.AddCoverage(9728)
				gas += gt.CreateBySuicide
			} else {
				fuzz_helper.AddCoverage(20787)
			}
		}
	} else {
		fuzz_helper.AddCoverage(14592)
	}
	fuzz_helper.AddCoverage(35072)

	if !evm.StateDB.HasSuicided(contract.Address()) {
		fuzz_helper.AddCoverage(1274)
		evm.StateDB.AddRefund(new(big.Int).SetUint64(params.SuicideRefundGas))
	} else {
		fuzz_helper.AddCoverage(56143)
	}
	fuzz_helper.AddCoverage(55654)
	return gas, nil
}

func gasDelegateCall(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(17283)
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.AddCoverage(20593)
		return 0, err
	} else {
		fuzz_helper.AddCoverage(26119)
	}
	fuzz_helper.AddCoverage(11455)
	var overflow bool
	if gas, overflow = math.SafeAdd(gas, gt.Calls); overflow {
		fuzz_helper.AddCoverage(38633)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(27520)
	}
	fuzz_helper.AddCoverage(38950)

	cg, err := callGas(gt, contract.Gas, gas, stack.Back(0))
	if err != nil {
		fuzz_helper.AddCoverage(54694)
		return 0, err
	} else {
		fuzz_helper.AddCoverage(7779)
	}
	fuzz_helper.AddCoverage(65052)

	stack.data[stack.len()-1] = new(big.Int).SetUint64(cg)

	if gas, overflow = math.SafeAdd(gas, cg); overflow {
		fuzz_helper.AddCoverage(50594)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(53346)
	}
	fuzz_helper.AddCoverage(35391)
	return gas, nil
}

func gasStaticCall(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(38099)
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		fuzz_helper.AddCoverage(55029)
		return 0, err
	} else {
		fuzz_helper.AddCoverage(24853)
	}
	fuzz_helper.AddCoverage(28290)
	var overflow bool
	if gas, overflow = math.SafeAdd(gas, gt.Calls); overflow {
		fuzz_helper.AddCoverage(8)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(39569)
	}
	fuzz_helper.AddCoverage(46467)

	cg, err := callGas(gt, contract.Gas, gas, stack.Back(0))
	if err != nil {
		fuzz_helper.AddCoverage(14681)
		return 0, err
	} else {
		fuzz_helper.AddCoverage(32228)
	}
	fuzz_helper.AddCoverage(2463)

	stack.data[stack.len()-1] = new(big.Int).SetUint64(cg)

	if gas, overflow = math.SafeAdd(gas, cg); overflow {
		fuzz_helper.AddCoverage(42372)
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.AddCoverage(15925)
	}
	fuzz_helper.AddCoverage(41764)
	return gas, nil
}

func gasPush(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(61802)
	return GasFastestStep, nil
}

func gasSwap(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(63150)
	return GasFastestStep, nil
}

func gasDup(gt params.GasTable, evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	fuzz_helper.AddCoverage(26030)
	return GasFastestStep, nil
}

var _ = fuzz_helper.AddCoverage
