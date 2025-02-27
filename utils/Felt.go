package utils

import (
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/internal/utils"
)

// Uint64ToFelt generates a new *felt.Felt from a given uint64 number.
//
// Parameters:
// - num: the uint64 number to convert to a *felt.Felt
// Returns:
// - *felt.Felt: a *felt.Felt
func Uint64ToFelt(num uint64) *felt.Felt {
	return utils.Uint64ToFelt(num)
}

// HexToFelt converts a hexadecimal string to a *felt.Felt object.
//
// Parameters:
// - hex: the input hexadecimal string to be converted.
// Returns:
// - *felt.Felt: a *felt.Felt object
// - error: if conversion fails
func HexToFelt(hex string) (*felt.Felt, error) {
	return utils.HexToFelt(hex)
}

// HexArrToFelt converts an array of hexadecimal strings to an array of felt objects.
//
// The function iterates over each element in the hexArr array and calls the HexToFelt function to convert each hexadecimal value to a felt object.
// If any error occurs during the conversion, the function will return nil and the corresponding error.
// Otherwise, it appends the converted felt object to the feltArr array.
// Finally, the function returns the feltArr array containing all the converted felt objects.
//
// Parameters:
// - hexArr: an array of strings representing hexadecimal values
// Returns:
// - []*felt.Felt: an array of *felt.Felt objects, or nil if there was
// - error: an error if any
func HexArrToFelt(hexArr []string) ([]*felt.Felt, error) {
	return utils.HexArrToFelt(hexArr)
}

// FeltToBigInt converts a Felt value to a *big.Int.
//
// Parameters:
// - f: the Felt value to convert
// Returns:
// - *big.Int: the converted value
func FeltToBigInt(f *felt.Felt) *big.Int {
	return utils.FeltToBigInt(f)
}

// BigIntToFelt converts a big integer to a felt.Felt.
//
// Parameters:
// - big: the big integer to convert
// Returns:
// - *felt.Felt: the converted value
func BigIntToFelt(big *big.Int) *felt.Felt {
	return utils.BigIntToFelt(big)
}

// FeltArrToBigIntArr converts an array of Felt objects to an array of big.Int objects.
//
// Parameters:
// - f: the array of Felt objects to convert
// Returns:
// - []*big.Int: the array of big.Int objects
func FeltArrToBigIntArr(f []*felt.Felt) []*big.Int {
	return utils.FeltArrToBigIntArr(f)
}

// FeltArrToStringArr converts an array of Felt objects to an array of string objects.
//
// Parameters:
// - f: the array of Felt objects to convert
// Returns:
// - []string: the array of string objects
func FeltArrToStringArr(f []*felt.Felt) []string {
	return utils.FeltArrToStringArr(f)
}

// StringToByteArrFelt converts string to array of Felt objects.
// The returned array of felts will be of the format
//
// [number of felts with 31 characters in length, 31 byte felts..., pending word with max size of 30 bytes, pending words bytes size]
//
// For further explanation, refer the [article]
//
// Parameters:
//
// - s: string/bytearray to convert
//
// Returns:
//
// - []*felt.Felt: the array of felt.Felt objects
//
// - error: an error, if any
//
// [article]: https://docs.starknet.io/architecture-and-concepts/smart-contracts/serialization-of-cairo-types/#serialization_of_byte_arrays
func StringToByteArrFelt(s string) ([]*felt.Felt, error) {
	return utils.StringToByteArrFelt(s)
}

// ByteArrFeltToString converts array of Felts to string.
// The input array of felts will be of the format
//
// [number of felts with 31 characters in length, 31 byte felts..., pending word with max size of 30 bytes, pending words bytes size]
//
// For further explanation, refer the [article]
//
// Parameters:
//
// - []*felt.Felt: the array of felt.Felt objects
//
// Returns:
//
// - s: string/bytearray
//
// - error: an error, if any
//
// [article]: https://docs.starknet.io/architecture-and-concepts/smart-contracts/serialization-of-cairo-types/#serialization_of_byte_arrays
func ByteArrFeltToString(arr []*felt.Felt) (string, error) {
	return utils.ByteArrFeltToString(arr)
}

// BigIntArrToFeltArr converts an array of big.Int objects to an array of Felt objects.
//
// Parameters:
// - bigArr: the array of big.Int objects to convert
// Returns:
// - []*felt.Felt: the array of Felt objects
func BigIntArrToFeltArr(bigArr []*big.Int) []*felt.Felt {
	return utils.BigIntArrToFeltArr(bigArr)
}
