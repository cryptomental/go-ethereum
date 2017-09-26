package vm

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"fmt"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

// Config are the configuration options for the Interpreter
type Config struct {
	// Debug enabled debugging Interpreter options
	Debug bool
	// EnableJit enabled the JIT VM
	EnableJit bool
	// ForceJit forces the JIT VM
	ForceJit bool
	// Tracer is the op code logger
	Tracer Tracer
	// NoRecursion disabled Interpreter call, callcode,
	// delegate call and create.
	NoRecursion bool
	// Disable gas metering
	DisableGasMetering bool
	// Enable recording of SHA3/keccak preimages
	EnablePreimageRecording bool
	// JumpTable contains the EVM instruction table. This
	// may be left uninitialised and will be set to the default
	// table.
	JumpTable [256]operation
}

// Interpreter is used to run Ethereum based contracts and will utilise the
// passed evmironment to query external sources for state information.
// The Interpreter will run the byte code VM or JIT VM based on the passed
// configuration.
type Interpreter struct {
	evm      *EVM
	cfg      Config
	gasTable params.GasTable
	intPool  *intPool

	readOnly   bool   // Whether to throw on stateful modifications
	returnData []byte // Last CALL's return data for subsequent reuse
}

// NewInterpreter returns a new instance of the Interpreter.
func NewInterpreter(evm *EVM, cfg Config) *Interpreter {
	fuzz_helper.AddCoverage(22588)

	if !cfg.JumpTable[STOP].valid {
		fuzz_helper.AddCoverage(5262)
		switch {
		case evm.ChainConfig().IsByzantium(evm.BlockNumber):
			fuzz_helper.AddCoverage(17878)
			cfg.JumpTable = byzantiumInstructionSet
		case evm.ChainConfig().IsHomestead(evm.BlockNumber):
			fuzz_helper.AddCoverage(45021)
			cfg.JumpTable = homesteadInstructionSet
		default:
			fuzz_helper.AddCoverage(39040)
			cfg.JumpTable = frontierInstructionSet
		}
	} else {
		fuzz_helper.AddCoverage(2095)
	}
	fuzz_helper.AddCoverage(44810)

	return &Interpreter{
		evm:      evm,
		cfg:      cfg,
		gasTable: evm.ChainConfig().GasTable(evm.BlockNumber),
		intPool:  newIntPool(),
	}
}

func (in *Interpreter) enforceRestrictions(op OpCode, operation operation, stack *Stack) error {
	fuzz_helper.AddCoverage(21668)
	if in.evm.chainRules.IsByzantium {
		fuzz_helper.AddCoverage(16619)
		if in.readOnly {
			fuzz_helper.AddCoverage(12692)

			if operation.writes || (op == CALL && stack.Back(2).BitLen() > 0) {
				fuzz_helper.AddCoverage(42483)
				return errWriteProtection
			} else {
				fuzz_helper.AddCoverage(6577)
			}
		} else {
			fuzz_helper.AddCoverage(17393)
		}
	} else {
		fuzz_helper.AddCoverage(64174)
	}
	fuzz_helper.AddCoverage(45213)
	return nil
}

// Run loops and evaluates the contract's code with the given input data and returns
// the return byte-slice and an error if one occurred.
//
// It's important to note that any errors returned by the interpreter should be
// considered a revert-and-consume-all-gas operation. No error specific checks
// should be handled to reduce complexity and errors further down the in.
func (in *Interpreter) Run(snapshot int, contract *Contract, input []byte) (ret []byte, err error) {
	fuzz_helper.AddCoverage(38740)

	in.evm.depth++
	defer func() { fuzz_helper.AddCoverage(49217); in.evm.depth-- }()
	fuzz_helper.AddCoverage(35657)

	in.returnData = nil

	if len(contract.Code) == 0 {
		fuzz_helper.AddCoverage(34511)
		return nil, nil
	} else {
		fuzz_helper.AddCoverage(64074)
	}
	fuzz_helper.AddCoverage(30358)

	codehash := contract.CodeHash
	if codehash == (common.Hash{}) {
		fuzz_helper.AddCoverage(28614)
		codehash = crypto.Keccak256Hash(contract.Code)
	} else {
		fuzz_helper.AddCoverage(39226)
	}
	fuzz_helper.AddCoverage(23294)

	var (
		op    OpCode        // current opcode
		mem   = NewMemory() // bound memory
		stack = newstack()  // local stack
		// For optimisation reason we're using uint64 as the program counter.
		// It's theoretically possible to go above 2^64. The YP defines the PC
		// to be uint256. Practically much less so feasible.
		pc   = uint64(0) // program counter
		cost uint64
	)
	contract.Input = input

	defer func() {
		fuzz_helper.AddCoverage(2297)
		if err != nil && in.cfg.Debug {
			fuzz_helper.AddCoverage(40870)
			in.cfg.Tracer.CaptureState(in.evm, pc, op, contract.Gas, cost, mem, stack, contract, in.evm.depth, err)
		} else {
			fuzz_helper.AddCoverage(52877)
		}
	}()
	fuzz_helper.AddCoverage(61639)

	for atomic.LoadInt32(&in.evm.abort) == 0 {
		fuzz_helper.AddCoverage(778)

		op = contract.GetOp(pc)

		operation := in.cfg.JumpTable[op]
		if err := in.enforceRestrictions(op, operation, stack); err != nil {
			fuzz_helper.AddCoverage(264)
			return nil, err
		} else {
			fuzz_helper.AddCoverage(3566)
		}
		fuzz_helper.AddCoverage(33340)

		if !operation.valid {
			fuzz_helper.AddCoverage(47636)
			return nil, fmt.Errorf("invalid opcode 0x%x", int(op))
		} else {
			fuzz_helper.AddCoverage(8730)
		}
		fuzz_helper.AddCoverage(15638)

		if err := operation.validateStack(stack); err != nil {
			fuzz_helper.AddCoverage(20539)
			return nil, err
		} else {
			fuzz_helper.AddCoverage(63931)
		}
		fuzz_helper.AddCoverage(45869)

		var memorySize uint64

		if operation.memorySize != nil {
			fuzz_helper.AddCoverage(19009)
			memSize, overflow := bigUint64(operation.memorySize(stack))
			if overflow {
				fuzz_helper.AddCoverage(50446)
				return nil, errGasUintOverflow
			} else {
				fuzz_helper.AddCoverage(18500)
			}
			fuzz_helper.AddCoverage(64748)

			if memorySize, overflow = math.SafeMul(toWordSize(memSize), 32); overflow {
				fuzz_helper.AddCoverage(52152)
				return nil, errGasUintOverflow
			} else {
				fuzz_helper.AddCoverage(17111)
			}
		} else {
			fuzz_helper.AddCoverage(9670)
		}
		fuzz_helper.AddCoverage(23368)

		if !in.cfg.DisableGasMetering {
			fuzz_helper.AddCoverage(55848)

			cost, err = operation.gasCost(in.gasTable, in.evm, contract, stack, mem, memorySize)
			if err != nil || !contract.UseGas(cost) {
				fuzz_helper.AddCoverage(50755)
				return nil, ErrOutOfGas
			} else {
				fuzz_helper.AddCoverage(912)
			}
		} else {
			fuzz_helper.AddCoverage(64631)
		}
		fuzz_helper.AddCoverage(12901)
		if memorySize > 0 {
			fuzz_helper.AddCoverage(15513)
			mem.Resize(memorySize)
		} else {
			fuzz_helper.AddCoverage(17300)
		}
		fuzz_helper.AddCoverage(12499)

		if in.cfg.Debug {
			fuzz_helper.AddCoverage(16403)
			in.cfg.Tracer.CaptureState(in.evm, pc, op, contract.Gas, cost, mem, stack, contract, in.evm.depth, err)
		} else {
			fuzz_helper.AddCoverage(40937)
		}
		fuzz_helper.AddCoverage(42993)

		res, err := operation.execute(&pc, in.evm, contract, mem, stack)

		if verifyPool {
			fuzz_helper.AddCoverage(33825)
			verifyIntegerPool(in.intPool)
		} else {
			fuzz_helper.AddCoverage(7237)
		}
		fuzz_helper.AddCoverage(30301)

		if operation.returns {
			fuzz_helper.AddCoverage(23248)
			in.returnData = res
		} else {
			fuzz_helper.AddCoverage(52715)
		}
		fuzz_helper.AddCoverage(45210)

		switch {
		case err != nil:
			fuzz_helper.AddCoverage(11389)
			return nil, err
		case operation.reverts:
			fuzz_helper.AddCoverage(60629)
			return res, errExecutionReverted
		case operation.halts:
			fuzz_helper.AddCoverage(23245)
			return res, nil
		case !operation.jumps:
			fuzz_helper.AddCoverage(5383)
			pc++
		default:
			fuzz_helper.AddCoverage(52957)
		}
	}
	fuzz_helper.AddCoverage(11162)
	return nil, nil
}

var _ = fuzz_helper.AddCoverage
