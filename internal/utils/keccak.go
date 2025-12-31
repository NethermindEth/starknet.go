package utils

import (
	"encoding/hex"
	"errors"
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
// Parameters:
//   - str: The UTF-8 string to convert to a big integer
//
// Returns:
//   - *big.Int: a pointer to a big.Int representing the converted value
func UTF8StrToBig(str string) *big.Int {
	hexStr := hex.EncodeToString([]byte(str))
	b, _ := new(big.Int).SetString(hexStr, 16) //nolint:mnd //set as hex

	return b
}

// StrToBig generates a *big.Int from a string representation.
//
// Parameters:
//   - str: The string to convert to a *big.Int
//
// Returns:
//   - *big.Int: a pointer to a big.Int representing the converted value
func StrToBig(str string) *big.Int {
	b, _ := new(big.Int).SetString(str, 10) //nolint:mnd //set as decimal

	return b
}

// StrToBig generates a hexadecimal from a string/number representation.
//
// Parameters:
//   - str: The string to convert to a hexadecimal
//
// Returns:
//   - hex: a string representing the converted value
func StrToHex(str string) string {
	if strings.HasPrefix(str, "0x") {
		return str
	}

	if bigNum, ok := new(big.Int).SetString(str, 0); ok {
		return "0x" + bigNum.Text(16) //nolint:mnd //turn to hex
	}

	return "0x" + fmt.Sprintf("%x", str)
}

// HexToShortStr converts a hexadecimal string to a short string (Starknet) representation.
//
// Parameters:
//   - hexStr: the hexadecimal string to convert to a short string
//
// Returns:
//   - string: a short string
func HexToShortStr(hexStr string) string {
	numStr := strings.ReplaceAll(hexStr, "0x", "")
	hb, _ := new(big.Int).SetString(numStr, 16) //nolint:mnd //set as hex

	return string(hb.Bytes())
}

// HexToBN converts a hexadecimal string to a big.Int.
// trim "0x" prefix(if exists)
//
// Parameters:
//   - hexString: the hexadecimal string to be converted
//
// Returns:
//   - *big.Int: the converted value
func HexToBN(hexString string) *big.Int {
	numStr := strings.ReplaceAll(hexString, "0x", "")

	// TODO: maybe make this func return this ignored bool value
	n, _ := new(big.Int).SetString(numStr, 16) //nolint:mnd //set as hex

	return n
}

// HexArrToBNArr converts a hexadecimal string array to a *big.Int array.
// Trim "0x" prefix(if exists)
//
// Parameters:
//   - hexArr: the hexadecimal string array to be converted
//
// Returns:
//   - *big.Int: the converted array
func HexArrToBNArr(hexArr []string) []*big.Int {
	bigNumArr := make([]*big.Int, len(hexArr))
	for i, hexStr := range hexArr {
		bigNumArr[i] = HexToBN(hexStr)
	}

	return bigNumArr
}

// HexToBytes converts a hexadecimal string to a byte slice.
// trim "0x" prefix(if exists)
//
// Parameters:
//   - hexString: the hexadecimal string to be converted
//
// Returns:
//   - []byte: the converted value
//   - error: an error if any
func HexToBytes(hexString string) ([]byte, error) {
	numStr := strings.ReplaceAll(hexString, "0x", "")
	if (len(numStr) % 2) != 0 {
		numStr = fmt.Sprintf("%s%s", "0", numStr)
	}

	return hex.DecodeString(numStr)
}

// BytesToBig converts a byte slice to a big.Int.
//
// Parameters:
//   - bytesVal: the byte slice to be converted
//
// Returns:
//   - *big.Int: the converted value
func BytesToBig(bytesVal []byte) *big.Int {
	return new(big.Int).SetBytes(bytesVal)
}

// BigToHex converts a big integer to its hexadecimal representation.
//
// Parameters:
//   - in: the big integer to be converted
//
// Returns:
//   - string: the hexadecimal representation
func BigToHex(in *big.Int) string {
	return fmt.Sprintf("0x%x", in)
}

// GetSelectorFromName generates a selector from a given function name.
// Ref: https://www.starknet.io/cairo-book/ch101-01-00-contract-storage.html#addresses-of-storage-variables
//
// Parameters:
//   - funcName: the name of the function
//
// Returns:
//   - *big.Int: the selector
//
//nolint:lll // The link would be unclickable if we break the line.
func GetSelectorFromName(funcName string) *big.Int {
	kec := Keccak256([]byte(funcName))

	maskedKec := MaskBits(250, 8, kec) //nolint:mnd // Values taken from the starknet documentation.

	return new(big.Int).SetBytes(maskedKec)
}

// GetSelectorFromNameFelt does the same as GetSelectorFromName, but returns
// the result as a *felt.Felt.
//
// Parameters:
//   - funcName: the name of the function
//
// Returns:
//   - *felt.Felt: the *felt.Felt
func GetSelectorFromNameFelt(funcName string) *felt.Felt {
	kec := Keccak256([]byte(funcName))

	maskedKec := MaskBits(250, 8, kec) //nolint:mnd // Values taken from the starknet documentation.

	return new(felt.Felt).SetBytes(maskedKec)
}

// Keccak256 returns the Keccak-256 hash of the input data.
// (ref: https://github.com/ethereum/go-ethereum/blob/master/crypto/crypto.go)
//
// Parameters:
//   - data: a variadic parameter of type []byte representing the input data
//
// Returns:
//   - []byte: a 32-byte hash output
func Keccak256(data ...[]byte) []byte {
	b := make([]byte, 32) //nolint:mnd // 32 bytes = 256 bits, necessary for keccak256
	d := NewKeccakState()
	for _, b := range data {
		_, err := d.Write(b)
		if err != nil {
			panic(err)
		}
	}

	if _, err := d.Read(b); err != nil {
		panic(err)
	}

	return b
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
	return sha3.NewLegacyKeccak256().(KeccakState)
}

// MaskBits masks (excess) bits in a slice of bytes based on the given mask and wordSize.
//
// Parameters:
//   - mask: an integer representing the number of bits to mask
//   - wordSize: an integer representing the size of each word in bits
//   - slice: a slice of bytes to mask
//
// Returns:
//   - ret: a slice of bytes with the masked bits
func MaskBits(mask, wordSize int, slice []byte) (ret []byte) {
	excess := len(slice)*wordSize - mask
	for _, by := range slice {
		if excess > 0 {
			if excess > wordSize {
				excess -= wordSize

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

// ComputeFact computes a cryptographic hash from a program hash and its outputs.
// It combines the program hash with a hash of all program outputs using Keccak256.
//
// Parameters:
//   - programHash: a pointer to a big.Int representing the program hash
//   - programOutputs: a slice of pointers to big.Int representing the program outputs
//
// Returns:
//   - *big.Int: a pointer to a big.Int representing the computed factorial
//
//nolint:mnd // 32 bytes = 256 bits, necessary for keccak256
func ComputeFact(programHash *big.Int, programOutputs []*big.Int) *big.Int {
	var progOutBuf []byte
	for _, programOutput := range programOutputs {
		inBuf := FmtKecBytes(
			programOutput,
			32,
		)
		progOutBuf = append(progOutBuf, inBuf...)
	}

	kecBuf := FmtKecBytes(
		programHash,
		32,
	)
	kecBuf = append(kecBuf, Keccak256(progOutBuf)...)

	return new(big.Int).SetBytes(Keccak256(kecBuf))
}

// SplitFactStr splits a given fact, with maximum 256 bits size, into two
// parts (felts): fact_low and fact_high.
//
// The function takes a fact string as input and converts it to a big number
// using the HexToBN function. It then converts the big number to bytes using
// the Bytes method. If the length of the bytes is less than 32, it pads the
// bytes with zeros using the bytes.Repeat method.
// The padded bytes are then appended to the original bytes.
// The function then extracts the low part of the bytes by taking the last 16
// bytes and converts it to a big number using the BytesToBig function.
// It also extracts the high part of the bytes by taking the first 16 bytes and
// converts it to a big number using the BytesToBig function.
// Finally, it converts the low and high big numbers to hexadecimal strings using
// the BigToHex function and returns them.
//
// Parameters:
//   - fact: The fact string to be split
//
// Return types:
//   - fact_low: The low part of the fact string in hexadecimal format
//   - fact_high: The high part of the fact string in hexadecimal format
//   - err: An error if any
//
//nolint:mnd // There's a comment explaining each magic number.
func SplitFactStr(fact string) (factLow, factHigh string, err error) {
	numStr := strings.ReplaceAll(fact, "0x", "")
	factBN, ok := new(big.Int).SetString(numStr, 16) // hex base
	if !ok {
		return "", "", errors.New("failed to convert fact string to big.Int")
	}
	if factBN.BitLen() > 256 { // max 256 bits
		return "", "", errors.New("fact string is too large")
	}
	factBytes := factBN.Bytes()
	lpadfactBytes := make([]byte, 32-len(factBytes)) // left pad with zeros
	factBytes = append(lpadfactBytes, factBytes...)  //nolint:makezero // we want the zeros
	high := BytesToBig(factBytes[:16])               // first 16 bytes
	low := BytesToBig(factBytes[16:])                // last 16 bytes

	return BigToHex(low), BigToHex(high), nil
}

// FmtKecBytes formats the given big.Int as a byte slice (Keccak hash) with
// a specified length.
//
// The function appends the bytes of the big.Int to a buffer and returns it.
// If the length of the buffer is less than the specified length, the function
// pads the buffer with zeros.
//
// Parameters:
//   - in: the big.Int to be formatted
//   - rolen: the length of the buffer
//
// Returns:
//   - buf: the formatted buffer
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
// Parameters:
//   - str: a string to be converted to *big.Int
//
// Returns:
//   - *big.Int: a pointer to a big.Int representing the converted value
func SNValToBN(str string) *big.Int {
	if strings.Contains(str, "0x") {
		return HexToBN(str)
	} else {
		return StrToBig(str)
	}
}
