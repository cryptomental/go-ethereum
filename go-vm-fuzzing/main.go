package main

import "C"

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
)

//export GoResetCoverage
func GoResetCoverage() {
    for i := 0; i < fuzz_helper.CoverSize; i++ {
        fuzz_helper.CoverTab[i] = 0
    }
}

//export GoCalcCoverage
func GoCalcCoverage() int {
    coverage := 0

    for i := 0; i < fuzz_helper.CoverSize; i++ {
        if fuzz_helper.CoverTab[i] != 0 {
            coverage += 1
        }
	}

    return coverage
}

type account struct{}

func (account) SubBalance(amount *big.Int)                          {}
func (account) AddBalance(amount *big.Int)                          {}
func (account) SetAddress(common.Address)                           {}
func (account) Value() *big.Int                                     { return nil }
func (account) SetBalance(*big.Int)                                 {}
func (account) SetNonce(uint64)                                     {}
func (account) Balance() *big.Int                                   { return nil }
func (account) Address() common.Address                             { return common.Address{} }
func (account) ReturnGas(*big.Int)                                  {}
func (account) SetCode(common.Hash, []byte)                         {}
func (account) ForEachStorage(cb func(key, value common.Hash) bool) {}

//export runVM
func runVM(input []byte) {

	db, _ := ethdb.NewMemDatabase()
	sdb := state.NewDatabase(db)
	statedb, _ := state.New(common.Hash{}, sdb)

	canTransfer := func(db vm.StateDB, address common.Address, amount *big.Int) bool {
        return false;
	}
	context := vm.Context{
		CanTransfer: canTransfer,
		GetHash:     vmTestBlockHash,
		BlockNumber: new(big.Int).SetUint64(1000),
		Time:   new(big.Int).SetUint64(1000),
		GasLimit:   new(big.Int).SetUint64(1000),
		Difficulty:   new(big.Int).SetUint64(1000),
		GasPrice:   new(big.Int).SetUint64(1000),
	}
	env := vm.NewEVM(context, statedb, params.TestChainConfig, vm.Config{})
	contract := vm.NewContract(account{}, account{}, big.NewInt(0), 10000)
	contract.Code = input

	_, _ = env.Interpreter().Run(0, contract, []byte{})
}

func vmTestBlockHash(n uint64) common.Hash {
	return common.BytesToHash(crypto.Keccak256([]byte(big.NewInt(int64(n)).String())))
}

func main() { }
