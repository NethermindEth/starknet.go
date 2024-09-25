package hash_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/hash"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/stretchr/testify/require"
)

// TestUnmarshalCasmClassHash is a test function that tests the unmarshaling of a CasmClass hash.
//
// It reads the content of the "./tests/hello_starknet_compiled.casm.json" file and unmarshals it into a contracts.CasmClass variable.
// The function returns an assertion error if there is an error reading the file or unmarshaling the content.
// Parameters:
// - t: A testing.T object used for running the test and reporting any failures.
// Returns:
//
//	none
func TestUnmarshalCasmClassHash(t *testing.T) {
	content, err := os.ReadFile("./tests/hello_starknet_compiled.casm.json")
	require.NoError(t, err)

	var class contracts.CasmClass
	err = json.Unmarshal(content, &class)
	require.NoError(t, err)
}

// TestClassHash is a test function that verifies the correctness of the ClassHash function.
//
// It reads the expected class hash from a specific source and compares it with the computed class hash.
// The function expects the class hash to be equal to the expected value.
//
// Parameters:
// - t: A testing.T object used for running the test and reporting any failures.
// Returns:
//
//	none
func TestClassHash(t *testing.T) {
	//https://github.com/software-mansion/starknet.py/blob/development/starknet_py/hash/class_hash_test.py
	expectedClasshash := "0x4ec2ecf58014bc2ffd7c84843c3525e5ecb0a2cac33c47e9c347f39fc0c0944"

	content, err := os.ReadFile("./tests/hello_starknet_compiled.sierra.json")
	require.NoError(t, err)

	var class rpc.ContractClass
	err = json.Unmarshal(content, &class)
	require.NoError(t, err)
	compClassHash := hash.ClassHash(class)
	require.Equal(t, expectedClasshash, compClassHash.String())
}

// TestCompiledClassHash is a test function that verifies the correctness of the CompiledClassHash function in the hash package.
//
// It reads the contents of a file, unmarshals it into a CasmClass object, computes the hash using the CompiledClassHash function,
// and asserts that the computed hash matches the expected hash.
//
// Parameters:
// - t: A testing.T object used for running the test and reporting any failures.
// Returns:
//
//	none
func TestCompiledClassHash(t *testing.T) {
	//https://github.com/software-mansion/starknet.py/blob/development/starknet_py/hash/casm_class_hash_test.py
	expectedHash := "0x785fa5f2bacf0bfe3bc413be5820a61e1ea63f2ec27ef00331ee9f46ad07603"

	content, err := os.ReadFile("./tests/hello_starknet_compiled.casm.json")
	require.NoError(t, err)

	var casmClass contracts.CasmClass
	err = json.Unmarshal(content, &casmClass)
	require.NoError(t, err)

	hash := hash.CompiledClassHash(casmClass)
	require.Equal(t, expectedHash, hash.String())
}
