package caigo

import (
	"fmt"
	"bytes"
	"hash"
	"math/big"
	"crypto/sha256"
)

func (sc StarkCurve) Verify(msgHash, r, s, pubX, pubY *big.Int) bool {
	w := sc.InvModCurveSize(s)

	if s.Cmp(big.NewInt(0)) != 1 || s.Cmp(sc.N) != -1 {
		return false
	}
	if r.Cmp(big.NewInt(0)) != 1 || r.BitLen() > 502 {
		return false
	}
	if w.Cmp(big.NewInt(0)) != 1 || w.BitLen() > 502 {
		return false
	}
	if msgHash.Cmp(big.NewInt(0)) != 1 || msgHash.BitLen() > 502 {
		return false
	}
	if !sc.IsOnCurve(pubX, pubY) {
		return false
	}

	zGx, zGy, err := sc.MimicEcMultAir(msgHash, sc.EcGenX, sc.EcGenY, sc.MinusShiftPointX, sc.MinusShiftPointY)
	if err != nil {
		return false
	}

	rQx, rQy, err := sc.MimicEcMultAir(r, pubX, pubY, sc.Gx, sc.Gy)
	if err != nil {
		return false
	}
	inX, inY := sc.Add(zGx, zGy, rQx, rQy)
	wBx, wBy, err := sc.MimicEcMultAir(w, inX, inY, sc.Gx, sc.Gy)
	if err != nil {
		return false
	}

	outX, _ := sc.Add(wBx, wBy, sc.MinusShiftPointX, sc.MinusShiftPointY)
	if r.Cmp(outX) == 0 {
		return true
	} else {
		altY := new(big.Int)
		altY = altY.Neg(pubY)
		altY = altY.Mod(altY, sc.P)
		pubY = altY

		zGx, zGy, err = sc.MimicEcMultAir(msgHash, sc.EcGenX, sc.EcGenY, sc.MinusShiftPointX, sc.MinusShiftPointY)
		if err != nil {
			return false
		}

		rQx, rQy, err = sc.MimicEcMultAir(r, pubX, pubY, sc.Gx, sc.Gy)
		if err != nil {
			return false
		}
		inX, inY = sc.Add(zGx, zGy, rQx, rQy)
		wBx, wBy, err = sc.MimicEcMultAir(w, inX, inY, sc.Gx, sc.Gy)
		if err != nil {
			return false
		}

		outX, _ = sc.Add(wBx, wBy, sc.MinusShiftPointX, sc.MinusShiftPointY)
		if r.Cmp(outX) == 0 {
			return true
		}
	}
	return false
}

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

func (sc StarkCurve) HashElements(elems []*big.Int) (hash *big.Int, err error) {
	if len(elems) < 2 {
		return hash, fmt.Errorf("must have element slice larger than 2 for hashing")
	}
	hash = big.NewInt(0)
	for _, h := range elems {
		hash, err = sc.PedersenHash([]*big.Int{hash, h})
		if err != nil {
			return hash, err
		}
	}
	return hash, err
}

func (sc StarkCurve) PedersenHash(elems []*big.Int) (hash *big.Int, err error) {
	if len(sc.ConstantPoints) == 0 {
		return hash, fmt.Errorf("must initiate precomputed constant points")
	}

	ptx := new(big.Int)
	pty := new(big.Int)
	ptx = ptx.Set(sc.Gx)
	pty = pty.Set(sc.Gy)
	for i, elem := range elems {
		x := new(big.Int)
		x = x.Set(elem)

		if x.Cmp(big.NewInt(0)) != 1 && x.Cmp(sc.P) != -1 {
			return hash, fmt.Errorf("invalid x: %v", x)
		}

		for j := 0; j < 252; j++ {
			idx := 2 + (i * 252) + j
			xin := new(big.Int)
			yin := new(big.Int)
			xin = xin.Set(sc.ConstantPoints[idx][0])
			yin = yin.Set(sc.ConstantPoints[idx][1])
			if xin.Cmp(ptx) == 0 {
				return hash, fmt.Errorf("constant point duplication: %v %v", ptx, xin)
			}
			if x.Bit(0) == 1 {
				ptx, pty = sc.Add(ptx, pty, xin, yin)
			}
			x = x.Rsh(x, 1)
		}
	}

	return ptx, nil
}

// implementation based on https://github.com/codahale/rfc6979/blob/master/rfc6979.go
func (sc StarkCurve) GenerateSecret(msgHash, privKey *big.Int, alg func() hash.Hash) (secret *big.Int) {
	holen := alg().Size()
	rolen := (sc.BitSize + 7) >> 3

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
		// TODO: implement seed here, final gating function
		if secret.Cmp(big.NewInt(0)) == 1 && secret.Cmp(sc.N) == -1 {
			return secret
		}
		k = mac(alg, k, append(v, 0x00), k)
		v = mac(alg, k, v, v)
	}

	return secret
}
