package caigo

import (
	"fmt"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"hash"
	"math/big"
)

func (sc StarkCurve) Sign(msgHash, privKey *big.Int) (x, y *big.Int, err error) {
	if msgHash.Cmp(big.NewInt(0)) != 1 || msgHash.BitLen() > 502 {
		return x, y, fmt.Errorf("invalid bit length")
	}

	invalidK := true 
	for invalidK {
		k := sc.GenerateSecret(msgHash, privKey, sha256.New)

		r := new(big.Int)
		kin := new(big.Int)
		kin = kin.Set(k)
		r, _ = sc.EcMult(kin, sc.EcGenX, sc.EcGenY)

		// DIFF: in classic ECDSA, we take int(x) % n.
		if r.Cmp(big.NewInt(0)) != 1 || r.BitLen() >= 502 {
			// Bad value. This fails with negligible probability.
			continue
		}

		agg := new(big.Int)
		agg = agg.Mul(r, privKey)
		agg = agg.Add(agg, msgHash)

		cagg := new(big.Int)
		cagg = cagg.Mod(agg, sc.N)
		if cagg.Cmp(big.NewInt(0)) == 0 {
			// Bad value. This fails with negligible probability.
			continue
		}

		w := new(big.Int)
		w = DivMod(k, agg, sc.N)
		if w.Cmp(big.NewInt(0)) != 1 || w.BitLen() >= 502 {
			// Bad value. This fails with negligible probability.
			continue
		}

		s := new(big.Int)
		s = sc.InvModCurveSize(w)
		return r, s, nil

	}

	return x, y, nil
}

// https://github.com/codahale/rfc6979/blob/master/rfc6979.go
func (sc StarkCurve) GenerateSecret(msgHash, privKey *big.Int, alg func() hash.Hash) (secret *big.Int) {
	holen := alg().Size() //32
	rolen := (sc.BitSize + 7) >> 3 //32

	by := append(int2octets(privKey, rolen), bits2octets(msgHash.Bytes(), sc.N, sc.BitSize, rolen)...)

	v := bytes.Repeat([]byte{0x01}, holen)
	
	k := bytes.Repeat([]byte{0x00}, holen)

	k = mac(alg, k, append(append(v, 0x00), by...), k)
	
	v = mac(alg, k, v, v)
	
	k = mac(alg, k, append(append(v, 0x01), by...), k)
	
	v = mac(alg, k, v, v)

	for {
		var t []byte

		for len(t) < sc.BitSize/8 {
			v = mac(alg, k, v, v)
			t = append(t, v...)
		}

		secret = bits2int(t, sc.BitSize)
		// TODO implement seed here, final gating function
		if secret.Cmp(big.NewInt(0)) == 1 && secret.Cmp(sc.N) == -1 {
			return secret
		}
		k = mac(alg, k, append(v, 0x00), k)
		v = mac(alg, k, v, v)
	}

	return secret
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
func bits2octets(in []byte, q *big.Int, qlen, rolen int) []byte {
	z1 := bits2int(in, qlen)
	z2 := new(big.Int).Sub(z1, q)
	if z2.Sign() < 0 {
		return int2octets(z1, rolen)
	}
	return int2octets(z2, rolen)
}

// https://tools.ietf.org/html/rfc6979#section-2.3.2
func bits2int(in []byte, qlen int) *big.Int {
	vlen := len(in) * 8
	v := new(big.Int).SetBytes(in)
	if vlen > qlen {
		v = new(big.Int).Rsh(v, uint(vlen-qlen))
	}
	return v
}

// mac returns an HMAC of the given key and message.
func mac(alg func() hash.Hash, k, m, buf []byte) []byte {
	h := hmac.New(alg, k)
	h.Write(m)
	return h.Sum(buf[:0])
}