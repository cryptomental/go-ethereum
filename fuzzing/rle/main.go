package main

import "C"
import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"
import rle "github.com/ethereum/go-ethereum/compression/rle"

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

//export run_rle
func run_rle(input []byte, mode int) {
	mode = mode % 2
	if mode == 0 {
		rle.Decompress(input)
	} else if mode == 1 {
		rle.Compress(input)
	}
}

/* No main() body because this file is compiled to a static archive */
func main() {
}
