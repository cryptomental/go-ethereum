package vm

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import "fmt"

// Memory implements a simple memory model for the ethereum virtual machine.
type Memory struct {
	store       []byte
	lastGasCost uint64
}

func NewMemory() *Memory {
	fuzz_helper.AddCoverage(22588)
	return &Memory{}
}

// Set sets offset + size to value
func (m *Memory) Set(offset, size uint64, value []byte) {
	fuzz_helper.AddCoverage(44810)

	if size > uint64(len(m.store)) {
		fuzz_helper.AddCoverage(17878)
		panic("INVALID memory: store empty")
	} else {
		fuzz_helper.AddCoverage(45021)
	}
	fuzz_helper.AddCoverage(5262)

	if size > 0 {
		fuzz_helper.AddCoverage(39040)
		copy(m.store[offset:offset+size], value)
	} else {
		fuzz_helper.AddCoverage(2095)
	}
}

// Resize resizes the memory to size
func (m *Memory) Resize(size uint64) {
	fuzz_helper.AddCoverage(21668)
	if uint64(m.Len()) < size {
		fuzz_helper.AddCoverage(45213)
		m.store = append(m.store, make([]byte, size-uint64(m.Len()))...)
	} else {
		fuzz_helper.AddCoverage(16619)
	}
}

// Get returns offset + size as a new slice
func (self *Memory) Get(offset, size int64) (cpy []byte) {
	fuzz_helper.AddCoverage(12692)
	if size == 0 {
		fuzz_helper.AddCoverage(17393)
		return nil
	} else {
		fuzz_helper.AddCoverage(64174)
	}
	fuzz_helper.AddCoverage(42483)

	if len(self.store) > int(offset) {
		fuzz_helper.AddCoverage(38740)
		cpy = make([]byte, size)
		copy(cpy, self.store[offset:offset+size])

		return
	} else {
		fuzz_helper.AddCoverage(35657)
	}
	fuzz_helper.AddCoverage(6577)

	return
}

// GetPtr returns the offset + size
func (self *Memory) GetPtr(offset, size int64) []byte {
	fuzz_helper.AddCoverage(30358)
	if size == 0 {
		fuzz_helper.AddCoverage(11162)
		return nil
	} else {
		fuzz_helper.AddCoverage(49217)
	}
	fuzz_helper.AddCoverage(23294)

	if len(self.store) > int(offset) {
		fuzz_helper.AddCoverage(34511)
		return self.store[offset : offset+size]
	} else {
		fuzz_helper.AddCoverage(64074)
	}
	fuzz_helper.AddCoverage(61639)

	return nil
}

// Len returns the length of the backing slice
func (m *Memory) Len() int {
	fuzz_helper.AddCoverage(28614)
	return len(m.store)
}

// Data returns the backing slice
func (m *Memory) Data() []byte {
	fuzz_helper.AddCoverage(39226)
	return m.store
}

func (m *Memory) Print() {
	fuzz_helper.AddCoverage(2297)
	fmt.Printf("### mem %d bytes ###\n", len(m.store))
	if len(m.store) > 0 {
		fuzz_helper.AddCoverage(52877)
		addr := 0
		for i := 0; i+32 <= len(m.store); i += 32 {
			fuzz_helper.AddCoverage(778)
			fmt.Printf("%03d: % x\n", addr, m.store[i:i+32])
			addr++
		}
	} else {
		fuzz_helper.AddCoverage(33340)
		fmt.Println("-- empty --")
	}
	fuzz_helper.AddCoverage(40870)
	fmt.Println("####################")
}

var _ = fuzz_helper.AddCoverage
