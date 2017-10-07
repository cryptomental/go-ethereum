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
	fuzz_helper.AddCoverage(55665)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	if contract.CodeAddr != nil {
		fuzz_helper.AddCoverage(65356)
		precompiles := PrecompiledContractsHomestead
		if evm.ChainConfig().IsByzantium(evm.BlockNumber) {
			fuzz_helper.AddCoverage(41896)
			precompiles = PrecompiledContractsByzantium
		} else {
			fuzz_helper.AddCoverage(5111)
		}
		fuzz_helper.AddCoverage(8275)
		if p := precompiles[*contract.CodeAddr]; p != nil {
			fuzz_helper.AddCoverage(2960)
			return RunPrecompiledContract(p, input, contract)
		} else {
			fuzz_helper.AddCoverage(35072)
		}
	} else {
		fuzz_helper.AddCoverage(55654)
	}
	fuzz_helper.AddCoverage(57927)
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
	fuzz_helper.AddCoverage(25547)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
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
	fuzz_helper.AddCoverage(6705)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	atomic.StoreInt32(&evm.abort, 1)
}

// Call executes the contract associated with the addr with the given input as
// parameters. It also handles any necessary value transfer required and takes
// the necessary steps to create accounts and reverses the state in case of an
// execution error or failed value transfer.
func (evm *EVM) Call(caller ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	fuzz_helper.AddCoverage(12502)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		fuzz_helper.AddCoverage(1274)
		return nil, gas, nil
	} else {
		fuzz_helper.AddCoverage(56143)
	}
	fuzz_helper.AddCoverage(5209)

	if evm.depth > int(params.CallCreateDepth) {
		fuzz_helper.AddCoverage(17283)
		return nil, gas, ErrDepth
	} else {
		fuzz_helper.AddCoverage(11455)
	}
	fuzz_helper.AddCoverage(3613)

	if !evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
		fuzz_helper.AddCoverage(38950)
		return nil, gas, ErrInsufficientBalance
	} else {
		fuzz_helper.AddCoverage(65052)
	}
	fuzz_helper.AddCoverage(9728)

	var (
		to       = AccountRef(addr)
		snapshot = evm.StateDB.Snapshot()
	)
	if !evm.StateDB.Exist(addr) {
		fuzz_helper.AddCoverage(35391)
		precompiles := PrecompiledContractsHomestead
		if evm.ChainConfig().IsByzantium(evm.BlockNumber) {
			fuzz_helper.AddCoverage(38633)
			precompiles = PrecompiledContractsByzantium
		} else {
			fuzz_helper.AddCoverage(27520)
		}
		fuzz_helper.AddCoverage(20593)
		if precompiles[addr] == nil && evm.ChainConfig().IsEIP158(evm.BlockNumber) && value.Sign() == 0 {
			fuzz_helper.AddCoverage(54694)
			return nil, gas, nil
		} else {
			fuzz_helper.AddCoverage(7779)
		}
		fuzz_helper.AddCoverage(26119)
		evm.StateDB.CreateAccount(addr)
	} else {
		fuzz_helper.AddCoverage(50594)
	}
	fuzz_helper.AddCoverage(20787)
	evm.Transfer(evm.StateDB, caller.Address(), to.Address(), value)

	contract := NewContract(caller, to, value, gas)
	contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr), evm.StateDB.GetCode(addr))

	ret, err = run(evm, snapshot, contract, input)

	if err != nil {
		fuzz_helper.AddCoverage(53346)
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			fuzz_helper.AddCoverage(38099)
			contract.UseGas(contract.Gas)
		} else {
			fuzz_helper.AddCoverage(28290)
		}
	} else {
		fuzz_helper.AddCoverage(46467)
	}
	fuzz_helper.AddCoverage(14592)
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
	fuzz_helper.AddCoverage(2463)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		fuzz_helper.AddCoverage(39569)
		return nil, gas, nil
	} else {
		fuzz_helper.AddCoverage(14681)
	}
	fuzz_helper.AddCoverage(41764)

	if evm.depth > int(params.CallCreateDepth) {
		fuzz_helper.AddCoverage(32228)
		return nil, gas, ErrDepth
	} else {
		fuzz_helper.AddCoverage(42372)
	}
	fuzz_helper.AddCoverage(55029)

	if !evm.CanTransfer(evm.StateDB, caller.Address(), value) {
		fuzz_helper.AddCoverage(15925)
		return nil, gas, ErrInsufficientBalance
	} else {
		fuzz_helper.AddCoverage(61802)
	}
	fuzz_helper.AddCoverage(24853)

	var (
		snapshot = evm.StateDB.Snapshot()
		to       = AccountRef(caller.Address())
	)

	contract := NewContract(caller, to, value, gas)
	contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr), evm.StateDB.GetCode(addr))

	ret, err = run(evm, snapshot, contract, input)
	if err != nil {
		fuzz_helper.AddCoverage(63150)
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			fuzz_helper.AddCoverage(26030)
			contract.UseGas(contract.Gas)
		} else {
			fuzz_helper.AddCoverage(36716)
		}
	} else {
		fuzz_helper.AddCoverage(21696)
	}
	fuzz_helper.AddCoverage(8)
	return ret, contract.Gas, err
}

// DelegateCall executes the contract associated with the addr with the given input
// as parameters. It reverses the state in case of an execution error.
//
// DelegateCall differs from CallCode in the sense that it executes the given address'
// code with the caller as context and the caller is set to the caller of the caller.
func (evm *EVM) DelegateCall(caller ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	fuzz_helper.AddCoverage(55609)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		fuzz_helper.AddCoverage(14325)
		return nil, gas, nil
	} else {
		fuzz_helper.AddCoverage(28666)
	}
	fuzz_helper.AddCoverage(554)

	if evm.depth > int(params.CallCreateDepth) {
		fuzz_helper.AddCoverage(35311)
		return nil, gas, ErrDepth
	} else {
		fuzz_helper.AddCoverage(46666)
	}
	fuzz_helper.AddCoverage(10614)

	var (
		snapshot = evm.StateDB.Snapshot()
		to       = AccountRef(caller.Address())
	)

	contract := NewContract(caller, to, nil, gas).AsDelegate()
	contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr), evm.StateDB.GetCode(addr))

	ret, err = run(evm, snapshot, contract, input)
	if err != nil {
		fuzz_helper.AddCoverage(53923)
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			fuzz_helper.AddCoverage(2311)
			contract.UseGas(contract.Gas)
		} else {
			fuzz_helper.AddCoverage(50649)
		}
	} else {
		fuzz_helper.AddCoverage(10532)
	}
	fuzz_helper.AddCoverage(61101)
	return ret, contract.Gas, err
}

// StaticCall executes the contract associated with the addr with the given input
// as parameters while disallowing any modifications to the state during the call.
// Opcodes that attempt to perform such modifications will result in exceptions
// instead of performing the modifications.
func (evm *EVM) StaticCall(caller ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	fuzz_helper.AddCoverage(18713)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		fuzz_helper.AddCoverage(12160)
		return nil, gas, nil
	} else {
		fuzz_helper.AddCoverage(2367)
	}
	fuzz_helper.AddCoverage(48086)

	if evm.depth > int(params.CallCreateDepth) {
		fuzz_helper.AddCoverage(53428)
		return nil, gas, ErrDepth
	} else {
		fuzz_helper.AddCoverage(47215)
	}
	fuzz_helper.AddCoverage(26434)

	if !evm.interpreter.readOnly {
		fuzz_helper.AddCoverage(31431)
		evm.interpreter.readOnly = true
		defer func() { fuzz_helper.AddCoverage(28516); evm.interpreter.readOnly = false }()
	} else {
		fuzz_helper.AddCoverage(54185)
	}
	fuzz_helper.AddCoverage(23521)

	var (
		to       = AccountRef(addr)
		snapshot = evm.StateDB.Snapshot()
	)

	contract := NewContract(caller, to, new(big.Int), gas)
	contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr), evm.StateDB.GetCode(addr))

	ret, err = run(evm, snapshot, contract, input)
	if err != nil {
		fuzz_helper.AddCoverage(50831)
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			fuzz_helper.AddCoverage(53403)
			contract.UseGas(contract.Gas)
		} else {
			fuzz_helper.AddCoverage(62616)
		}
	} else {
		fuzz_helper.AddCoverage(47149)
	}
	fuzz_helper.AddCoverage(11617)
	return ret, contract.Gas, err
}

// Create creates a new contract using code as deployment code.
func (evm *EVM) Create(caller ContractRef, code []byte, gas uint64, value *big.Int) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error) {
	fuzz_helper.AddCoverage(41214)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()

	if evm.depth > int(params.CallCreateDepth) {
		fuzz_helper.AddCoverage(39209)
		return nil, common.Address{}, gas, ErrDepth
	} else {
		fuzz_helper.AddCoverage(16335)
	}
	fuzz_helper.AddCoverage(64295)
	if !evm.CanTransfer(evm.StateDB, caller.Address(), value) {
		fuzz_helper.AddCoverage(59818)
		return nil, common.Address{}, gas, ErrInsufficientBalance
	} else {
		fuzz_helper.AddCoverage(51577)
	}
	fuzz_helper.AddCoverage(18054)

	nonce := evm.StateDB.GetNonce(caller.Address())
	evm.StateDB.SetNonce(caller.Address(), nonce+1)

	contractAddr = crypto.CreateAddress(caller.Address(), nonce)
	contractHash := evm.StateDB.GetCodeHash(contractAddr)
	if evm.StateDB.GetNonce(contractAddr) != 0 || (contractHash != (common.Hash{}) && contractHash != emptyCodeHash) {
		fuzz_helper.AddCoverage(60851)
		return nil, common.Address{}, 0, ErrContractAddressCollision
	} else {
		fuzz_helper.AddCoverage(10511)
	}
	fuzz_helper.AddCoverage(20109)

	snapshot := evm.StateDB.Snapshot()
	evm.StateDB.CreateAccount(contractAddr)
	if evm.ChainConfig().IsEIP158(evm.BlockNumber) {
		fuzz_helper.AddCoverage(27377)
		evm.StateDB.SetNonce(contractAddr, 1)
	} else {
		fuzz_helper.AddCoverage(15456)
	}
	fuzz_helper.AddCoverage(38750)
	evm.Transfer(evm.StateDB, caller.Address(), contractAddr, value)

	contract := NewContract(caller, AccountRef(contractAddr), value, gas)
	contract.SetCallCode(&contractAddr, crypto.Keccak256Hash(code), code)

	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		fuzz_helper.AddCoverage(52683)
		return nil, contractAddr, gas, nil
	} else {
		fuzz_helper.AddCoverage(33609)
	}
	fuzz_helper.AddCoverage(8966)
	ret, err = run(evm, snapshot, contract, nil)

	maxCodeSizeExceeded := evm.ChainConfig().IsEIP158(evm.BlockNumber) && len(ret) > params.MaxCodeSize

	if err == nil && !maxCodeSizeExceeded {
		fuzz_helper.AddCoverage(11535)
		createDataGas := uint64(len(ret)) * params.CreateDataGas
		if contract.UseGas(createDataGas) {
			fuzz_helper.AddCoverage(26893)
			evm.StateDB.SetCode(contractAddr, ret)
		} else {
			fuzz_helper.AddCoverage(5129)
			err = ErrCodeStoreOutOfGas
		}
	} else {
		fuzz_helper.AddCoverage(46588)
	}
	fuzz_helper.AddCoverage(21967)

	if maxCodeSizeExceeded || (err != nil && (evm.ChainConfig().IsHomestead(evm.BlockNumber) || err != ErrCodeStoreOutOfGas)) {
		fuzz_helper.AddCoverage(27525)
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			fuzz_helper.AddCoverage(58523)
			contract.UseGas(contract.Gas)
		} else {
			fuzz_helper.AddCoverage(25670)
		}
	} else {
		fuzz_helper.AddCoverage(26793)
	}
	fuzz_helper.AddCoverage(57184)

	if maxCodeSizeExceeded && err == nil {
		fuzz_helper.AddCoverage(27366)
		err = errMaxCodeSizeExceeded
	} else {
		fuzz_helper.AddCoverage(62659)
	}
	fuzz_helper.AddCoverage(9620)
	return ret, contractAddr, contract.Gas, err
}

// ChainConfig returns the evmironment's chain configuration
func (evm *EVM) ChainConfig() *params.ChainConfig {
	fuzz_helper.AddCoverage(61108)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return evm.chainConfig
}

// Interpreter returns the EVM interpreter
func (evm *EVM) Interpreter() *Interpreter {
	fuzz_helper.AddCoverage(24008)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return evm.interpreter
}

var _ = fuzz_helper.AddCoverage
