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
	tests.RunTestOn(t, tests.TestnetEnv, tests.MainnetEnv)
	t.Parallel()

	testConfig := BeforeEach(t, false)

	// setup provider with a spy
	provider := testConfig.Provider
	spy := tests.NewRPCSpy(provider.c, false)
	provider.c = spy

	t.Run("No multiplier", func(t *testing.T) {
		t.Parallel()
		estimatedTip, err := EstimateTip(t.Context(), provider, 0)
		require.NoError(t, err)

		// get from the spy the latest block used to estimate the tip
		rawBlock := spy.LastResponse()
		var block Block
		require.NoError(t, json.Unmarshal(rawBlock, &block))

		averageTip := getTipAverageFromBlock(t, &block)

		// compare the estimated tips
		currentTip := U64("0x" + strconv.FormatUint(averageTip, 16))
		assert.Equal(t, currentTip, estimatedTip)
	})
	t.Run("With multiplier 1.5", func(t *testing.T) {
		t.Parallel()
		estimatedTip, err := EstimateTip(t.Context(), provider, 1.5)
		require.NoError(t, err)

		// get from the spy the latest block used to estimate the tip
		rawBlock := spy.LastResponse()
		var block Block
		require.NoError(t, json.Unmarshal(rawBlock, &block))

		averageTip := getTipAverageFromBlock(t, &block)

		// compare the estimated tips
		currentTip := U64("0x" + strconv.FormatUint(uint64(float64(averageTip)*1.5), 16))
		assert.Equal(t, currentTip, estimatedTip)
	})
	t.Run("With negative multiplier", func(t *testing.T) {
		t.Parallel()
		estimatedTip, err := EstimateTip(t.Context(), provider, -1.5)
		require.NoError(t, err)

		// get from the spy the latest block used to estimate the tip
		rawBlock := spy.LastResponse()
		var block Block
		require.NoError(t, json.Unmarshal(rawBlock, &block))

		averageTip := getTipAverageFromBlock(t, &block)

		// compare the estimated tips
		// (no multiplier is applied for negative multipliers)
		currentTip := U64("0x" + strconv.FormatUint(averageTip, 16))
		assert.Equal(t, currentTip, estimatedTip)
	})
	t.Run("With multiplier less than 1", func(t *testing.T) {
		t.Parallel()
		estimatedTip, err := EstimateTip(t.Context(), provider, 0.5)
		require.NoError(t, err)

		// get from the spy the latest block used to estimate the tip
		rawBlock := spy.LastResponse()
		var block Block
		require.NoError(t, json.Unmarshal(rawBlock, &block))

		averageTip := getTipAverageFromBlock(t, &block)
		if averageTip == 0 {
			assert.Equal(t, U64("0x0"), estimatedTip)

			return
		}

		// compare the estimated tips
		currentTip := (uint64(float64(averageTip) * 0.5))
		assert.Equal(t, U64("0x"+strconv.FormatUint(currentTip, 16)), estimatedTip)
		assert.Less(t, currentTip, averageTip)
	})
}

// getTipAverageFromBlock returns the average of the tips from all transactions in the block
func getTipAverageFromBlock(t *testing.T, block *Block) uint64 {
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

	return tipCounter / uint64(len(block.Transactions))
}
