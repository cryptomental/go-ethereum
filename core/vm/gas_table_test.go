package vm

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import "testing"

func TestMemoryGasCost(t *testing.T) {
	fuzz_helper.AddCoverage(13222)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()

	size := uint64(0xffffffffe0)
	v, err := memoryGasCost(&Memory{}, size)
	if err != nil {
		fuzz_helper.AddCoverage(3204)
		t.Error("didn't expect error:", err)
	} else {
		fuzz_helper.AddCoverage(34140)
	}
	fuzz_helper.AddCoverage(17465)
	if v != 36028899963961341 {
		fuzz_helper.AddCoverage(39557)
		t.Errorf("Expected: 36028899963961341, got %d", v)
	} else {
		fuzz_helper.AddCoverage(31052)
	}
	fuzz_helper.AddCoverage(10456)

	_, err = memoryGasCost(&Memory{}, size+1)
	if err == nil {
		fuzz_helper.AddCoverage(6234)
		t.Error("expected error")
	} else {
		fuzz_helper.AddCoverage(18961)
	}
}

var _ = fuzz_helper.AddCoverage
