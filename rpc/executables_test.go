package rpc

import (
	"context"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/stretchr/testify/require"
)

func TestCompiledCasm(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		Description    string
		ClassHash      *felt.Felt
		ExpectedResult *CasmCompiledContractClass
		ExpectedError  error
	}
	testSet := map[string][]testSetType{
		"mock":   {},
		"devnet": {},
		"testnet": {
			{
				Description:    "normal call, only required field class_hash",
				ClassHash:      utils.TestHexToFelt(t, "0x00d764f235da1c654c4ca14c47bfc2a54ccd4c0c56b3f4570cd241bd638db448"),
				ExpectedResult: utils.UnmarshallFileToType[CasmCompiledContractClass](t, "./tests/compiledCasm.json", true),
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		require := require.New(t)
		result, err := testConfig.provider.CompiledCasm(context.Background(), test.ClassHash)
		require.NoError(err)
		require.NotNil(result, "should return a nonce")
		require.Equal(test.ExpectedResult, result)
	}
}
