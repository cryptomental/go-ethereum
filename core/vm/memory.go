package vm

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import "fmt"

// Memory implements a simple memory model for the ethereum virtual machine.
type Memory struct {
	store       []byte
	lastGasCost uint64
}

func NewMemory() *Memory {
	fuzz_helper.CoverTab[22588]++
	return &Memory{}
}

// Set sets offset + size to value
func (m *Memory) Set(offset, size uint64, value []byte) {
	fuzz_helper.CoverTab[44810]++

	if size > uint64(len(m.store)) {
		fuzz_helper.CoverTab[17878]++
		panic("INVALID memory: store empty")
	} else {
		fuzz_helper.CoverTab[45021]++
	}
	fuzz_helper.CoverTab[5262]++

	if size > 0 {
		fuzz_helper.CoverTab[39040]++
		copy(m.store[offset:offset+size], value)
	} else {
		fuzz_helper.CoverTab[2095]++
	}
}

// Resize resizes the memory to size
func (m *Memory) Resize(size uint64) {
	fuzz_helper.CoverTab[21668]++
	if uint64(m.Len()) < size {
		fuzz_helper.CoverTab[45213]++
		m.store = append(m.store, make([]byte, size-uint64(m.Len()))...)
	} else {
		fuzz_helper.CoverTab[16619]++
	}
}

// Get returns offset + size as a new slice
func (self *Memory) Get(offset, size int64) (cpy []byte) {
	fuzz_helper.CoverTab[12692]++
	if size == 0 {
		fuzz_helper.CoverTab[17393]++
		return nil
	} else {
		fuzz_helper.CoverTab[64174]++
	}
	fuzz_helper.CoverTab[42483]++

	if len(self.store) > int(offset) {
		fuzz_helper.CoverTab[38740]++
		cpy = make([]byte, size)
		copy(cpy, self.store[offset:offset+size])

		return
	} else {
		fuzz_helper.CoverTab[35657]++
	}
	fuzz_helper.CoverTab[6577]++

	return
}

// GetPtr returns the offset + size
func (self *Memory) GetPtr(offset, size int64) []byte {
	fuzz_helper.CoverTab[30358]++
	if size == 0 {
		fuzz_helper.CoverTab[11162]++
		return nil
	} else {
		fuzz_helper.CoverTab[49217]++
	}
	fuzz_helper.CoverTab[23294]++

	if len(self.store) > int(offset) {
		fuzz_helper.CoverTab[34511]++
		return self.store[offset : offset+size]
	} else {
		fuzz_helper.CoverTab[64074]++
	}
	fuzz_helper.CoverTab[61639]++

	return nil
}

// Len returns the length of the backing slice
func (m *Memory) Len() int {
	fuzz_helper.CoverTab[28614]++
	return len(m.store)
}

// Data returns the backing slice
func (m *Memory) Data() []byte {
	fuzz_helper.CoverTab[39226]++
	return m.store
}

func (m *Memory) Print() {
	fuzz_helper.CoverTab[2297]++
	fmt.Printf("### mem %d bytes ###\n", len(m.store))
	if len(m.store) > 0 {
		fuzz_helper.CoverTab[52877]++
		addr := 0
		for i := 0; i+32 <= len(m.store); i += 32 {
			fuzz_helper.CoverTab[778]++
			fmt.Printf("%03d: % x\n", addr, m.store[i:i+32])
			addr++
		}
	} else {
		fuzz_helper.CoverTab[33340]++
		fmt.Println("-- empty --")
	}
	fuzz_helper.CoverTab[40870]++
	fmt.Println("####################")
}

var _ = fuzz_helper.CoverTab
