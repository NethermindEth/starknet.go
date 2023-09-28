package hash_test

import (
	"encoding/json"
	"testing"

	"github.com/NethermindEth/starknet.go/artifacts"
	"github.com/NethermindEth/starknet.go/hash"
	newcontract "github.com/NethermindEth/starknet.go/newcontracts"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/test-go/testify/require"
)

func TestUnmarshalCompiledClassHash(t *testing.T) {
	// todo get ClassHash hash from sequencer
	expectedClasshash := "0x_todo_insert_"

	compiledClass := artifacts.HelloWorldSierra
	var class rpc.ContractClass
	err := json.Unmarshal(compiledClass, &class)
	require.NoError(t, err)
	compClassHash, err := hash.ClassHash(class)
	require.NoError(t, err)
	require.Equal(t, expectedClasshash, compClassHash.String())
}

func TestUnmarshalCasmClassHash(t *testing.T) {
	compiledClass := artifacts.HelloWorldCasm
	var class newcontract.CasmClass
	err := json.Unmarshal(compiledClass, &class)
	require.NoError(t, err)
}

func TestCompiledClassHash(t *testing.T) {
	expectedHash := "0x_todo_"

	casmClass, err := newcontract.UnmarshalCasmClass("../artifacts/starknet_hello_world_Balance.casm.json")
	require.NoError(t, err)

	hash, err := hash.CompiledClassHash(*casmClass)
	require.NoError(t, err)
	require.Equal(t, expectedHash, hash.String())
}
