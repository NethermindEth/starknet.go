package contracts_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/test-go/testify/assert"
	"github.com/test-go/testify/require"
)

// TestUnmarshalContractClass is a test function to unmarshal a contract class.
//
// It reads the content of a file and unmarshals it into a ContractClass struct.
// Then it asserts the equality of certain values within the struct.
//
// Parameters:
// - t: The testing.T instance for running the test
// Returns:
//   none
func TestUnmarshalContractClass(t *testing.T) {
	content, err := os.ReadFile("./tests/hello_starknet_compiled.sierra.json")
	require.NoError(t, err)

	var class rpc.ContractClass
	err = json.Unmarshal(content, &class)
	require.NoError(t, err)
	assert.Equal(t, class.SierraProgram[0].String(), "0x1")
	assert.Equal(t, class.SierraProgram[1].String(), "0x3")
}

// TestUnmarshalCasmClass tests the UnmarshalCasmClass function.
//
// It reads the content of a file and unmarshals it into a CasmClass struct.
// The function uses the require.NoError and assert.Equal functions from the github.com/stretchr/testify/assert package to perform the assertions.
// It is a test function and is meant to be used with the Go testing framework.
//
// Parameters:
// - t: The testing.T instance for running the test
// Returns:
//   none
func TestUnmarshalCasmClass(t *testing.T) {
	casmClass, err := contracts.UnmarshalCasmClass("./tests/hello_starknet_compiled.casm.json")
	require.NoError(t, err)
	assert.Equal(t, casmClass.Prime, "0x800000000000011000000000000000000000000000000000000000000000001")
	assert.Equal(t, casmClass.Version, "2.1.0")
	assert.Equal(t, casmClass.EntryPointByType.External[0].Selector.String(), "0x362398bec32bc0ebb411203221a35a0301193a96f317ebe5e40be9f60d15320")
	assert.Equal(t, casmClass.EntryPointByType.External[1].Offset, 130)
	assert.Equal(t, casmClass.EntryPointByType.External[1].Builtins[0], "range_check")
}
