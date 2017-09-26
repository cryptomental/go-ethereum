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
	fuzz_helper.AddCoverage(22588)
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
		fuzz_helper.AddCoverage(5262)
		enc.Stack = make([]*math.HexOrDecimal256, len(s.Stack))
		for k, v := range s.Stack {
			fuzz_helper.AddCoverage(17878)
			enc.Stack[k] = (*math.HexOrDecimal256)(v)
		}
	} else {
		fuzz_helper.AddCoverage(45021)
	}
	fuzz_helper.AddCoverage(44810)
	enc.Storage = s.Storage
	enc.Depth = s.Depth
	enc.Err = s.Err
	enc.OpName = s.OpName()
	return json.Marshal(&enc)
}

func (s *StructLog) UnmarshalJSON(input []byte) error {
	fuzz_helper.AddCoverage(39040)
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
		fuzz_helper.AddCoverage(30358)
		return err
	} else {
		fuzz_helper.AddCoverage(23294)
	}
	fuzz_helper.AddCoverage(2095)
	if dec.Pc != nil {
		fuzz_helper.AddCoverage(61639)
		s.Pc = *dec.Pc
	} else {
		fuzz_helper.AddCoverage(11162)
	}
	fuzz_helper.AddCoverage(21668)
	if dec.Op != nil {
		fuzz_helper.AddCoverage(49217)
		s.Op = *dec.Op
	} else {
		fuzz_helper.AddCoverage(34511)
	}
	fuzz_helper.AddCoverage(45213)
	if dec.Gas != nil {
		fuzz_helper.AddCoverage(64074)
		s.Gas = uint64(*dec.Gas)
	} else {
		fuzz_helper.AddCoverage(28614)
	}
	fuzz_helper.AddCoverage(16619)
	if dec.GasCost != nil {
		fuzz_helper.AddCoverage(39226)
		s.GasCost = uint64(*dec.GasCost)
	} else {
		fuzz_helper.AddCoverage(2297)
	}
	fuzz_helper.AddCoverage(12692)
	if dec.Memory != nil {
		fuzz_helper.AddCoverage(40870)
		s.Memory = dec.Memory
	} else {
		fuzz_helper.AddCoverage(52877)
	}
	fuzz_helper.AddCoverage(42483)
	if dec.MemorySize != nil {
		fuzz_helper.AddCoverage(778)
		s.MemorySize = *dec.MemorySize
	} else {
		fuzz_helper.AddCoverage(33340)
	}
	fuzz_helper.AddCoverage(6577)
	if dec.Stack != nil {
		fuzz_helper.AddCoverage(15638)
		s.Stack = make([]*big.Int, len(dec.Stack))
		for k, v := range dec.Stack {
			fuzz_helper.AddCoverage(45869)
			s.Stack[k] = (*big.Int)(v)
		}
	} else {
		fuzz_helper.AddCoverage(23368)
	}
	fuzz_helper.AddCoverage(17393)
	if dec.Storage != nil {
		fuzz_helper.AddCoverage(12901)
		s.Storage = dec.Storage
	} else {
		fuzz_helper.AddCoverage(12499)
	}
	fuzz_helper.AddCoverage(64174)
	if dec.Depth != nil {
		fuzz_helper.AddCoverage(42993)
		s.Depth = *dec.Depth
	} else {
		fuzz_helper.AddCoverage(30301)
	}
	fuzz_helper.AddCoverage(38740)
	if dec.Err != nil {
		fuzz_helper.AddCoverage(45210)
		s.Err = *dec.Err
	} else {
		fuzz_helper.AddCoverage(264)
	}
	fuzz_helper.AddCoverage(35657)
	return nil
}

var _ = fuzz_helper.AddCoverage
