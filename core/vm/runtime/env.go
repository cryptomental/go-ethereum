package runtime

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
)

func NewEnv(cfg *Config) *vm.EVM {
	fuzz_helper.AddCoverage(11986)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	context := vm.Context{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		GetHash:     func(uint64) common.Hash { fuzz_helper.AddCoverage(45867); return common.Hash{} },

		Origin:      cfg.Origin,
		Coinbase:    cfg.Coinbase,
		BlockNumber: cfg.BlockNumber,
		Time:        cfg.Time,
		Difficulty:  cfg.Difficulty,
		GasLimit:    new(big.Int).SetUint64(cfg.GasLimit),
		GasPrice:    cfg.GasPrice,
	}
	fuzz_helper.AddCoverage(11299)

	return vm.NewEVM(context, cfg.State, cfg.ChainConfig, cfg.EVMConfig)
}

var _ = fuzz_helper.AddCoverage
