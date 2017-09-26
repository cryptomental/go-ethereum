package runtime

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

// Fuzz is the basic entry point for the go-fuzz tool
//
// This returns 1 for valid parsable/runable code, 0
// for invalid opcode.
func Fuzz(input []byte) int {
	fuzz_helper.AddCoverage(22588)
	_, _, err := Execute(input, input, &Config{
		GasLimit: 3000000,
	})

	if err != nil && len(err.Error()) > 6 && string(err.Error()[:7]) == "invalid" {
		fuzz_helper.AddCoverage(5262)
		return 0
	} else {
		fuzz_helper.AddCoverage(17878)
	}
	fuzz_helper.AddCoverage(44810)

	return 1
}

var _ = fuzz_helper.AddCoverage
