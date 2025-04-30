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
//   - num: the uint64 number to convert to a *felt.Felt
// Returns:
//   - *felt.Felt: a *felt.Felt
func Uint64ToFelt(num uint64) *felt.Felt {
	return new(felt.Felt).SetUint64(num)
}

// HexToFelt converts a hexadecimal string to a *felt.Felt object.
//
// Parameters:
//   - hex: the input hexadecimal string to be converted.
// Returns:
//   - *felt.Felt: a *felt.Felt object
//   - error: if conversion fails
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
//   - hexArr: an array of strings representing hexadecimal values
// Returns:
//   - []*felt.Felt: an array of *felt.Felt objects, or nil if there was
//   - error: an error if any
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
//   - f: the Felt value to convert
// Returns:
//   - *big.Int: the converted value
func FeltToBigInt(f *felt.Felt) *big.Int {
	tmp := f.Bytes()
	return new(big.Int).SetBytes(tmp[:])
}

// BigIntToFelt converts a big integer to a felt.Felt.
//
// Parameters:
//   - big: the big integer to convert
// Returns:
//   - *felt.Felt: the converted value
func BigIntToFelt(big *big.Int) *felt.Felt {
	return new(felt.Felt).SetBytes(big.Bytes())
}

// FeltArrToBigIntArr converts an array of Felt objects to an array of big.Int objects.
//
// Parameters:
//   - f: the array of Felt objects to convert
// Returns:
//   - []*big.Int: the array of big.Int objects
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
//   - f: the array of Felt objects to convert
// Returns:
//   - []string: the array of string objects
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
//   - s: string/bytearray to convert
//
// Returns:
//
//   - []*felt.Felt: the array of felt.Felt objects
//
//   - error: an error, if any
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
//   - []*felt.Felt: the array of felt.Felt objects
//
// Returns:
//
//   - s: string/bytearray
//
//   - error: an error, if any
//
// [article]: https://docs.starknet.io/architecture-and-concepts/smart-contracts/serialization-of-cairo-types/#serialization_of_byte_arrays
func ByteArrFeltToString(arr []*felt.Felt) (string, error) {
	const SHORT_LENGTH = 31

	if len(arr) < 3 {
		return "", fmt.Errorf("invalid felt array, require atleast 3 elements in array")
	}

	count := arr[0].Uint64()
	pendingWordLength := arr[len(arr)-1].Uint64()

	// pending word length is in the range [0, SHORT_LENGTH-1]
	if pendingWordLength > SHORT_LENGTH-1 {
		return "", fmt.Errorf("invalid felt array, invalid pending word length")
	}

	if uint64(len(arr)) != 3+count {
		return "", fmt.Errorf("invalid felt array, invalid length got %d expected %d", len(arr), 3+count)
	}

	var res []string
	if pendingWordLength == 0 {
		res = make([]string, count)
	} else {
		res = make([]string, count+1)
	}

	for index := range count {
		res[index] = bytesFeltToString(arr[1+index], SHORT_LENGTH)
	}

	if pendingWordLength != 0 {
		res[count] = bytesFeltToString(arr[1+count], int(pendingWordLength))
	}

	return strings.Join(res, ""), nil
}

func bytesFeltToString(f *felt.Felt, length int) string {
	b := f.Bytes()
	return string(b[len(b)-length:])
}

// BigIntArrToFeltArr converts an array of big.Int objects to an array of Felt objects.
//
// Parameters:
//   - bigArr: the array of big.Int objects to convert
// Returns:
//   - []*felt.Felt: the array of Felt objects
func BigIntArrToFeltArr(bigArr []*big.Int) []*felt.Felt {
	var feltArr []*felt.Felt
	for _, big := range bigArr {
		feltArr = append(feltArr, BigIntToFelt(big))
	}
	return feltArr
}

// HexToU256Felt converts a hexadecimal string to a Cairo u256 representation.
// The Cairo u256 is represented as two felt.Felt values:
//   - The first felt.Felt contains the 128 least significant bits (low part)
//   - The second felt.Felt contains the 128 most significant bits (high part)
//
// Parameters:
//   - hexStr: the hexadecimal string to convert to a Cairo u256
// Returns:
//   - []*felt.Felt: a slice containing two felt.Felt values [low, high]
//   - error: if conversion fails
func HexToU256Felt(hexStr string) ([]*felt.Felt, error) {
	// Ensure the hex string has the 0x prefix
	if !strings.HasPrefix(hexStr, "0x") && !strings.HasPrefix(hexStr, "0X") {
		hexStr = "0x" + hexStr
	}

	// Parse the hex string to a big.Int
	value := new(big.Int)
	value, success := value.SetString(hexStr[2:], 16)
	if !success {
		return nil, fmt.Errorf("failed to parse hex string: %s", hexStr)
	}

	// Create a mask for the low 128 bits (2^128 - 1)
	mask128 := new(big.Int).Sub(
		new(big.Int).Lsh(big.NewInt(1), 128),
		big.NewInt(1),
	)

	// Extract the low 128 bits
	lowBits := new(big.Int).And(value, mask128)

	// Extract the high bits by shifting right by 128
	highBits := new(big.Int).Rsh(value, 128)

	// Convert to felt.Felt values
	lowFelt := BigIntToFelt(lowBits)
	highFelt := BigIntToFelt(highBits)

	// Return as a slice [low, high]
	return []*felt.Felt{lowFelt, highFelt}, nil
}

// U256FeltToHex converts a Cairo u256 representation (two felt.Felt values) back to a hexadecimal string.
// The Cairo u256 is represented as two felt.Felt values:
//   - The first felt.Felt contains the 128 least significant bits (low part)
//   - The second felt.Felt contains the 128 most significant bits (high part)
//
// Parameters:
//   - u256: a slice containing two felt.Felt values [low, high]
// Returns:
//   - string: the hexadecimal representation of the combined value
//   - error: if conversion fails
func U256FeltToHex(u256 []*felt.Felt) (string, error) {
	// Check if the input is valid
	if len(u256) != 2 {
		return "", fmt.Errorf("expected 2 felt values for u256, got %d", len(u256))
	}

	// Extract low and high parts
	lowFelt, highFelt := u256[0], u256[1]

	// Convert to big.Int
	lowBits := FeltToBigInt(lowFelt)
	highBits := FeltToBigInt(highFelt)

	// Combine the parts: result = highBits << 128 + lowBits
	result := new(big.Int).Lsh(highBits, 128)  // Shift high bits left by 128 bits
	result = new(big.Int).Add(result, lowBits) // Add low bits

	// Convert to hex string with "0x" prefix
	hexStr := fmt.Sprintf("%#x", result)
	return hexStr, nil
}
