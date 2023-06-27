package utils

import "github.com/NethermindEth/juno/core/felt"

func HexToFelt(hex string) (*felt.Felt, error) {
	return new(felt.Felt).SetString(hex)
}
