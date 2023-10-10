package utils

import (
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
)

// Uint64ToFelt generates a new *felt.Felt from a given uint64 number.
//
// It takes a uint64 parameter called 'num' and returns a *felt.Felt.
func Uint64ToFelt(num uint64) *felt.Felt {
	return new(felt.Felt).SetUint64(num)
}

// HexToFelt converts a hexadecimal string to a *felt.Felt object.
//
// hex: the input hexadecimal string to be converted.
// Returns a *felt.Felt object and an error if conversion fails.
func HexToFelt(hex string) (*felt.Felt, error) {
	return new(felt.Felt).SetString(hex)
}

// HexArrToFelt converts an array of hexadecimal strings to an array of felt objects.
//
// It takes an input parameter hexArr, which is an array of strings representing hexadecimal values.
// The function iterates over each element in the hexArr array and calls the HexToFelt function to convert each hexadecimal value to a felt object.
// If any error occurs during the conversion, the function will return nil and the corresponding error.
// Otherwise, it appends the converted felt object to the feltArr array.
// Finally, the function returns the feltArr array containing all the converted felt objects.
//
// The function has the following return type: []*felt.Felt, error.
func HexArrToFelt(hexArr []string) ([]*felt.Felt, error) {
	var feltArr []*felt.Felt
	for _, hex := range hexArr {
		felt, err := HexToFelt(hex)
		if err != nil {
			return nil, err
		}
		feltArr = append(feltArr, felt)
	}
	return feltArr, nil
}

// FeltToBigInt converts a Felt value to a *big.Int.
//
// It takes a pointer to a Felt value as a parameter and returns a pointer to a big.Int.
func FeltToBigInt(f *felt.Felt) *big.Int {
	tmp := f.Bytes()
	return new(big.Int).SetBytes(tmp[:])
}

// BigIntToFelt converts a big integer to a felt.Felt.
//
// It takes a pointer to a big.Int as a parameter and returns a pointer to a felt.Felt.
func BigIntToFelt(big *big.Int) *felt.Felt {
	return new(felt.Felt).SetBytes(big.Bytes())
}

// FeltArrToBigIntArr converts an array of Felt objects to an array of big.Int objects.
//
// f - the array of Felt objects to convert.
// Returns an array of big.Int objects.
func FeltArrToBigIntArr(f []*felt.Felt) []*big.Int {
	var bigArr []*big.Int
	for _, felt := range f {
		bigArr = append(bigArr, FeltToBigInt(felt))
	}
	return bigArr
}
