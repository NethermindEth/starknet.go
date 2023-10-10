package utils

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

// UTF8StrToBig converts a UTF-8 string to a big integer.
//
// str: The UTF-8 string to convert.
// Returns a pointer to a big.Int representing the converted value.
func UTF8StrToBig(str string) *big.Int {
	hexStr := hex.EncodeToString([]byte(str))
	b, _ := new(big.Int).SetString(hexStr, 16)

	return b
}

// StrToBig generates a *big.Int from a string representation.
//
// It takes a string parameter and returns a *big.Int.
func StrToBig(str string) *big.Int {
	b, _ := new(big.Int).SetString(str, 10)

	return b
}

// HexToShortStr converts a hexadecimal string to a short string (Starknet) representation.
//
// It takes a hexadecimal string as input and returns a short string representation.
// The input hexadecimal string must start with "0x".
// The return value is a string.
func HexToShortStr(hexStr string) string {
	numStr := strings.Replace(hexStr, "0x", "", -1)
	hb, _ := new(big.Int).SetString(numStr, 16)

	return string(hb.Bytes())
}

// HexToBN converts a hexadecimal string to a big.Int.
// trim "0x" prefix(if exists)
//
// It takes a hexString parameter, which is the hexadecimal string to be converted.
// It returns a *big.Int, which is the converted value.
func HexToBN(hexString string) *big.Int {
	numStr := strings.Replace(hexString, "0x", "", -1)

	n, _ := new(big.Int).SetString(numStr, 16)
	return n
}

// HexToBytes converts a hexadecimal string to a byte slice.
// trim "0x" prefix(if exists) 
//
// It takes a hexString parameter which is the string representation of the hexadecimal number.
// The function returns a byte slice and an error.
func HexToBytes(hexString string) ([]byte, error) {
	numStr := strings.Replace(hexString, "0x", "", -1)
	if (len(numStr) % 2) != 0 {
		numStr = fmt.Sprintf("%s%s", "0", numStr)
	}

	return hex.DecodeString(numStr)
}

// BytesToBig converts a byte slice to a big.Int.
//
// It takes a byte slice as input and returns a pointer to a big.Int.
func BytesToBig(bytes []byte) *big.Int {
	return new(big.Int).SetBytes(bytes)
}

// BigToHex converts a big integer to its hexadecimal representation.
//
// It takes a pointer to a big.Int as input.
// It returns a string containing the hexadecimal representation of the input big integer.
func BigToHex(in *big.Int) string {
	return fmt.Sprintf("0x%x", in)
}

// GetSelectorFromName generates a selector from a given function name.
//
// It takes a string parameter `funcName` which represents the name of the function.
// It returns a pointer to a `big.Int` type.
//
// TODO: this is used by the signer. Should it return a felt?
func GetSelectorFromName(funcName string) *big.Int {
	kec := Keccak256([]byte(funcName))

	maskedKec := MaskBits(250, 8, kec)

	return new(big.Int).SetBytes(maskedKec)
}

// GetSelectorFromNameFelt returns a *felt.Felt based on the given function name.
//
// It takes a string parameter, funcName, which represents the name of the function.
//
// The function returns a *felt.Felt, which is calculated based on the funcName parameter.
func GetSelectorFromNameFelt(funcName string) *felt.Felt {
	kec := Keccak256([]byte(funcName))

	maskedKec := MaskBits(250, 8, kec)

	return new(felt.Felt).SetBytes(maskedKec)
}

// Keccak256 returns the Keccak-256 hash of the input data.
// (ref: https://github.com/ethereum/go-ethereum/blob/master/crypto/crypto.go)
//
// It accepts a variadic parameter of type []byte, representing the input data.
// It returns a []byte, which is the 32-byte hash output.
func Keccak256(data ...[]byte) []byte {
	b := make([]byte, 32)
	d := NewKeccakState()
	for _, b := range data {
		d.Write(b)
	}
	d.Read(b)
	return b
}

// NewKeccakState returns a new instance of KeccakState.
// (ref: https://github.com/ethereum/go-ethereum/blob/master/crypto/crypto.go)
//
// It does not take any parameters.
// It returns a value of type KeccakState.
func NewKeccakState() KeccakState {
	return sha3.NewLegacyKeccak256().(KeccakState)
}

// MaskBits masks (excess) bits in a slice of bytes based on the given mask and wordSize.
//
// Parameters:
// - mask: an integer representing the number of bits to mask.
// - wordSize: an integer representing the size of each word in bits.
// - slice: a slice of bytes to mask.
//
// Returns:
// - ret: a slice of bytes with the masked bits.
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

// ComputeFact computes the factorial of a given number.
//
// It takes two parameters:
// - programHash: a pointer to a big.Int representing the program hash.
// - programOutputs: a slice of pointers to big.Int representing the program outputs.
//
// It returns a pointer to a big.Int representing the computed factorial.
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

// SplitFactStr splits a given fact string into two parts (felts): fact_low and fact_high.
//
// The function takes a fact string as input and converts it to a big number using the HexToBN function.
// It then converts the big number to bytes using the Bytes method.
// If the length of the bytes is less than 32, it pads the bytes with zeros using the bytes.Repeat method.
// The padded bytes are then appended to the original bytes.
// The function then extracts the low part of the bytes by taking the last 16 bytes and converts it to a big number using the BytesToBig function.
// It also extracts the high part of the bytes by taking the first 16 bytes and converts it to a big number using the BytesToBig function.
// Finally, it converts the low and high big numbers to hexadecimal strings using the BigToHex function and returns them.
//
// Parameters:
// - fact: The fact string to be split.
//
// Return types:
// - fact_low: The low part of the fact string in hexadecimal format.
// - fact_high: The high part of the fact string in hexadecimal format.
func SplitFactStr(fact string) (fact_low, fact_high string) {
	factBN := HexToBN(fact)
	factBytes := factBN.Bytes()
	lpadfactBytes := bytes.Repeat([]byte{0x00}, 32-len(factBytes))
	factBytes = append(lpadfactBytes, factBytes...)
	low := BytesToBig(factBytes[16:])
	high := BytesToBig(factBytes[:16])
	return BigToHex(low), BigToHex(high)
}

// FmtKecBytes formats the given big.Int as a byte slice (Keccak hash) with a specified length.
//
// It takes in a pointer to a big.Int and an integer representing the desired length of the resulting byte slice.
// The function appends the bytes of the big.Int to a buffer and returns it.
// If the length of the buffer is less than the specified length, the function pads the buffer with zeros.
//
// The function returns a byte slice.
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

// SNValToBN converts a given string to a *big.Int by checking if the string contains "0x" prefix.
// used in string conversions when interfacing with the APIs
//
// It takes a single parameter:
// - str: a string to be converted to *big.Int
//
// It returns a *big.Int.
func SNValToBN(str string) *big.Int {
	if strings.Contains(str, "0x") {
		return HexToBN(str)
	} else {
		return StrToBig(str)
	}
}
