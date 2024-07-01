package utils

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"regexp"

	"github.com/NethermindEth/juno/core/felt"
)

// Uint64ToFelt generates a new *felt.Felt from a given uint64 number.
//
// Parameters:
// - num: the uint64 number to convert to a *felt.Felt
// Returns:
// - *felt.Felt: a *felt.Felt
func Uint64ToFelt(num uint64) *felt.Felt {
	return new(felt.Felt).SetUint64(num)
}

// HexToFelt converts a hexadecimal string to a *felt.Felt object.
//
// Parameters:
// - hex: the input hexadecimal string to be converted.
// Returns:
// - *felt.Felt: a *felt.Felt object
// - error: if conversion fails
func HexToFelt(hex string) (*felt.Felt, error) {
	return new(felt.Felt).SetString(hex)
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

	feltArr := make([]*felt.Felt, len(hexArr))
	for i, e := range hexArr {
		felt, err := HexToFelt(e)
		if err != nil {
			return nil, err
		}
		feltArr[i] = felt
	}
	return feltArr, nil

}

// FeltToBigInt converts a Felt value to a *big.Int.
//
// Parameters:
// - f: the Felt value to convert
// Returns:
// - *big.Int: the converted value
func FeltToBigInt(f *felt.Felt) *big.Int {
	tmp := f.Bytes()
	return new(big.Int).SetBytes(tmp[:])
}

// BigIntToFelt converts a big integer to a felt.Felt.
//
// Parameters:
// - big: the big integer to convert
// Returns:
// - *felt.Felt: the converted value
func BigIntToFelt(big *big.Int) *felt.Felt {
	return new(felt.Felt).SetBytes(big.Bytes())
}

// FeltArrToBigIntArr converts an array of Felt objects to an array of big.Int objects.
//
// Parameters:
// - f: the array of Felt objects to convert
// Returns:
// - []*big.Int: the array of big.Int objects
func FeltArrToBigIntArr(f []*felt.Felt) []*big.Int {
	var bigArr []*big.Int
	for _, felt := range f {
		bigArr = append(bigArr, FeltToBigInt(felt))
	}
	return bigArr
}

// StringToByteArrFelt converts string to array of Felt objects.
//
// Parameters:
// - s: string/bytearray to convert
// Returns:
// - []*felt.Felt: the array of felt.Felt objects
// - error: an error, if any
func StringToByteArrFelt(s string) ([]*felt.Felt, error) {
	const SHORT_LENGTH = 31
	exp := fmt.Sprintf(".{1,%d}", SHORT_LENGTH)
	r := regexp.MustCompile(exp)

	arr := r.FindAllString(s, -1)
	if len(arr) == 0 {
		return []*felt.Felt{}, fmt.Errorf("invalid string no matches found, s: %s", s)
	}

	hexarr := []string{}
	var (
		count uint64
		size  uint64
	)

	for _, val := range arr {
		if len(val) == SHORT_LENGTH {
			count += 1
		}
		size = uint64(len(val))
		hexarr = append(hexarr, hex.EncodeToString([]byte(val)))
	}

	harr, err := HexArrToFelt(hexarr)
	if err != nil {
		return nil, err
	}

	harr = append(harr, new(felt.Felt).SetUint64(size))
	return append([]*felt.Felt{new(felt.Felt).SetUint64(count)}, harr...), nil
}
