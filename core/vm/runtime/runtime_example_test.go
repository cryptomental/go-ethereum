package runtime_test

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm/runtime"
)

func ExampleExecute() {
	fuzz_helper.AddCoverage(15655)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	ret, _, err := runtime.Execute(common.Hex2Bytes("6060604052600a8060106000396000f360606040526008565b00"), nil, nil)
	if err != nil {
		fuzz_helper.AddCoverage(64465)
		fmt.Println(err)
	} else {
		fuzz_helper.AddCoverage(28899)
	}
	fuzz_helper.AddCoverage(24272)
	fmt.Println(ret)

}

var _ = fuzz_helper.AddCoverage
