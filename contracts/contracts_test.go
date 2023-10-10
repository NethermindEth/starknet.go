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
// It checks if the UnmarshalCasmClass function correctly unmarshals the contents of the casmClass file.
// The function takes a file path as a parameter and returns the unmarshalled CasmClass object and an error.
// The function asserts that the Prime field of the CasmClass object is equal to "0x800000000000011000000000000000000000000000000000000000000000001",
// the Version field is equal to "2.1.0", the Selector field of the first External entry point is equal to "0x362398bec32bc0ebb411203221a35a0301193a96f317ebe5e40be9f60d15320",
// the Offset field of the second External entry point is equal to 130, and the first Builtin of the second External entry point is equal to "range_check".
//
// The function uses the require.NoError and assert.Equal functions from the github.com/stretchr/testify/assert package to perform the assertions.
// It is a test function and is meant to be used with the Go testing framework.
func TestUnmarshalCasmClass(t *testing.T) {
	casmClass, err := contracts.UnmarshalCasmClass("./tests/hello_starknet_compiled.casm.json")
	require.NoError(t, err)
	assert.Equal(t, casmClass.Prime, "0x800000000000011000000000000000000000000000000000000000000000001")
	assert.Equal(t, casmClass.Version, "2.1.0")
	assert.Equal(t, casmClass.EntryPointByType.External[0].Selector.String(), "0x362398bec32bc0ebb411203221a35a0301193a96f317ebe5e40be9f60d15320")
	assert.Equal(t, casmClass.EntryPointByType.External[1].Offset, 130)
	assert.Equal(t, casmClass.EntryPointByType.External[1].Builtins[0], "range_check")
}
