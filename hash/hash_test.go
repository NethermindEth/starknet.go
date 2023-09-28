package hash_test

import (
	"encoding/json"
	"testing"

	"github.com/NethermindEth/starknet.go/artifacts"
	"github.com/NethermindEth/starknet.go/hash"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/test-go/testify/require"
)

func TestCompiledClassHash(t *testing.T) {
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
