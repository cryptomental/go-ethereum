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
	fuzz_helper.AddCoverage(25005)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	cpy := make(Storage)
	for key, value := range self {
		fuzz_helper.AddCoverage(45659)
		cpy[key] = value
	}
	fuzz_helper.AddCoverage(36849)

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
	fuzz_helper.AddCoverage(16063)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
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
	fuzz_helper.AddCoverage(35235)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	logger := &StructLogger{
		changedValues: make(map[common.Address]Storage),
	}
	if cfg != nil {
		fuzz_helper.AddCoverage(17202)
		logger.cfg = *cfg
	} else {
		fuzz_helper.AddCoverage(26553)
	}
	fuzz_helper.AddCoverage(40694)
	return logger
}

func (l *StructLogger) CaptureState(env *EVM, pc uint64, op OpCode, gas, cost uint64, memory *Memory, stack *Stack, contract *Contract, depth int, err error) error {
	fuzz_helper.AddCoverage(19645)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()

	if l.cfg.Limit != 0 && l.cfg.Limit <= len(l.logs) {
		fuzz_helper.AddCoverage(54184)
		return ErrTraceLimitReached
	} else {
		fuzz_helper.AddCoverage(15617)
	}
	fuzz_helper.AddCoverage(23691)

	if l.changedValues[contract.Address()] == nil {
		fuzz_helper.AddCoverage(23241)
		l.changedValues[contract.Address()] = make(Storage)
	} else {
		fuzz_helper.AddCoverage(9541)
	}
	fuzz_helper.AddCoverage(24413)

	if op == SSTORE && stack.len() >= 2 {
		fuzz_helper.AddCoverage(37781)
		var (
			value   = common.BigToHash(stack.data[stack.len()-2])
			address = common.BigToHash(stack.data[stack.len()-1])
		)
		l.changedValues[contract.Address()][address] = value
	} else {
		fuzz_helper.AddCoverage(27899)
	}
	fuzz_helper.AddCoverage(59115)

	var mem []byte
	if !l.cfg.DisableMemory {
		fuzz_helper.AddCoverage(4930)
		mem = make([]byte, len(memory.Data()))
		copy(mem, memory.Data())
	} else {
		fuzz_helper.AddCoverage(15533)
	}
	fuzz_helper.AddCoverage(11919)

	var stck []*big.Int
	if !l.cfg.DisableStack {
		fuzz_helper.AddCoverage(45727)
		stck = make([]*big.Int, len(stack.Data()))
		for i, item := range stack.Data() {
			fuzz_helper.AddCoverage(50494)
			stck[i] = new(big.Int).Set(item)
		}
	} else {
		fuzz_helper.AddCoverage(22109)
	}
	fuzz_helper.AddCoverage(64269)

	var storage Storage
	if !l.cfg.DisableStorage {
		fuzz_helper.AddCoverage(57486)
		if l.cfg.FullStorage {
			fuzz_helper.AddCoverage(36021)
			storage = make(Storage)

			env.StateDB.ForEachStorage(contract.Address(), func(key, value common.Hash) bool {
				fuzz_helper.AddCoverage(51702)
				storage[key] = value

				return true
			})
		} else {
			fuzz_helper.AddCoverage(64906)

			storage = l.changedValues[contract.Address()].Copy()
		}
	} else {
		fuzz_helper.AddCoverage(37749)
	}
	fuzz_helper.AddCoverage(51530)

	log := StructLog{pc, op, gas, cost, mem, memory.Len(), stck, storage, depth, err}

	l.logs = append(l.logs, log)
	return nil
}

func (l *StructLogger) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) error {
	fuzz_helper.AddCoverage(49487)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	fmt.Printf("0x%x", output)
	if err != nil {
		fuzz_helper.AddCoverage(48118)
		fmt.Printf(" error: %v\n", err)
	} else {
		fuzz_helper.AddCoverage(65279)
	}
	fuzz_helper.AddCoverage(7661)
	return nil
}

func (l *StructLogger) StructLogs() []StructLog {
	fuzz_helper.AddCoverage(15949)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return l.logs
}

func WriteTrace(writer io.Writer, logs []StructLog) {
	fuzz_helper.AddCoverage(8398)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	for _, log := range logs {
		fuzz_helper.AddCoverage(63898)
		fmt.Fprintf(writer, "%-16spc=%08d gas=%v cost=%v", log.Op, log.Pc, log.Gas, log.GasCost)
		if log.Err != nil {
			fuzz_helper.AddCoverage(38067)
			fmt.Fprintf(writer, " ERROR: %v", log.Err)
		} else {
			fuzz_helper.AddCoverage(34912)
		}
		fuzz_helper.AddCoverage(45683)
		fmt.Fprintln(writer)

		if len(log.Stack) > 0 {
			fuzz_helper.AddCoverage(4382)
			fmt.Fprintln(writer, "Stack:")
			for i := len(log.Stack) - 1; i >= 0; i-- {
				fuzz_helper.AddCoverage(16441)
				fmt.Fprintf(writer, "%08d  %x\n", len(log.Stack)-i-1, math.PaddedBigBytes(log.Stack[i], 32))
			}
		} else {
			fuzz_helper.AddCoverage(26948)
		}
		fuzz_helper.AddCoverage(35150)
		if len(log.Memory) > 0 {
			fuzz_helper.AddCoverage(42589)
			fmt.Fprintln(writer, "Memory:")
			fmt.Fprint(writer, hex.Dump(log.Memory))
		} else {
			fuzz_helper.AddCoverage(26675)
		}
		fuzz_helper.AddCoverage(42912)
		if len(log.Storage) > 0 {
			fuzz_helper.AddCoverage(16788)
			fmt.Fprintln(writer, "Storage:")
			for h, item := range log.Storage {
				fuzz_helper.AddCoverage(42712)
				fmt.Fprintf(writer, "%x: %x\n", h, item)
			}
		} else {
			fuzz_helper.AddCoverage(64109)
		}
		fuzz_helper.AddCoverage(23188)
		fmt.Fprintln(writer)
	}
}

func WriteLogs(writer io.Writer, logs []*types.Log) {
	fuzz_helper.AddCoverage(55235)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	for _, log := range logs {
		fuzz_helper.AddCoverage(15730)
		fmt.Fprintf(writer, "LOG%d: %x bn=%d txi=%x\n", len(log.Topics), log.Address, log.BlockNumber, log.TxIndex)

		for i, topic := range log.Topics {
			fuzz_helper.AddCoverage(37055)
			fmt.Fprintf(writer, "%08d  %x\n", i, topic)
		}
		fuzz_helper.AddCoverage(37737)

		fmt.Fprint(writer, hex.Dump(log.Data))
		fmt.Fprintln(writer)
	}
}

var _ = fuzz_helper.AddCoverage
