package vm

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"fmt"
	"math/big"
)

// stack is an object for basic stack operations. Items popped to the stack are
// expected to be changed and modified. stack does not take care of adding newly
// initialised objects.
type Stack struct {
	data []*big.Int
}

func newstack() *Stack {
	fuzz_helper.AddCoverage(43618)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return &Stack{data: make([]*big.Int, 0, 1024)}
}

func (st *Stack) Data() []*big.Int {
	fuzz_helper.AddCoverage(40159)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return st.data
}

func (st *Stack) push(d *big.Int) {
	fuzz_helper.AddCoverage(19222)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()

	st.data = append(st.data, d)
}
func (st *Stack) pushN(ds ...*big.Int) {
	fuzz_helper.AddCoverage(7173)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	st.data = append(st.data, ds...)
}

func (st *Stack) pop() (ret *big.Int) {
	fuzz_helper.AddCoverage(53159)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	ret = st.data[len(st.data)-1]
	st.data = st.data[:len(st.data)-1]
	return
}

func (st *Stack) len() int {
	fuzz_helper.AddCoverage(12977)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return len(st.data)
}

func (st *Stack) swap(n int) {
	fuzz_helper.AddCoverage(4281)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	st.data[st.len()-n], st.data[st.len()-1] = st.data[st.len()-1], st.data[st.len()-n]
}

func (st *Stack) dup(pool *intPool, n int) {
	fuzz_helper.AddCoverage(19349)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	st.push(pool.get().Set(st.data[st.len()-n]))
}

func (st *Stack) peek() *big.Int {
	fuzz_helper.AddCoverage(29465)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return st.data[st.len()-1]
}

// Back returns the n'th item in stack
func (st *Stack) Back(n int) *big.Int {
	fuzz_helper.AddCoverage(46863)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return st.data[st.len()-n-1]
}

func (st *Stack) require(n int) error {
	fuzz_helper.AddCoverage(2065)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	if st.len() < n {
		fuzz_helper.AddCoverage(58769)
		return fmt.Errorf("stack underflow (%d <=> %d)", len(st.data), n)
	} else {
		fuzz_helper.AddCoverage(45067)
	}
	fuzz_helper.AddCoverage(689)
	return nil
}

func (st *Stack) Print() {
	fuzz_helper.AddCoverage(17061)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	fmt.Println("### stack ###")
	if len(st.data) > 0 {
		fuzz_helper.AddCoverage(26638)
		for i, val := range st.data {
			fuzz_helper.AddCoverage(33507)
			fmt.Printf("%-3d  %v\n", i, val)
		}
	} else {
		fuzz_helper.AddCoverage(3939)
		fmt.Println("-- empty --")
	}
	fuzz_helper.AddCoverage(57324)
	fmt.Println("#############")
}

var _ = fuzz_helper.AddCoverage
