package hash_test

import (
	"encoding/json"
	"testing"

	"github.com/NethermindEth/starknet.go/artifacts"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/hash"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/test-go/testify/require"
)

func TestUnmarshalCasmClassHash(t *testing.T) {
	compiledClass := artifacts.ExampleWorldCasm
	var class contracts.CasmClass
	err := json.Unmarshal(compiledClass, &class)
	require.NoError(t, err)
}

func TestClassHash(t *testing.T) {
	//https://github.com/software-mansion/starknet.py/blob/development/starknet_py/hash/class_hash_test.py
	expectedClasshash := "0x4ec2ecf58014bc2ffd7c84843c3525e5ecb0a2cac33c47e9c347f39fc0c0944"

	compiledClass := artifacts.ExampleWorldSierra
	var class rpc.ContractClass
	err := json.Unmarshal(compiledClass, &class)
	require.NoError(t, err)
	compClassHash, err := hash.ClassHash(class)
	require.NoError(t, err)
	require.Equal(t, expectedClasshash, compClassHash.String())
}

func TestCompiledClassHash(t *testing.T) {
	//https://github.com/software-mansion/starknet.py/blob/development/starknet_py/hash/casm_class_hash_test.py
	expectedHash := "0x785fa5f2bacf0bfe3bc413be5820a61e1ea63f2ec27ef00331ee9f46ad07603"

	compiledClass := artifacts.ExampleWorldCasm
	var casmClass contracts.CasmClass
	err := json.Unmarshal(compiledClass, &casmClass)
	require.NoError(t, err)

	hash := hash.CompiledClassHash(casmClass)
	require.Equal(t, expectedHash, hash.String())
}