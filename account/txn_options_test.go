package account

import (
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/stretchr/testify/assert"
)

// TestTxnOptions tests the methods of the TxnOptions struct,
// verifying that default values are set correctly and edge cases are handled properly
//

func TestTxnOptions(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv)
	t.Parallel()

	t.Run("FmtFeeMultiplier", func(t *testing.T) {
		t.Parallel()
		testcases := []struct {
			name               string
			opts               *TxnOptions
			expectedMultiplier float64
		}{
			{
				name:               "Default value (nil)",
				opts:               nil,
				expectedMultiplier: 1.5,
			},
			{
				name: "Zero multiplier",
				opts: &TxnOptions{
					FeeMultiplier: 0,
				},
				expectedMultiplier: 1.5,
			},
			{
				name: "Negative multiplier",
				opts: &TxnOptions{
					FeeMultiplier: -1.0,
				},
				expectedMultiplier: 1.5,
			},
			{
				name: "Custom multiplier",
				opts: &TxnOptions{
					FeeMultiplier: 2.0,
				},
				expectedMultiplier: 2.0,
			},
			{
				name: "Custom multiplier below 1",
				opts: &TxnOptions{
					FeeMultiplier: 0.5,
				},
				expectedMultiplier: 0.5,
			},
		}

		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				if tt.opts == nil {
					tt.opts = new(TxnOptions)
				}
				assert.Equal(t, tt.expectedMultiplier, tt.opts.FmtFeeMultiplier())
			})
		}
	})

	t.Run("FmtTipMultiplier", func(t *testing.T) {
		t.Parallel()
		testcases := []struct {
			name               string
			opts               *TxnOptions
			expectedMultiplier float64
		}{
			{
				name:               "Default value (nil)",
				opts:               nil,
				expectedMultiplier: 1.0,
			},
			{
				name: "Zero multiplier",
				opts: &TxnOptions{
					TipMultiplier: 0,
				},
				expectedMultiplier: 1.0,
			},
			{
				name: "Negative multiplier",
				opts: &TxnOptions{
					TipMultiplier: -1.0,
				},
				expectedMultiplier: 1.0,
			},
			{
				name: "Custom multiplier",
				opts: &TxnOptions{
					TipMultiplier: 2.0,
				},
				expectedMultiplier: 2.0,
			},
			{
				name: "Custom multiplier below 1",
				opts: &TxnOptions{
					TipMultiplier: 0.5,
				},
				expectedMultiplier: 0.5,
			},
		}

		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				if tt.opts == nil {
					tt.opts = new(TxnOptions)
				}
				assert.Equal(t, tt.expectedMultiplier, tt.opts.FmtTipMultiplier())
			})
		}
	})

	t.Run("BlockID", func(t *testing.T) {
		t.Parallel()
		testcases := []struct {
			name            string
			opts            *TxnOptions
			expectedBlockID rpc.BlockID
		}{
			{
				name:            "Default value (nil)",
				opts:            nil,
				expectedBlockID: rpc.WithBlockTag(rpc.BlockTagPreConfirmed),
			},
			{
				name: "latest set to true",
				opts: &TxnOptions{
					UseLatest: true,
				},
				expectedBlockID: rpc.WithBlockTag(rpc.BlockTagLatest),
			},
			{
				name: "latest set to false",
				opts: &TxnOptions{
					UseLatest: false,
				},
				expectedBlockID: rpc.WithBlockTag(rpc.BlockTagPreConfirmed),
			},
		}

		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				if tt.opts == nil {
					tt.opts = new(TxnOptions)
				}
				assert.Equal(t, tt.expectedBlockID, tt.opts.BlockID())
			})
		}
	})

	t.Run("SimulationFlags", func(t *testing.T) {
		t.Parallel()
		testcases := []struct {
			name             string
			opts             *TxnOptions
			expectedSimFlags []rpc.SimulationFlag
		}{
			{
				name:             "Default value (nil)",
				opts:             nil,
				expectedSimFlags: []rpc.SimulationFlag{},
			},
			{
				name:             "Empty simulation flag",
				opts:             &TxnOptions{SimulationFlag: ""},
				expectedSimFlags: []rpc.SimulationFlag{},
			},
			{
				name: "SKIP_VALIDATE flag",
				opts: &TxnOptions{
					SimulationFlag: rpc.SkipValidate,
				},
				expectedSimFlags: []rpc.SimulationFlag{rpc.SkipValidate},
			},
			{
				name: "SKIP_FEE_CHARGE flag",
				opts: &TxnOptions{
					SimulationFlag: rpc.SkipFeeCharge,
				},
				expectedSimFlags: []rpc.SimulationFlag{rpc.SkipFeeCharge},
			},
		}

		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				if tt.opts == nil {
					tt.opts = new(TxnOptions)
				}
				assert.Equal(t, tt.expectedSimFlags, tt.opts.SimulationFlags())
			})
		}
	})
}
