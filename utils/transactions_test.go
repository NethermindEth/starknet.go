package utils

import (
	"fmt"
	"math/big"
	"testing"

	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResBoundsMapToOverallFee(t *testing.T) {
	tests := []struct {
		name        string
		resBounds   rpc.ResourceBoundsMapping
		multiplier  float64
		expectedRes string
		expectedErr string
	}{
		{
			name: "Basic calculation",
			resBounds: rpc.ResourceBoundsMapping{
				L1Gas: rpc.ResourceBounds{
					MaxAmount:       "0x64", // 100
					MaxPricePerUnit: "0x64", // 100
				},
				L1DataGas: rpc.ResourceBounds{
					MaxAmount:       "0x64", // 100
					MaxPricePerUnit: "0x64", // 100
				},
				L2Gas: rpc.ResourceBounds{
					MaxAmount:       "0x64", // 100
					MaxPricePerUnit: "0x64", // 100
				},
			},
			multiplier: 1.0,
			// Expected: (100*100) + (100*100) + (100*100) = 10000 + 10000 + 10000 = 30000
			expectedRes: "30000",
		},
		{
			name: "Zero values",
			resBounds: rpc.ResourceBoundsMapping{
				L1Gas: rpc.ResourceBounds{
					MaxAmount:       "0x0",
					MaxPricePerUnit: "0x0",
				},
				L1DataGas: rpc.ResourceBounds{
					MaxAmount:       "0x0",
					MaxPricePerUnit: "0x0",
				},
				L2Gas: rpc.ResourceBounds{
					MaxAmount:       "0x0",
					MaxPricePerUnit: "0x0",
				},
			},
			multiplier:  1.5,
			expectedRes: "0",
		},
		{
			name: "With multiplier 1.5",
			resBounds: rpc.ResourceBoundsMapping{
				L1Gas: rpc.ResourceBounds{
					MaxAmount:       "0x64", // 100
					MaxPricePerUnit: "0x64", // 100
				},
				L1DataGas: rpc.ResourceBounds{
					MaxAmount:       "0x64", // 100
					MaxPricePerUnit: "0x64", // 100
				},
				L2Gas: rpc.ResourceBounds{
					MaxAmount:       "0x64", // 100
					MaxPricePerUnit: "0x64", // 100
				},
			},
			multiplier: 1.5,
			// Expected: ((100*100) + (100*100) + (100*100)) * 1.5 = 30000 * 1.5 = 45000
			expectedRes: "45000",
		},
		{
			name: "Negative multiplier",
			resBounds: rpc.ResourceBoundsMapping{
				L1Gas: rpc.ResourceBounds{
					MaxAmount:       "0x64", // 100
					MaxPricePerUnit: "0x64", // 100
				},
				L1DataGas: rpc.ResourceBounds{
					MaxAmount:       "0x64", // 100
					MaxPricePerUnit: "0x64", // 100
				},
				L2Gas: rpc.ResourceBounds{
					MaxAmount:       "0x64", // 100
					MaxPricePerUnit: "0x64", // 100
				},
			},
			multiplier:  -1.0,
			expectedErr: "multiplier cannot be negative",
		},
		{
			name: "Multiplier less than 1",
			resBounds: rpc.ResourceBoundsMapping{
				L1Gas: rpc.ResourceBounds{
					MaxAmount:       "0x64", // 100
					MaxPricePerUnit: "0x64", // 100
				},
				L1DataGas: rpc.ResourceBounds{
					MaxAmount:       "0x64", // 100
					MaxPricePerUnit: "0x64", // 100
				},
				L2Gas: rpc.ResourceBounds{
					MaxAmount:       "0x64", // 100
					MaxPricePerUnit: "0x64", // 100
				},
			},
			multiplier: 0.5,
			// Expected: ((100*100) + (100*100) + (100*100)) * 0.5 = 30000 * 0.5 = 15000
			expectedRes: "15000",
		},
		{
			name: "Extremely large values",
			resBounds: rpc.ResourceBoundsMapping{
				L1Gas: rpc.ResourceBounds{
					MaxAmount:       "0x38d7ea4c68000",            // 1,000,000,000,000,000
					MaxPricePerUnit: "0x204fce5e3e25026110000000", // 10,000,000,000,000,000,000,000,000,000
				},
				L1DataGas: rpc.ResourceBounds{
					MaxAmount:       "0x38d7ea4c68000",            // 1,000,000,000,000,000
					MaxPricePerUnit: "0x204fce5e3e25026110000000", // 10,000,000,000,000,000,000,000,000,000
				},
				L2Gas: rpc.ResourceBounds{
					MaxAmount:       "0x38d7ea4c68000",            // 1,000,000,000,000,000
					MaxPricePerUnit: "0x204fce5e3e25026110000000", // 10,000,000,000,000,000,000,000,000,000
				},
			},
			multiplier:  1.5,
			expectedRes: "45000000000000000000000000000000000000000000",
		},
		{
			name: "Invalid resource bounds values",
			resBounds: rpc.ResourceBoundsMapping{
				L1Gas: rpc.ResourceBounds{
					MaxAmount:       "invalidValue", // Invalid format
					MaxPricePerUnit: "0xa",          // 10
				},
			},
			expectedErr: "invalid resource bounds: 'invalidValue' is not a valid big.Int",
		},
		{
			name: "Empty fields",
			resBounds: rpc.ResourceBoundsMapping{
				L1Gas: rpc.ResourceBounds{
					MaxAmount:       "",
					MaxPricePerUnit: "0x64", // 100
				},
				L1DataGas: rpc.ResourceBounds{
					MaxAmount:       "0x64", // 100
					MaxPricePerUnit: "0x64", // 100
				},
				L2Gas: rpc.ResourceBounds{
					MaxAmount:       "0x64", // 100
					MaxPricePerUnit: "0x64", // 100
				},
			},
			multiplier:  1.0,
			expectedErr: "invalid resource bounds: '' is not a valid big.Int",
		},
		{
			name: "Overflow",
			resBounds: rpc.ResourceBoundsMapping{
				L1Gas: rpc.ResourceBounds{
					MaxAmount:       "0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
					MaxPricePerUnit: "0x64", // 100
				},
				L1DataGas: rpc.ResourceBounds{
					MaxAmount:       "0x64", // 100
					MaxPricePerUnit: "0x64", // 100
				},
				L2Gas: rpc.ResourceBounds{
					MaxAmount:       "0x64", // 100
					MaxPricePerUnit: "0x64", // 100
				},
			},
			multiplier:  1.0,
			expectedErr: "can't fit in felt: 0x64000000000000000000000000000000000000000000000000000000000000000004dbc",
		},
		{
			name: "Underflow",
			resBounds: rpc.ResourceBoundsMapping{
				L1Gas: rpc.ResourceBounds{
					MaxAmount:       "-0x64",
					MaxPricePerUnit: "0x64", // 100
				},
			},
			multiplier:  1.0,
			expectedErr: "resource bounds cannot be negative, got '-0x64'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResBoundsMapToOverallFee(&tt.resBounds, tt.multiplier)

			if tt.expectedErr != "" {
				require.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr)

				return
			}
			require.NoError(t, err)

			expectedBigInt, ok := new(big.Int).SetString(tt.expectedRes, 0)
			require.True(t, ok)
			assert.Equal(t, fmt.Sprintf("%#x", expectedBigInt), got.String(), "ResBoundsMapToOverallFee() returned incorrect value")
		})
	}
}

//nolint:dupl
func TestFeeEstToResBoundsMap(t *testing.T) {
	tests := []struct {
		name          string
		feeEstimation rpc.FeeEstimation
		multiplier    float64
		expected      rpc.ResourceBoundsMapping
	}{
		{
			name: "Basic calculation with multiplier 1.0",
			feeEstimation: rpc.FeeEstimation{
				L1GasPrice:        BigIntToFelt(big.NewInt(10)),
				L1GasConsumed:     BigIntToFelt(big.NewInt(100)),
				L1DataGasPrice:    BigIntToFelt(big.NewInt(5)),
				L1DataGasConsumed: BigIntToFelt(big.NewInt(50)),
				L2GasPrice:        BigIntToFelt(big.NewInt(3)),
				L2GasConsumed:     BigIntToFelt(big.NewInt(200)),
			},
			multiplier: 1.0,
			expected: rpc.ResourceBoundsMapping{
				L1Gas: rpc.ResourceBounds{
					MaxAmount:       "0x64", // 100
					MaxPricePerUnit: "0xa",  // 10
				},
				L1DataGas: rpc.ResourceBounds{
					MaxAmount:       "0x32", // 50
					MaxPricePerUnit: "0x5",  // 5
				},
				L2Gas: rpc.ResourceBounds{
					MaxAmount:       "0xc8", // 200
					MaxPricePerUnit: "0x3",  // 3
				},
			},
		},
		{
			name: "Multiplier less than 1",
			feeEstimation: rpc.FeeEstimation{
				L1GasPrice:        internalUtils.TestHexToFelt(t, "0xabcdef1234567890abcdef1234567"),         // valid uint128
				L1GasConsumed:     internalUtils.TestHexToFelt(t, "0x8b2c3d4e5f607182"),                      // valid uint64
				L1DataGasPrice:    internalUtils.TestHexToFelt(t, "0xa2ffe1d2c3b4a5968778695a4b3c2d15"),      // valid uint128
				L1DataGasConsumed: internalUtils.TestHexToFelt(t, "0xac2b3c4d5e6f7a8b"),                      // valid uint64
				L2GasPrice:        internalUtils.TestHexToFelt(t, "0x123456789abcdef0123456789abcdabcdabcd"), // invalid uint128
				L2GasConsumed:     internalUtils.TestHexToFelt(t, "0x123456789abcdef0123456789abcdabcdabcd"), // invalid uint64
			},
			multiplier: 0.5,
			expected: rpc.ResourceBoundsMapping{
				L1Gas: rpc.ResourceBounds{
					// 10028457877064151426 * 0.5 ~= 5014228938532075713
					MaxAmount: "0x45961ea72fb038c1",
					// 55753724871440480815496793359074663 * 0.5 ~= 27876862435720240407748396679537331
					MaxPricePerUnit: "0x55e6f7891a2b3c4855e6f7891a2b3",
				},
				L1DataGas: rpc.ResourceBounds{
					// 12406075901516675723 * 0.5 ~= 6203037950758337861
					MaxAmount: "0x56159e26af37bd45",
					// 216663551256725667606984177334664047893 * 0.5 ~= 108331775628362833803492088667332023946
					MaxPricePerUnit: "0x517ff0e961da52cb43bc34ad259e168a",
				},
				L2Gas: rpc.ResourceBounds{
					// As these inputs overflow, the multiplier will be applied to the max values
					// 18446744073709551615 * 0.5 ~= 9223372036854775807
					MaxAmount: "0x7fffffffffffffff",
					// 340282366920938463463374607431768211455 * 0.5 ~= 170141183460469231731687303715884105727
					MaxPricePerUnit: "0x7fffffffffffffffffffffffffffffff",
				},
			},
		},
		{
			name: "With multiplier 1.5",
			feeEstimation: rpc.FeeEstimation{
				L1GasPrice:        BigIntToFelt(big.NewInt(10)),
				L1GasConsumed:     BigIntToFelt(big.NewInt(100)),
				L1DataGasPrice:    BigIntToFelt(big.NewInt(5)),
				L1DataGasConsumed: BigIntToFelt(big.NewInt(50)),
				L2GasPrice:        BigIntToFelt(big.NewInt(3)),
				L2GasConsumed:     BigIntToFelt(big.NewInt(200)),
			},
			multiplier: 1.5,
			expected: rpc.ResourceBoundsMapping{
				L1Gas: rpc.ResourceBounds{
					MaxAmount:       "0x96", // 150 (100 * 1.5)
					MaxPricePerUnit: "0xf",  // 15 (10 * 1.5)
				},
				L1DataGas: rpc.ResourceBounds{
					MaxAmount:       "0x4b", // 75 (50 * 1.5)
					MaxPricePerUnit: "0x7",  // 7 (5 * 1.5 = 7.5, truncated to 7)
				},
				L2Gas: rpc.ResourceBounds{
					MaxAmount:       "0x12c", // 300 (200 * 1.5)
					MaxPricePerUnit: "0x4",   // 4 (3 * 1.5 = 4.5, truncated to 4)
				},
			},
		},
		{
			name: "Very large fractional values, within the uint128 and uint64 ranges",
			feeEstimation: rpc.FeeEstimation{
				L1GasPrice:        internalUtils.TestHexToFelt(t, "0xabcdef1234567890abcdef1234567"), // 55753724871440480815496793359074663
				L1GasConsumed:     internalUtils.TestHexToFelt(t, "0x8b2c3d4e5f607182"),              // 10028457877064151426
				L1DataGasPrice:    internalUtils.TestHexToFelt(t, "0xf0e1d2c3b4a5968778695a4b3c2d1"), // 78170717918204611383717257769370321
				L1DataGasConsumed: internalUtils.TestHexToFelt(t, "0x1a2b3c4d5e6f7a8b"),              // 1885667171979197067
				L2GasPrice:        internalUtils.TestHexToFelt(t, "0x123456789abcdef0123456789abcd"), // 5907679981266292691599931071900621
				L2GasConsumed:     internalUtils.TestHexToFelt(t, "0xfedcba98765432"),                // 71737338064426034
			},
			multiplier: 1.7,
			expected: rpc.ResourceBoundsMapping{
				L1Gas: rpc.ResourceBounds{
					// 10028457877064151426 * 1.7 ~= 17048378391009056979
					MaxAmount: "0xec9801d2088a58d3",
					// 55753724871440480815496793359074663 * 1.7 ~= 94781332281448814910381786274848222
					MaxPricePerUnit: "0x12411499ef292fe035de10f5ddddde",
				},
				L1DataGas: rpc.ResourceBounds{
					// 1885667171979197067 * 1.7 ~= 3205634192364634930
					MaxAmount: "0x2c7cb35053bd8332",
					// 78170717918204611383717257769370321 * 1.7 ~= 132890220460947835880842102837167790
					MaxPricePerUnit: "0x1997fe64cb3197ce37a10a73dd46ae",
				},
				L2Gas: rpc.ResourceBounds{
					// 71737338064426034 * 1.7 ~= 121953474709524254
					MaxAmount: "0x1b1440a032f8f1e",
					// 5907679981266292691599931071900621 * 1.7 ~= 10043055968152697313366189329472992
					MaxPricePerUnit: "0x1ef293003a41145dddddddddddde0",
				},
			},
		},
		{
			name: "Zero values",
			feeEstimation: rpc.FeeEstimation{
				L1GasPrice:        BigIntToFelt(big.NewInt(0)),
				L1GasConsumed:     BigIntToFelt(big.NewInt(0)),
				L1DataGasPrice:    BigIntToFelt(big.NewInt(0)),
				L1DataGasConsumed: BigIntToFelt(big.NewInt(0)),
				L2GasPrice:        BigIntToFelt(big.NewInt(0)),
				L2GasConsumed:     BigIntToFelt(big.NewInt(0)),
			},
			multiplier: 1.0,
			expected: rpc.ResourceBoundsMapping{
				L1Gas: rpc.ResourceBounds{
					MaxAmount:       "0x0",
					MaxPricePerUnit: "0x0",
				},
				L1DataGas: rpc.ResourceBounds{
					MaxAmount:       "0x0",
					MaxPricePerUnit: "0x0",
				},
				L2Gas: rpc.ResourceBounds{
					MaxAmount:       "0x0",
					MaxPricePerUnit: "0x0",
				},
			},
		},
		{
			name: "Overflow",
			feeEstimation: rpc.FeeEstimation{
				L1GasPrice:        internalUtils.TestHexToFelt(t, "0xabcdef1234567890abcdef1234567"),         // valid uint128
				L1GasConsumed:     internalUtils.TestHexToFelt(t, "0x8b2c3d4e5f607182"),                      // valid uint64
				L1DataGasPrice:    internalUtils.TestHexToFelt(t, "0xa2ffe1d2c3b4a5968778695a4b3c2d15"),      // valid uint128
				L1DataGasConsumed: internalUtils.TestHexToFelt(t, "0xac2b3c4d5e6f7a8b"),                      // valid uint64
				L2GasPrice:        internalUtils.TestHexToFelt(t, "0x123456789abcdef0123456789abcdabcdabcd"), // invalid uint128
				L2GasConsumed:     internalUtils.TestHexToFelt(t, "0x123456789abcdef0123456789abcdabcdabcd"), // invalid uint64
			},
			multiplier: 1.7,
			expected: rpc.ResourceBoundsMapping{
				L1Gas: rpc.ResourceBounds{
					// 10028457877064151426 * 1.7 ~= 17048378391009056979
					MaxAmount: "0xec9801d2088a58d3",
					// 55753724871440480815496793359074663 * 1.7 ~= 94781332281448814910381786274848222
					MaxPricePerUnit: "0x12411499ef292fe035de10f5ddddde",
				},
				L1DataGas: rpc.ResourceBounds{
					// 12406075901516675723 * 1.7 ~= 21090329032578348178
					// This result is too large to fit in a uint64, so the function returns the max uint64 value
					MaxAmount: rpc.U64(fmt.Sprintf("%#x", maxUint64)),
					// 216663551256725667606984177334664047893 * 1.7 ~= 368328037136433625310078573378144594674
					// This result is too large to fit in a uint128, so the function returns the max uint128 value
					MaxPricePerUnit: rpc.U128(maxUint128),
				},
				L2Gas: rpc.ResourceBounds{
					// The inputs overflow, so should the output
					MaxAmount:       rpc.U64(fmt.Sprintf("%#x", maxUint64)),
					MaxPricePerUnit: rpc.U128(maxUint128),
				},
			},
		},
		{
			name: "Negative multiplier",
			feeEstimation: rpc.FeeEstimation{
				L1GasPrice:        internalUtils.TestHexToFelt(t, "0xabcdef1234567890abcdef1234567"), // 55753724871440480815496793359074663
				L1GasConsumed:     internalUtils.TestHexToFelt(t, "0x8b2c3d4e5f607182"),              // 10028457877064151426
				L1DataGasPrice:    internalUtils.TestHexToFelt(t, "0xf0e1d2c3b4a5968778695a4b3c2d1"), // 78170717918204611383717257769370321
				L1DataGasConsumed: internalUtils.TestHexToFelt(t, "0x1a2b3c4d5e6f7a8b"),              // 1885667171979197067
				L2GasPrice:        internalUtils.TestHexToFelt(t, "0x123456789abcdef0123456789abcd"), // 5907679981266292691599931071900621
				L2GasConsumed:     internalUtils.TestHexToFelt(t, "0xfedcba98765432"),                // 71737338064426034
			},
			multiplier: -1.7,
			expected: rpc.ResourceBoundsMapping{
				// when multiplier is negative, the max amount and max price per unit should be 0
				L1Gas: rpc.ResourceBounds{
					MaxAmount:       "0x0",
					MaxPricePerUnit: "0x0",
				},
				L1DataGas: rpc.ResourceBounds{
					MaxAmount:       "0x0",
					MaxPricePerUnit: "0x0",
				},
				L2Gas: rpc.ResourceBounds{
					MaxAmount:       "0x0",
					MaxPricePerUnit: "0x0",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FeeEstToResBoundsMap(tt.feeEstimation, tt.multiplier)

			// Compare each field individually for better error messages
			assert.Equal(t, tt.expected.L1Gas.MaxAmount, got.L1Gas.MaxAmount,
				"L1Gas.MaxAmount mismatch")
			assert.Equal(t, tt.expected.L1Gas.MaxPricePerUnit, got.L1Gas.MaxPricePerUnit,
				"L1Gas.MaxPricePerUnit mismatch")

			assert.Equal(t, tt.expected.L1DataGas.MaxAmount, got.L1DataGas.MaxAmount,
				"L1DataGas.MaxAmount mismatch")
			assert.Equal(t, tt.expected.L1DataGas.MaxPricePerUnit, got.L1DataGas.MaxPricePerUnit,
				"L1DataGas.MaxPricePerUnit mismatch")

			assert.Equal(t, tt.expected.L2Gas.MaxAmount, got.L2Gas.MaxAmount,
				"L2Gas.MaxAmount mismatch")
			assert.Equal(t, tt.expected.L2Gas.MaxPricePerUnit, got.L2Gas.MaxPricePerUnit,
				"L2Gas.MaxPricePerUnit mismatch")
		})
	}
}

// TestTxnOptions tests the ApplyOptions method of the TxnOptions struct,
// testing whether the method sets the default values and checks for edge cases
func TestTxnOptions(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name            string
		opts            *TxnOptions
		expectedTip     rpc.U64
		expectedVersion rpc.TransactionVersion
	}{
		{
			name:            "Default values",
			opts:            nil,
			expectedTip:     "0x0",
			expectedVersion: rpc.TransactionV3,
		},
		{
			name: "WithQueryBitVersion true",
			opts: &TxnOptions{
				UseQueryBit: true,
			},
			expectedTip:     "0x0",
			expectedVersion: rpc.TransactionV3WithQueryBit,
		},
		{
			name: "Tip set",
			opts: &TxnOptions{
				Tip: "0x1234567890",
			},
			expectedTip:     "0x1234567890",
			expectedVersion: rpc.TransactionV3,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			opts := tt.opts

			assert.Equal(t, tt.expectedTip, opts.SafeTip())
			assert.Equal(t, tt.expectedVersion, opts.TxnVersion())
		})
	}
}
