package caigo

import (
	"fmt"
	"math/big"
	"crypto/ecdsa"
)

// N_ELEMENT_BITS_ECDSA = math.floor(math.log(FIELD_PRIME, 2))
// assert N_ELEMENT_BITS_ECDSA == 251

// N_ELEMENT_BITS_HASH = FIELD_PRIME.bit_length()
// assert N_ELEMENT_BITS_HASH == 252

func Verify(msgHash, r, s *big.Int, pub ecdsa.PublicKey, sc StarkCurve) bool {
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
	if !sc.IsOnCurve(pub.X, pub.Y) {
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

	rQx, rQy, err := sc.MimicEcMultAir(r, pub.X, pub.Y, sc.Gx, sc.Gy)
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
		altY = altY.Neg(pub.Y)
		altY = altY.Mod(altY, sc.P)
		pub.Y = altY

		rSigIn := new(big.Int)
		rSigIn = rSigIn.Set(rSig)

		zGx, zGy, err = sc.MimicEcMultAir(holdHash, sc.EcGenX, sc.EcGenY, sc.MinusShiftPointX, sc.MinusShiftPointY)
		if err != nil {
			return false
		}

		rQx, rQy, err = sc.MimicEcMultAir(rSig, pub.X, pub.Y, sc.Gx, sc.Gy)
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

func hashElements(elems []*big.Int) (hash *big.Int, err error) {
	return hash, err
}

func pedersenHash(x, y *big.Int) (hash *big.Int, err error) {
	curv := SC()

	if x.Cmp(big.NewInt(0)) != 1 && x.Cmp(curv.P) != -1 {
		return hash, fmt.Errorf("invalid x: %v", x)
	}
	if y.Cmp(big.NewInt(0)) != 1 && y.Cmp(curv.P) != -1 {
		return hash, fmt.Errorf("invalid y: %v", y)
	}

	for i := 0; i < 252; i++ {
		ptx := new(big.Int)
		pty := new(big.Int)
		fmt.Println("PTS: ", ptx, pty)
		// ptx := ptx.
		
	}
	return hash, err
}