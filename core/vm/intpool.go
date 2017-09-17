package vm

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import "math/big"

var checkVal = big.NewInt(-42)

const poolLimit = 256

// intPool is a pool of big integers that
// can be reused for all big.Int operations.
type intPool struct {
	pool *Stack
}

func newIntPool() *intPool {
	fuzz_helper.CoverTab[22588]++
	return &intPool{pool: newstack()}
}

func (p *intPool) get() *big.Int {
	fuzz_helper.CoverTab[44810]++
	if p.pool.len() > 0 {
		fuzz_helper.CoverTab[17878]++
		return p.pool.pop()
	} else {
		fuzz_helper.CoverTab[45021]++
	}
	fuzz_helper.CoverTab[5262]++
	return new(big.Int)
}
func (p *intPool) put(is ...*big.Int) {
	fuzz_helper.CoverTab[39040]++
	if len(p.pool.data) > poolLimit {
		fuzz_helper.CoverTab[21668]++
		return
	} else {
		fuzz_helper.CoverTab[45213]++
	}
	fuzz_helper.CoverTab[2095]++

	for _, i := range is {
		fuzz_helper.CoverTab[16619]++

		if verifyPool {
			fuzz_helper.CoverTab[42483]++
			i.Set(checkVal)
		} else {
			fuzz_helper.CoverTab[6577]++
		}
		fuzz_helper.CoverTab[12692]++

		p.pool.push(i)
	}
}

var _ = fuzz_helper.CoverTab
