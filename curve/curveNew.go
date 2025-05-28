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

func GetYCoordinate(starkX *big.Int) *big.Int {
	// ref: https://github.com/NethermindEth/juno/blob/7d64642de90b6957c40a3b3ea75e6ad548a37f39/core/crypto/ecdsa.go#L26
	xEl := new(fp.Element).SetBigInt(starkX)

	var ySquared fp.Element
	ySquared.Mul(xEl, xEl).Mul(&ySquared, xEl) // x^3
	ySquared.Add(&ySquared, xEl)               // + x

	_, b := starkcurve.CurveCoefficients()
	ySquared.Add(&ySquared, &b) // ySquared equals to (x^3 + x + b)
	return ySquared.Sqrt(&ySquared).BigInt(new(big.Int))
}

func Verify(msgHash, r, s, pubX, pubY *big.Int) (bool, error) {
	pubKey := crypto.NewPublicKey(new(felt.Felt).SetBigInt(pubX))
	msgHashFelt := new(felt.Felt).SetBigInt(msgHash)
	rFelt := new(felt.Felt).SetBigInt(r)
	sFelt := new(felt.Felt).SetBigInt(s)

	return pubKey.Verify(&crypto.Signature{R: *rFelt, S: *sFelt}, msgHashFelt)
}

func Sign(msgHash, privKey *big.Int, seed ...*big.Int) (r, s *big.Int, err error) {
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

func GetRandomPrivateKey() (*big.Int, error) {
	priv, err := ecdsa.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	privBytes := priv.Bytes()           // a 64 bytes array containing both public (compressed) and private keys
	finalPrivKeyBytes := privBytes[32:] // the remaining 32 bytes are the private key

	finalPrivKey := new(big.Int).SetBytes(finalPrivKeyBytes)

	return finalPrivKey, nil
}

func PrivateToPoint(privKey *big.Int) (x, y *big.Int, err error) {
	g1a := new(starkcurve.G1Affine)
	res := g1a.ScalarMultiplicationBase(privKey)
	return res.X.BigInt(new(big.Int)), res.Y.BigInt(new(big.Int)), nil
}
