package vm

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func NoopCanTransfer(db StateDB, from common.Address, balance *big.Int) bool {
	fuzz_helper.CoverTab[22588]++
	return true
}
func NoopTransfer(db StateDB, from, to common.Address, amount *big.Int) { fuzz_helper.CoverTab[44810]++ }

type NoopEVMCallContext struct{}

func (NoopEVMCallContext) Call(caller ContractRef, addr common.Address, data []byte, gas, value *big.Int) ([]byte, error) {
	fuzz_helper.CoverTab[5262]++
	return nil, nil
}
func (NoopEVMCallContext) CallCode(caller ContractRef, addr common.Address, data []byte, gas, value *big.Int) ([]byte, error) {
	fuzz_helper.CoverTab[17878]++
	return nil, nil
}
func (NoopEVMCallContext) Create(caller ContractRef, data []byte, gas, value *big.Int) ([]byte, common.Address, error) {
	fuzz_helper.CoverTab[45021]++
	return nil, common.Address{}, nil
}
func (NoopEVMCallContext) DelegateCall(me ContractRef, addr common.Address, data []byte, gas *big.Int) ([]byte, error) {
	fuzz_helper.CoverTab[39040]++
	return nil, nil
}

type NoopStateDB struct{}

func (NoopStateDB) CreateAccount(common.Address)        { fuzz_helper.CoverTab[2095]++ }
func (NoopStateDB) SubBalance(common.Address, *big.Int) { fuzz_helper.CoverTab[21668]++ }
func (NoopStateDB) AddBalance(common.Address, *big.Int) { fuzz_helper.CoverTab[45213]++ }
func (NoopStateDB) GetBalance(common.Address) *big.Int  { fuzz_helper.CoverTab[16619]++; return nil }
func (NoopStateDB) GetNonce(common.Address) uint64      { fuzz_helper.CoverTab[12692]++; return 0 }
func (NoopStateDB) SetNonce(common.Address, uint64)     { fuzz_helper.CoverTab[42483]++ }
func (NoopStateDB) GetCodeHash(common.Address) common.Hash {
	fuzz_helper.CoverTab[6577]++
	return common.Hash{}
}
func (NoopStateDB) GetCode(common.Address) []byte  { fuzz_helper.CoverTab[17393]++; return nil }
func (NoopStateDB) SetCode(common.Address, []byte) { fuzz_helper.CoverTab[64174]++ }
func (NoopStateDB) GetCodeSize(common.Address) int { fuzz_helper.CoverTab[38740]++; return 0 }
func (NoopStateDB) AddRefund(*big.Int)             { fuzz_helper.CoverTab[35657]++ }
func (NoopStateDB) GetRefund() *big.Int            { fuzz_helper.CoverTab[30358]++; return nil }
func (NoopStateDB) GetState(common.Address, common.Hash) common.Hash {
	fuzz_helper.CoverTab[23294]++
	return common.Hash{}
}
func (NoopStateDB) SetState(common.Address, common.Hash, common.Hash) { fuzz_helper.CoverTab[61639]++ }
func (NoopStateDB) Suicide(common.Address) bool                       { fuzz_helper.CoverTab[11162]++; return false }
func (NoopStateDB) HasSuicided(common.Address) bool                   { fuzz_helper.CoverTab[49217]++; return false }
func (NoopStateDB) Exist(common.Address) bool                         { fuzz_helper.CoverTab[34511]++; return false }
func (NoopStateDB) Empty(common.Address) bool                         { fuzz_helper.CoverTab[64074]++; return false }
func (NoopStateDB) RevertToSnapshot(int)                              { fuzz_helper.CoverTab[28614]++ }
func (NoopStateDB) Snapshot() int                                     { fuzz_helper.CoverTab[39226]++; return 0 }
func (NoopStateDB) AddLog(*types.Log)                                 { fuzz_helper.CoverTab[2297]++ }
func (NoopStateDB) AddPreimage(common.Hash, []byte)                   { fuzz_helper.CoverTab[40870]++ }
func (NoopStateDB) ForEachStorage(common.Address, func(common.Hash, common.Hash) bool) {
	fuzz_helper.CoverTab[52877]++
}

var _ = fuzz_helper.CoverTab
