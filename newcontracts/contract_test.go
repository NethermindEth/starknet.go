package newcontract_test

import (
	"encoding/json"
	"testing"

	"github.com/NethermindEth/starknet.go/artifacts"
	"github.com/NethermindEth/starknet.go/hash"
	newcontract "github.com/NethermindEth/starknet.go/newcontracts"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/test-go/testify/require"
)

func TestUnmarshalContractClass(t *testing.T) {
	compiledClass := artifacts.HelloWorldSierra
	var class rpc.ContractClass
	err := json.Unmarshal(compiledClass, &class)
	require.NoError(t, err)
	require.Equal(t, class.SierraProgram[0].String(), "0x1")
	require.Equal(t, class.SierraProgram[1].String(), "0x3")
}

func TestUnmarshalCasmClass(t *testing.T) {
	casmClass, err := newcontract.UnmarshalCasmClass("../artifacts/starknet_hello_world_Balance.casm.json")
	require.NoError(t, err)
	require.Equal(t, casmClass.Prime, "0x800000000000011000000000000000000000000000000000000000000000001")
	require.Equal(t, casmClass.Version, "2.2.0")
	require.Equal(t, casmClass.EntryPointByType.External[0].Selector.String(), "0x17c00f03de8b5bd58d2016b59d251c13056b989171c5852949903bc043bc27")
	require.Equal(t, casmClass.EntryPointByType.External[0].Selector.String(), "0x17c00f03de8b5bd58d2016b59d251c13056b989171c5852949903bc043bc27")
	require.Equal(t, casmClass.EntryPointByType.External[1].Offset, 111)
	require.Equal(t, casmClass.EntryPointByType.External[1].Builtins[0], "range_check")
}

func TestCompiledClassHash(t *testing.T) {
	expectedHash := "0x_todo_"

	casmClass, err := newcontract.UnmarshalCasmClass("../artifacts/starknet_hello_world_Balance.casm.json")
	require.NoError(t, err)

	hash, err := hash.CompiledClassHash(*casmClass)
	require.NoError(t, err)
	require.Equal(t, expectedHash, hash.String())
}
