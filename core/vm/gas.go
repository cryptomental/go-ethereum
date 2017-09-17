package vm

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"math/big"

	"github.com/ethereum/go-ethereum/params"
)

const (
	GasQuickStep   uint64 = 2
	GasFastestStep uint64 = 3
	GasFastStep    uint64 = 5
	GasMidStep     uint64 = 8
	GasSlowStep    uint64 = 10
	GasExtStep     uint64 = 20

	GasReturn       uint64 = 0
	GasStop         uint64 = 0
	GasContractByte uint64 = 200
)

// calcGas returns the actual gas cost of the call.
//
// The cost of gas was changed during the homestead price change HF. To allow for EIP150
// to be implemented. The returned gas is gas - base * 63 / 64.
func callGas(gasTable params.GasTable, availableGas, base uint64, callCost *big.Int) (uint64, error) {
	fuzz_helper.CoverTab[22588]++
	if gasTable.CreateBySuicide > 0 {
		fuzz_helper.CoverTab[17878]++
		availableGas = availableGas - base
		gas := availableGas - availableGas/64

		if callCost.BitLen() > 64 || gas < callCost.Uint64() {
			fuzz_helper.CoverTab[45021]++
			return gas, nil
		} else {
			fuzz_helper.CoverTab[39040]++
		}
	} else {
		fuzz_helper.CoverTab[2095]++
	}
	fuzz_helper.CoverTab[44810]++
	if callCost.BitLen() > 64 {
		fuzz_helper.CoverTab[21668]++
		return 0, errGasUintOverflow
	} else {
		fuzz_helper.CoverTab[45213]++
	}
	fuzz_helper.CoverTab[5262]++

	return callCost.Uint64(), nil
}

var _ = fuzz_helper.CoverTab
