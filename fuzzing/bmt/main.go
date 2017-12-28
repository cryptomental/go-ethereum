package main

import "C"
import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"
import bmt "github.com/ethereum/go-ethereum/bmt_instrumented"
import "github.com/ethereum/go-ethereum/crypto/sha3"

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

//export run_bmt
func run_bmt(input []byte, mode int) {
	hasher := sha3.NewKeccak256
	pool := bmt.NewTreePool(hasher, 128, 1)
	bmt := bmt.New(pool)
	bmt.Reset()
	bmt.Write(input)
	bmt.Sum(nil)
}

/* No main() body because this file is compiled to a static archive */
func main() {
}
