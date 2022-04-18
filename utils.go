package caigo

import (
	"crypto/hmac"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"hash"
	"math/big"
	"strings"

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

// convert hex string to StarkNet 'short string'
func HexToShortStr(hexStr string) string {
	numStr := strings.Replace(hexStr, "0x", "", -1)
	hb, _ := new(big.Int).SetString(numStr, 16)

	return string(hb.Bytes())
}

// trim "0x" prefix(if exists) and converts hexidecimal string to big int
func HexToBN(hexString string) *big.Int {
	numStr := strings.Replace(hexString, "0x", "", -1)

	n, _ := new(big.Int).SetString(numStr, 16)
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

func BytesToBig(bytes []byte) *big.Int {
	return new(big.Int).SetBytes(bytes)
}

// convert big int to hexidecimal string
func BigToHex(in *big.Int) string {
	return fmt.Sprintf("0x%x", in)
}

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

func GetSelectorFromName(funcName string) *big.Int {
	kec := Keccak256([]byte(funcName))

	maskedKec := MaskBits(250, 8, kec)

	return new(big.Int).SetBytes(maskedKec)
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

// compute the keccack fact given the program hash and outputs
func ComputeFact(programHash *big.Int, programOutputs []*big.Int) *big.Int {
	var progOutBuf []byte
	for _, programOutput := range programOutputs {
		inBuf := FmtKecBytes(programOutput, 32)
		progOutBuf = append(progOutBuf[:], inBuf...)
	}

	kecBuf := FmtKecBytes(programHash, 32)
	kecBuf = append(kecBuf[:], Keccak256(progOutBuf)...)

	return new(big.Int).SetBytes(Keccak256(kecBuf))
}

// split a fact into two felts
func SplitFactStr(fact string) (fact_low, fact_high string) {
	factBN := HexToBN(fact)
	factBytes := factBN.Bytes()
	low := BytesToBig(factBytes[16:])
	high := BytesToBig(factBytes[:16])
	return BigToHex(low), BigToHex(high)
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

// used in string conversions when interfacing with the APIs
func SNValToBN(str string) *big.Int {
	if strings.Contains(str, "0x") {
		return HexToBN(str)
	} else {
		return StrToBig(str)
	}
}
