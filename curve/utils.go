package curve

import (
	"crypto/hmac"
	"fmt"
	"hash"
	"math/big"
)

// Finds a nonnegative integer 0 <= x < p such that (m * x) % p == n
//
// (ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/crypto/signature/math_utils.py#L50)
func DivMod(n, m, p *big.Int) *big.Int {
	q := new(big.Int)
	gx := new(big.Int)
	gy := new(big.Int)
	q.GCD(gx, gy, m, p)

	r := new(big.Int).Mul(n, gx)
	r = r.Mod(r, p)
	fmt.Println(r)
	return r
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

// mask excess bits
func MaskBits(mask, wordSize int, slice []byte) (ret []byte) {
	excess := len(slice)*wordSize - mask
	for _, by := range slice {
		if excess > 0 {
			if excess > wordSize {
				excess = excess - wordSize
				continue
			}
			by <<= excess
			by >>= excess
			excess = 0
		}
		ret = append(ret, by)
	}
	return ret
}

// format the bytes in Keccak hash
func FmtKecBytes(in *big.Int, rolen int) (buf []byte) {
	buf = append(buf, in.Bytes()...)

	// pad with zeros if too short
	if len(buf) < rolen {
		padded := make([]byte, rolen)
		copy(padded[rolen-len(buf):], buf)

		return padded
	}

	return buf
}
