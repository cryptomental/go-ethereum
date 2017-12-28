package main

import "C"
import "bytes"
import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"
import ecies "github.com/ethereum/go-ethereum/crypto/ecies_instrumented"
import "github.com/ethereum/go-ethereum/crypto"
import "crypto/rand"

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

func testEncryptDecrypt(input []byte) {
    var n int
    if len(input) >= 32 {
        n = 32
    } else if len(input) == 0 {
        return
    } else {
        n = len(input) / 2
    }
    prv1, err := crypto.ToECDSA(input[:n])
    if err != nil {
        return
    }
    message := input[n:]

    ecies_prv1 := ecies.ImportECDSA(prv1);
    ecies_pub1 := ecies.ImportECDSAPublic(&prv1.PublicKey)

    ciphertext, err := ecies.Encrypt(rand.Reader, ecies_pub1, message, nil, nil)
    if err != nil {
        return
    }

    plaintext, err := ecies_prv1.Decrypt(rand.Reader, ciphertext, nil, nil)
    if err != nil {
        return
    }

    if !bytes.Equal(plaintext, message) {
        panic("Plaintext not equal to original input")
    }
}

func testEncryptDecryptGenerated(input []byte) {
	prv1, err := ecies.GenerateKey(rand.Reader, ecies.DefaultCurve, nil)
	if err != nil {
        return
	}

	prv2, err := ecies.GenerateKey(rand.Reader, ecies.DefaultCurve, nil)
	if err != nil {
        return
	}

	ct, err := ecies.Encrypt(rand.Reader, &prv2.PublicKey, input, nil, nil)
	if err != nil {
        return
	}

	pt, err := prv2.Decrypt(rand.Reader, ct, nil, nil)
	if err != nil {
        return
	}

	if !bytes.Equal(pt, input) {
        panic("Plaintext not equal to original input")
	}

	_, err = prv1.Decrypt(rand.Reader, ct, nil, nil)
	if err == nil {
		panic("encryption should not have succeeded")
	}
}

func testShared(input []byte) {
    if len(input) < 64 {
        return
    }

    prv1, err := crypto.ToECDSA(input[:32])
    if err != nil {
        return
    }
    prv2, err := crypto.ToECDSA(input[32:64])
    if err != nil {
        return
    }

    ecies_prv1 := ecies.ImportECDSA(prv1);
    //ecies_pub1 := ecies.ImportECDSAPublic(&prv1.PublicKey)
    ecies_prv2 := ecies.ImportECDSA(prv2);
    //ecies_pub2 := ecies.ImportECDSAPublic(&prv2.PublicKey)

	skLen := ecies.MaxSharedKeyLength(&ecies_prv1.PublicKey) / 2

    sk1, err := ecies_prv1.GenerateShared(&ecies_prv2.PublicKey, skLen, skLen)
    if err != nil {
        return
    }

    sk2, err := ecies_prv2.GenerateShared(&ecies_prv1.PublicKey, skLen, skLen)
    if err != nil {
        return
    }

	if !bytes.Equal(sk1, sk2) {
        panic("Shared keys not equal")
	}
}

func testTooBigShared(input []byte) {
    if len(input) < 64 {
        return
    }

    prv1, err := crypto.ToECDSA(input[:32])
    if err != nil {
        return
    }
    prv2, err := crypto.ToECDSA(input[32:64])
    if err != nil {
        return
    }

    ecies_prv1 := ecies.ImportECDSA(prv1);
    ecies_prv2 := ecies.ImportECDSA(prv2);

	skLen := 32

    _, err = ecies_prv1.GenerateShared(&ecies_prv2.PublicKey, skLen, skLen)
    if err != ecies.ErrSharedKeyTooBig {
        panic("GenerateShared() did not return ErrSharedKeyTooBig")
    }

    _, err = ecies_prv2.GenerateShared(&ecies_prv1.PublicKey, skLen, skLen)
    if err != ecies.ErrSharedKeyTooBig {
        panic("GenerateShared() did not return ErrSharedKeyTooBig")
    }
}

func testDecrypt(input []byte) {
    if len(input) < (32+1+1) {
        return
    }

    prv1, err := crypto.ToECDSA(input[:32])
    if err != nil {
        return
    }
    ecies_prv1 := ecies.ImportECDSA(prv1);

    input = input[32:]


    var withs1 bool
    var withs2 bool
    if input[0] % 2 == 0 {
        withs2 = true
    } else {
        withs2 = false
    }
    if input[1] % 2 == 0 {
        withs2 = true
    } else {
        withs2 = false
    }

    input = input[2:]
    ciphertext := input[:len(input)/3]
    var s1 []byte = nil
    var s2 []byte = nil
    if withs1 {
        s1 = input[:len(input)/3]
    }
    if withs2 {
        s2 = input[:len(input)/3]
    }

    _, err = ecies_prv1.Decrypt(rand.Reader, ciphertext, s1, s2)
}

//export run_ecies
func run_ecies(input []byte, mode int) {
    testEncryptDecrypt(input)
    testEncryptDecryptGenerated(input)
    testShared(input)
    testTooBigShared(input)
    testDecrypt(input)
}

/* No main() body because this file is compiled to a static archive */
func main() {
}
