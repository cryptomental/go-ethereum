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
	fuzz_helper.AddCoverage(15763)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return &intPool{pool: newstack()}
}

func (p *intPool) get() *big.Int {
	fuzz_helper.AddCoverage(48468)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	if p.pool.len() > 0 {
		fuzz_helper.AddCoverage(64618)
		return p.pool.pop()
	} else {
		fuzz_helper.AddCoverage(64058)
	}
	fuzz_helper.AddCoverage(45453)
	return new(big.Int)
}
func (p *intPool) put(is ...*big.Int) {
	fuzz_helper.AddCoverage(61123)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	if len(p.pool.data) > poolLimit {
		fuzz_helper.AddCoverage(9751)
		return
	} else {
		fuzz_helper.AddCoverage(13627)
	}
	fuzz_helper.AddCoverage(51006)

	for _, i := range is {
		fuzz_helper.AddCoverage(55286)

		if verifyPool {
			fuzz_helper.AddCoverage(36956)
			i.Set(checkVal)
		} else {
			fuzz_helper.AddCoverage(19191)
		}
		fuzz_helper.AddCoverage(47586)

		p.pool.push(i)
	}
}

var _ = fuzz_helper.AddCoverage
