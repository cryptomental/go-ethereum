package vm

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"fmt"

	"github.com/ethereum/go-ethereum/params"
)

func makeStackFunc(pop, push int) stackValidationFunc {
	fuzz_helper.AddCoverage(22588)
	return func(stack *Stack) error {
		fuzz_helper.AddCoverage(44810)
		if err := stack.require(pop); err != nil {
			fuzz_helper.AddCoverage(45021)
			return err
		} else {
			fuzz_helper.AddCoverage(39040)
		}
		fuzz_helper.AddCoverage(5262)

		if stack.len()+push-pop > int(params.StackLimit) {
			fuzz_helper.AddCoverage(2095)
			return fmt.Errorf("stack limit reached %d (%d)", stack.len(), params.StackLimit)
		} else {
			fuzz_helper.AddCoverage(21668)
		}
		fuzz_helper.AddCoverage(17878)
		return nil
	}
}

func makeDupStackFunc(n int) stackValidationFunc {
	fuzz_helper.AddCoverage(45213)
	return makeStackFunc(n, n+1)
}

func makeSwapStackFunc(n int) stackValidationFunc {
	fuzz_helper.AddCoverage(16619)
	return makeStackFunc(n, n)
}

var _ = fuzz_helper.AddCoverage
