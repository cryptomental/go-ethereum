package vm

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"crypto/sha256"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/bn256"
	"github.com/ethereum/go-ethereum/params"
	"golang.org/x/crypto/ripemd160"
)

// PrecompiledContract is the basic interface for native Go contracts. The implementation
// requires a deterministic gas count based on the input size of the Run method of the
// contract.
type PrecompiledContract interface {
	RequiredGas(input []byte) uint64  // RequiredPrice calculates the contract gas use
	Run(input []byte) ([]byte, error) // Run runs the precompiled contract
}

// PrecompiledContractsHomestead contains the default set of pre-compiled Ethereum
// contracts used in the Frontier and Homestead releases.
var PrecompiledContractsHomestead = map[common.Address]PrecompiledContract{
	common.BytesToAddress([]byte{1}): &ecrecover{},
	common.BytesToAddress([]byte{2}): &sha256hash{},
	common.BytesToAddress([]byte{3}): &ripemd160hash{},
	common.BytesToAddress([]byte{4}): &dataCopy{},
}

// PrecompiledContractsByzantium contains the default set of pre-compiled Ethereum
// contracts used in the Byzantium release.
var PrecompiledContractsByzantium = map[common.Address]PrecompiledContract{
	common.BytesToAddress([]byte{1}): &ecrecover{},
	common.BytesToAddress([]byte{2}): &sha256hash{},
	common.BytesToAddress([]byte{3}): &ripemd160hash{},
	common.BytesToAddress([]byte{4}): &dataCopy{},
	common.BytesToAddress([]byte{5}): &bigModExp{},
	common.BytesToAddress([]byte{6}): &bn256Add{},
	common.BytesToAddress([]byte{7}): &bn256ScalarMul{},
	common.BytesToAddress([]byte{8}): &bn256Pairing{},
}

// RunPrecompiledContract runs and evaluates the output of a precompiled contract.
func RunPrecompiledContract(p PrecompiledContract, input []byte, contract *Contract) (ret []byte, err error) {
	fuzz_helper.AddCoverage(15513)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	gas := p.RequiredGas(input)
	if contract.UseGas(gas) {
		fuzz_helper.AddCoverage(16403)
		return p.Run(input)
	} else {
		fuzz_helper.AddCoverage(40937)
	}
	fuzz_helper.AddCoverage(17300)
	return nil, ErrOutOfGas
}

// ECRECOVER implemented as a native contract.
type ecrecover struct{}

func (c *ecrecover) RequiredGas(input []byte) uint64 {
	fuzz_helper.AddCoverage(33825)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return params.EcrecoverGas
}

func (c *ecrecover) Run(input []byte) ([]byte, error) {
	fuzz_helper.AddCoverage(7237)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	const ecRecoverInputLength = 128

	input = common.RightPadBytes(input, ecRecoverInputLength)

	r := new(big.Int).SetBytes(input[64:96])
	s := new(big.Int).SetBytes(input[96:128])
	v := input[63] - 27

	if !allZero(input[32:63]) || !crypto.ValidateSignatureValues(v, r, s, false) {
		fuzz_helper.AddCoverage(11389)
		return nil, nil
	} else {
		fuzz_helper.AddCoverage(60629)
	}
	fuzz_helper.AddCoverage(23248)

	pubKey, err := crypto.Ecrecover(input[:32], append(input[64:128], v))

	if err != nil {
		fuzz_helper.AddCoverage(23245)
		return nil, nil
	} else {
		fuzz_helper.AddCoverage(5383)
	}
	fuzz_helper.AddCoverage(52715)

	return common.LeftPadBytes(crypto.Keccak256(pubKey[1:])[12:], 32), nil
}

// SHA256 implemented as a native contract.
type sha256hash struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
//
// This method does not require any overflow checking as the input size gas costs
// required for anything significant is so high it's impossible to pay for.
func (c *sha256hash) RequiredGas(input []byte) uint64 {
	fuzz_helper.AddCoverage(52957)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return uint64(len(input)+31)/32*params.Sha256PerWordGas + params.Sha256BaseGas
}
func (c *sha256hash) Run(input []byte) ([]byte, error) {
	fuzz_helper.AddCoverage(6211)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	h := sha256.Sum256(input)
	return h[:], nil
}

// RIPMED160 implemented as a native contract.
type ripemd160hash struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
//
// This method does not require any overflow checking as the input size gas costs
// required for anything significant is so high it's impossible to pay for.
func (c *ripemd160hash) RequiredGas(input []byte) uint64 {
	fuzz_helper.AddCoverage(49245)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return uint64(len(input)+31)/32*params.Ripemd160PerWordGas + params.Ripemd160BaseGas
}
func (c *ripemd160hash) Run(input []byte) ([]byte, error) {
	fuzz_helper.AddCoverage(15785)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	ripemd := ripemd160.New()
	ripemd.Write(input)
	return common.LeftPadBytes(ripemd.Sum(nil), 32), nil
}

// data copy implemented as a native contract.
type dataCopy struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
//
// This method does not require any overflow checking as the input size gas costs
// required for anything significant is so high it's impossible to pay for.
func (c *dataCopy) RequiredGas(input []byte) uint64 {
	fuzz_helper.AddCoverage(9735)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return uint64(len(input)+31)/32*params.IdentityPerWordGas + params.IdentityBaseGas
}
func (c *dataCopy) Run(in []byte) ([]byte, error) {
	fuzz_helper.AddCoverage(45823)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return in, nil
}

// bigModExp implements a native big integer exponential modular operation.
type bigModExp struct{}

var (
	big1      = big.NewInt(1)
	big4      = big.NewInt(4)
	big8      = big.NewInt(8)
	big16     = big.NewInt(16)
	big32     = big.NewInt(32)
	big64     = big.NewInt(64)
	big96     = big.NewInt(96)
	big480    = big.NewInt(480)
	big1024   = big.NewInt(1024)
	big3072   = big.NewInt(3072)
	big199680 = big.NewInt(199680)
)

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *bigModExp) RequiredGas(input []byte) uint64 {
	fuzz_helper.AddCoverage(48647)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	var (
		baseLen = new(big.Int).SetBytes(getData(input, 0, 32))
		expLen  = new(big.Int).SetBytes(getData(input, 32, 32))
		modLen  = new(big.Int).SetBytes(getData(input, 64, 32))
	)
	if len(input) > 96 {
		fuzz_helper.AddCoverage(63449)
		input = input[96:]
	} else {
		fuzz_helper.AddCoverage(47485)
		input = input[:0]
	}
	fuzz_helper.AddCoverage(24978)
	// Retrieve the head 32 bytes of exp for the adjusted exponent length
	var expHead *big.Int
	if big.NewInt(int64(len(input))).Cmp(baseLen) <= 0 {
		fuzz_helper.AddCoverage(60075)
		expHead = new(big.Int)
	} else {
		fuzz_helper.AddCoverage(21817)
		if expLen.Cmp(big32) > 0 {
			fuzz_helper.AddCoverage(57682)
			expHead = new(big.Int).SetBytes(getData(input, baseLen.Uint64(), 32))
		} else {
			fuzz_helper.AddCoverage(45496)
			expHead = new(big.Int).SetBytes(getData(input, baseLen.Uint64(), expLen.Uint64()))
		}
	}
	fuzz_helper.AddCoverage(61755)
	// Calculate the adjusted exponent length
	var msb int
	if bitlen := expHead.BitLen(); bitlen > 0 {
		fuzz_helper.AddCoverage(3661)
		msb = bitlen - 1
	} else {
		fuzz_helper.AddCoverage(22210)
	}
	fuzz_helper.AddCoverage(19607)
	adjExpLen := new(big.Int)
	if expLen.Cmp(big32) > 0 {
		fuzz_helper.AddCoverage(4417)
		adjExpLen.Sub(expLen, big32)
		adjExpLen.Mul(big8, adjExpLen)
	} else {
		fuzz_helper.AddCoverage(5093)
	}
	fuzz_helper.AddCoverage(28743)
	adjExpLen.Add(adjExpLen, big.NewInt(int64(msb)))

	gas := new(big.Int).Set(math.BigMax(modLen, baseLen))
	switch {
	case gas.Cmp(big64) <= 0:
		fuzz_helper.AddCoverage(56274)
		gas.Mul(gas, gas)
	case gas.Cmp(big1024) <= 0:
		fuzz_helper.AddCoverage(1404)
		gas = new(big.Int).Add(
			new(big.Int).Div(new(big.Int).Mul(gas, gas), big4),
			new(big.Int).Sub(new(big.Int).Mul(big96, gas), big3072),
		)
	default:
		fuzz_helper.AddCoverage(34815)
		gas = new(big.Int).Add(
			new(big.Int).Div(new(big.Int).Mul(gas, gas), big16),
			new(big.Int).Sub(new(big.Int).Mul(big480, gas), big199680),
		)
	}
	fuzz_helper.AddCoverage(8832)
	gas.Mul(gas, math.BigMax(adjExpLen, big1))
	gas.Div(gas, new(big.Int).SetUint64(params.ModExpQuadCoeffDiv))

	if gas.BitLen() > 64 {
		fuzz_helper.AddCoverage(58334)
		return math.MaxUint64
	} else {
		fuzz_helper.AddCoverage(46899)
	}
	fuzz_helper.AddCoverage(40052)
	return gas.Uint64()
}

func (c *bigModExp) Run(input []byte) ([]byte, error) {
	fuzz_helper.AddCoverage(40738)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	var (
		baseLen = new(big.Int).SetBytes(getData(input, 0, 32)).Uint64()
		expLen  = new(big.Int).SetBytes(getData(input, 32, 32)).Uint64()
		modLen  = new(big.Int).SetBytes(getData(input, 64, 32)).Uint64()
	)
	if len(input) > 96 {
		fuzz_helper.AddCoverage(12109)
		input = input[96:]
	} else {
		fuzz_helper.AddCoverage(29966)
		input = input[:0]
	}
	fuzz_helper.AddCoverage(10840)

	if baseLen == 0 && modLen == 0 {
		fuzz_helper.AddCoverage(27153)
		return []byte{}, nil
	} else {
		fuzz_helper.AddCoverage(51386)
	}
	fuzz_helper.AddCoverage(43066)
	// Retrieve the operands and execute the exponentiation
	var (
		base = new(big.Int).SetBytes(getData(input, 0, baseLen))
		exp  = new(big.Int).SetBytes(getData(input, baseLen, expLen))
		mod  = new(big.Int).SetBytes(getData(input, baseLen+expLen, modLen))
	)
	if mod.BitLen() == 0 {
		fuzz_helper.AddCoverage(4676)

		return common.LeftPadBytes([]byte{}, int(modLen)), nil
	} else {
		fuzz_helper.AddCoverage(61379)
	}
	fuzz_helper.AddCoverage(14816)
	return common.LeftPadBytes(base.Exp(base, exp, mod).Bytes(), int(modLen)), nil
}

var (
	// errNotOnCurve is returned if a point being unmarshalled as a bn256 elliptic
	// curve point is not on the curve.
	errNotOnCurve = errors.New("point not on elliptic curve")

	// errInvalidCurvePoint is returned if a point being unmarshalled as a bn256
	// elliptic curve point is invalid.
	errInvalidCurvePoint = errors.New("invalid elliptic curve point")
)

// newCurvePoint unmarshals a binary blob into a bn256 elliptic curve point,
// returning it, or an error if the point is invalid.
func newCurvePoint(blob []byte) (*bn256.G1, error) {
	fuzz_helper.AddCoverage(16873)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	p, onCurve := new(bn256.G1).Unmarshal(blob)
	if !onCurve {
		fuzz_helper.AddCoverage(12762)
		return nil, errNotOnCurve
	} else {
		fuzz_helper.AddCoverage(17951)
	}
	fuzz_helper.AddCoverage(20099)
	gx, gy, _, _ := p.CurvePoints()
	if gx.Cmp(bn256.P) >= 0 || gy.Cmp(bn256.P) >= 0 {
		fuzz_helper.AddCoverage(53729)
		return nil, errInvalidCurvePoint
	} else {
		fuzz_helper.AddCoverage(64454)
	}
	fuzz_helper.AddCoverage(61712)
	return p, nil
}

// newTwistPoint unmarshals a binary blob into a bn256 elliptic curve point,
// returning it, or an error if the point is invalid.
func newTwistPoint(blob []byte) (*bn256.G2, error) {
	fuzz_helper.AddCoverage(7160)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	p, onCurve := new(bn256.G2).Unmarshal(blob)
	if !onCurve {
		fuzz_helper.AddCoverage(55434)
		return nil, errNotOnCurve
	} else {
		fuzz_helper.AddCoverage(35217)
	}
	fuzz_helper.AddCoverage(7189)
	x2, y2, _, _ := p.CurvePoints()
	if x2.Real().Cmp(bn256.P) >= 0 || x2.Imag().Cmp(bn256.P) >= 0 ||
		y2.Real().Cmp(bn256.P) >= 0 || y2.Imag().Cmp(bn256.P) >= 0 {
		fuzz_helper.AddCoverage(51703)
		return nil, errInvalidCurvePoint
	} else {
		fuzz_helper.AddCoverage(39545)
	}
	fuzz_helper.AddCoverage(54641)
	return p, nil
}

// bn256Add implements a native elliptic curve point addition.
type bn256Add struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *bn256Add) RequiredGas(input []byte) uint64 {
	fuzz_helper.AddCoverage(19417)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return params.Bn256AddGas
}

func (c *bn256Add) Run(input []byte) ([]byte, error) {
	fuzz_helper.AddCoverage(47979)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	x, err := newCurvePoint(getData(input, 0, 64))
	if err != nil {
		fuzz_helper.AddCoverage(29885)
		return nil, err
	} else {
		fuzz_helper.AddCoverage(17574)
	}
	fuzz_helper.AddCoverage(22809)
	y, err := newCurvePoint(getData(input, 64, 64))
	if err != nil {
		fuzz_helper.AddCoverage(29431)
		return nil, err
	} else {
		fuzz_helper.AddCoverage(31201)
	}
	fuzz_helper.AddCoverage(20851)
	res := new(bn256.G1)
	res.Add(x, y)
	return res.Marshal(), nil
}

// bn256ScalarMul implements a native elliptic curve scalar multiplication.
type bn256ScalarMul struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *bn256ScalarMul) RequiredGas(input []byte) uint64 {
	fuzz_helper.AddCoverage(25315)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return params.Bn256ScalarMulGas
}

func (c *bn256ScalarMul) Run(input []byte) ([]byte, error) {
	fuzz_helper.AddCoverage(56379)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	p, err := newCurvePoint(getData(input, 0, 64))
	if err != nil {
		fuzz_helper.AddCoverage(33884)
		return nil, err
	} else {
		fuzz_helper.AddCoverage(49673)
	}
	fuzz_helper.AddCoverage(7887)
	res := new(bn256.G1)
	res.ScalarMult(p, new(big.Int).SetBytes(getData(input, 64, 32)))
	return res.Marshal(), nil
}

var (
	// true32Byte is returned if the bn256 pairing check succeeds.
	true32Byte = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}

	// false32Byte is returned if the bn256 pairing check fails.
	false32Byte = make([]byte, 32)

	// errBadPairingInput is returned if the bn256 pairing input is invalid.
	errBadPairingInput = errors.New("bad elliptic curve pairing size")
)

// bn256Pairing implements a pairing pre-compile for the bn256 curve
type bn256Pairing struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *bn256Pairing) RequiredGas(input []byte) uint64 {
	fuzz_helper.AddCoverage(4765)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()
	return params.Bn256PairingBaseGas + uint64(len(input)/192)*params.Bn256PairingPerPointGas
}

func (c *bn256Pairing) Run(input []byte) ([]byte, error) {
	fuzz_helper.AddCoverage(30464)
	fuzz_helper.IncrementStack()
	defer fuzz_helper.DecrementStack()

	if len(input)%192 > 0 {
		fuzz_helper.AddCoverage(26509)
		return nil, errBadPairingInput
	} else {
		fuzz_helper.AddCoverage(62902)
	}
	fuzz_helper.AddCoverage(62973)
	// Convert the input into a set of coordinates
	var (
		cs []*bn256.G1
		ts []*bn256.G2
	)
	for i := 0; i < len(input); i += 192 {
		fuzz_helper.AddCoverage(2802)
		c, err := newCurvePoint(input[i : i+64])
		if err != nil {
			fuzz_helper.AddCoverage(60446)
			return nil, err
		} else {
			fuzz_helper.AddCoverage(58133)
		}
		fuzz_helper.AddCoverage(28688)
		t, err := newTwistPoint(input[i+64 : i+192])
		if err != nil {
			fuzz_helper.AddCoverage(20790)
			return nil, err
		} else {
			fuzz_helper.AddCoverage(63383)
		}
		fuzz_helper.AddCoverage(16762)
		cs = append(cs, c)
		ts = append(ts, t)
	}
	fuzz_helper.AddCoverage(21243)

	if bn256.PairingCheck(cs, ts) {
		fuzz_helper.AddCoverage(58404)
		return true32Byte, nil
	} else {
		fuzz_helper.AddCoverage(54412)
	}
	fuzz_helper.AddCoverage(29150)
	return false32Byte, nil
}

var _ = fuzz_helper.AddCoverage
