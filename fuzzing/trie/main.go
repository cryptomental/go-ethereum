package main

import "C"
//import "fmt"
import "bytes"
import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"
import "github.com/ethereum/go-ethereum/ethdb"
import trie "github.com/ethereum/go-ethereum/trie_instrumented"
import "github.com/ethereum/go-ethereum/common"

//export SetInstrumentationType
func SetInstrumentationType(t int) {
    fuzz_helper.SetInstrumentationType(t)
}
//export GoResetCoverage
func GoResetCoverage() {
    fuzz_helper.ResetCoverage()
}

//export GoCalcCoverage
func GoCalcCoverage() uint64 {
    return fuzz_helper.CalcCoverage()
}

/*
func run_trie(input []byte, mode int) {
	mode = mode % 5
	if mode == 0 {
		trie.HexToCompact(input)
	} else if mode == 1 {
		trie.CompactToHex(input)
	} else if mode == 2 {
		trie.KeybytesToHex(input)
	} else if mode == 3 {
        if len(input)&1 == 0 {
            trie.HexToKeybytes(input)
        }
	} else if mode == 4 {
		a := input[:len(input)]
		b := input[len(input):]
		trie.PrefixLen(a, b)
	}
}
*/

func updateString(trie *trie.Trie, k, v string) {
	trie.Update([]byte(k), []byte(v))
}
//export run_trie
func run_trie(input []byte, mode int) {
	db, _ := ethdb.NewMemDatabase()
	xtrie, _ := trie.New(common.Hash{}, db)
    separator := make([]byte, 1)
    separator[0] = byte(',')
    parts := bytes.Split(input, separator)
    lenparts := len(parts)
    var sub1lenparts int
    var sub2lenparts int
    if lenparts > 10 {
        sub1lenparts = 10
        sub2lenparts = lenparts - 10
    } else {
        sub1lenparts = lenparts
        sub2lenparts = 0
    }
    if sub1lenparts & 1 == 1 {
        sub1lenparts--
    }
    if sub2lenparts & 1 == 1 {
        sub2lenparts--
    }
    for i := 0; i < sub1lenparts; i = i + 2 {
        xtrie.Update(parts[i+0], parts[i+1])
    }
    root, _ := xtrie.Commit()
	xtrie, _ = trie.New(root, db)
    for i := 0; i < sub2lenparts; i = i + 2 {
        xtrie.Get(parts[sub1lenparts+i+0])
        xtrie.Delete(parts[sub1lenparts+i+1])
    }
    checktr, _ := trie.New(common.Hash{}, nil)
    it := trie.NewIterator(xtrie.NodeIterator(nil))
    for it.Next() {
        checktr.Update(it.Key, it.Value)
    }
    if xtrie.Hash() != checktr.Hash() {
        panic("hash mismatch")
    }
    xtrie.Hash()
    xtrie.Commit()
}
/* No main() body because this file is compiled to a static archive */
func main() {
}
