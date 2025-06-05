package account_test

import (
	"testing"

	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/stretchr/testify/assert"
)

// TestTxnOptions tests the methods of the TxnOptions struct,
// verifying that default values are set correctly and edge cases are handled properly
func TestTxnOptions(t *testing.T) {
	t.Parallel()

	t.Run("SafeMultiplier", func(t *testing.T) {
		t.Parallel()
		testcases := []struct {
			name               string
			opts               *account.TxnOptions
			expectedMultiplier float64
		}{
			{
				name:               "Default value (nil)",
				opts:               nil,
				expectedMultiplier: 1.5,
			},
			{
				name: "Zero multiplier",
				opts: &account.TxnOptions{
					Multiplier: 0,
				},
				expectedMultiplier: 1.5,
			},
			{
				name: "Negative multiplier",
				opts: &account.TxnOptions{
					Multiplier: -1.0,
				},
				expectedMultiplier: 1.5,
			},
			{
				name: "Custom multiplier",
				opts: &account.TxnOptions{
					Multiplier: 2.0,
				},
				expectedMultiplier: 2.0,
			},
			{
				name: "Custom multiplier below 1",
				opts: &account.TxnOptions{
					Multiplier: 0.5,
				},
				expectedMultiplier: 0.5,
			},
		}

		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.Equal(t, tt.expectedMultiplier, tt.opts.SafeMultiplier())
			})
		}
	})

	t.Run("BlockID", func(t *testing.T) {
		t.Parallel()
		testcases := []struct {
			name            string
			opts            *account.TxnOptions
			expectedBlockID rpc.BlockID
		}{
			{
				name:            "Default value (nil)",
				opts:            nil,
				expectedBlockID: rpc.WithBlockTag(rpc.BlockTagPending),
			},
			{
				name: "Empty block tag",
				opts: &account.TxnOptions{
					EstimationBlockTag: "",
				},
				expectedBlockID: rpc.WithBlockTag(rpc.BlockTagPending),
			},
			{
				name: "latest block tag",
				opts: &account.TxnOptions{
					EstimationBlockTag: rpc.BlockTagLatest,
				},
				expectedBlockID: rpc.WithBlockTag(rpc.BlockTagLatest),
			},
			{
				name: "pending block tag",
				opts: &account.TxnOptions{
					EstimationBlockTag: rpc.BlockTagPending,
				},
				expectedBlockID: rpc.WithBlockTag(rpc.BlockTagPending),
			},
		}

		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.Equal(t, tt.expectedBlockID, tt.opts.BlockID())
			})
		}
	})

	t.Run("SimulationFlags", func(t *testing.T) {
		t.Parallel()
		testcases := []struct {
			name             string
			opts             *account.TxnOptions
			expectedSimFlags []rpc.SimulationFlag
		}{
			{
				name:             "Default value (nil)",
				opts:             nil,
				expectedSimFlags: []rpc.SimulationFlag{},
			},
			{
				name:             "Empty simulation flag",
				opts:             &account.TxnOptions{SimulationFlag: ""},
				expectedSimFlags: []rpc.SimulationFlag{},
			},
			{
				name: "SKIP_VALIDATE flag",
				opts: &account.TxnOptions{
					SimulationFlag: rpc.SKIP_VALIDATE,
				},
				expectedSimFlags: []rpc.SimulationFlag{rpc.SKIP_VALIDATE},
			},
			{
				name: "SKIP_FEE_CHARGE flag",
				opts: &account.TxnOptions{
					SimulationFlag: rpc.SKIP_FEE_CHARGE,
				},
				expectedSimFlags: []rpc.SimulationFlag{rpc.SKIP_FEE_CHARGE},
			},
			{
				name: "SKIP_EXECUTE flag",
				opts: &account.TxnOptions{
					SimulationFlag: rpc.SKIP_EXECUTE,
				},
				expectedSimFlags: []rpc.SimulationFlag{rpc.SKIP_EXECUTE},
			},
		}

		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.Equal(t, tt.expectedSimFlags, tt.opts.SimulationFlags())
			})
		}
	})
}
