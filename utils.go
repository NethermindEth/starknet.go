package starknetgo

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GetRandomPrivateKey generates a random primary key on Stark elliptic curve
// NOTE: to be used for testing purposes
//
// It returns the generated private key as *big.Int and an error if any.
func (sc StarkCurve) GetRandomPrivateKey() (priv *big.Int, err error) {
	max := new(big.Int).Sub(sc.Max, big.NewInt(1))

	priv, err = rand.Int(rand.Reader, max)
	if err != nil {
		return priv, err
	}

	x, y, err := sc.PrivateToPoint(priv)
	if err != nil {
		return priv, err
	}

	if !sc.IsOnCurve(x, y) {
		return priv, fmt.Errorf("key gen is not on stark cruve")
	}

	return priv, nil
}

// PrivateToPoint calculates the coordinates of a point on the Stark elliptic curve given a private key.
//
// The function takes a private key as a parameter and returns the x and y coordinates of the
// corresponding point on the curve. If the private key is not within the range of the curve,
// an error is returned.
//
// Parameters:
// - privKey: A pointer to a big.Int representing the private key.
//
// Returns:
// - x: A pointer to a big.Int representing the x-coordinate of the point.
// - y: A pointer to a big.Int representing the y-coordinate of the point.
// - err: An error indicating if the private key is not within the curve range.
func (sc StarkCurve) PrivateToPoint(privKey *big.Int) (x, y *big.Int, err error) {
	if privKey.Cmp(big.NewInt(0)) != 1 || privKey.Cmp(sc.N) != -1 {
		return x, y, fmt.Errorf("private key not in curve range")
	}
	x, y = sc.EcMult(privKey, sc.EcGenX, sc.EcGenY)
	return x, y, nil
}
