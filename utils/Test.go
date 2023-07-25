package utils

import (
	"math/big"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/test-go/testify/require"
)

func TestHexToFelt(t testing.TB, hex string) *felt.Felt {
	t.Helper()
	f, err := HexToFelt(hex)
	require.NoError(t, err)
	return f
}

func TestBigIntToFelt(t testing.TB, big *big.Int) *felt.Felt {
	t.Helper()
	felt, err := BigIntToFelt(big)
	require.NoError(t, err)
	return felt
}

func TestHexArrToFelt(t testing.TB, hexArr []string) []*felt.Felt {
	t.Helper()
	feltArr, err := HexArrToFelt(hexArr)
	require.NoError(t, err)
	return feltArr
}
