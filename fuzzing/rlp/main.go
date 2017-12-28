package main

import "C"
import "bytes"
import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"
import rlp "github.com/ethereum/go-ethereum/rlp_instrumented"

//export SetInstrumentationType
func SetInstrumentationType(t int) {
    fuzz_helper.SetInstrumentationType(t)
}
//export GoResetCoverage
func GoResetCoverage() {
    fuzz_helper.ResetCoverage()
}

//export GoCalcCoverage
func GoCalcCoverage() uint64 {
    return fuzz_helper.CalcCoverage()
}

//export decode_rlp
func decode_rlp(input []byte, mode int) {
	mode = mode % 5
	if mode == 0 {
		msg := new(interface{})
		err := rlp.DecodeBytes(input, msg)
		if err == nil {
			rlp.EncodeToBytes(msg)
		}
	} else if mode == 1 && len(input) > 0 {
		rlp.Split(input)
	} else if mode == 2 && len(input) > 0 {
		elems, _, err := rlp.SplitList(input)
		if err == nil {
			rlp.CountValues(elems)
		}
	} else if mode == 3 {
		s := rlp.NewStream(bytes.NewReader(input), 0)
		msg := new(interface{})
		s.Decode(msg)
	} else if mode == 4 {
		type AllTypes struct{
			Int uint
			String string
			Bytes []byte
			Bool bool
			Raw rlp.RawValue
			Slice []*AllTypes
			Array [3]*AllTypes
			Iface []interface{}
		}

		var v AllTypes
        err := rlp.DecodeBytes(input, &v)
		if err == nil {
			rlp.EncodeToBytes(v)
		}
    }
}

/* No main() body because this file is compiled to a static archive */
func main() {
}
