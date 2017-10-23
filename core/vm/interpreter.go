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
	fuzz_helper.AddCoverage(53076)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()

	if !cfg.JumpTable[STOP].valid {
		fuzz_helper.AddCoverage(49904)
		switch {
		case evm.ChainConfig().IsByzantium(evm.BlockNumber):
			fuzz_helper.AddCoverage(53754)
			cfg.JumpTable = byzantiumInstructionSet
		case evm.ChainConfig().IsHomestead(evm.BlockNumber):
			fuzz_helper.AddCoverage(14461)
			cfg.JumpTable = homesteadInstructionSet
		default:
			fuzz_helper.AddCoverage(48783)
			cfg.JumpTable = frontierInstructionSet
		}
	} else {
		fuzz_helper.AddCoverage(25041)
	}
	fuzz_helper.AddCoverage(30267)

	return &Interpreter{
		evm:      evm,
		cfg:      cfg,
		gasTable: evm.ChainConfig().GasTable(evm.BlockNumber),
		intPool:  newIntPool(),
	}
}

func (in *Interpreter) enforceRestrictions(op OpCode, operation operation, stack *Stack) error {
	fuzz_helper.AddCoverage(5456)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	if in.evm.chainRules.IsByzantium {
		fuzz_helper.AddCoverage(57329)
		if in.readOnly {
			fuzz_helper.AddCoverage(41491)

			if operation.writes || (op == CALL && stack.Back(2).BitLen() > 0) {
				fuzz_helper.AddCoverage(28583)
				return errWriteProtection
			} else {
				fuzz_helper.AddCoverage(60967)
			}
		} else {
			fuzz_helper.AddCoverage(11887)
		}
	} else {
		fuzz_helper.AddCoverage(29635)
	}
	fuzz_helper.AddCoverage(9809)
	return nil
}

// Run loops and evaluates the contract's code with the given input data and returns
// the return byte-slice and an error if one occurred.
//
// It's important to note that any errors returned by the interpreter should be
// considered a revert-and-consume-all-gas operation. No error specific checks
// should be handled to reduce complexity and errors further down the in.
func (in *Interpreter) Run(snapshot int, contract *Contract, input []byte) (ret []byte, err error) {
	fuzz_helper.AddCoverage(23400)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()

    logged := false
	in.evm.depth++
	defer func() { fuzz_helper.AddCoverage(25685); in.evm.depth-- }()
	fuzz_helper.AddCoverage(56930)

	in.returnData = nil

	if len(contract.Code) == 0 {
		fuzz_helper.AddCoverage(25085)
		return nil, nil
	} else {
		fuzz_helper.AddCoverage(288)
	}
	fuzz_helper.AddCoverage(1335)

	codehash := contract.CodeHash
	if codehash == (common.Hash{}) {
		fuzz_helper.AddCoverage(65508)
		codehash = crypto.Keccak256Hash(contract.Code)
	} else {
		fuzz_helper.AddCoverage(19803)
	}
	fuzz_helper.AddCoverage(8190)

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
		fuzz_helper.AddCoverage(53187)
		if err != nil && in.cfg.Debug {
			fuzz_helper.AddCoverage(4556)
            if logged == false {
                in.cfg.Tracer.CaptureState(in.evm, pc, op, contract.Gas, cost, mem, stack, contract, in.evm.depth, err)
            }
		} else {
			fuzz_helper.AddCoverage(53452)
		}
	}()
	fuzz_helper.AddCoverage(14908)

	for atomic.LoadInt32(&in.evm.abort) == 0 {
		fuzz_helper.AddCoverage(900)
        if pc >= uint64(len(contract.Code)) {
            break;
        }
        logged = false

		op = contract.GetOp(pc)

		operation := in.cfg.JumpTable[op]
		fuzz_helper.AddCoverage(52932)

		if !operation.valid {
			fuzz_helper.AddCoverage(60391)
			return nil, fmt.Errorf("invalid opcode 0x%x", int(op))
		} else {
			fuzz_helper.AddCoverage(43288)
		}
		fuzz_helper.AddCoverage(33045)

		if err := operation.validateStack(stack); err != nil {
			fuzz_helper.AddCoverage(64823)
			return nil, err
		} else {
			fuzz_helper.AddCoverage(61834)
		}
        // If the operation is valid, enforce and write restrictions
        if err := in.enforceRestrictions(op, operation, stack); err != nil {
			fuzz_helper.AddCoverage(61835)
            return nil, err
        }
		fuzz_helper.AddCoverage(42112)

		var memorySize uint64

		if operation.memorySize != nil {
			fuzz_helper.AddCoverage(32414)
			memSize, overflow := bigUint64(operation.memorySize(stack))
			if overflow {
				fuzz_helper.AddCoverage(43146)
				return nil, errGasUintOverflow
			} else {
				fuzz_helper.AddCoverage(52720)
			}
			fuzz_helper.AddCoverage(38787)

			if memorySize, overflow = math.SafeMul(toWordSize(memSize), 32); overflow {
				fuzz_helper.AddCoverage(58681)
				return nil, errGasUintOverflow
			} else {
				fuzz_helper.AddCoverage(61396)
			}
		} else {
			fuzz_helper.AddCoverage(2227)
		}
		fuzz_helper.AddCoverage(15716)

		if !in.cfg.DisableGasMetering {
			fuzz_helper.AddCoverage(46657)

			cost, err = operation.gasCost(in.gasTable, in.evm, contract, stack, mem, memorySize)
			if err != nil || !contract.UseGas(cost) {
				fuzz_helper.AddCoverage(12511)
				return nil, ErrOutOfGas
			} else {
				fuzz_helper.AddCoverage(35870)
			}
		} else {
			fuzz_helper.AddCoverage(27555)
		}
		fuzz_helper.AddCoverage(15856)
		if memorySize > 0 {
			fuzz_helper.AddCoverage(41224)
			mem.Resize(memorySize)
		} else {
			fuzz_helper.AddCoverage(64354)
		}
		fuzz_helper.AddCoverage(33163)

        oldpc := pc
		res, err := operation.execute(&pc, in.evm, contract, mem, stack)

		if in.cfg.Debug {
			fuzz_helper.AddCoverage(21358)
            if err == nil {
                in.cfg.Tracer.CaptureState(in.evm, oldpc, op, contract.Gas, cost, mem, stack, contract, in.evm.depth, err)
            }
		} else {
			fuzz_helper.AddCoverage(8695)
		}
		fuzz_helper.AddCoverage(42178)


		if verifyPool {
			fuzz_helper.AddCoverage(15580)
			verifyIntegerPool(in.intPool)
		} else {
			fuzz_helper.AddCoverage(52456)
		}
		fuzz_helper.AddCoverage(10038)

		if operation.returns {
			fuzz_helper.AddCoverage(32513)
			in.returnData = res
		} else {
			fuzz_helper.AddCoverage(64046)
		}
		fuzz_helper.AddCoverage(16517)

		switch {
		case err != nil:
			fuzz_helper.AddCoverage(18339)
			return nil, err
		case operation.reverts:
			fuzz_helper.AddCoverage(34754)
			return res, errExecutionReverted
		case operation.halts:
			fuzz_helper.AddCoverage(59259)
			return res, nil
		case !operation.jumps:
			fuzz_helper.AddCoverage(32239)
			pc++
		default:
			fuzz_helper.AddCoverage(6719)
		}
	}
	fuzz_helper.AddCoverage(34767)
	return nil, nil
}

var _ = fuzz_helper.AddCoverage
