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
	fuzz_helper.CoverTab[22588]++

	if !cfg.JumpTable[STOP].valid {
		fuzz_helper.CoverTab[5262]++
		switch {
		case evm.ChainConfig().IsByzantium(evm.BlockNumber):
			fuzz_helper.CoverTab[17878]++
			cfg.JumpTable = byzantiumInstructionSet
		case evm.ChainConfig().IsHomestead(evm.BlockNumber):
			fuzz_helper.CoverTab[45021]++
			cfg.JumpTable = homesteadInstructionSet
		default:
			fuzz_helper.CoverTab[39040]++
			cfg.JumpTable = frontierInstructionSet
		}
	} else {
		fuzz_helper.CoverTab[2095]++
	}
	fuzz_helper.CoverTab[44810]++

	return &Interpreter{
		evm:      evm,
		cfg:      cfg,
		gasTable: evm.ChainConfig().GasTable(evm.BlockNumber),
		intPool:  newIntPool(),
	}
}

func (in *Interpreter) enforceRestrictions(op OpCode, operation operation, stack *Stack) error {
	fuzz_helper.CoverTab[21668]++
	if in.evm.chainRules.IsByzantium {
		fuzz_helper.CoverTab[16619]++
		if in.readOnly {
			fuzz_helper.CoverTab[12692]++

			if operation.writes || (op == CALL && stack.Back(2).BitLen() > 0) {
				fuzz_helper.CoverTab[42483]++
				return errWriteProtection
			} else {
				fuzz_helper.CoverTab[6577]++
			}
		} else {
			fuzz_helper.CoverTab[17393]++
		}
	} else {
		fuzz_helper.CoverTab[64174]++
	}
	fuzz_helper.CoverTab[45213]++
	return nil
}

// Run loops and evaluates the contract's code with the given input data and returns
// the return byte-slice and an error if one occurred.
//
// It's important to note that any errors returned by the interpreter should be
// considered a revert-and-consume-all-gas operation. No error specific checks
// should be handled to reduce complexity and errors further down the in.
func (in *Interpreter) Run(snapshot int, contract *Contract, input []byte) (ret []byte, err error) {
	fuzz_helper.CoverTab[38740]++

	in.evm.depth++
	defer func() { fuzz_helper.CoverTab[49217]++; in.evm.depth-- }()
	fuzz_helper.CoverTab[35657]++

	in.returnData = nil

	if len(contract.Code) == 0 {
		fuzz_helper.CoverTab[34511]++
		return nil, nil
	} else {
		fuzz_helper.CoverTab[64074]++
	}
	fuzz_helper.CoverTab[30358]++

	codehash := contract.CodeHash
	if codehash == (common.Hash{}) {
		fuzz_helper.CoverTab[28614]++
		codehash = crypto.Keccak256Hash(contract.Code)
	} else {
		fuzz_helper.CoverTab[39226]++
	}
	fuzz_helper.CoverTab[23294]++

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
		fuzz_helper.CoverTab[2297]++
		if err != nil && in.cfg.Debug {
			fuzz_helper.CoverTab[40870]++
			in.cfg.Tracer.CaptureState(in.evm, pc, op, contract.Gas, cost, mem, stack, contract, in.evm.depth, err)
		} else {
			fuzz_helper.CoverTab[52877]++
		}
	}()
	fuzz_helper.CoverTab[61639]++

	for atomic.LoadInt32(&in.evm.abort) == 0 {
		fuzz_helper.CoverTab[778]++

		op = contract.GetOp(pc)

		operation := in.cfg.JumpTable[op]
		if err := in.enforceRestrictions(op, operation, stack); err != nil {
			fuzz_helper.CoverTab[264]++
			return nil, err
		} else {
			fuzz_helper.CoverTab[3566]++
		}
		fuzz_helper.CoverTab[33340]++

		if !operation.valid {
			fuzz_helper.CoverTab[47636]++
			return nil, fmt.Errorf("invalid opcode 0x%x", int(op))
		} else {
			fuzz_helper.CoverTab[8730]++
		}
		fuzz_helper.CoverTab[15638]++

		if err := operation.validateStack(stack); err != nil {
			fuzz_helper.CoverTab[20539]++
			return nil, err
		} else {
			fuzz_helper.CoverTab[63931]++
		}
		fuzz_helper.CoverTab[45869]++

		var memorySize uint64

		if operation.memorySize != nil {
			fuzz_helper.CoverTab[19009]++
			memSize, overflow := bigUint64(operation.memorySize(stack))
			if overflow {
				fuzz_helper.CoverTab[50446]++
				return nil, errGasUintOverflow
			} else {
				fuzz_helper.CoverTab[18500]++
			}
			fuzz_helper.CoverTab[64748]++

			if memorySize, overflow = math.SafeMul(toWordSize(memSize), 32); overflow {
				fuzz_helper.CoverTab[52152]++
				return nil, errGasUintOverflow
			} else {
				fuzz_helper.CoverTab[17111]++
			}
		} else {
			fuzz_helper.CoverTab[9670]++
		}
		fuzz_helper.CoverTab[23368]++

		if !in.cfg.DisableGasMetering {
			fuzz_helper.CoverTab[55848]++

			cost, err = operation.gasCost(in.gasTable, in.evm, contract, stack, mem, memorySize)
			if err != nil || !contract.UseGas(cost) {
				fuzz_helper.CoverTab[50755]++
				return nil, ErrOutOfGas
			} else {
				fuzz_helper.CoverTab[912]++
			}
		} else {
			fuzz_helper.CoverTab[64631]++
		}
		fuzz_helper.CoverTab[12901]++
		if memorySize > 0 {
			fuzz_helper.CoverTab[15513]++
			mem.Resize(memorySize)
		} else {
			fuzz_helper.CoverTab[17300]++
		}
		fuzz_helper.CoverTab[12499]++

		if in.cfg.Debug {
			fuzz_helper.CoverTab[16403]++
			in.cfg.Tracer.CaptureState(in.evm, pc, op, contract.Gas, cost, mem, stack, contract, in.evm.depth, err)
		} else {
			fuzz_helper.CoverTab[40937]++
		}
		fuzz_helper.CoverTab[42993]++

		res, err := operation.execute(&pc, in.evm, contract, mem, stack)

		if verifyPool {
			fuzz_helper.CoverTab[33825]++
			verifyIntegerPool(in.intPool)
		} else {
			fuzz_helper.CoverTab[7237]++
		}
		fuzz_helper.CoverTab[30301]++

		if operation.returns {
			fuzz_helper.CoverTab[23248]++
			in.returnData = res
		} else {
			fuzz_helper.CoverTab[52715]++
		}
		fuzz_helper.CoverTab[45210]++

		switch {
		case err != nil:
			fuzz_helper.CoverTab[11389]++
			return nil, err
		case operation.reverts:
			fuzz_helper.CoverTab[60629]++
			return res, errExecutionReverted
		case operation.halts:
			fuzz_helper.CoverTab[23245]++
			return res, nil
		case !operation.jumps:
			fuzz_helper.CoverTab[5383]++
			pc++
		default:
			fuzz_helper.CoverTab[52957]++
		}
	}
	fuzz_helper.CoverTab[11162]++
	return nil, nil
}

var _ = fuzz_helper.CoverTab
