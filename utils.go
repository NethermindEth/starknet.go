package caigo

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// obtain random primary key on stark curve
// NOTE: to be used for testing purposes
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

// obtain public key coordinates from stark curve given the private key
func (sc StarkCurve) PrivateToPoint(privKey *big.Int) (x, y *big.Int, err error) {
	return sc.privateToPoint(privKey, sc.EcMult)
}

// ec multiplication fn
type EcMultiFn func(m, x1, y1 *big.Int) (x, y *big.Int)

// obtain public key coordinates from stark curve given the private key
// NOTICE: configurable ec multiplication fn, used for testing
func (sc StarkCurve) privateToPoint(privKey *big.Int, ecMulti EcMultiFn) (x, y *big.Int, err error) {
	if privKey.Cmp(big.NewInt(0)) != 1 || privKey.Cmp(sc.N) != -1 {
		return x, y, fmt.Errorf("private key not in curve range")
	}

	x, y = ecMulti(privKey, sc.EcGenX, sc.EcGenY)
	return x, y, nil
}
