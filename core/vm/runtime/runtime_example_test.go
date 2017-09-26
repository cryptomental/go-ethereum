package runtime_test

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm/runtime"
)

func ExampleExecute() {
	fuzz_helper.AddCoverage(22588)
	ret, _, err := runtime.Execute(common.Hex2Bytes("6060604052600a8060106000396000f360606040526008565b00"), nil, nil)
	if err != nil {
		fuzz_helper.AddCoverage(5262)
		fmt.Println(err)
	} else {
		fuzz_helper.AddCoverage(17878)
	}
	fuzz_helper.AddCoverage(44810)
	fmt.Println(ret)

}

var _ = fuzz_helper.AddCoverage
