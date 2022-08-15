package curve

import (
	"crypto/hmac"
	"crypto/rand"
	"fmt"
	"hash"
	"math/big"
)

// KeccakState wraps sha3.state. In addition to the usual hash methods, it also supports
// Read to get a variable amount of data from the hash state. Read is faster than Sum
// because it doesn't copy the internal state, but also modifies the internal state.
type KeccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

// given x will find corresponding public key coordinate on curve
func XToPubKey(x string) (*big.Int, *big.Int) {
	xin := HexToBN(x)

	yout := GetYCoordinate(xin)

	return xin, yout
}

// obtain random primary key on stark curve
// NOTE: to be used for testing purposes
func GetRandomPrivateKey() (priv *big.Int, err error) {
	max := new(big.Int).Sub(sc.Max, big.NewInt(1))

	priv, err = rand.Int(rand.Reader, max)
	if err != nil {
		return priv, err
	}

	x, y, err := PrivateToPoint(priv)
	if err != nil {
		return priv, err
	}

	if !sc.IsOnCurve(x, y) {
		return priv, fmt.Errorf("key gen is not on stark cruve")
	}

	return priv, nil
}

// obtain public key coordinates from stark curve given the private key
func PrivateToPoint(privKey *big.Int) (x, y *big.Int, err error) {
	if privKey.Cmp(big.NewInt(0)) != 1 || privKey.Cmp(sc.N) != -1 {
		return x, y, fmt.Errorf("private key not in curve range")
	}
	x, y = EcMult(privKey, sc.EcGenX, sc.EcGenY)
	return x, y, nil
}

// https://tools.ietf.org/html/rfc6979#section-2.3.3
func int2octets(v *big.Int, rolen int) []byte {
	out := v.Bytes()

	// pad with zeros if it's too short
	if len(out) < rolen {
		out2 := make([]byte, rolen)
		copy(out2[rolen-len(out):], out)
		return out2
	}

	// drop most significant bytes if it's too long
	if len(out) > rolen {
		out2 := make([]byte, rolen)
		copy(out2, out[len(out)-rolen:])
		return out2
	}

	return out
}

// https://tools.ietf.org/html/rfc6979#section-2.3.4
func bits2octets(in, q *big.Int, qlen, rolen int) []byte {
	z1 := bits2int(in, qlen)
	z2 := new(big.Int).Sub(z1, q)
	if z2.Sign() < 0 {
		return int2octets(z1, rolen)
	}
	return int2octets(z2, rolen)
}

// https://tools.ietf.org/html/rfc6979#section-2.3.2
func bits2int(in *big.Int, qlen int) *big.Int {
	blen := len(in.Bytes()) * 8

	if blen > qlen {

		return new(big.Int).Rsh(in, uint(blen-qlen))
	}
	return in
}

// mac returns an HMAC of the given key and message.
func mac(alg func() hash.Hash, k, m, buf []byte) []byte {
	h := hmac.New(alg, k)
	h.Write(m)
	return h.Sum(buf[:0])
}
