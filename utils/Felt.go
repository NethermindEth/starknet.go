package utils

import (
	"errors"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
)

func Uint64ToFelt(num uint64) *felt.Felt {
	return new(felt.Felt).SetUint64(num)
}

func HexToFelt(hex string) (*felt.Felt, error) {
	return new(felt.Felt).SetString(hex)
}

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

func FeltToBigInt(f *felt.Felt) *big.Int {
	tmp := f.Bytes()
	return new(big.Int).SetBytes(tmp[:])
}

func BigIntToFelt(big *big.Int) *felt.Felt {
	return new(felt.Felt).SetBytes(big.Bytes())
}

func FeltArrToBigIntArr(f []*felt.Felt) ([]*big.Int, error) {
	var bigArr []*big.Int
	for _, felt := range f {
		bigArr = append(bigArr, FeltToBigInt(felt))
	}
	return bigArr, nil
}
func FeltArrToBigIntArr(f []*felt.Felt) (*[]*big.Int, error) {
	var bigArr []*big.Int
	for _, felt := range f {
		bint, ok := FeltToBigInt(felt)
		if !ok {
			return nil, errors.New("Failed to convert felt to big.Int")
		}
		bigArr = append(bigArr, bint)
	}
	return &bigArr, nil
}
