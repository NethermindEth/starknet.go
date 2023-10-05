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

func TestUnmarshalContractClass(t *testing.T) {
	content, err := os.ReadFile("./tests/hello_starknet_compiled.sierra.json")
	require.NoError(t, err)

	var class rpc.ContractClass
	err = json.Unmarshal(content, &class)
	require.NoError(t, err)
	assert.Equal(t, class.SierraProgram[0].String(), "0x1")
	assert.Equal(t, class.SierraProgram[1].String(), "0x3")
}

func TestUnmarshalCasmClass(t *testing.T) {
	casmClass, err := contracts.UnmarshalCasmClass("./tests/hello_starknet_compiled.casm.json")
	require.NoError(t, err)
	assert.Equal(t, casmClass.Prime, "0x800000000000011000000000000000000000000000000000000000000000001")
	assert.Equal(t, casmClass.Version, "2.1.0")
	assert.Equal(t, casmClass.EntryPointByType.External[0].Selector.String(), "0x362398bec32bc0ebb411203221a35a0301193a96f317ebe5e40be9f60d15320")
	assert.Equal(t, casmClass.EntryPointByType.External[1].Offset, 130)
	assert.Equal(t, casmClass.EntryPointByType.External[1].Builtins[0], "range_check")
}
