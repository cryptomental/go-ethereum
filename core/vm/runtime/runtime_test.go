package runtime

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
)

func TestDefaults(t *testing.T) {
	fuzz_helper.AddCoverage(12711)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	cfg := new(Config)
	setDefaults(cfg)

	if cfg.Difficulty == nil {
		fuzz_helper.AddCoverage(14830)
		t.Error("expected difficulty to be non nil")
	} else {
		fuzz_helper.AddCoverage(14455)
	}
	fuzz_helper.AddCoverage(49630)

	if cfg.Time == nil {
		fuzz_helper.AddCoverage(25572)
		t.Error("expected time to be non nil")
	} else {
		fuzz_helper.AddCoverage(54637)
	}
	fuzz_helper.AddCoverage(57933)
	if cfg.GasLimit == 0 {
		fuzz_helper.AddCoverage(37407)
		t.Error("didn't expect gaslimit to be zero")
	} else {
		fuzz_helper.AddCoverage(13027)
	}
	fuzz_helper.AddCoverage(51460)
	if cfg.GasPrice == nil {
		fuzz_helper.AddCoverage(57235)
		t.Error("expected time to be non nil")
	} else {
		fuzz_helper.AddCoverage(47801)
	}
	fuzz_helper.AddCoverage(1152)
	if cfg.Value == nil {
		fuzz_helper.AddCoverage(3138)
		t.Error("expected time to be non nil")
	} else {
		fuzz_helper.AddCoverage(52298)
	}
	fuzz_helper.AddCoverage(60988)
	if cfg.GetHashFn == nil {
		fuzz_helper.AddCoverage(34497)
		t.Error("expected time to be non nil")
	} else {
		fuzz_helper.AddCoverage(11490)
	}
	fuzz_helper.AddCoverage(18447)
	if cfg.BlockNumber == nil {
		fuzz_helper.AddCoverage(45028)
		t.Error("expected block number to be non nil")
	} else {
		fuzz_helper.AddCoverage(35800)
	}
}

func TestEVM(t *testing.T) {
	fuzz_helper.AddCoverage(7078)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	defer func() {
		fuzz_helper.AddCoverage(34511)
		if r := recover(); r != nil {
			fuzz_helper.AddCoverage(57903)
			t.Fatalf("crashed with: %v", r)
		} else {
			fuzz_helper.AddCoverage(23350)
		}
	}()
	fuzz_helper.AddCoverage(4451)

	Execute([]byte{
		byte(vm.DIFFICULTY),
		byte(vm.TIMESTAMP),
		byte(vm.GASLIMIT),
		byte(vm.PUSH1),
		byte(vm.ORIGIN),
		byte(vm.BLOCKHASH),
		byte(vm.COINBASE),
	}, nil, nil)
}

func TestExecute(t *testing.T) {
	fuzz_helper.AddCoverage(28480)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	ret, _, err := Execute([]byte{
		byte(vm.PUSH1), 10,
		byte(vm.PUSH1), 0,
		byte(vm.MSTORE),
		byte(vm.PUSH1), 32,
		byte(vm.PUSH1), 0,
		byte(vm.RETURN),
	}, nil, nil)
	if err != nil {
		fuzz_helper.AddCoverage(20295)
		t.Fatal("didn't expect error", err)
	} else {
		fuzz_helper.AddCoverage(40052)
	}
	fuzz_helper.AddCoverage(7377)

	num := new(big.Int).SetBytes(ret)
	if num.Cmp(big.NewInt(10)) != 0 {
		fuzz_helper.AddCoverage(27533)
		t.Error("Expected 10, got", num)
	} else {
		fuzz_helper.AddCoverage(18600)
	}
}

func TestCall(t *testing.T) {
	fuzz_helper.AddCoverage(11893)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	db, _ := ethdb.NewMemDatabase()
	state, _ := state.New(common.Hash{}, state.NewDatabase(db))
	address := common.HexToAddress("0x0a")
	state.SetCode(address, []byte{
		byte(vm.PUSH1), 10,
		byte(vm.PUSH1), 0,
		byte(vm.MSTORE),
		byte(vm.PUSH1), 32,
		byte(vm.PUSH1), 0,
		byte(vm.RETURN),
	})

	ret, _, err := Call(address, nil, &Config{State: state})
	if err != nil {
		fuzz_helper.AddCoverage(25369)
		t.Fatal("didn't expect error", err)
	} else {
		fuzz_helper.AddCoverage(59744)
	}
	fuzz_helper.AddCoverage(52991)

	num := new(big.Int).SetBytes(ret)
	if num.Cmp(big.NewInt(10)) != 0 {
		fuzz_helper.AddCoverage(31823)
		t.Error("Expected 10, got", num)
	} else {
		fuzz_helper.AddCoverage(37955)
	}
}

func BenchmarkCall(b *testing.B) {
	fuzz_helper.AddCoverage(63557)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	var definition = `[{"constant":true,"inputs":[],"name":"seller","outputs":[{"name":"","type":"address"}],"type":"function"},{"constant":false,"inputs":[],"name":"abort","outputs":[],"type":"function"},{"constant":true,"inputs":[],"name":"value","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":false,"inputs":[],"name":"refund","outputs":[],"type":"function"},{"constant":true,"inputs":[],"name":"buyer","outputs":[{"name":"","type":"address"}],"type":"function"},{"constant":false,"inputs":[],"name":"confirmReceived","outputs":[],"type":"function"},{"constant":true,"inputs":[],"name":"state","outputs":[{"name":"","type":"uint8"}],"type":"function"},{"constant":false,"inputs":[],"name":"confirmPurchase","outputs":[],"type":"function"},{"inputs":[],"type":"constructor"},{"anonymous":false,"inputs":[],"name":"Aborted","type":"event"},{"anonymous":false,"inputs":[],"name":"PurchaseConfirmed","type":"event"},{"anonymous":false,"inputs":[],"name":"ItemReceived","type":"event"},{"anonymous":false,"inputs":[],"name":"Refunded","type":"event"}]`

	var code = common.Hex2Bytes("6060604052361561006c5760e060020a600035046308551a53811461007457806335a063b4146100865780633fa4f245146100a6578063590e1ae3146100af5780637150d8ae146100cf57806373fac6f0146100e1578063c19d93fb146100fe578063d696069714610112575b610131610002565b610133600154600160a060020a031681565b610131600154600160a060020a0390811633919091161461015057610002565b61014660005481565b610131600154600160a060020a039081163391909116146102d557610002565b610133600254600160a060020a031681565b610131600254600160a060020a0333811691161461023757610002565b61014660025460ff60a060020a9091041681565b61013160025460009060ff60a060020a9091041681146101cc57610002565b005b600160a060020a03166060908152602090f35b6060908152602090f35b60025460009060a060020a900460ff16811461016b57610002565b600154600160a060020a03908116908290301631606082818181858883f150506002805460a060020a60ff02191660a160020a179055506040517f72c874aeff0b183a56e2b79c71b46e1aed4dee5e09862134b8821ba2fddbf8bf9250a150565b80546002023414806101dd57610002565b6002805460a060020a60ff021973ffffffffffffffffffffffffffffffffffffffff1990911633171660a060020a1790557fd5d55c8a68912e9a110618df8d5e2e83b8d83211c57a8ddd1203df92885dc881826060a15050565b60025460019060a060020a900460ff16811461025257610002565b60025460008054600160a060020a0390921691606082818181858883f150508354604051600160a060020a0391821694503090911631915082818181858883f150506002805460a060020a60ff02191660a160020a179055506040517fe89152acd703c9d8c7d28829d443260b411454d45394e7995815140c8cbcbcf79250a150565b60025460019060a060020a900460ff1681146102f057610002565b6002805460008054600160a060020a0390921692909102606082818181858883f150508354604051600160a060020a0391821694503090911631915082818181858883f150506002805460a060020a60ff02191660a160020a179055506040517f8616bbbbad963e4e65b1366f1d75dfb63f9e9704bbbf91fb01bec70849906cf79250a15056")

	abi, err := abi.JSON(strings.NewReader(definition))
	if err != nil {
		fuzz_helper.AddCoverage(49540)
		b.Fatal(err)
	} else {
		fuzz_helper.AddCoverage(54865)
	}
	fuzz_helper.AddCoverage(54431)

	cpurchase, err := abi.Pack("confirmPurchase")
	if err != nil {
		fuzz_helper.AddCoverage(39156)
		b.Fatal(err)
	} else {
		fuzz_helper.AddCoverage(46278)
	}
	fuzz_helper.AddCoverage(50612)
	creceived, err := abi.Pack("confirmReceived")
	if err != nil {
		fuzz_helper.AddCoverage(37366)
		b.Fatal(err)
	} else {
		fuzz_helper.AddCoverage(11612)
	}
	fuzz_helper.AddCoverage(51992)
	refund, err := abi.Pack("refund")
	if err != nil {
		fuzz_helper.AddCoverage(22109)
		b.Fatal(err)
	} else {
		fuzz_helper.AddCoverage(9639)
	}
	fuzz_helper.AddCoverage(26141)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fuzz_helper.AddCoverage(56132)
		for j := 0; j < 400; j++ {
			fuzz_helper.AddCoverage(36761)
			Execute(code, cpurchase, nil)
			Execute(code, creceived, nil)
			Execute(code, refund, nil)
		}
	}
}

var _ = fuzz_helper.AddCoverage
