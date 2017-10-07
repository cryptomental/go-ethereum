package vm

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// ContractRef is a reference to the contract's backing object
type ContractRef interface {
	Address() common.Address
}

// AccountRef implements ContractRef.
//
// Account references are used during EVM initialisation and
// it's primary use is to fetch addresses. Removing this object
// proves difficult because of the cached jump destinations which
// are fetched from the parent contract (i.e. the caller), which
// is a ContractRef.
type AccountRef common.Address

// Address casts AccountRef to a Address
func (ar AccountRef) Address() common.Address {
	fuzz_helper.AddCoverage(42993)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return (common.Address)(ar)
}

// Contract represents an ethereum contract in the state database. It contains
// the the contract code, calling arguments. Contract implements ContractRef
type Contract struct {
	// CallerAddress is the result of the caller which initialised this
	// contract. However when the "call method" is delegated this value
	// needs to be initialised to that of the caller's caller.
	CallerAddress common.Address
	caller        ContractRef
	self          ContractRef

	jumpdests destinations // result of JUMPDEST analysis.

	Code     []byte
	CodeHash common.Hash
	CodeAddr *common.Address
	Input    []byte

	Gas   uint64
	value *big.Int

	Args []byte

	DelegateCall bool
}

// NewContract returns a new contract environment for the execution of EVM.
func NewContract(caller ContractRef, object ContractRef, value *big.Int, gas uint64) *Contract {
	fuzz_helper.AddCoverage(30301)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	c := &Contract{CallerAddress: caller.Address(), caller: caller, self: object, Args: nil}

	if parent, ok := caller.(*Contract); ok {
		fuzz_helper.AddCoverage(264)

		c.jumpdests = parent.jumpdests
	} else {
		fuzz_helper.AddCoverage(3566)
		c.jumpdests = make(destinations)
	}
	fuzz_helper.AddCoverage(45210)

	c.Gas = gas

	c.value = value

	return c
}

// AsDelegate sets the contract to be a delegate call and returns the current
// contract (for chaining calls)
func (c *Contract) AsDelegate() *Contract {
	fuzz_helper.AddCoverage(47636)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	c.DelegateCall = true

	parent := c.caller.(*Contract)
	c.CallerAddress = parent.CallerAddress
	c.value = parent.value

	return c
}

// GetOp returns the n'th element in the contract's byte array
func (c *Contract) GetOp(n uint64) OpCode {
	fuzz_helper.AddCoverage(8730)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return OpCode(c.GetByte(n))
}

// GetByte returns the n'th byte in the contract's byte array
func (c *Contract) GetByte(n uint64) byte {
	fuzz_helper.AddCoverage(20539)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	if n < uint64(len(c.Code)) {
		fuzz_helper.AddCoverage(19009)
		return c.Code[n]
	} else {
		fuzz_helper.AddCoverage(64748)
	}
	fuzz_helper.AddCoverage(63931)

	return 0
}

// Caller returns the caller of the contract.
//
// Caller will recursively call caller when the contract is a delegate
// call, including that of caller's caller.
func (c *Contract) Caller() common.Address {
	fuzz_helper.AddCoverage(50446)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return c.CallerAddress
}

// UseGas attempts the use gas and subtracts it and returns true on success
func (c *Contract) UseGas(gas uint64) (ok bool) {
	fuzz_helper.AddCoverage(18500)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	if c.Gas < gas {
		fuzz_helper.AddCoverage(17111)
		return false
	} else {
		fuzz_helper.AddCoverage(9670)
	}
	fuzz_helper.AddCoverage(52152)
	c.Gas -= gas
	return true
}

// Address returns the contracts address
func (c *Contract) Address() common.Address {
	fuzz_helper.AddCoverage(55848)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return c.self.Address()
}

// Value returns the contracts value (sent to it from it's caller)
func (c *Contract) Value() *big.Int {
	fuzz_helper.AddCoverage(50755)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return c.value
}

// SetCode sets the code to the contract
func (self *Contract) SetCode(hash common.Hash, code []byte) {
	fuzz_helper.AddCoverage(912)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	self.Code = code
	self.CodeHash = hash
}

// SetCallCode sets the code of the contract and address of the backing data
// object
func (self *Contract) SetCallCode(addr *common.Address, hash common.Hash, code []byte) {
	fuzz_helper.AddCoverage(64631)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	self.Code = code
	self.CodeHash = hash
	self.CodeAddr = addr
}

var _ = fuzz_helper.AddCoverage
