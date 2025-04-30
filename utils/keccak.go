package utils

import (
	"hash"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
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
// Parameters:
//   - str: The UTF-8 string to convert to a big integer
// Returns:
//   - *big.Int: a pointer to a big.Int representing the converted value
func UTF8StrToBig(str string) *big.Int {
	return internalUtils.UTF8StrToBig(str)
}

// StrToBig generates a *big.Int from a string representation.
//
// Parameters:
//   - str: The string to convert to a *big.Int
// Returns:
//   - *big.Int: a pointer to a big.Int representing the converted value
func StrToBig(str string) *big.Int {
	return internalUtils.StrToBig(str)
}

// StrToBig generates a hexadecimal from a string/number representation.
//
// Parameters:
//   - str: The string to convert to a hexadecimal
// Returns:
//   - hex: a string representing the converted value
func StrToHex(str string) string {
	return internalUtils.StrToHex(str)
}

// HexToShortStr converts a hexadecimal string to a short string (Starknet) representation.
//
// Parameters:
//   - hexStr: the hexadecimal string to convert to a short string
// Returns:
//   - string: a short string
func HexToShortStr(hexStr string) string {
	return internalUtils.HexToShortStr(hexStr)
}

// HexToBN converts a hexadecimal string to a big.Int.
// trim "0x" prefix(if exists)
//
// Parameters:
//   - hexString: the hexadecimal string to be converted
// Returns:
//   - *big.Int: the converted value
func HexToBN(hexString string) *big.Int {
	return internalUtils.HexToBN(hexString)
}

// HexArrToBNArr converts a hexadecimal string array to a *big.Int array.
// Trim "0x" prefix(if exists)
//
// Parameters:
//   - hexArr: the hexadecimal string array to be converted
// Returns:
//   - *big.Int: the converted array
func HexArrToBNArr(hexArr []string) []*big.Int {
	return internalUtils.HexArrToBNArr(hexArr)
}

// HexToBytes converts a hexadecimal string to a byte slice.
// trim "0x" prefix(if exists)
//
// Parameters:
//   - hexString: the hexadecimal string to be converted
// Returns:
//   - []byte: the converted value
//   - error: an error if any
func HexToBytes(hexString string) ([]byte, error) {
	return internalUtils.HexToBytes(hexString)
}

// BytesToBig converts a byte slice to a big.Int.
//
// Parameters:
//   - bytes: the byte slice to be converted
// Returns:
//   - *big.Int: the converted value
func BytesToBig(bytes []byte) *big.Int {
	return internalUtils.BytesToBig(bytes)
}

// BigToHex converts a big integer to its hexadecimal representation.
//
// Parameters:
//   - in: the big integer to be converted
// Returns:
//   - string: the hexadecimal representation
func BigToHex(in *big.Int) string {
	return internalUtils.BigToHex(in)
}

// GetSelectorFromName generates a selector from a given function name.
//
// Parameters:
//   - funcName: the name of the function
// Returns:
//   - *big.Int: the selector
func GetSelectorFromName(funcName string) *big.Int {
	return internalUtils.GetSelectorFromName(funcName)
}

// GetSelectorFromNameFelt returns a *felt.Felt based on the given function name.
//
// Parameters:
//   - funcName: the name of the function
// Returns:
//   - *felt.Felt: the *felt.Felt
func GetSelectorFromNameFelt(funcName string) *felt.Felt {
	return internalUtils.GetSelectorFromNameFelt(funcName)
}

// Keccak256 returns the Keccak-256 hash of the input data.
// (ref: https://github.com/ethereum/go-ethereum/blob/master/crypto/crypto.go)
//
// Parameters:
//   - data: a variadic parameter of type []byte representing the input data
// Returns:
//   - []byte: a 32-byte hash output
func Keccak256(data ...[]byte) []byte {
	return internalUtils.Keccak256(data...)
}

// NewKeccakState returns a new instance of KeccakState.
// (ref: https://github.com/ethereum/go-ethereum/blob/master/crypto/crypto.go)
//
// Parameters:
//
//	none
//
// Returns:
//   - KeccakState: a new instance of KeccakState
func NewKeccakState() KeccakState {
	return internalUtils.NewKeccakState()
}

// MaskBits masks (excess) bits in a slice of bytes based on the given mask and wordSize.
//
// Parameters:
//   - mask: an integer representing the number of bits to mask
//   - wordSize: an integer representing the size of each word in bits
//   - slice: a slice of bytes to mask
// Returns:
//   - ret: a slice of bytes with the masked bits
func MaskBits(mask, wordSize int, slice []byte) (ret []byte) {
	return internalUtils.MaskBits(mask, wordSize, slice)
}

// ComputeFact computes the factorial of a given number.
//
// Parameters:
//   - programHash: a pointer to a big.Int representing the program hash
//   - programOutputs: a slice of pointers to big.Int representing the program outputs
// Returns:
//   - *big.Int: a pointer to a big.Int representing the computed factorial
func ComputeFact(programHash *big.Int, programOutputs []*big.Int) *big.Int {
	return internalUtils.ComputeFact(programHash, programOutputs)
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
//   - fact: The fact string to be split
// Return types:
//   - fact_low: The low part of the fact string in hexadecimal format
//   - fact_high: The high part of the fact string in hexadecimal format
func SplitFactStr(fact string) (fact_low, fact_high string) {
	return internalUtils.SplitFactStr(fact)
}

// FmtKecBytes formats the given big.Int as a byte slice (Keccak hash) with a specified length.
//
// The function appends the bytes of the big.Int to a buffer and returns it.
// If the length of the buffer is less than the specified length, the function pads the buffer with zeros.
//
// Parameters:
//   - in: the big.Int to be formatted
//   - rolen: the length of the buffer
// Returns:
// buf: the formatted buffer
func FmtKecBytes(in *big.Int, rolen int) (buf []byte) {
	return internalUtils.FmtKecBytes(in, rolen)
}

// SNValToBN converts a given string to a *big.Int by checking if the string contains "0x" prefix.
// used in string conversions when interfacing with the APIs
//
// Parameters:
//   - str: a string to be converted to *big.Int
// Returns:
//   - *big.Int: a pointer to a big.Int representing the converted value
func SNValToBN(str string) *big.Int {
	return internalUtils.SNValToBN(str)
}
