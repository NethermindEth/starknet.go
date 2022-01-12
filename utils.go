package caigo

import (
	"crypto/hmac"
	"encoding/hex"
	"fmt"
	"hash"
	"math/big"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/sha3"
)

// KeccakState wraps sha3.state. In addition to the usual hash methods, it also supports
// Read to get a variable amount of data from the hash state. Read is faster than Sum
// because it doesn't copy the internal state, but also modifies the internal state.
type KeccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

// given x will find corresponding public key coordinate on curve
func (sc StarkCurve) XToPubKey(x string) (*big.Int, *big.Int) {
	xin := HexToBN(x)

	yout := sc.GetYCoordinate(xin)

	return xin, yout
}

// convert utf8 string to big int
func UTF8StrToBig(str string) *big.Int {
	hexStr := hex.EncodeToString([]byte(str))
	b, _ := new(big.Int).SetString(hexStr, 16)

	return b
}

// convert decimal string to big int
func StrToBig(str string) *big.Int {
	b, _ := new(big.Int).SetString(str, 10)

	return b
}

// trim "0x" prefix(if exists) and converts hexidecimal string to big int
func HexToBN(hexString string) (n *big.Int) {
	numStr := strings.Replace(hexString, "0x", "", -1)

	n = new(big.Int)
	n.SetString(numStr, 16)
	return n
}

// trim "0x" prefix(if exists) and converts hexidecimal string to byte slice
func HexToBytes(hexString string) ([]byte, error) {
	numStr := strings.Replace(hexString, "0x", "", -1)
	if (len(numStr) % 2) != 0 {
		numStr = fmt.Sprintf("%s%s", "0", numStr)
	}

	return hex.DecodeString(numStr)
}

// convert big int to hexidecimal string
func BigToHex(in *big.Int) string {
	return fmt.Sprintf("0x%x", in)
}

// obtain random primary key on stark curve
// NOTE: to be used for testing purposes
func (sc StarkCurve) GetRandomPrivateKey() *big.Int {
	max := new(big.Int)
	max = max.Sub(sc.N, big.NewInt(1))

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	priv := new(big.Int)
	priv = priv.Rand(r, max)
	return priv
}

// obtain public key coordinates from stark curve given the private key
func (sc StarkCurve) PrivateToPoint(privKey *big.Int) (x, y *big.Int, err error) {
	if privKey.Cmp(big.NewInt(0)) != 1 || privKey.Cmp(sc.N) != -1 {
		return x, y, fmt.Errorf("private key not in curve range")
	}
	x, y = sc.EcMult(privKey, sc.EcGenX, sc.EcGenY)
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


func getSelectorFromName(funcName string) *big.Int {
	kec := Keccak256([]byte(funcName))

	maskedKec := MaskBits(250, 8, kec)
	ret := new(big.Int)
	ret = ret.SetBytes(maskedKec)

	return ret
}

// Keccak256 calculates and returns the Keccak256 hash of the input data.
// (ref: https://github.com/ethereum/go-ethereum/blob/master/crypto/crypto.go)
func Keccak256(data ...[]byte) []byte {
	b := make([]byte, 32)
	d := NewKeccakState()
	for _, b := range data {
		d.Write(b)
	}
	d.Read(b)
	return b
}

// NewKeccakState creates a new KeccakState
// (ref: https://github.com/ethereum/go-ethereum/blob/master/crypto/crypto.go)
func NewKeccakState() KeccakState {
	return sha3.NewLegacyKeccak256().(KeccakState)
}

func MaskBits(mask, wordSize int, slice []byte) (ret []byte) {
	excess := len(slice) * wordSize - mask
	for _, by := range slice {
		if excess > 0  {
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