package utils

import (
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
	var qwe []*felt.Felt
	for _, hex := range hexArr {
		felt, err := HexToFelt(hex)
		if err != nil {
			return nil, err
		}
		qwe = append(qwe, felt)
	}
	return qwe, nil
}

func FeltToBigInt(f *felt.Felt) (*big.Int, bool) {
	return new(big.Int).SetString(f.String(), 0)
}

func BigIntToFelt(big *big.Int) (*felt.Felt, error) {
	return new(felt.Felt).SetString(big.String())
}
