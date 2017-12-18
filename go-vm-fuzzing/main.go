package main

import "C"

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"
import vmlogger "github.com/ethereum/go-ethereum/core/vm_instrumented"

import (
	"math/big"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm_instrumented"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/params"
)
var no_tracer bool;

//export WriteSymcov
func WriteSymcov(filename string) {
    fuzz_helper.WriteSymcov(filename)
}
//export EnableSymcovWriter
func EnableSymcovWriter() {
    fuzz_helper.EnableSymcovWriter()
}

//export SetInstrumentationType
func SetInstrumentationType(t int) {
    if t == 1 || t == 2 {
        no_tracer = true;
    }
    fuzz_helper.SetInstrumentationType(t)
}

//export DisableTracer
func DisableTracer() {
    no_tracer = true;
}

//export GoResetCoverage
func GoResetCoverage() {
    fuzz_helper.ResetCoverage()
}

//export GoCalcCoverage
func GoCalcCoverage() uint64 {
    return fuzz_helper.CalcCoverage()
}

//export MergeMode
func MergeMode() {
    fuzz_helper.MergeMode()
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
var g_msizes = make([]uint64, 0)
var g_trace_idx int;
var g_gastrace_idx int;
var g_msizetrace_idx int;

type AccountData struct {
    address uint64
    balance uint64
    code []byte
}

var g_accounts = make([]AccountData, 0)

var g_stack = make([](big.Int), 0);
var g_stack_idx int;

/* This function is called by the fuzzer to retrieve the execution specifics
   after a run.
*/
//export getTrace
func getTrace(finished *int, address *uint64, opcode *uint64 ) {
    if g_trace_idx >= len(g_addresses) {
        /* Reset to 0 so getTrace may be called again */
        g_trace_idx = 0

        /* Signal to the fuzzer that it has retrieved all trace items */
        *finished = 1
        return
    }

    *address = g_addresses[g_trace_idx]
    *opcode = g_opcodes[g_trace_idx]

    *finished = 0
    g_trace_idx++
}

/* This function is called by the fuzzer to retrieve the gas trace
   after a run.
*/
//export getGasTrace
func getGasTrace(finished *int, gas *uint64 ) {
    if g_gastrace_idx >= len(g_gases) {
        /* Reset to 0 so getTrace may be called again */
        g_gastrace_idx = 0

        /* Signal to the fuzzer that it has retrieved all trace items */
        *finished = 1
        return
    }

    *gas = g_gases[g_gastrace_idx]

    *finished = 0
    g_gastrace_idx++
}

/* This function is called by the fuzzer to retrieve the gas trace
   after a run.
*/
//export getMSizeTrace
func getMSizeTrace(finished *int, msize *uint64 ) {
    if g_msizetrace_idx >= len(g_msizes) {
        /* Reset to 0 so getTrace may be called again */
        g_msizetrace_idx = 0

        /* Signal to the fuzzer that it has retrieved all trace items */
        *finished = 1
        return
    }

    *msize = g_msizes[g_msizetrace_idx]

    *finished = 0
    g_msizetrace_idx++
}

/* This function is called by the fuzzer to retrieve the final stack state
   after a run
*/
//export getStack
func getStack(finished *int, stackitem []byte) {
    if g_stack_idx >= len(g_stack) {
        /* Reset to 0 so getStack may be called again */
        g_stack_idx = 0

        /* Signal to the fuzzer that it has retrieved all stack items */
        *finished = 1
        return
    }

    /* Prevent a buffer overwrite */
    stackitem_len := len(g_stack[g_stack_idx].Bytes())
    if stackitem_len > 32 {
        panic("stackitem too long")
    }

    copy(stackitem, g_stack[g_stack_idx].Bytes())

    *finished = 0
    g_stack_idx++
}

//export setAccounts
func setAccounts(address uint64, balance uint64, code []byte) {
    account := AccountData{
        address: address,
        balance: balance,
        code : code,
    }

    g_accounts = append(g_accounts, account)
}

//export runVM
func runVM(
    input []byte,
    success *int,
    do_trace int,
    gas uint64,
    blocknumber uint64,
    time uint64,
    gaslimit uint64,
    difficulty uint64,
    gasprice uint64,
    c_balance int64) {

    g_addresses = nil
    g_opcodes = nil
    g_gases = nil
    g_msizes = nil
    g_trace_idx = 0
    g_gastrace_idx = 0
    g_msizetrace_idx = 0

    g_stack = nil
    g_stack_idx = 0

	db, _ := ethdb.NewMemDatabase()
	sdb := state.NewDatabase(db)
	statedb, _ := state.New(common.Hash{}, sdb)

    /*
    addr := common.HexToAddress("0x1")
    balance := new(big.Int).SetUint64(1);
    statedb.SetBalance(addr, balance)

    addr = common.HexToAddress("0x2")
    statedb.SetBalance(addr, balance)

    addr = common.HexToAddress("0x3")
    statedb.SetBalance(addr, balance)

    addr = common.HexToAddress("0x4")
    statedb.SetBalance(addr, balance)
    */

    for _, acc := range(g_accounts) {
        addr := common.BigToAddress( new(big.Int).SetUint64(acc.address) )
        statedb.SetBalance(addr, new(big.Int).SetUint64(acc.balance))
        statedb.SetCode(addr, acc.code)
    }

    g_accounts = nil

	root, _ := statedb.CommitTo(db, false)
	statedb, _ = state.New(root, sdb)

    /* Helper functions required for correct functioning of the VM */
	canTransfer := func(db vm.StateDB, address common.Address, amount *big.Int) bool {
        return db.GetBalance(address).Cmp(amount) >= 0
	}
    transfer := func(db vm.StateDB, sender, recipient common.Address, amount *big.Int) {
        db.SubBalance(sender, amount)
        db.AddBalance(recipient, amount)
    }
    vmTestBlockHash := func(n uint64) common.Hash {
        i := new(big.Int).SetUint64(0)
        return common.BigToHash(i)
    }

	context := vm.Context{
		CanTransfer: canTransfer,
		Transfer: transfer,
		GetHash:     vmTestBlockHash,
		BlockNumber: new(big.Int).SetUint64(blocknumber),
		Time:   new(big.Int).SetUint64(time),
		Coinbase:   common.Address{},
		GasLimit:   new(big.Int).SetUint64(gaslimit),
		Difficulty:   new(big.Int).SetUint64(difficulty),
		GasPrice:   new(big.Int).SetUint64(gasprice),
	}

    logStack := false
    if do_trace != 0 {
        logStack = true
    }
	logger_config := &vm.LogConfig{
        DisableStack: false,
        LogStack: logStack,
        DisableStorage: true,
        FullStorage: false,
        Limit: 0,
        DisableMemory: false,
	}
    tracer := vm.NewStructLogger(logger_config)
    vm_config := vm.Config{Debug: true, Tracer: tracer}
    if no_tracer == true  {
        vm_config = vm.Config{}
    }
    env := vm.NewEVM(context, statedb, params.TestnetChainConfig, vm_config)
	contract := vm.NewContract(account{}, account{}, big.NewInt(c_balance), gas)
	contract.Code = input

    vmlogger.LastStack = nil
    vmlogger.PrevLastStack = nil

    /* Execute the byte code */
    _, err := env.Interpreter().Run(0, contract, []byte{})

    if err != nil {
        errStr := fmt.Sprintf("%v", err)
        if errStr == "evm: execution reverted" {
            err = nil
        }
    }
    if err == nil {
        *success = 1
    } else {
        /* Determine whether the error is caused by a REVERT */
        if do_trace != 0 {
            fmt.Printf("err is %v\n", err);
        }
        *success = 0
    }

    if no_tracer == false {
        logs := tracer.StructLogs()
        /* This loop stores the variables address, opcode, gas, msize at every step
        of the execution as well as the final stack state, for later
        retrieval by the fuzzer.
        */
        for _, t := range logs {
             var o = uint64(t.Op)
             g_opcodes = append(g_opcodes, o)
             g_addresses = append(g_addresses, t.Pc)
             g_gases = append(g_gases, t.Gas)
             g_msizes = append(g_msizes, uint64(t.MemorySize))
        }

        /* Set g_stack to the final stack state */
        for _, s := range vmlogger.PrevLastStack {
            g_stack = append(g_stack, *s)
        }

        /* Print address, opcode, gas at every step of the execution
        if the fuzzer is run with --trace
        */
        if do_trace != 0 {
            execution_num := 1
            for _, t := range logs {
                fmt.Printf("[%v] %v : %v\n", execution_num, t.Pc, t.Op)
                fmt.Printf("Stack: %v\n", t.Stack)
                fmt.Printf("Gas: %v\n", t.Gas)
                fmt.Printf("Depth: %v\n", t.Depth)
                fmt.Printf("Memory size: %v\n", t.MemorySize)

                execution_num++;
            }
        }
    }
}

/* No main() body because this file is compiled to a static archive */
func main() { }
