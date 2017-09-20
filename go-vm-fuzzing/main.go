package main

import "C"

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"math/big"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
//	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/common/math"
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
func (account) Address() common.Address                             { a := new(big.Int).SetUint64(0x155); return common.BigToAddress(a); }
func (account) ReturnGas(*big.Int)                                  {}
func (account) SetCode(common.Hash, []byte)                         {}
func (account) ForEachStorage(cb func(key, value common.Hash) bool) {}

type StructLogRes struct {
	Pc      uint64            `json:"pc"`
	Op      string            `json:"op"`
	Gas     uint64            `json:"gas"`
	GasCost uint64            `json:"gasCost"`
	Depth   int               `json:"depth"`
	Error   error             `json:"error"`
	Stack   []string          `json:"stack"`
	Memory  []string          `json:"memory"`
	Storage map[string]string `json:"storage"`
}
func FormatLogs(structLogs []vm.StructLog) []StructLogRes {
	formattedStructLogs := make([]StructLogRes, len(structLogs))
	for index, trace := range structLogs {
		formattedStructLogs[index] = StructLogRes{
			Pc:      trace.Pc,
			Op:      trace.Op.String(),
			Gas:     trace.Gas,
			GasCost: trace.GasCost,
			Depth:   trace.Depth,
			Error:   trace.Err,
			Stack:   make([]string, len(trace.Stack)),
			Storage: make(map[string]string),
		}

		for i, stackValue := range trace.Stack {
			formattedStructLogs[index].Stack[i] = fmt.Sprintf("%x", math.PaddedBigBytes(stackValue, 32))
		}

		for i := 0; i+32 <= len(trace.Memory); i += 32 {
			formattedStructLogs[index].Memory = append(formattedStructLogs[index].Memory, fmt.Sprintf("%x", trace.Memory[i:i+32]))
		}

		for i, storageValue := range trace.Storage {
			formattedStructLogs[index].Storage[fmt.Sprintf("%x", i)] = fmt.Sprintf("%x", storageValue)
		}
	}
	return formattedStructLogs
}

var g_addresses = make([]uint64, 0)
var g_opcodes = make([]uint64, 0)
var g_gases = make([]uint64, 0)
var g_trace_idx int;

//export getTrace
func getTrace(finished *int, address *uint64, opcode *uint64, gas *uint64 ) {
    if g_trace_idx >= len(g_addresses) {
        *finished = 1
        return
    }
    *address = g_addresses[g_trace_idx]
    *opcode = g_opcodes[g_trace_idx]
    *gas= g_gases[g_trace_idx]
    *finished = 0
    g_trace_idx++
}

//export runVM
func runVM(input []byte, success *int, do_trace int, gas uint64) {

    g_addresses = nil
    g_opcodes = nil
    g_gases = nil
    g_trace_idx = 0

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
		Coinbase:   common.Address{},
		GasLimit:   new(big.Int).SetUint64(1000),
		Difficulty:   new(big.Int).SetUint64(1000),
		GasPrice:   new(big.Int).SetUint64(1000),
	}
    tracer := vm.NewStructLogger(nil)
    env := vm.NewEVM(context, statedb, params.TestChainConfig, vm.Config{Debug: true, Tracer: tracer})
	contract := vm.NewContract(account{}, account{}, big.NewInt(0), gas)
	contract.Code = input

    _, err := env.Interpreter().Run(0, contract, []byte{})
    if err == nil {
        *success = 1
    } else {
        if do_trace != 0 {
            fmt.Printf("err is %v\n", err);
        }
        *success = 0
    }
    for _, t := range tracer.StructLogs() {
        g_addresses = append(g_addresses, t.Pc)
        var o = uint64(t.Op)
        g_opcodes = append(g_opcodes, o)
        g_gases = append(g_gases, t.Gas)
    }
    if do_trace != 0 {
        for _, t := range tracer.StructLogs() {
            fmt.Printf("%v : %v\n", t.Pc, t.Op)
            fmt.Printf("Stack: %v\n", t.Stack)
            fmt.Printf("Gas: %v\n", t.Gas)
            /*
            fmt.Printf("Op: %v\n", t.Op)
            fmt.Printf("Gas: %v\n", t.Gas)
            fmt.Printf("GasCost: %v\n", t.GasCost)
            fmt.Printf("Depth: %v\n", t.Depth)
            fmt.Printf("Stack: %v\n", t.Stack)
            fmt.Printf("Storage: %v\n", t.Storage)
            fmt.Printf("\n")
            */
        }
    }
}

func vmTestBlockHash(n uint64) common.Hash {
    i := new(big.Int).SetUint64(123)
    return common.BigToHash(i)
}

func main() { }
