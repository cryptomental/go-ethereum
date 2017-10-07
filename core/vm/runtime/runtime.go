package runtime

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
)

// Config is a basic type specifying certain configuration flags for running
// the EVM.
type Config struct {
	ChainConfig *params.ChainConfig
	Difficulty  *big.Int
	Origin      common.Address
	Coinbase    common.Address
	BlockNumber *big.Int
	Time        *big.Int
	GasLimit    uint64
	GasPrice    *big.Int
	Value       *big.Int
	DisableJit  bool // "disable" so it's enabled by default
	Debug       bool
	EVMConfig   vm.Config

	State     *state.StateDB
	GetHashFn func(n uint64) common.Hash
}

// sets defaults on the config
func setDefaults(cfg *Config) {
	fuzz_helper.AddCoverage(18098)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	if cfg.ChainConfig == nil {
		fuzz_helper.AddCoverage(29662)
		cfg.ChainConfig = &params.ChainConfig{
			ChainId:        big.NewInt(1),
			HomesteadBlock: new(big.Int),
			DAOForkBlock:   new(big.Int),
			DAOForkSupport: false,
			EIP150Block:    new(big.Int),
			EIP155Block:    new(big.Int),
			EIP158Block:    new(big.Int),
		}
	} else {
		fuzz_helper.AddCoverage(37203)
	}
	fuzz_helper.AddCoverage(23681)

	if cfg.Difficulty == nil {
		fuzz_helper.AddCoverage(60325)
		cfg.Difficulty = new(big.Int)
	} else {
		fuzz_helper.AddCoverage(4473)
	}
	fuzz_helper.AddCoverage(10981)
	if cfg.Time == nil {
		fuzz_helper.AddCoverage(18212)
		cfg.Time = big.NewInt(time.Now().Unix())
	} else {
		fuzz_helper.AddCoverage(46492)
	}
	fuzz_helper.AddCoverage(8103)
	if cfg.GasLimit == 0 {
		fuzz_helper.AddCoverage(63372)
		cfg.GasLimit = math.MaxUint64
	} else {
		fuzz_helper.AddCoverage(18996)
	}
	fuzz_helper.AddCoverage(14995)
	if cfg.GasPrice == nil {
		fuzz_helper.AddCoverage(56685)
		cfg.GasPrice = new(big.Int)
	} else {
		fuzz_helper.AddCoverage(7152)
	}
	fuzz_helper.AddCoverage(34345)
	if cfg.Value == nil {
		fuzz_helper.AddCoverage(30914)
		cfg.Value = new(big.Int)
	} else {
		fuzz_helper.AddCoverage(61203)
	}
	fuzz_helper.AddCoverage(60137)
	if cfg.BlockNumber == nil {
		fuzz_helper.AddCoverage(45477)
		cfg.BlockNumber = new(big.Int)
	} else {
		fuzz_helper.AddCoverage(61123)
	}
	fuzz_helper.AddCoverage(907)
	if cfg.GetHashFn == nil {
		fuzz_helper.AddCoverage(3505)
		cfg.GetHashFn = func(n uint64) common.Hash {
			fuzz_helper.AddCoverage(27510)
			return common.BytesToHash(crypto.Keccak256([]byte(new(big.Int).SetUint64(n).String())))
		}
	} else {
		fuzz_helper.AddCoverage(17408)
	}
}

// Execute executes the code using the input as call data during the execution.
// It returns the EVM's return value, the new state and an error if it failed.
//
// Executes sets up a in memory, temporarily, environment for the execution of
// the given code. It enabled the JIT by default and make sure that it's restored
// to it's original state afterwards.
func Execute(code, input []byte, cfg *Config) ([]byte, *state.StateDB, error) {
	fuzz_helper.AddCoverage(29199)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	if cfg == nil {
		fuzz_helper.AddCoverage(16497)
		cfg = new(Config)
	} else {
		fuzz_helper.AddCoverage(33471)
	}
	fuzz_helper.AddCoverage(5881)
	setDefaults(cfg)

	if cfg.State == nil {
		fuzz_helper.AddCoverage(43595)
		db, _ := ethdb.NewMemDatabase()
		cfg.State, _ = state.New(common.Hash{}, state.NewDatabase(db))
	} else {
		fuzz_helper.AddCoverage(34531)
	}
	fuzz_helper.AddCoverage(4435)
	var (
		address = common.StringToAddress("contract")
		vmenv   = NewEnv(cfg)
		sender  = vm.AccountRef(cfg.Origin)
	)
	cfg.State.CreateAccount(address)

	cfg.State.SetCode(address, code)

	ret, _, err := vmenv.Call(
		sender,
		common.StringToAddress("contract"),
		input,
		cfg.GasLimit,
		cfg.Value,
	)

	return ret, cfg.State, err
}

// Create executes the code using the EVM create method
func Create(input []byte, cfg *Config) ([]byte, common.Address, uint64, error) {
	fuzz_helper.AddCoverage(8882)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	if cfg == nil {
		fuzz_helper.AddCoverage(16595)
		cfg = new(Config)
	} else {
		fuzz_helper.AddCoverage(56682)
	}
	fuzz_helper.AddCoverage(64691)
	setDefaults(cfg)

	if cfg.State == nil {
		fuzz_helper.AddCoverage(9660)
		db, _ := ethdb.NewMemDatabase()
		cfg.State, _ = state.New(common.Hash{}, state.NewDatabase(db))
	} else {
		fuzz_helper.AddCoverage(56029)
	}
	fuzz_helper.AddCoverage(51357)
	var (
		vmenv  = NewEnv(cfg)
		sender = vm.AccountRef(cfg.Origin)
	)

	code, address, leftOverGas, err := vmenv.Create(
		sender,
		input,
		cfg.GasLimit,
		cfg.Value,
	)
	return code, address, leftOverGas, err
}

// Call executes the code given by the contract's address. It will return the
// EVM's return value or an error if it failed.
//
// Call, unlike Execute, requires a config and also requires the State field to
// be set.
func Call(address common.Address, input []byte, cfg *Config) ([]byte, uint64, error) {
	fuzz_helper.AddCoverage(42980)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	setDefaults(cfg)

	vmenv := NewEnv(cfg)

	sender := cfg.State.GetOrNewStateObject(cfg.Origin)

	ret, leftOverGas, err := vmenv.Call(
		sender,
		address,
		input,
		cfg.GasLimit,
		cfg.Value,
	)

	return ret, leftOverGas, err
}

var _ = fuzz_helper.AddCoverage
