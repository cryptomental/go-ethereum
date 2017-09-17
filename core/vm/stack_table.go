package vm

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"fmt"

	"github.com/ethereum/go-ethereum/params"
)

func makeStackFunc(pop, push int) stackValidationFunc {
	fuzz_helper.CoverTab[22588]++
	return func(stack *Stack) error {
		fuzz_helper.CoverTab[44810]++
		if err := stack.require(pop); err != nil {
			fuzz_helper.CoverTab[45021]++
			return err
		} else {
			fuzz_helper.CoverTab[39040]++
		}
		fuzz_helper.CoverTab[5262]++

		if stack.len()+push-pop > int(params.StackLimit) {
			fuzz_helper.CoverTab[2095]++
			return fmt.Errorf("stack limit reached %d (%d)", stack.len(), params.StackLimit)
		} else {
			fuzz_helper.CoverTab[21668]++
		}
		fuzz_helper.CoverTab[17878]++
		return nil
	}
}

func makeDupStackFunc(n int) stackValidationFunc {
	fuzz_helper.CoverTab[45213]++
	return makeStackFunc(n, n+1)
}

func makeSwapStackFunc(n int) stackValidationFunc {
	fuzz_helper.CoverTab[16619]++
	return makeStackFunc(n, n)
}

var _ = fuzz_helper.CoverTab
