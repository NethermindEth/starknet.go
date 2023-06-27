package utils

import "github.com/NethermindEth/caigo/types/felt"

func HexToFelt(hex string) (*felt.Felt, error) {
	return new(felt.Felt).SetString(hex)
}
