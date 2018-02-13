package abi_fuzzing

import "C"
import "math/big"

var g_calldataloads = make([]*big.Int, 0)

var (
    Enabled = false
)

func EnableABIFuzzing() {
    Enabled = true
}

func AddCallDataLoad(val *big.Int) {
    num := new(big.Int).Set(val)
    g_calldataloads = append(g_calldataloads, num)
}

func ResetCallDataLoads() {
    g_calldataloads = nil
}

func GetCallDataLoads() []*big.Int {
    return g_calldataloads
}
