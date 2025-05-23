package curve

/*
	Although the library adheres to the 'elliptic/curve' interface.
	All testing has been done against library function explicity.
	It is recommended to use in the same way(i.e. `curve.Sign` and not `ecdsa.Sign`).
*/
import (
	"crypto/rand"
	_ "embed"
	"math/big"

	"github.com/NethermindEth/juno/core/crypto"
	"github.com/NethermindEth/juno/core/felt"
	starkcurve "github.com/consensys/gnark-crypto/ecc/stark-curve"
	"github.com/consensys/gnark-crypto/ecc/stark-curve/ecdsa"
	"github.com/consensys/gnark-crypto/ecc/stark-curve/fp"
)

var CurveNew StarkCurveNew

type StarkCurveNew starkcurve.G1Affine

func (sc StarkCurveNew) Add(x1, y1, x2, y2 *big.Int) (x, y *big.Int) {
	x1Felt := new(felt.Felt).SetBigInt(x1)
	y1Felt := new(felt.Felt).SetBigInt(y1)
	g1a1 := starkcurve.G1Affine{X: *x1Felt.Impl(), Y: *y1Felt.Impl()}

	x2Felt := new(felt.Felt).SetBigInt(x2)
	y2Felt := new(felt.Felt).SetBigInt(y2)
	g1a2 := starkcurve.G1Affine{X: *x2Felt.Impl(), Y: *y2Felt.Impl()}

	curve := starkcurve.G1Affine(sc)
	curve.Add(&g1a1, &g1a2)

	return curve.X.BigInt(new(big.Int)), curve.Y.BigInt(new(big.Int))
}

func (sc StarkCurveNew) GetYCoordinate(starkX *big.Int) *big.Int {
	// ref: https://github.com/NethermindEth/juno/blob/7d64642de90b6957c40a3b3ea75e6ad548a37f39/core/crypto/ecdsa.go#L26
	xEl := new(fp.Element).SetBigInt(starkX)

	var ySquared fp.Element
	ySquared.Mul(xEl, xEl).Mul(&ySquared, xEl) // x^3
	ySquared.Add(&ySquared, xEl)               // + x

	_, b := starkcurve.CurveCoefficients()
	ySquared.Add(&ySquared, &b) // ySquared equals to (x^3 + x + b)
	return ySquared.Sqrt(&ySquared).BigInt(new(big.Int))
}

func (sc StarkCurveNew) Verify(msgHash, r, s, pubX, pubY *big.Int) (bool, error) {
	pubKey := crypto.NewPublicKey(new(felt.Felt).SetBigInt(pubX))
	msgHashFelt := new(felt.Felt).SetBigInt(msgHash)
	rFelt := new(felt.Felt).SetBigInt(r)
	sFelt := new(felt.Felt).SetBigInt(s)

	return pubKey.Verify(&crypto.Signature{R: *rFelt, S: *sFelt}, msgHashFelt)
}

func (sc StarkCurveNew) Sign(msgHash, privKey *big.Int, seed ...*big.Int) (r, s *big.Int, err error) {
	// generating pub and priv keys
	g1a := new(starkcurve.G1Affine).ScalarMultiplicationBase(privKey)

	var pubKeyStruct ecdsa.PublicKey
	pubKeyBytes := g1a.Bytes()
	_, err = pubKeyStruct.SetBytes(pubKeyBytes[:])
	if err != nil {
		return nil, nil, err
	}

	var privKeyStruct ecdsa.PrivateKey
	privKeyBytes := privKey.Bytes()
	privKeyInput := append(pubKeyStruct.Bytes(), privKeyBytes...)
	_, err = privKeyStruct.SetBytes(privKeyInput[:])
	if err != nil {
		return nil, nil, err
	}

	// signing
	_, r, s, err = privKeyStruct.SignForRecover(msgHash.Bytes(), nil)

	return r, s, err
}

func (sc StarkCurveNew) GetRandomPrivateKey() (string, error) {
	// TODO: this is a test
	priv, _ := ecdsa.GenerateKey(rand.Reader)
	pub := priv.PublicKey

	var temp any = pub.A.X

	print(temp)
	return pub.A.X.String(), nil
}

func (sc StarkCurveNew) PrivateToPoint(privKey *big.Int) (x, y *big.Int, err error) {
	g1a := starkcurve.G1Affine(sc)
	res := g1a.ScalarMultiplicationBase(privKey)
	return res.X.BigInt(new(big.Int)), res.Y.BigInt(new(big.Int)), nil
}
