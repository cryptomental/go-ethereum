package vm

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import "testing"

func TestMemoryGasCost(t *testing.T) {
	fuzz_helper.AddCoverage(22588)

	size := uint64(0xffffffffe0)
	v, err := memoryGasCost(&Memory{}, size)
	if err != nil {
		fuzz_helper.AddCoverage(17878)
		t.Error("didn't expect error:", err)
	} else {
		fuzz_helper.AddCoverage(45021)
	}
	fuzz_helper.AddCoverage(44810)
	if v != 36028899963961341 {
		fuzz_helper.AddCoverage(39040)
		t.Errorf("Expected: 36028899963961341, got %d", v)
	} else {
		fuzz_helper.AddCoverage(2095)
	}
	fuzz_helper.AddCoverage(5262)

	_, err = memoryGasCost(&Memory{}, size+1)
	if err == nil {
		fuzz_helper.AddCoverage(21668)
		t.Error("expected error")
	} else {
		fuzz_helper.AddCoverage(45213)
	}
}

var _ = fuzz_helper.AddCoverage
