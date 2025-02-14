package utils

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/stretchr/testify/require"
)

// Just a random felt variable to use when needed. The value is "0xdeadbeef"
var RANDOM_FELT = new(felt.Felt).SetUint64(3735928559)

// TestHexToFelt generates a felt.Felt from a hexadecimal string.
//
// Parameters:
// - t: the testing.TB object for test logging and reporting
// - hex: the hexadecimal string to convert to a felt.Felt
// Returns:
// - *felt.Felt: the generated felt.Felt object
func TestHexToFelt(t testing.TB, hex string) *felt.Felt {
	t.Helper()
	f, err := HexToFelt(hex)
	require.NoError(t, err)
	return f
}

// TestHexArrToFelt generates a slice of *felt.Felt from a slice of strings representing hexadecimal values.
//
// Parameters:
// - t: A testing.TB interface used for test logging and error reporting
// - hexArr: A slice of strings representing hexadecimal values
// Returns:
// - []*felt.Felt: a slice of *felt.Felt
func TestHexArrToFelt(t testing.TB, hexArr []string) []*felt.Felt {
	t.Helper()
	feltArr, err := HexArrToFelt(hexArr)
	require.NoError(t, err)
	return feltArr
}

// TestUnmarshallFileToType reads a JSON file at the given path and unmarshals it into the specified type T.
// If any error occurs during file reading or unmarshalling, it fails the test.
//
// Parameters:
// - t: testing.TB interface for test logging and error reporting
// - filePath: string path to the JSON file
// - isRPCResp: boolean indicating if the JSON file is in JSON-RPC response format
// Returns:
// - T: the unmarshalled data of type T
func TestUnmarshallFileToType[T any](t testing.TB, filePath string, isRPCResp bool) *T {
	t.Helper()
	var result T

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", filePath, err)
	}

	if isRPCResp {
		var rpcResponse struct {
			Result json.RawMessage `json:"result"`
		}
		if err := json.Unmarshal(data, &rpcResponse); err != nil {
			t.Fatalf("failed to unmarshal RPC response from file %s: %v", filePath, err)
		}
		data = rpcResponse.Result
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal file %s: %v", filePath, err)
	}

	return &result
}
