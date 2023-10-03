package contracts_test

import (
	"encoding/json"
	"testing"

	"github.com/NethermindEth/starknet.go/artifacts"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/test-go/testify/require"
)

func TestUnmarshalContractClass(t *testing.T) {
	compiledClass := artifacts.ExampleWorldSierra
	var class rpc.ContractClass
	err := json.Unmarshal(compiledClass, &class)
	require.NoError(t, err)
	require.Equal(t, class.SierraProgram[0].String(), "0x1")
	require.Equal(t, class.SierraProgram[1].String(), "0x3")
}

func TestUnmarshalCasmClass(t *testing.T) {
	casmClass, err := contracts.UnmarshalCasmClass("../artifacts/hello_starknet_compiled.casm.json")
	require.NoError(t, err)
	require.Equal(t, casmClass.Prime, "0x800000000000011000000000000000000000000000000000000000000000001")
	require.Equal(t, casmClass.Version, "2.1.0")
	require.Equal(t, casmClass.EntryPointByType.External[0].Selector.String(), "0x362398bec32bc0ebb411203221a35a0301193a96f317ebe5e40be9f60d15320")
	require.Equal(t, casmClass.EntryPointByType.External[1].Offset, 130)
	require.Equal(t, casmClass.EntryPointByType.External[1].Builtins[0], "range_check")
}
