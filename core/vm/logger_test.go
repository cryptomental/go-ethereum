package vm

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

type dummyContractRef struct {
	calledForEach bool
}

func (dummyContractRef) ReturnGas(*big.Int) { fuzz_helper.AddCoverage(22588) }
func (dummyContractRef) Address() common.Address {
	fuzz_helper.AddCoverage(44810)
	return common.Address{}
}
func (dummyContractRef) Value() *big.Int             { fuzz_helper.AddCoverage(5262); return new(big.Int) }
func (dummyContractRef) SetCode(common.Hash, []byte) { fuzz_helper.AddCoverage(17878) }
func (d *dummyContractRef) ForEachStorage(callback func(key, value common.Hash) bool) {
	fuzz_helper.AddCoverage(45021)
	d.calledForEach = true
}
func (d *dummyContractRef) SubBalance(amount *big.Int) { fuzz_helper.AddCoverage(39040) }
func (d *dummyContractRef) AddBalance(amount *big.Int) { fuzz_helper.AddCoverage(2095) }
func (d *dummyContractRef) SetBalance(*big.Int)        { fuzz_helper.AddCoverage(21668) }
func (d *dummyContractRef) SetNonce(uint64)            { fuzz_helper.AddCoverage(45213) }
func (d *dummyContractRef) Balance() *big.Int          { fuzz_helper.AddCoverage(16619); return new(big.Int) }

type dummyStateDB struct {
	NoopStateDB
	ref *dummyContractRef
}

func TestStoreCapture(t *testing.T) {
	fuzz_helper.AddCoverage(12692)
	var (
		env      = NewEVM(Context{}, nil, params.TestChainConfig, Config{EnableJit: false, ForceJit: false})
		logger   = NewStructLogger(nil)
		mem      = NewMemory()
		stack    = newstack()
		contract = NewContract(&dummyContractRef{}, &dummyContractRef{}, new(big.Int), 0)
	)
	stack.push(big.NewInt(1))
	stack.push(big.NewInt(0))

	var index common.Hash

	logger.CaptureState(env, 0, SSTORE, 0, 0, mem, stack, contract, 0, nil)
	if len(logger.changedValues[contract.Address()]) == 0 {
		fuzz_helper.AddCoverage(6577)
		t.Fatalf("expected exactly 1 changed value on address %x, got %d", contract.Address(), len(logger.changedValues[contract.Address()]))
	} else {
		fuzz_helper.AddCoverage(17393)
	}
	fuzz_helper.AddCoverage(42483)

	exp := common.BigToHash(big.NewInt(1))
	if logger.changedValues[contract.Address()][index] != exp {
		fuzz_helper.AddCoverage(64174)
		t.Errorf("expected %x, got %x", exp, logger.changedValues[contract.Address()][index])
	} else {
		fuzz_helper.AddCoverage(38740)
	}
}

func TestStorageCapture(t *testing.T) {
	fuzz_helper.AddCoverage(35657)
	t.Skip("implementing this function is difficult. it requires all sort of interfaces to be implemented which isn't trivial. The value (the actual test) isn't worth it")
	var (
		ref      = &dummyContractRef{}
		contract = NewContract(ref, ref, new(big.Int), 0)
		env      = NewEVM(Context{}, dummyStateDB{ref: ref}, params.TestChainConfig, Config{EnableJit: false, ForceJit: false})
		logger   = NewStructLogger(nil)
		mem      = NewMemory()
		stack    = newstack()
	)

	logger.CaptureState(env, 0, STOP, 0, 0, mem, stack, contract, 0, nil)
	if ref.calledForEach {
		fuzz_helper.AddCoverage(23294)
		t.Error("didn't expect for each to be called")
	} else {
		fuzz_helper.AddCoverage(61639)
	}
	fuzz_helper.AddCoverage(30358)

	logger = NewStructLogger(&LogConfig{FullStorage: true})
	logger.CaptureState(env, 0, STOP, 0, 0, mem, stack, contract, 0, nil)
	if !ref.calledForEach {
		fuzz_helper.AddCoverage(11162)
		t.Error("expected for each to be called")
	} else {
		fuzz_helper.AddCoverage(49217)
	}
}

var _ = fuzz_helper.AddCoverage
