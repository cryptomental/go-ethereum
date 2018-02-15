package abi_fuzzing

import "C"
import "math/big"

type CallDataVars struct {
    Opcode uint64
    Size *big.Int
    Offset *big.Int
    To *big.Int
}

var g_calldatavars = make([]CallDataVars, 0)

var (
    Enabled = false
)

func EnableABIFuzzing() {
    Enabled = true
}

func AddCallDataOp(opcode uint64, size *big.Int, offset *big.Int, to *big.Int) {
    size_copy := new(big.Int).Set(size)
    offset_copy := new(big.Int).Set(offset)
    to_copy := new(big.Int)
    if to != nil {
        to_copy.Set(to)
    }
    g_calldatavars = append(g_calldatavars, CallDataVars{opcode, size_copy, offset_copy, to_copy})
}

func ResetCallDataLoads() {
    g_calldatavars = nil
}

func GetCallDataLoads() []CallDataVars {
    return g_calldatavars
}
