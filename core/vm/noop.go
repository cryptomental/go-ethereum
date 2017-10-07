package vm

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func NoopCanTransfer(db StateDB, from common.Address, balance *big.Int) bool {
	fuzz_helper.AddCoverage(12095)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return true
}
func NoopTransfer(db StateDB, from, to common.Address, amount *big.Int) {
	fuzz_helper.AddCoverage(14350)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
}

type NoopEVMCallContext struct{}

func (NoopEVMCallContext) Call(caller ContractRef, addr common.Address, data []byte, gas, value *big.Int) ([]byte, error) {
	fuzz_helper.AddCoverage(23405)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return nil, nil
}
func (NoopEVMCallContext) CallCode(caller ContractRef, addr common.Address, data []byte, gas, value *big.Int) ([]byte, error) {
	fuzz_helper.AddCoverage(2467)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return nil, nil
}
func (NoopEVMCallContext) Create(caller ContractRef, data []byte, gas, value *big.Int) ([]byte, common.Address, error) {
	fuzz_helper.AddCoverage(29828)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return nil, common.Address{}, nil
}
func (NoopEVMCallContext) DelegateCall(me ContractRef, addr common.Address, data []byte, gas *big.Int) ([]byte, error) {
	fuzz_helper.AddCoverage(13291)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return nil, nil
}

type NoopStateDB struct{}

func (NoopStateDB) CreateAccount(common.Address) {
	fuzz_helper.AddCoverage(16133)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
}
func (NoopStateDB) SubBalance(common.Address, *big.Int) {
	fuzz_helper.AddCoverage(46085)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
}
func (NoopStateDB) AddBalance(common.Address, *big.Int) {
	fuzz_helper.AddCoverage(9125)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
}
func (NoopStateDB) GetBalance(common.Address) *big.Int {
	fuzz_helper.AddCoverage(22394)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return nil
}
func (NoopStateDB) GetNonce(common.Address) uint64 {
	fuzz_helper.AddCoverage(63521)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return 0
}
func (NoopStateDB) SetNonce(common.Address, uint64) {
	fuzz_helper.AddCoverage(60963)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
}
func (NoopStateDB) GetCodeHash(common.Address) common.Hash {
	fuzz_helper.AddCoverage(47938)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return common.Hash{}
}
func (NoopStateDB) GetCode(common.Address) []byte {
	fuzz_helper.AddCoverage(57932)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return nil
}
func (NoopStateDB) SetCode(common.Address, []byte) {
	fuzz_helper.AddCoverage(47610)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
}
func (NoopStateDB) GetCodeSize(common.Address) int {
	fuzz_helper.AddCoverage(9424)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return 0
}
func (NoopStateDB) AddRefund(*big.Int) {
	fuzz_helper.AddCoverage(3615)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
}
func (NoopStateDB) GetRefund() *big.Int {
	fuzz_helper.AddCoverage(48133)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return nil
}
func (NoopStateDB) GetState(common.Address, common.Hash) common.Hash {
	fuzz_helper.AddCoverage(47314)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return common.Hash{}
}
func (NoopStateDB) SetState(common.Address, common.Hash, common.Hash) {
	fuzz_helper.AddCoverage(12873)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
}
func (NoopStateDB) Suicide(common.Address) bool {
	fuzz_helper.AddCoverage(33227)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return false
}
func (NoopStateDB) HasSuicided(common.Address) bool {
	fuzz_helper.AddCoverage(15626)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return false
}
func (NoopStateDB) Exist(common.Address) bool {
	fuzz_helper.AddCoverage(56132)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return false
}
func (NoopStateDB) Empty(common.Address) bool {
	fuzz_helper.AddCoverage(37300)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return false
}
func (NoopStateDB) RevertToSnapshot(int) {
	fuzz_helper.AddCoverage(18889)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
}
func (NoopStateDB) Snapshot() int {
	fuzz_helper.AddCoverage(55662)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return 0
}
func (NoopStateDB) AddLog(*types.Log) {
	fuzz_helper.AddCoverage(33806)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
}
func (NoopStateDB) AddPreimage(common.Hash, []byte) {
	fuzz_helper.AddCoverage(9009)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
}
func (NoopStateDB) ForEachStorage(common.Address, func(common.Hash, common.Hash) bool) {
	fuzz_helper.AddCoverage(6760)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
}

var _ = fuzz_helper.AddCoverage
