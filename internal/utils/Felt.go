package utils

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"regexp"
	"strings"

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

// FeltArrToStringArr converts an array of Felt objects to an array of string objects.
//
// Parameters:
// - f: the array of Felt objects to convert
// Returns:
// - []string: the array of string objects
func FeltArrToStringArr(f []*felt.Felt) []string {
	stringArr := make([]string, len(f))
	for i, felt := range f {
		stringArr[i] = felt.String()
	}
	return stringArr
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
	const SHORT_LENGTH = 31
	exp := fmt.Sprintf(".{1,%d}", SHORT_LENGTH)
	r := regexp.MustCompile(exp)

	arr := r.FindAllString(s, -1)
	if len(arr) == 0 {
		return []*felt.Felt{&felt.Zero, &felt.Zero, &felt.Zero}, nil
	}

	hexarr := []string{}
	var count, size uint64

	for _, val := range arr {
		if len(val) == SHORT_LENGTH {
			count += 1
		} else {
			size = uint64(len(val))
		}
		hexarr = append(hexarr, "0x"+hex.EncodeToString([]byte(val)))
	}

	harr, err := HexArrToFelt(hexarr)
	if err != nil {
		return nil, err
	}

	if size == 0 {
		harr = append(harr, new(felt.Felt).SetUint64(0))
	}

	harr = append(harr, new(felt.Felt).SetUint64(size))
	return append([]*felt.Felt{new(felt.Felt).SetUint64(count)}, harr...), nil
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
	if len(arr) < 3 {
		return "", fmt.Errorf("invalid felt array, require atleast 3 elements in array")
	}

	count := FeltToBigInt(arr[0]).Uint64()
	var index uint64
	var res []string
	for index = 0; index < count; index++ {
		f := arr[1+index]
		s, err := feltToString(f)
		if err != nil {
			return "", err
		}
		res = append(res, s)
	}

	pendingWordLength := arr[len(arr)-1]
	if pendingWordLength.IsZero() {
		return strings.Join(res, ""), nil
	}

	pendingWordFelt := arr[1+index]
	s, err := feltToString(pendingWordFelt)
	if err != nil {
		return "", fmt.Errorf("invalid pending word")
	}

	res = append(res, s)
	return strings.Join(res, ""), nil
}

func feltToString(f *felt.Felt) (string, error) {
	b, err := hex.DecodeString(f.String()[2:])
	if err != nil {
		return "", fmt.Errorf("unable to decode to string")
	}
	return string(b), nil
}

// BigIntArrToFeltArr converts an array of big.Int objects to an array of Felt objects.
//
// Parameters:
// - bigArr: the array of big.Int objects to convert
// Returns:
// - []*felt.Felt: the array of Felt objects
func BigIntArrToFeltArr(bigArr []*big.Int) []*felt.Felt {
	var feltArr []*felt.Felt
	for _, big := range bigArr {
		feltArr = append(feltArr, BigIntToFelt(big))
	}
	return feltArr
}
