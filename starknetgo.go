package starknetgo

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"hash"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
)

// Verify verifies the validity of the stark curve signature of a message hash using the Stark elliptic curve used to sign the message.
// (ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/crypto/starkware/crypto/signature/signature.py)
//
// Parameters:
// - msgHash: The hash of the message to be verified.
// - r: The r component of the signature.
// - s: The s component of the signature.
// - pubX: The x-coordinate of the public key.
// - pubY: The y-coordinate of the public key.
//
// Returns:
// - bool: True if the signature is valid, false otherwise.
func (sc StarkCurve) Verify(msgHash, r, s, pubX, pubY *big.Int) bool {
	w := sc.InvModCurveSize(s)

	if s.Cmp(big.NewInt(0)) != 1 || s.Cmp(sc.N) != -1 {
		return false
	}
	if r.Cmp(big.NewInt(0)) != 1 || r.Cmp(sc.Max) != -1 {
		return false
	}
	if w.Cmp(big.NewInt(0)) != 1 || w.Cmp(sc.Max) != -1 {
		return false
	}
	if msgHash.Cmp(big.NewInt(0)) != 1 || msgHash.Cmp(sc.Max) != -1 {
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
		altY := new(big.Int).Neg(pubY)

		zGx, zGy, err = sc.MimicEcMultAir(msgHash, sc.EcGenX, sc.EcGenY, sc.MinusShiftPointX, sc.MinusShiftPointY)
		if err != nil {
			return false
		}

		rQx, rQy, err = sc.MimicEcMultAir(r, pubX, new(big.Int).Set(altY), sc.Gx, sc.Gy)
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

// Sign signs the hash value of contents with a private key using the Stark elliptic curve algorithm.
// Secret is generated using a golang implementation of RFC 6979.
// (ref: https://datatracker.ietf.org/doc/html/rfc6979)
//
// Implementation does not yet include "extra entropy" or "retry gen".
//
// Parameters:
// - msgHash: The message hash to be signed.
// - privKey: The private key used for signing.
// - seed: (optional) Additional seed values for generating the secret.
//
// Returns:
// - x: The x-coordinate of the resulting point on the curve.
// - y: The y-coordinate of the resulting point on the curve.
// - err: An error if any occurred during the signing process.
func (sc StarkCurve) Sign(msgHash, privKey *big.Int, seed ...*big.Int) (x, y *big.Int, err error) {
	if msgHash == nil {
		return x, y, fmt.Errorf("nil msgHash")
	}
	if privKey == nil {
		return x, y, fmt.Errorf("nil privKey")
	}
	if msgHash.Cmp(big.NewInt(0)) != 1 || msgHash.Cmp(sc.Max) != -1 {
		return x, y, fmt.Errorf("invalid bit length")
	}

	inSeed := big.NewInt(0)
	if len(seed) == 1 && inSeed != nil {
		inSeed = seed[0]
	}
	for {
		k := sc.GenerateSecret(big.NewInt(0).Set(msgHash), big.NewInt(0).Set(privKey), big.NewInt(0).Set(inSeed))
		// In case r is rejected k shall be generated with new seed
		inSeed = inSeed.Add(inSeed, big.NewInt(1))

		r, _ := sc.EcMult(k, sc.EcGenX, sc.EcGenY)

		// DIFF: in classic ECDSA, we take int(x) % n.
		if r.Cmp(big.NewInt(0)) != 1 || r.Cmp(sc.Max) != -1 {
			// Bad value. This fails with negligible probability.
			continue
		}

		agg := new(big.Int).Mul(r, privKey)
		agg = agg.Add(agg, msgHash)

		if new(big.Int).Mod(agg, sc.N).Cmp(big.NewInt(0)) == 0 {
			// Bad value. This fails with negligible probability.
			continue
		}

		w := DivMod(k, agg, sc.N)
		if w.Cmp(big.NewInt(0)) != 1 || w.Cmp(sc.Max) != -1 {
			// Bad value. This fails with negligible probability.
			continue
		}

		s := sc.InvModCurveSize(w)
		return r, s, nil
	}

	return x, y, nil
}

// SignFelt signs a message hash with a private key using the Stark elliptic curve.
// See Sign. SignFelt just wraps Sign.
//
// Parameters:
// - msgHash: The message hash to be signed, of type *felt.Felt.
// - privKey: The private key to sign the message hash, of type *felt.Felt.
//
// Returns:
// - xFelt: The x-coordinate of the generated signature, of type *felt.Felt.
// - yFelt: The y-coordinate of the generated signature, of type *felt.Felt.
// - error: Any error that occurred during the signing process.
func (sc StarkCurve) SignFelt(msgHash, privKey *felt.Felt) (*felt.Felt, *felt.Felt, error) {
	msgHashInt := msgHash.BigInt(new(big.Int))
	privKeyInt := privKey.BigInt(new(big.Int))
	x, y, err := sc.Sign(msgHashInt, privKeyInt)
	if err != nil {
		return nil, nil, err
	}
	xFelt := felt.NewFelt(new(felt.Felt).Impl().SetBigInt(x))
	yFelt := felt.NewFelt(new(felt.Felt).Impl().SetBigInt(y))
	return xFelt, yFelt, nil

}

// HashElements calculates the hash value of the contents of a given array using a golang Pedersen Hash implementation.
// (ref: https://github.com/seanjameshan/starknet.js/blob/main/src/utils/ellipticCurve.ts)
//
// It takes a slice of big integers as input and returns the resulting hash value and an error, if any.
func (sc StarkCurve) HashElements(elems []*big.Int) (hash *big.Int, err error) {
	if len(elems) == 0 {
		elems = append(elems, big.NewInt(0))
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

// ComputeHashOnElements computes the hash of a given array with its size using a golang Pedersen Hash implementation.
// (ref: https://github.com/starkware-libs/cairo-lang/blob/13cef109cd811474de114925ee61fd5ac84a25eb/src/starkware/cairo/common/hash_state.py#L6)
//
// It takes a slice of big.Int pointers as the parameter `elems`.
// It returns a big.Int pointer `hash` and an error `err`.
func (sc StarkCurve) ComputeHashOnElements(elems []*big.Int) (hash *big.Int, err error) {
	elems = append(elems, big.NewInt(int64(len(elems))))
	return Curve.HashElements((elems))
}

// PedersenHash calculates the Pedersen hash of the given array of big integers.
// NOTE: This function assumes the curve has been initialized with contant points
// (ref: https://github.com/seanjameshan/starknet.js/blob/main/src/utils/ellipticCurve.ts)
//
// Parameters:
// - elems: a slice of big integers representing the elements to hash.
//
// Returns:
// - hash: the resulting hash as a big integer.
// - err: an error if the precomputed constant points are not initiated or if there is an invalid x value or a constant point duplication.
func (sc StarkCurve) PedersenHash(elems []*big.Int) (hash *big.Int, err error) {
	if len(sc.ConstantPoints) == 0 {
		return hash, fmt.Errorf("must initiate precomputed constant points")
	}

	ptx := new(big.Int).Set(sc.Gx)
	pty := new(big.Int).Set(sc.Gy)
	for i, elem := range elems {
		x := new(big.Int).Set(elem)

		if x.Cmp(big.NewInt(0)) != -1 && x.Cmp(sc.P) != -1 {
			return ptx, fmt.Errorf("invalid x: %v", x)
		}

		for j := 0; j < 252; j++ {
			idx := 2 + (i * 252) + j
			xin := new(big.Int).Set(sc.ConstantPoints[idx][0])
			yin := new(big.Int).Set(sc.ConstantPoints[idx][1])
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

// GenerateSecret generates a secret using the StarkCurve struct in Go.
// implementation based on https://github.com/codahale/rfc6979/blob/master/rfc6979.go
// for the specification, see https://tools.ietf.org/html/rfc6979#section-3.2
//
// Parameters:
// - msgHash: The hash of the message as a big.Int pointer.
// - privKey: The private key as a big.Int pointer.
// - seed: The seed as a big.Int pointer.
//
// Return:
// - secret: The generated secret as a big.Int pointer.
func (sc StarkCurve) GenerateSecret(msgHash, privKey, seed *big.Int) (secret *big.Int) {
	alg := sha256.New
	holen := alg().Size()
	rolen := (sc.BitSize + 7) >> 3

	if msgHash.BitLen()%8 <= 4 && msgHash.BitLen() >= 248 {
		msgHash = msgHash.Mul(msgHash, big.NewInt(16))
	}

	by := append(int2octets(privKey, rolen), bits2octets(msgHash, sc.N, sc.BitSize, rolen)...)

	if seed.Cmp(big.NewInt(0)) == 1 {
		by = append(by, seed.Bytes()...)
	}

	v := bytes.Repeat([]byte{0x01}, holen)

	k := bytes.Repeat([]byte{0x00}, holen)

	k = mac(alg, k, append(append(v, 0x00), by...), k)

	v = mac(alg, k, v, v)

	k = mac(alg, k, append(append(v, 0x01), by...), k)

	v = mac(alg, k, v, v)

	for {
		var t []byte

		for len(t) < rolen {
			v = mac(alg, k, v, v)
			t = append(t, v...)
		}

		secret = bits2int(new(big.Int).SetBytes(t), sc.BitSize)
		// TODO: implement seed here, final gating function
		if secret.Cmp(big.NewInt(0)) == 1 && secret.Cmp(sc.N) == -1 {
			return secret
		}
		k = mac(alg, k, append(v, 0x00), k)
		v = mac(alg, k, v, v)
	}
}

// int2octets converts a big integer to a byte slice of specified length.
// https://tools.ietf.org/html/rfc6979#section-2.3.3
//
// v: The big integer to convert.
// rolen: The desired length of the byte slice.
// Returns a byte slice representing the big integer, padded with zeros if it's shorter than the desired length,
// or truncated from the most significant bytes if it's longer.
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

// bits2octets generates octets from bits.
// https://tools.ietf.org/html/rfc6979#section-2.3.4
//
// It takes in two parameters:
// - in: a pointer to a big.Int, the input value in bits.
// - q: a pointer to a big.Int, the modulus value in bits.
//
// It also takes in two integers as parameters:
// - qlen: the length of the modulus value in bits.
// - rolen: the length of the output value in octets.
//
// It returns a byte slice representing the octets.
func bits2octets(in, q *big.Int, qlen, rolen int) []byte {
	z1 := bits2int(in, qlen)
	z2 := new(big.Int).Sub(z1, q)
	if z2.Sign() < 0 {
		return int2octets(z1, rolen)
	}
	return int2octets(z2, rolen)
}

// bits2int returns a new big.Int that is the result of converting the given big.Int to an integer with a length of qlen bits.
// https://tools.ietf.org/html/rfc6979#section-2.3.2
//
// It takes the following parameter(s):
// - in: a pointer to a big.Int, the input value to convert.
// - qlen: an integer, the desired length of the resulting integer in bits.
//
// It returns a pointer to a big.Int, the converted integer.
func bits2int(in *big.Int, qlen int) *big.Int {
	blen := len(in.Bytes()) * 8

	if blen > qlen {

		return new(big.Int).Rsh(in, uint(blen-qlen))
	}
	return in
}

// mac calculates the message authentication code (MAC) using the specified hash algorithm, key, message, and buffer.
//
// The alg parameter is a function that returns a hash.Hash implementation, which will be used to calculate the MAC.
// The k parameter is the key used for the MAC calculation.
// The m parameter is the message for which the MAC is being calculated.
// The buf parameter is a buffer used to store the result of the MAC calculation.
//
// The function returns the calculated MAC as a byte slice.
func mac(alg func() hash.Hash, k, m, buf []byte) []byte {
	h := hmac.New(alg, k)
	h.Write(m)
	return h.Sum(buf[:0])
}


// MaskBits masks the bits in a byte slice based on the given mask and word size.
//
// It takes three parameters: 
//   - mask: the number of bits to mask
//   - wordSize: the size of each word in bits
//   - slice: the byte slice to mask
//
// It returns a new byte slice with the masked bits.
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

// FmtKecBytes formats the given big integer as a byte slice with a specified length.
// format the bytes in Keccak hash
//
// It takes in a big integer 'in' and an integer 'rolen' which specifies the desired length of the byte slice.
// It returns a byte slice 'buf' which contains the formatted bytes of the big integer.
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
