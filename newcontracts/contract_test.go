package newcontract

import (
	"encoding/json"
	"testing"

	"github.com/NethermindEth/starknet.go/artifacts"
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

func TestCompiledClassHash(t *testing.T) {
	expectedClasshash := "0x_todo_insert_"

	compiledClass := artifacts.HelloWorldSierra
	var class rpc.ContractClass
	err := json.Unmarshal(compiledClass, &class)
	require.NoError(t, err)
	compClassHash, err := ClassHash(class)
	require.NoError(t, err)
	require.Equal(t, expectedClasshash, compClassHash.String())
}
