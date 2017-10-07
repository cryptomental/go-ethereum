package vm

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
)

func (s StructLog) MarshalJSON() ([]byte, error) {
	fuzz_helper.AddCoverage(25636)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	type StructLog struct {
		Pc         uint64                      `json:"pc"`
		Op         OpCode                      `json:"op"`
		Gas        math.HexOrDecimal64         `json:"gas"`
		GasCost    math.HexOrDecimal64         `json:"gasCost"`
		Memory     hexutil.Bytes               `json:"memory"`
		MemorySize int                         `json:"memSize"`
		Stack      []*math.HexOrDecimal256     `json:"stack"`
		Storage    map[common.Hash]common.Hash `json:"-"`
		Depth      int                         `json:"depth"`
		Err        error                       `json:"error"`
		OpName     string                      `json:"opName"`
	}
	var enc StructLog
	enc.Pc = s.Pc
	enc.Op = s.Op
	enc.Gas = math.HexOrDecimal64(s.Gas)
	enc.GasCost = math.HexOrDecimal64(s.GasCost)
	enc.Memory = s.Memory
	enc.MemorySize = s.MemorySize
	if s.Stack != nil {
		fuzz_helper.AddCoverage(6081)
		enc.Stack = make([]*math.HexOrDecimal256, len(s.Stack))
		for k, v := range s.Stack {
			fuzz_helper.AddCoverage(54249)
			enc.Stack[k] = (*math.HexOrDecimal256)(v)
		}
	} else {
		fuzz_helper.AddCoverage(25709)
	}
	fuzz_helper.AddCoverage(12742)
	enc.Storage = s.Storage
	enc.Depth = s.Depth
	enc.Err = s.Err
	enc.OpName = s.OpName()
	return json.Marshal(&enc)
}

func (s *StructLog) UnmarshalJSON(input []byte) error {
	fuzz_helper.AddCoverage(10433)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	type StructLog struct {
		Pc         *uint64                     `json:"pc"`
		Op         *OpCode                     `json:"op"`
		Gas        *math.HexOrDecimal64        `json:"gas"`
		GasCost    *math.HexOrDecimal64        `json:"gasCost"`
		Memory     hexutil.Bytes               `json:"memory"`
		MemorySize *int                        `json:"memSize"`
		Stack      []*math.HexOrDecimal256     `json:"stack"`
		Storage    map[common.Hash]common.Hash `json:"-"`
		Depth      *int                        `json:"depth"`
		Err        *error                      `json:"error"`
	}
	var dec StructLog
	if err := json.Unmarshal(input, &dec); err != nil {
		fuzz_helper.AddCoverage(58192)
		return err
	} else {
		fuzz_helper.AddCoverage(6264)
	}
	fuzz_helper.AddCoverage(53609)
	if dec.Pc != nil {
		fuzz_helper.AddCoverage(47208)
		s.Pc = *dec.Pc
	} else {
		fuzz_helper.AddCoverage(34399)
	}
	fuzz_helper.AddCoverage(27059)
	if dec.Op != nil {
		fuzz_helper.AddCoverage(41921)
		s.Op = *dec.Op
	} else {
		fuzz_helper.AddCoverage(28285)
	}
	fuzz_helper.AddCoverage(34057)
	if dec.Gas != nil {
		fuzz_helper.AddCoverage(17422)
		s.Gas = uint64(*dec.Gas)
	} else {
		fuzz_helper.AddCoverage(56637)
	}
	fuzz_helper.AddCoverage(45816)
	if dec.GasCost != nil {
		fuzz_helper.AddCoverage(39095)
		s.GasCost = uint64(*dec.GasCost)
	} else {
		fuzz_helper.AddCoverage(28179)
	}
	fuzz_helper.AddCoverage(44693)
	if dec.Memory != nil {
		fuzz_helper.AddCoverage(3887)
		s.Memory = dec.Memory
	} else {
		fuzz_helper.AddCoverage(51012)
	}
	fuzz_helper.AddCoverage(29622)
	if dec.MemorySize != nil {
		fuzz_helper.AddCoverage(17785)
		s.MemorySize = *dec.MemorySize
	} else {
		fuzz_helper.AddCoverage(39631)
	}
	fuzz_helper.AddCoverage(29487)
	if dec.Stack != nil {
		fuzz_helper.AddCoverage(22481)
		s.Stack = make([]*big.Int, len(dec.Stack))
		for k, v := range dec.Stack {
			fuzz_helper.AddCoverage(36927)
			s.Stack[k] = (*big.Int)(v)
		}
	} else {
		fuzz_helper.AddCoverage(41592)
	}
	fuzz_helper.AddCoverage(60994)
	if dec.Storage != nil {
		fuzz_helper.AddCoverage(63414)
		s.Storage = dec.Storage
	} else {
		fuzz_helper.AddCoverage(34780)
	}
	fuzz_helper.AddCoverage(5594)
	if dec.Depth != nil {
		fuzz_helper.AddCoverage(18579)
		s.Depth = *dec.Depth
	} else {
		fuzz_helper.AddCoverage(15389)
	}
	fuzz_helper.AddCoverage(35504)
	if dec.Err != nil {
		fuzz_helper.AddCoverage(63602)
		s.Err = *dec.Err
	} else {
		fuzz_helper.AddCoverage(27256)
	}
	fuzz_helper.AddCoverage(57201)
	return nil
}

var _ = fuzz_helper.AddCoverage
