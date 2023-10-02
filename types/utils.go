package types

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"hash"
	"math/big"
	"strings"

	"github.com/NethermindEth/juno/core/felt"
	"golang.org/x/crypto/sha3"
)

// KeccakState wraps sha3.state. In addition to the usual hash methods, it also supports
// Read to get a variable amount of data from the hash state. Read is faster than Sum
// because it doesn't copy the internal state, but also modifies the internal state.
type KeccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

// UTF8StrToBig converts a UTF-8 string to a big.Int.
//
// It takes a string as a parameter and returns a pointer to a big.Int.
func UTF8StrToBig(str string) *big.Int {
	hexStr := hex.EncodeToString([]byte(str))
	b, _ := new(big.Int).SetString(hexStr, 16)

	return b
}

// StrToBig converts a string to a *big.Int.
//
// It takes a string as input and returns a *big.Int.
func StrToBig(str string) *big.Int {
	b, _ := new(big.Int).SetString(str, 10)

	return b
}

// HexToShortStr converts a hexadecimal string to a Starknet short string.
//
// hexStr: The hexadecimal string to be converted.
// Returns: The converted short string.
func HexToShortStr(hexStr string) string {
	numStr := strings.Replace(hexStr, "0x", "", -1)
	hb, _ := new(big.Int).SetString(numStr, 16)

	return string(hb.Bytes())
}

// HexToBN converts a hexadecimal string to a big.Int, trims "0x" prefix(if exists).
//
// It takes a hexString as input parameter and returns a pointer to a big.Int.
func HexToBN(hexString string) *big.Int {
	numStr := strings.Replace(hexString, "0x", "", -1)

	n, _ := new(big.Int).SetString(numStr, 16)
	return n
}

// HexToBytes converts a hexadecimal string to a byte array, trims "0x" prefix(if exists).
//
// It takes a hexString parameter, which is the hexadecimal string to be converted.
// The function returns a byte array and an error.
func HexToBytes(hexString string) ([]byte, error) {
	numStr := strings.Replace(hexString, "0x", "", -1)
	if (len(numStr) % 2) != 0 {
		numStr = fmt.Sprintf("%s%s", "0", numStr)
	}

	return hex.DecodeString(numStr)
}

// BytesToBig converts a byte slice to a *big.Int.
//
// It takes a parameter bytes, which is a byte slice.
// It returns a *big.Int.
func BytesToBig(bytes []byte) *big.Int {
	return new(big.Int).SetBytes(bytes)
}

// BigToHex returns a hexadecimal representation of the given big integer.
//
// It takes a pointer to a big.Int as its parameter and returns a string.
func BigToHex(in *big.Int) string {
	return fmt.Sprintf("0x%x", in)
}

// GetSelectorFromName returns the selector from the given function name.
// TODO: this is used by the signer. Should it return a Felt?
//
// The function takes a string parameter `funcName` which represents the name of the function.
// It returns a pointer to a `big.Int` type.
func GetSelectorFromName(funcName string) *big.Int {
	kec := Keccak256([]byte(funcName))

	maskedKec := MaskBits(250, 8, kec)

	return new(big.Int).SetBytes(maskedKec)
}

// GetSelectorFromNameFelt generates a selector from a given function name using the Felt package.
//
// Parameters:
// - funcName: the name of the function.
//
// Return:
// - *felt.Felt: the generated selector.
func GetSelectorFromNameFelt(funcName string) *felt.Felt {
	kec := Keccak256([]byte(funcName))

	maskedKec := MaskBits(250, 8, kec)

	return new(felt.Felt).SetBytes(maskedKec)
}

// Keccak256 calculates the Keccak-256 hash of the given data.
// (ref: https://github.com/ethereum/go-ethereum/blob/master/crypto/crypto.go)
//
// It takes one or more byte slices as input and returns a 32-byte hash.
func Keccak256(data ...[]byte) []byte {
	b := make([]byte, 32)
	d := NewKeccakState()
	for _, b := range data {
		d.Write(b)
	}
	d.Read(b)
	return b
}

// NewKeccakState creates a new instance of the KeccakState struct.
// (ref: https://github.com/ethereum/go-ethereum/blob/master/crypto/crypto.go)
//
// It does not take any parameters.
// It returns a KeccakState object.
func NewKeccakState() KeccakState {
	return sha3.NewLegacyKeccak256().(KeccakState)
}

// MaskBits masks the excess bits in a byte slice based on the given mask and word size.
//
// Parameters:
// - mask: an integer representing the mask value.
// - wordSize: an integer representing the size of each word.
// - slice: a byte slice containing the data to mask.
//
// Returns:
// - ret: a byte slice with the masked bits.
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

// ComputeFact calculates the factorial of a given programHash and programOutputs.
//
// The function takes in a programHash, which is a pointer to a big.Int, and
// programOutputs, which is a slice of pointers to big.Int. It calculates the
// factorial by iterating over the programOutputs and appending the formatted
// KecBytes to the progOutBuf. It then appends the formatted KecBytes of the
// programHash to the kecBuf. Finally, it returns a new big.Int set with the
// Keccak256 of the kecBuf.
//
// The function returns a pointer to a big.Int.
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

// SplitFactStr splits the given fact string into two Felts: fact_low and fact_high.
//
// Parameters:
// - fact: The fact string to be split.
//
// Return types:
// - fact_low: The lower part of the split fact string.
// - fact_high: The higher part of the split fact string.
func SplitFactStr(fact string) (fact_low, fact_high string) {
	factBN := HexToBN(fact)
	factBytes := factBN.Bytes()
	lpadfactBytes := bytes.Repeat([]byte{0x00}, 32-len(factBytes))
	factBytes = append(lpadfactBytes, factBytes...)
	low := BytesToBig(factBytes[16:])
	high := BytesToBig(factBytes[:16])
	return BigToHex(low), BigToHex(high)
}

// FmtKecBytes formats the given big.Int as a byte slice with a specified length.
//
// It takes in a pointer to a big.Int 'in' and an integer 'rolen' representing
// the desired length of the resulting byte slice.
//
// It returns a byte slice 'buf' containing the formatted bytes.
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

// SNValToBN converts a string representation of a number to a big integer.
// used in string conversions when interfacing with the APIs
//
// It takes a string parameter `str` which represents the number to be converted.
// The function returns a pointer to a big.Int.
func SNValToBN(str string) *big.Int {
	if strings.Contains(str, "0x") {
		return HexToBN(str)
	} else {
		return StrToBig(str)
	}
}
