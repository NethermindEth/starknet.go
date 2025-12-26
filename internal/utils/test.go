package utils

import (
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/stretchr/testify/require"
)

// Just a random felt variable to use when needed. The value is "0xdeadbeef"
var DeadBeef = new(felt.Felt).SetUint64(3735928559) //nolint:mnd // 0xdeadbeef

// TestHexToFelt generates a felt.Felt from a hexadecimal string.
//
// Parameters:
//   - t: the testing.TB object for test logging and reporting
//   - hex: the hexadecimal string to convert to a felt.Felt
//
// Returns:
//   - *felt.Felt: the generated felt.Felt object
func TestHexToFelt(t testing.TB, hex string) *felt.Felt {
	t.Helper()
	f, err := HexToFelt(hex)
	require.NoError(t, err)

	return f
}

// TestHexArrToFelt generates a slice of *felt.Felt from a slice of strings
// representing hexadecimal values.
//
// Parameters:
//   - t: A testing.TB interface used for test logging and error reporting
//   - hexArr: A slice of strings representing hexadecimal values
//
// Returns:
//   - []*felt.Felt: a slice of *felt.Felt
func TestHexArrToFelt(t testing.TB, hexArr []string) []*felt.Felt {
	t.Helper()
	feltArr, err := HexArrToFelt(hexArr)
	require.NoError(t, err)

	return feltArr
}

// TestUnmarshalJSONFileToType reads a JSON file at the given path and unmarshals it
// into the specified type T.
// If any error occurs during file reading or unmarshalling, it fails the test.
//
// Parameters:
//   - t: testing.TB interface for test logging and error reporting
//   - filePath: string path to the JSON file
//   - subfield: string subfield to unmarshal from the JSON file
//
// Returns:
//   - T: the unmarshalled data of type T
func TestUnmarshalJSONFileToType[T any](t testing.TB, filePath, subfield string) T {
	result, err := UnmarshalJSONFileToType[T](filePath, subfield)
	require.NoError(t, err)

	return result
}
