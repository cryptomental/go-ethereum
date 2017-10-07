package vm

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import "fmt"

// Memory implements a simple memory model for the ethereum virtual machine.
type Memory struct {
	store       []byte
	lastGasCost uint64
}

func NewMemory() *Memory {
	fuzz_helper.AddCoverage(28742)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return &Memory{}
}

// Set sets offset + size to value
func (m *Memory) Set(offset, size uint64, value []byte) {
	fuzz_helper.AddCoverage(12959)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()

	if size > uint64(len(m.store)) {
		fuzz_helper.AddCoverage(31678)
		panic("INVALID memory: store empty")
	} else {
		fuzz_helper.AddCoverage(47067)
	}
	fuzz_helper.AddCoverage(31799)

	if size > 0 {
		fuzz_helper.AddCoverage(46092)
		copy(m.store[offset:offset+size], value)
	} else {
		fuzz_helper.AddCoverage(63014)
	}
}

// Resize resizes the memory to size
func (m *Memory) Resize(size uint64) {
	fuzz_helper.AddCoverage(57608)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	if uint64(m.Len()) < size {
		fuzz_helper.AddCoverage(23568)
		m.store = append(m.store, make([]byte, size-uint64(m.Len()))...)
	} else {
		fuzz_helper.AddCoverage(30723)
	}
}

// Get returns offset + size as a new slice
func (self *Memory) Get(offset, size int64) (cpy []byte) {
	fuzz_helper.AddCoverage(20813)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	if size == 0 {
		fuzz_helper.AddCoverage(52839)
		return nil
	} else {
		fuzz_helper.AddCoverage(13689)
	}
	fuzz_helper.AddCoverage(6688)

	if len(self.store) > int(offset) {
		fuzz_helper.AddCoverage(35924)
		cpy = make([]byte, size)
		copy(cpy, self.store[offset:offset+size])

		return
	} else {
		fuzz_helper.AddCoverage(63364)
	}
	fuzz_helper.AddCoverage(10587)

	return
}

// GetPtr returns the offset + size
func (self *Memory) GetPtr(offset, size int64) []byte {
	fuzz_helper.AddCoverage(15305)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	if size == 0 {
		fuzz_helper.AddCoverage(19801)
		return nil
	} else {
		fuzz_helper.AddCoverage(15849)
	}
	fuzz_helper.AddCoverage(53653)

	if len(self.store) > int(offset) {
		fuzz_helper.AddCoverage(2277)
		return self.store[offset : offset+size]
	} else {
		fuzz_helper.AddCoverage(22274)
	}
	fuzz_helper.AddCoverage(16995)

	return nil
}

// Len returns the length of the backing slice
func (m *Memory) Len() int {
	fuzz_helper.AddCoverage(16478)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return len(m.store)
}

// Data returns the backing slice
func (m *Memory) Data() []byte {
	fuzz_helper.AddCoverage(20517)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return m.store
}

func (m *Memory) Print() {
	fuzz_helper.AddCoverage(24160)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	fmt.Printf("### mem %d bytes ###\n", len(m.store))
	if len(m.store) > 0 {
		fuzz_helper.AddCoverage(44402)
		addr := 0
		for i := 0; i+32 <= len(m.store); i += 32 {
			fuzz_helper.AddCoverage(38896)
			fmt.Printf("%03d: % x\n", addr, m.store[i:i+32])
			addr++
		}
	} else {
		fuzz_helper.AddCoverage(19578)
		fmt.Println("-- empty --")
	}
	fuzz_helper.AddCoverage(42380)
	fmt.Println("####################")
}

var _ = fuzz_helper.AddCoverage
