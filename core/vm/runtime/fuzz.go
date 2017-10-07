package runtime

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

// Fuzz is the basic entry point for the go-fuzz tool
//
// This returns 1 for valid parsable/runable code, 0
// for invalid opcode.
func Fuzz(input []byte) int {
	fuzz_helper.AddCoverage(63359)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	_, _, err := Execute(input, input, &Config{
		GasLimit: 3000000,
	})

	if err != nil && len(err.Error()) > 6 && string(err.Error()[:7]) == "invalid" {
		fuzz_helper.AddCoverage(63)
		return 0
	} else {
		fuzz_helper.AddCoverage(13937)
	}
	fuzz_helper.AddCoverage(62421)

	return 1
}

var _ = fuzz_helper.AddCoverage
