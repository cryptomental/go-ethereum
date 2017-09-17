package vm

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
)

type Storage map[common.Hash]common.Hash

func (self Storage) Copy() Storage {
	fuzz_helper.CoverTab[22588]++
	cpy := make(Storage)
	for key, value := range self {
		fuzz_helper.CoverTab[5262]++
		cpy[key] = value
	}
	fuzz_helper.CoverTab[44810]++

	return cpy
}

type LogConfig struct {
	DisableMemory  bool
	DisableStack   bool
	DisableStorage bool
	FullStorage    bool
	Limit          int
}

//go:generate gencodec -type StructLog -field-override structLogMarshaling -out gen_structlog.go

type StructLog struct {
	Pc         uint64                      `json:"pc"`
	Op         OpCode                      `json:"op"`
	Gas        uint64                      `json:"gas"`
	GasCost    uint64                      `json:"gasCost"`
	Memory     []byte                      `json:"memory"`
	MemorySize int                         `json:"memSize"`
	Stack      []*big.Int                  `json:"stack"`
	Storage    map[common.Hash]common.Hash `json:"-"`
	Depth      int                         `json:"depth"`
	Err        error                       `json:"error"`
}

type structLogMarshaling struct {
	Stack   []*math.HexOrDecimal256
	Gas     math.HexOrDecimal64
	GasCost math.HexOrDecimal64
	Memory  hexutil.Bytes
	OpName  string `json:"opName"`
}

func (s *StructLog) OpName() string {
	fuzz_helper.CoverTab[17878]++
	return s.Op.String()
}

type Tracer interface {
	CaptureState(env *EVM, pc uint64, op OpCode, gas, cost uint64, memory *Memory, stack *Stack, contract *Contract, depth int, err error) error
	CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) error
}

type StructLogger struct {
	cfg LogConfig

	logs          []StructLog
	changedValues map[common.Address]Storage
}

func NewStructLogger(cfg *LogConfig) *StructLogger {
	fuzz_helper.CoverTab[45021]++
	logger := &StructLogger{
		changedValues: make(map[common.Address]Storage),
	}
	if cfg != nil {
		fuzz_helper.CoverTab[2095]++
		logger.cfg = *cfg
	} else {
		fuzz_helper.CoverTab[21668]++
	}
	fuzz_helper.CoverTab[39040]++
	return logger
}

func (l *StructLogger) CaptureState(env *EVM, pc uint64, op OpCode, gas, cost uint64, memory *Memory, stack *Stack, contract *Contract, depth int, err error) error {
	fuzz_helper.CoverTab[45213]++

	if l.cfg.Limit != 0 && l.cfg.Limit <= len(l.logs) {
		fuzz_helper.CoverTab[38740]++
		return ErrTraceLimitReached
	} else {
		fuzz_helper.CoverTab[35657]++
	}
	fuzz_helper.CoverTab[16619]++

	if l.changedValues[contract.Address()] == nil {
		fuzz_helper.CoverTab[30358]++
		l.changedValues[contract.Address()] = make(Storage)
	} else {
		fuzz_helper.CoverTab[23294]++
	}
	fuzz_helper.CoverTab[12692]++

	if op == SSTORE && stack.len() >= 2 {
		fuzz_helper.CoverTab[61639]++
		var (
			value   = common.BigToHash(stack.data[stack.len()-2])
			address = common.BigToHash(stack.data[stack.len()-1])
		)
		l.changedValues[contract.Address()][address] = value
	} else {
		fuzz_helper.CoverTab[11162]++
	}
	fuzz_helper.CoverTab[42483]++

	var mem []byte
	if !l.cfg.DisableMemory {
		fuzz_helper.CoverTab[49217]++
		mem = make([]byte, len(memory.Data()))
		copy(mem, memory.Data())
	} else {
		fuzz_helper.CoverTab[34511]++
	}
	fuzz_helper.CoverTab[6577]++

	var stck []*big.Int
	if !l.cfg.DisableStack {
		fuzz_helper.CoverTab[64074]++
		stck = make([]*big.Int, len(stack.Data()))
		for i, item := range stack.Data() {
			fuzz_helper.CoverTab[28614]++
			stck[i] = new(big.Int).Set(item)
		}
	} else {
		fuzz_helper.CoverTab[39226]++
	}
	fuzz_helper.CoverTab[17393]++

	var storage Storage
	if !l.cfg.DisableStorage {
		fuzz_helper.CoverTab[2297]++
		if l.cfg.FullStorage {
			fuzz_helper.CoverTab[40870]++
			storage = make(Storage)

			env.StateDB.ForEachStorage(contract.Address(), func(key, value common.Hash) bool {
				fuzz_helper.CoverTab[52877]++
				storage[key] = value

				return true
			})
		} else {
			fuzz_helper.CoverTab[778]++

			storage = l.changedValues[contract.Address()].Copy()
		}
	} else {
		fuzz_helper.CoverTab[33340]++
	}
	fuzz_helper.CoverTab[64174]++

	log := StructLog{pc, op, gas, cost, mem, memory.Len(), stck, storage, depth, err}

	l.logs = append(l.logs, log)
	return nil
}

func (l *StructLogger) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) error {
	fuzz_helper.CoverTab[15638]++
	fmt.Printf("0x%x", output)
	if err != nil {
		fuzz_helper.CoverTab[23368]++
		fmt.Printf(" error: %v\n", err)
	} else {
		fuzz_helper.CoverTab[12901]++
	}
	fuzz_helper.CoverTab[45869]++
	return nil
}

func (l *StructLogger) StructLogs() []StructLog {
	fuzz_helper.CoverTab[12499]++
	return l.logs
}

func WriteTrace(writer io.Writer, logs []StructLog) {
	fuzz_helper.CoverTab[42993]++
	for _, log := range logs {
		fuzz_helper.CoverTab[30301]++
		fmt.Fprintf(writer, "%-16spc=%08d gas=%v cost=%v", log.Op, log.Pc, log.Gas, log.GasCost)
		if log.Err != nil {
			fuzz_helper.CoverTab[8730]++
			fmt.Fprintf(writer, " ERROR: %v", log.Err)
		} else {
			fuzz_helper.CoverTab[20539]++
		}
		fuzz_helper.CoverTab[45210]++
		fmt.Fprintln(writer)

		if len(log.Stack) > 0 {
			fuzz_helper.CoverTab[63931]++
			fmt.Fprintln(writer, "Stack:")
			for i := len(log.Stack) - 1; i >= 0; i-- {
				fuzz_helper.CoverTab[19009]++
				fmt.Fprintf(writer, "%08d  %x\n", len(log.Stack)-i-1, math.PaddedBigBytes(log.Stack[i], 32))
			}
		} else {
			fuzz_helper.CoverTab[64748]++
		}
		fuzz_helper.CoverTab[264]++
		if len(log.Memory) > 0 {
			fuzz_helper.CoverTab[50446]++
			fmt.Fprintln(writer, "Memory:")
			fmt.Fprint(writer, hex.Dump(log.Memory))
		} else {
			fuzz_helper.CoverTab[18500]++
		}
		fuzz_helper.CoverTab[3566]++
		if len(log.Storage) > 0 {
			fuzz_helper.CoverTab[52152]++
			fmt.Fprintln(writer, "Storage:")
			for h, item := range log.Storage {
				fuzz_helper.CoverTab[17111]++
				fmt.Fprintf(writer, "%x: %x\n", h, item)
			}
		} else {
			fuzz_helper.CoverTab[9670]++
		}
		fuzz_helper.CoverTab[47636]++
		fmt.Fprintln(writer)
	}
}

func WriteLogs(writer io.Writer, logs []*types.Log) {
	fuzz_helper.CoverTab[55848]++
	for _, log := range logs {
		fuzz_helper.CoverTab[50755]++
		fmt.Fprintf(writer, "LOG%d: %x bn=%d txi=%x\n", len(log.Topics), log.Address, log.BlockNumber, log.TxIndex)

		for i, topic := range log.Topics {
			fuzz_helper.CoverTab[64631]++
			fmt.Fprintf(writer, "%08d  %x\n", i, topic)
		}
		fuzz_helper.CoverTab[912]++

		fmt.Fprint(writer, hex.Dump(log.Data))
		fmt.Fprintln(writer)
	}
}

var _ = fuzz_helper.CoverTab
