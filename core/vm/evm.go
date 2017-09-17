package vm

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

// emptyCodeHash is used by create to ensure deployment is disallowed to already
// deployed contract addresses (relevant after the account abstraction).
var emptyCodeHash = crypto.Keccak256Hash(nil)

type (
	CanTransferFunc func(StateDB, common.Address, *big.Int) bool
	TransferFunc    func(StateDB, common.Address, common.Address, *big.Int)
	// GetHashFunc returns the nth block hash in the blockchain
	// and is used by the BLOCKHASH EVM op code.
	GetHashFunc func(uint64) common.Hash
)

// run runs the given contract and takes care of running precompiles with a fallback to the byte code interpreter.
func run(evm *EVM, snapshot int, contract *Contract, input []byte) ([]byte, error) {
	fuzz_helper.CoverTab[22588]++
	if contract.CodeAddr != nil {
		fuzz_helper.CoverTab[5262]++
		precompiles := PrecompiledContractsHomestead
		if evm.ChainConfig().IsByzantium(evm.BlockNumber) {
			fuzz_helper.CoverTab[45021]++
			precompiles = PrecompiledContractsByzantium
		} else {
			fuzz_helper.CoverTab[39040]++
		}
		fuzz_helper.CoverTab[17878]++
		if p := precompiles[*contract.CodeAddr]; p != nil {
			fuzz_helper.CoverTab[2095]++
			return RunPrecompiledContract(p, input, contract)
		} else {
			fuzz_helper.CoverTab[21668]++
		}
	} else {
		fuzz_helper.CoverTab[45213]++
	}
	fuzz_helper.CoverTab[44810]++
	return evm.interpreter.Run(snapshot, contract, input)
}

// Context provides the EVM with auxiliary information. Once provided
// it shouldn't be modified.
type Context struct {
	// CanTransfer returns whether the account contains
	// sufficient ether to transfer the value
	CanTransfer CanTransferFunc
	// Transfer transfers ether from one account to the other
	Transfer TransferFunc
	// GetHash returns the hash corresponding to n
	GetHash GetHashFunc

	// Message information
	Origin   common.Address // Provides information for ORIGIN
	GasPrice *big.Int       // Provides information for GASPRICE

	// Block information
	Coinbase    common.Address // Provides information for COINBASE
	GasLimit    *big.Int       // Provides information for GASLIMIT
	BlockNumber *big.Int       // Provides information for NUMBER
	Time        *big.Int       // Provides information for TIME
	Difficulty  *big.Int       // Provides information for DIFFICULTY
}

// EVM is the Ethereum Virtual Machine base object and provides
// the necessary tools to run a contract on the given state with
// the provided context. It should be noted that any error
// generated through any of the calls should be considered a
// revert-state-and-consume-all-gas operation, no checks on
// specific errors should ever be performed. The interpreter makes
// sure that any errors generated are to be considered faulty code.
//
// The EVM should never be reused and is not thread safe.
type EVM struct {
	// Context provides auxiliary blockchain related information
	Context
	// StateDB gives access to the underlying state
	StateDB StateDB
	// Depth is the current call stack
	depth int

	// chainConfig contains information about the current chain
	chainConfig *params.ChainConfig
	// chain rules contains the chain rules for the current epoch
	chainRules params.Rules
	// virtual machine configuration options used to initialise the
	// evm.
	vmConfig Config
	// global (to this context) ethereum virtual machine
	// used throughout the execution of the tx.
	interpreter *Interpreter
	// abort is used to abort the EVM calling operations
	// NOTE: must be set atomically
	abort int32
}

// NewEVM retutrns a new EVM . The returned EVM is not thread safe and should
// only ever be used *once*.
func NewEVM(ctx Context, statedb StateDB, chainConfig *params.ChainConfig, vmConfig Config) *EVM {
	fuzz_helper.CoverTab[16619]++
	evm := &EVM{
		Context:     ctx,
		StateDB:     statedb,
		vmConfig:    vmConfig,
		chainConfig: chainConfig,
		chainRules:  chainConfig.Rules(ctx.BlockNumber),
	}

	evm.interpreter = NewInterpreter(evm, vmConfig)
	return evm
}

// Cancel cancels any running EVM operation. This may be called concurrently and
// it's safe to be called multiple times.
func (evm *EVM) Cancel() {
	fuzz_helper.CoverTab[12692]++
	atomic.StoreInt32(&evm.abort, 1)
}

// Call executes the contract associated with the addr with the given input as
// parameters. It also handles any necessary value transfer required and takes
// the necessary steps to create accounts and reverses the state in case of an
// execution error or failed value transfer.
func (evm *EVM) Call(caller ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	fuzz_helper.CoverTab[42483]++
	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		fuzz_helper.CoverTab[30358]++
		return nil, gas, nil
	} else {
		fuzz_helper.CoverTab[23294]++
	}
	fuzz_helper.CoverTab[6577]++

	if evm.depth > int(params.CallCreateDepth) {
		fuzz_helper.CoverTab[61639]++
		return nil, gas, ErrDepth
	} else {
		fuzz_helper.CoverTab[11162]++
	}
	fuzz_helper.CoverTab[17393]++

	if !evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
		fuzz_helper.CoverTab[49217]++
		return nil, gas, ErrInsufficientBalance
	} else {
		fuzz_helper.CoverTab[34511]++
	}
	fuzz_helper.CoverTab[64174]++

	var (
		to       = AccountRef(addr)
		snapshot = evm.StateDB.Snapshot()
	)
	if !evm.StateDB.Exist(addr) {
		fuzz_helper.CoverTab[64074]++
		precompiles := PrecompiledContractsHomestead
		if evm.ChainConfig().IsByzantium(evm.BlockNumber) {
			fuzz_helper.CoverTab[2297]++
			precompiles = PrecompiledContractsByzantium
		} else {
			fuzz_helper.CoverTab[40870]++
		}
		fuzz_helper.CoverTab[28614]++
		if precompiles[addr] == nil && evm.ChainConfig().IsEIP158(evm.BlockNumber) && value.Sign() == 0 {
			fuzz_helper.CoverTab[52877]++
			return nil, gas, nil
		} else {
			fuzz_helper.CoverTab[778]++
		}
		fuzz_helper.CoverTab[39226]++
		evm.StateDB.CreateAccount(addr)
	} else {
		fuzz_helper.CoverTab[33340]++
	}
	fuzz_helper.CoverTab[38740]++
	evm.Transfer(evm.StateDB, caller.Address(), to.Address(), value)

	contract := NewContract(caller, to, value, gas)
	contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr), evm.StateDB.GetCode(addr))

	ret, err = run(evm, snapshot, contract, input)

	if err != nil {
		fuzz_helper.CoverTab[15638]++
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			fuzz_helper.CoverTab[45869]++
			contract.UseGas(contract.Gas)
		} else {
			fuzz_helper.CoverTab[23368]++
		}
	} else {
		fuzz_helper.CoverTab[12901]++
	}
	fuzz_helper.CoverTab[35657]++
	return ret, contract.Gas, err
}

// CallCode executes the contract associated with the addr with the given input
// as parameters. It also handles any necessary value transfer required and takes
// the necessary steps to create accounts and reverses the state in case of an
// execution error or failed value transfer.
//
// CallCode differs from Call in the sense that it executes the given address'
// code with the caller as context.
func (evm *EVM) CallCode(caller ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	fuzz_helper.CoverTab[12499]++
	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		fuzz_helper.CoverTab[3566]++
		return nil, gas, nil
	} else {
		fuzz_helper.CoverTab[47636]++
	}
	fuzz_helper.CoverTab[42993]++

	if evm.depth > int(params.CallCreateDepth) {
		fuzz_helper.CoverTab[8730]++
		return nil, gas, ErrDepth
	} else {
		fuzz_helper.CoverTab[20539]++
	}
	fuzz_helper.CoverTab[30301]++

	if !evm.CanTransfer(evm.StateDB, caller.Address(), value) {
		fuzz_helper.CoverTab[63931]++
		return nil, gas, ErrInsufficientBalance
	} else {
		fuzz_helper.CoverTab[19009]++
	}
	fuzz_helper.CoverTab[45210]++

	var (
		snapshot = evm.StateDB.Snapshot()
		to       = AccountRef(caller.Address())
	)

	contract := NewContract(caller, to, value, gas)
	contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr), evm.StateDB.GetCode(addr))

	ret, err = run(evm, snapshot, contract, input)
	if err != nil {
		fuzz_helper.CoverTab[64748]++
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			fuzz_helper.CoverTab[50446]++
			contract.UseGas(contract.Gas)
		} else {
			fuzz_helper.CoverTab[18500]++
		}
	} else {
		fuzz_helper.CoverTab[52152]++
	}
	fuzz_helper.CoverTab[264]++
	return ret, contract.Gas, err
}

// DelegateCall executes the contract associated with the addr with the given input
// as parameters. It reverses the state in case of an execution error.
//
// DelegateCall differs from CallCode in the sense that it executes the given address'
// code with the caller as context and the caller is set to the caller of the caller.
func (evm *EVM) DelegateCall(caller ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	fuzz_helper.CoverTab[17111]++
	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		fuzz_helper.CoverTab[912]++
		return nil, gas, nil
	} else {
		fuzz_helper.CoverTab[64631]++
	}
	fuzz_helper.CoverTab[9670]++

	if evm.depth > int(params.CallCreateDepth) {
		fuzz_helper.CoverTab[15513]++
		return nil, gas, ErrDepth
	} else {
		fuzz_helper.CoverTab[17300]++
	}
	fuzz_helper.CoverTab[55848]++

	var (
		snapshot = evm.StateDB.Snapshot()
		to       = AccountRef(caller.Address())
	)

	contract := NewContract(caller, to, nil, gas).AsDelegate()
	contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr), evm.StateDB.GetCode(addr))

	ret, err = run(evm, snapshot, contract, input)
	if err != nil {
		fuzz_helper.CoverTab[16403]++
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			fuzz_helper.CoverTab[40937]++
			contract.UseGas(contract.Gas)
		} else {
			fuzz_helper.CoverTab[33825]++
		}
	} else {
		fuzz_helper.CoverTab[7237]++
	}
	fuzz_helper.CoverTab[50755]++
	return ret, contract.Gas, err
}

// StaticCall executes the contract associated with the addr with the given input
// as parameters while disallowing any modifications to the state during the call.
// Opcodes that attempt to perform such modifications will result in exceptions
// instead of performing the modifications.
func (evm *EVM) StaticCall(caller ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	fuzz_helper.CoverTab[23248]++
	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		fuzz_helper.CoverTab[5383]++
		return nil, gas, nil
	} else {
		fuzz_helper.CoverTab[52957]++
	}
	fuzz_helper.CoverTab[52715]++

	if evm.depth > int(params.CallCreateDepth) {
		fuzz_helper.CoverTab[6211]++
		return nil, gas, ErrDepth
	} else {
		fuzz_helper.CoverTab[49245]++
	}
	fuzz_helper.CoverTab[11389]++

	if !evm.interpreter.readOnly {
		fuzz_helper.CoverTab[15785]++
		evm.interpreter.readOnly = true
		defer func() { fuzz_helper.CoverTab[9735]++; evm.interpreter.readOnly = false }()
	} else {
		fuzz_helper.CoverTab[45823]++
	}
	fuzz_helper.CoverTab[60629]++

	var (
		to       = AccountRef(addr)
		snapshot = evm.StateDB.Snapshot()
	)

	contract := NewContract(caller, to, new(big.Int), gas)
	contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr), evm.StateDB.GetCode(addr))

	ret, err = run(evm, snapshot, contract, input)
	if err != nil {
		fuzz_helper.CoverTab[48647]++
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			fuzz_helper.CoverTab[24978]++
			contract.UseGas(contract.Gas)
		} else {
			fuzz_helper.CoverTab[61755]++
		}
	} else {
		fuzz_helper.CoverTab[19607]++
	}
	fuzz_helper.CoverTab[23245]++
	return ret, contract.Gas, err
}

// Create creates a new contract using code as deployment code.
func (evm *EVM) Create(caller ContractRef, code []byte, gas uint64, value *big.Int) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error) {
	fuzz_helper.CoverTab[28743]++

	if evm.depth > int(params.CallCreateDepth) {
		fuzz_helper.CoverTab[3661]++
		return nil, common.Address{}, gas, ErrDepth
	} else {
		fuzz_helper.CoverTab[22210]++
	}
	fuzz_helper.CoverTab[8832]++
	if !evm.CanTransfer(evm.StateDB, caller.Address(), value) {
		fuzz_helper.CoverTab[4417]++
		return nil, common.Address{}, gas, ErrInsufficientBalance
	} else {
		fuzz_helper.CoverTab[5093]++
	}
	fuzz_helper.CoverTab[40052]++

	nonce := evm.StateDB.GetNonce(caller.Address())
	evm.StateDB.SetNonce(caller.Address(), nonce+1)

	contractAddr = crypto.CreateAddress(caller.Address(), nonce)
	contractHash := evm.StateDB.GetCodeHash(contractAddr)
	if evm.StateDB.GetNonce(contractAddr) != 0 || (contractHash != (common.Hash{}) && contractHash != emptyCodeHash) {
		fuzz_helper.CoverTab[56274]++
		return nil, common.Address{}, 0, ErrContractAddressCollision
	} else {
		fuzz_helper.CoverTab[1404]++
	}
	fuzz_helper.CoverTab[63449]++

	snapshot := evm.StateDB.Snapshot()
	evm.StateDB.CreateAccount(contractAddr)
	if evm.ChainConfig().IsEIP158(evm.BlockNumber) {
		fuzz_helper.CoverTab[34815]++
		evm.StateDB.SetNonce(contractAddr, 1)
	} else {
		fuzz_helper.CoverTab[58334]++
	}
	fuzz_helper.CoverTab[47485]++
	evm.Transfer(evm.StateDB, caller.Address(), contractAddr, value)

	contract := NewContract(caller, AccountRef(contractAddr), value, gas)
	contract.SetCallCode(&contractAddr, crypto.Keccak256Hash(code), code)

	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		fuzz_helper.CoverTab[46899]++
		return nil, contractAddr, gas, nil
	} else {
		fuzz_helper.CoverTab[40738]++
	}
	fuzz_helper.CoverTab[60075]++
	ret, err = run(evm, snapshot, contract, nil)

	maxCodeSizeExceeded := evm.ChainConfig().IsEIP158(evm.BlockNumber) && len(ret) > params.MaxCodeSize

	if err == nil && !maxCodeSizeExceeded {
		fuzz_helper.CoverTab[10840]++
		createDataGas := uint64(len(ret)) * params.CreateDataGas
		if contract.UseGas(createDataGas) {
			fuzz_helper.CoverTab[43066]++
			evm.StateDB.SetCode(contractAddr, ret)
		} else {
			fuzz_helper.CoverTab[14816]++
			err = ErrCodeStoreOutOfGas
		}
	} else {
		fuzz_helper.CoverTab[12109]++
	}
	fuzz_helper.CoverTab[21817]++

	if maxCodeSizeExceeded || (err != nil && (evm.ChainConfig().IsHomestead(evm.BlockNumber) || err != ErrCodeStoreOutOfGas)) {
		fuzz_helper.CoverTab[29966]++
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			fuzz_helper.CoverTab[27153]++
			contract.UseGas(contract.Gas)
		} else {
			fuzz_helper.CoverTab[51386]++
		}
	} else {
		fuzz_helper.CoverTab[4676]++
	}
	fuzz_helper.CoverTab[57682]++

	if maxCodeSizeExceeded && err == nil {
		fuzz_helper.CoverTab[61379]++
		err = errMaxCodeSizeExceeded
	} else {
		fuzz_helper.CoverTab[16873]++
	}
	fuzz_helper.CoverTab[45496]++
	return ret, contractAddr, contract.Gas, err
}

// ChainConfig returns the evmironment's chain configuration
func (evm *EVM) ChainConfig() *params.ChainConfig {
	fuzz_helper.CoverTab[20099]++
	return evm.chainConfig
}

// Interpreter returns the EVM interpreter
func (evm *EVM) Interpreter() *Interpreter { fuzz_helper.CoverTab[61712]++; return evm.interpreter }

var _ = fuzz_helper.CoverTab
