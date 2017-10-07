package vm

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"fmt"

	"github.com/ethereum/go-ethereum/params"
)

func makeStackFunc(pop, push int) stackValidationFunc {
	fuzz_helper.AddCoverage(9543)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return func(stack *Stack) error {
		fuzz_helper.AddCoverage(50572)
		if err := stack.require(pop); err != nil {
			fuzz_helper.AddCoverage(28042)
			return err
		} else {
			fuzz_helper.AddCoverage(18347)
		}
		fuzz_helper.AddCoverage(6005)

		if stack.len()+push-pop > int(params.StackLimit) {
			fuzz_helper.AddCoverage(16061)
			return fmt.Errorf("stack limit reached %d (%d)", stack.len(), params.StackLimit)
		} else {
			fuzz_helper.AddCoverage(38767)
		}
		fuzz_helper.AddCoverage(22317)
		return nil
	}
}

func makeDupStackFunc(n int) stackValidationFunc {
	fuzz_helper.AddCoverage(51553)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return makeStackFunc(n, n+1)
}

func makeSwapStackFunc(n int) stackValidationFunc {
	fuzz_helper.AddCoverage(26366)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return makeStackFunc(n, n)
}

var _ = fuzz_helper.AddCoverage
