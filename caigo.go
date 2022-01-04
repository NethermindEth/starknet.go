package caigo

import (
	"fmt"
	"math/big"
)

// N_ELEMENT_BITS_ECDSA = math.floor(math.log(FIELD_PRIME, 2))
// assert N_ELEMENT_BITS_ECDSA == 251

// N_ELEMENT_BITS_HASH = FIELD_PRIME.bit_length()
// assert N_ELEMENT_BITS_HASH == 252

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
	holdHash := new(big.Int)
	holdHash = holdHash.Set(msgHash)
	wHold := new(big.Int)
	wHold = wHold.Set(w)

	rSig := new(big.Int)
	rSig = rSig.Set(r)

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
	if rSig.Cmp(outX) == 0 {
		return true
	} else {
		altY := new(big.Int)
		altY = altY.Neg(pubY)
		altY = altY.Mod(altY, sc.P)
		pubY = altY

		rSigIn := new(big.Int)
		rSigIn = rSigIn.Set(rSig)

		zGx, zGy, err = sc.MimicEcMultAir(holdHash, sc.EcGenX, sc.EcGenY, sc.MinusShiftPointX, sc.MinusShiftPointY)
		if err != nil {
			return false
		}

		rQx, rQy, err = sc.MimicEcMultAir(rSig, pubX, pubY, sc.Gx, sc.Gy)
		if err != nil {
			return false
		}
		inX, inY = sc.Add(zGx, zGy, rQx, rQy)
		wBx, wBy, err = sc.MimicEcMultAir(wHold, inX, inY, sc.Gx, sc.Gy)
		if err != nil {
			return false
		}

		outX, _ = sc.Add(wBx, wBy, sc.MinusShiftPointX, sc.MinusShiftPointY)
		if rSigIn.Cmp(outX) == 0 {
			return true
		}
	}
	return false
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
	// TODO: implement/test the starkware fast pedersen hash 
	if (len(sc.ConstantPoints) == 0) {
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
