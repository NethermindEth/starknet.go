package rpc

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	"github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// IsCompatible method is already tested in the provider_test.go file,
// in the TestVersionCompatibility test

// TestProvider_EstimateTip tests the EstimateTip method by
// calculating the estimated tip again and comparing it with
// the estimated tip returned by the provider.
func TestProvider_EstimateTip(t *testing.T) {
	tests.RunTestOn(t, tests.TestnetEnv)
	testConfig := BeforeEach(t, false)

	// setup provider with a spy
	provider := testConfig.Provider
	spy := tests.NewJSONRPCSpy(provider.c, false)
	provider.c = spy

	estimatedTip, err := EstimateTip(t.Context(), provider)
	require.NoError(t, err)

	// get from the spy the latest block used to estimate the tip
	rawBlock := spy.LastResponse()
	var block Block
	require.NoError(t, json.Unmarshal(rawBlock, &block))

	var tipCounter uint64
	for _, tnx := range block.Transactions {
		// get the tip from the transaction
		rawTxn, err := json.Marshal(tnx.Transaction)
		require.NoError(t, err)
		var txnMap map[string]json.RawMessage
		require.NoError(t, json.Unmarshal(rawTxn, &txnMap))
		tip, err := utils.GetAndUnmarshalJSONFromMap[U64](txnMap, "tip")
		require.NoError(t, err)

		// convert the tip to uint64 and add it to the counter
		tipUint, err := tip.ToUint64()
		require.NoError(t, err)
		tipCounter += tipUint
	}

	// compare the estimated tips
	currentTip := U64(strconv.FormatUint(tipCounter/uint64(len(block.Transactions)), 16))
	assert.Equal(t, currentTip, estimatedTip)
}
