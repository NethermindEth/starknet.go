package utils

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/NethermindEth/starknet.go/contracts"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResBoundsMapToOverallFee(t *testing.T) {
	t.Parallel()

	zeroTip := rpc.U64("0x0")
	tests := []struct {
		name        string
		resBounds   rpc.ResourceBoundsMapping
		multiplier  float64
		tip         rpc.U64
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
			tip:        zeroTip,
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
			tip:         zeroTip,
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
			tip:        zeroTip,
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
			tip:         zeroTip,
			expectedErr: "multiplier must be greater than 0",
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
			tip:        zeroTip,
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
			tip:         zeroTip,
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
			multiplier:  1.0,
			tip:         zeroTip,
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
			tip:         zeroTip,
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
			tip:         zeroTip,
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
			tip:         zeroTip,
			expectedErr: "resource bounds cannot be negative, got '-0x64'",
		},
		{
			name: "Real values",
			resBounds: rpc.ResourceBoundsMapping{
				L1Gas: rpc.ResourceBounds{
					MaxAmount:       "0x0",
					MaxPricePerUnit: "0x1925a36320fc",
				},
				L1DataGas: rpc.ResourceBounds{
					MaxAmount:       "0x80",   // 128
					MaxPricePerUnit: "0x6c01", // 27649
				},
				L2Gas: rpc.ResourceBounds{
					MaxAmount:       "0xc25b1",    // 796081
					MaxPricePerUnit: "0xb2d05e00", // 3000000000
				},
			},
			multiplier: 1.0,
			tip:        rpc.U64("0x1000000"), // 16777216
			// Expected: 0 + (128*27649) + ((3000000000+16777216)*796081) =
			// 0 + 3539072 + 2401599022890496 = 2401599026429568
			expectedRes: "2401599026429568", // 0x8883dd8dcfe80
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ResBoundsMapToOverallFee(&tt.resBounds, tt.multiplier, tt.tip)

			if tt.expectedErr != "" {
				require.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr)

				return
			}
			require.NoError(t, err)

			expectedBigInt, ok := new(big.Int).SetString(tt.expectedRes, 0)
			require.True(t, ok)
			assert.Equal(
				t,
				fmt.Sprintf("%#x", expectedBigInt),
				got.String(),
				"ResBoundsMapToOverallFee() returned incorrect value",
			)
		})
	}
}

//nolint:tparallel // Run sequentially to avoid race conditions with the `tests` variable.
func TestFeeEstToResBoundsMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		feeEstimation rpc.FeeEstimation
		multiplier    float64
		expected      rpc.ResourceBoundsMapping
		feeLimit      FeeLimits // Only used in the `CustomFeeEstToResBoundsMap` test.
	}{
		{
			name: "Basic calculation with multiplier 1.0",
			feeEstimation: rpc.FeeEstimation{
				FeeEstimationCommon: rpc.FeeEstimationCommon{
					L1GasPrice:        BigIntToFelt(big.NewInt(10)),
					L1GasConsumed:     BigIntToFelt(big.NewInt(100)),
					L1DataGasPrice:    BigIntToFelt(big.NewInt(5)),
					L1DataGasConsumed: BigIntToFelt(big.NewInt(50)),
					L2GasPrice:        BigIntToFelt(big.NewInt(3)),
					L2GasConsumed:     BigIntToFelt(big.NewInt(200)),
				},
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
				FeeEstimationCommon: rpc.FeeEstimationCommon{
					L1GasPrice: internalUtils.TestHexToFelt(
						t,
						"0xabcdef1234567890abcdef1234567",
					), // valid uint128
					L1GasConsumed: internalUtils.TestHexToFelt(
						t,
						"0x8b2c3d4e5f607182",
					), // valid uint64
					L1DataGasPrice: internalUtils.TestHexToFelt(
						t,
						"0xa2ffe1d2c3b4a5968778695a4b3c2d15",
					), // valid uint128
					L1DataGasConsumed: internalUtils.TestHexToFelt(
						t,
						"0xac2b3c4d5e6f7a8b",
					), // valid uint64
					L2GasPrice: internalUtils.TestHexToFelt(
						t,
						"0x123456789abcdef0123456789abcdabcdabcd",
					), // invalid uint128
					L2GasConsumed: internalUtils.TestHexToFelt(
						t,
						"0x2faf0800",
					), // valid uint64, within L2 gas amount limit
				},
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
					// 800_000_000 * 0.5 ~= 400_000_000
					MaxAmount: "0x17d78400",
					// As the output overflows, the max value is used
					// 340282366920938463463374607431768211455 * 0.5 ~= 170141183460469231731687303715884105727
					MaxPricePerUnit: rpc.U128(maxUint128),
				},
			},
		},
		{
			name: "With multiplier 1.5",
			feeEstimation: rpc.FeeEstimation{
				FeeEstimationCommon: rpc.FeeEstimationCommon{
					L1GasPrice:        BigIntToFelt(big.NewInt(10)),
					L1GasConsumed:     BigIntToFelt(big.NewInt(100)),
					L1DataGasPrice:    BigIntToFelt(big.NewInt(5)),
					L1DataGasConsumed: BigIntToFelt(big.NewInt(50)),
					L2GasPrice:        BigIntToFelt(big.NewInt(3)),
					L2GasConsumed:     BigIntToFelt(big.NewInt(200)),
				},
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
				FeeEstimationCommon: rpc.FeeEstimationCommon{
					L1GasPrice: internalUtils.TestHexToFelt(
						t,
						"0xabcdef1234567890abcdef1234567",
					), // 55753724871440480815496793359074663
					L1GasConsumed: internalUtils.TestHexToFelt(
						t,
						"0x8b2c3d4e5f607182",
					), // 10028457877064151426
					L1DataGasPrice: internalUtils.TestHexToFelt(
						t,
						"0xf0e1d2c3b4a5968778695a4b3c2d1",
					), // 78170717918204611383717257769370321
					L1DataGasConsumed: internalUtils.TestHexToFelt(
						t,
						"0x1a2b3c4d5e6f7a8b",
					), // 1885667171979197067
					L2GasPrice: internalUtils.TestHexToFelt(
						t,
						"0x123456789abcdef0123456789abcd",
					), // 5907679981266292691599931071900621
					L2GasConsumed: internalUtils.TestHexToFelt(
						t,
						"0xfedcba98765432",
					), // 71737338064426034
				},
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
					// The result is bigger than the max L2 gas amount limit, so the function
					// should return the max L2 gas amount instead
					MaxAmount: rpc.U64(maxL2GasAmount),
					// 5907679981266292691599931071900621 * 1.7 ~= 10043055968152697313366189329472992
					MaxPricePerUnit: "0x1ef293003a41145dddddddddddde0",
				},
			},
			feeLimit: starknetLimits, // For the CustomFeeEstToResBoundsMap test.
		},
		{
			name: "Zero values",
			feeEstimation: rpc.FeeEstimation{
				FeeEstimationCommon: rpc.FeeEstimationCommon{
					L1GasPrice:        BigIntToFelt(big.NewInt(0)),
					L1GasConsumed:     BigIntToFelt(big.NewInt(0)),
					L1DataGasPrice:    BigIntToFelt(big.NewInt(0)),
					L1DataGasConsumed: BigIntToFelt(big.NewInt(0)),
					L2GasPrice:        BigIntToFelt(big.NewInt(0)),
					L2GasConsumed:     BigIntToFelt(big.NewInt(0)),
				},
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
				FeeEstimationCommon: rpc.FeeEstimationCommon{
					L1GasPrice: internalUtils.TestHexToFelt(
						t,
						"0xabcdef1234567890abcdef1234567",
					), // valid uint128
					L1GasConsumed: internalUtils.TestHexToFelt(
						t,
						"0x8b2c3d4e5f607182",
					), // valid uint64
					L1DataGasPrice: internalUtils.TestHexToFelt(
						t,
						"0xa2ffe1d2c3b4a5968778695a4b3c2d15",
					), // valid uint128
					L1DataGasConsumed: internalUtils.TestHexToFelt(
						t,
						"0xac2b3c4d5e6f7a8b",
					), // valid uint64
					L2GasPrice: internalUtils.TestHexToFelt(
						t,
						"0x123456789abcdef0123456789abcdabcdabcd",
					), // invalid uint128
					L2GasConsumed: internalUtils.TestHexToFelt(
						t,
						"0x123456789abcdef0123456789abcdabcdabcd",
					), // invalid uint64
				},
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
					MaxAmount: rpc.U64(maxUint64),
					// 216663551256725667606984177334664047893 * 1.7 ~= 368328037136433625310078573378144594674
					// This result is too large to fit in a uint128, so the function returns the max uint128 value
					MaxPricePerUnit: rpc.U128(maxUint128),
				},
				// The inputs overflow, so the output should be the max values
				L2Gas: rpc.ResourceBounds{
					// Default max L2 gas amount limit
					MaxAmount:       rpc.U64(maxL2GasAmount),
					MaxPricePerUnit: rpc.U128(maxUint128),
				},
			},
			feeLimit: starknetLimits, // For the CustomFeeEstToResBoundsMap test.
		},
		{
			name: "Negative multiplier",
			feeEstimation: rpc.FeeEstimation{
				FeeEstimationCommon: rpc.FeeEstimationCommon{
					L1GasPrice: internalUtils.TestHexToFelt(
						t,
						"0xabcdef1234567890abcdef1234567",
					), // 55753724871440480815496793359074663
					L1GasConsumed: internalUtils.TestHexToFelt(
						t,
						"0x8b2c3d4e5f607182",
					), // 10028457877064151426
					L1DataGasPrice: internalUtils.TestHexToFelt(
						t,
						"0xf0e1d2c3b4a5968778695a4b3c2d1",
					), // 78170717918204611383717257769370321
					L1DataGasConsumed: internalUtils.TestHexToFelt(
						t,
						"0x1a2b3c4d5e6f7a8b",
					), // 1885667171979197067
					L2GasPrice: internalUtils.TestHexToFelt(
						t,
						"0x123456789abcdef0123456789abcd",
					), // 5907679981266292691599931071900621
					L2GasConsumed: internalUtils.TestHexToFelt(
						t,
						"0xfedcba98765432",
					), // 71737338064426034
				},
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

	t.Run("Test FeeEstToResBoundsMap", func(t *testing.T) {
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
				assert.Equal(
					t,
					tt.expected.L1DataGas.MaxPricePerUnit,
					got.L1DataGas.MaxPricePerUnit,
					"L1DataGas.MaxPricePerUnit mismatch",
				)

				assert.Equal(t, tt.expected.L2Gas.MaxAmount, got.L2Gas.MaxAmount,
					"L2Gas.MaxAmount mismatch")
				assert.Equal(t, tt.expected.L2Gas.MaxPricePerUnit, got.L2Gas.MaxPricePerUnit,
					"L2Gas.MaxPricePerUnit mismatch")
			})
		}
	})

	maxUint128BigInt, ok := new(big.Int).SetString(maxUint128, 0)
	require.True(t, ok)

	// All previous tests + new test with custom fee limits.
	tests = append(tests, []struct {
		name          string
		feeEstimation rpc.FeeEstimation
		multiplier    float64
		expected      rpc.ResourceBoundsMapping
		feeLimit      FeeLimits
	}{
		{
			name: "With fee limit + multiplier 1.5",
			feeEstimation: rpc.FeeEstimation{
				FeeEstimationCommon: rpc.FeeEstimationCommon{
					L1GasPrice:        BigIntToFelt(big.NewInt(1_000_000)),
					L1GasConsumed:     BigIntToFelt(big.NewInt(500_000)),
					L1DataGasPrice:    BigIntToFelt(big.NewInt(1_000_000_000)),
					L1DataGasConsumed: BigIntToFelt(big.NewInt(1_000_000)),
					L2GasPrice:        BigIntToFelt(big.NewInt(800_000)),
					L2GasConsumed:     BigIntToFelt(big.NewInt(200_000_000)),
				},
			},
			feeLimit: FeeLimits{
				L1GasPriceLimit:      rpc.U128("0x124f80"),    // 1_200_000
				L1GasAmountLimit:     rpc.U64("0xf4240"),      // 1_000_000
				L1DataGasPriceLimit:  rpc.U128("0xf4240"),     // 1_000_000
				L1DataGasAmountLimit: rpc.U64("0xe8d4a51000"), // 1_000_000_000_000
				L2GasPriceLimit:      rpc.U128("0x7a120"),     // 500_000
				L2GasAmountLimit:     rpc.U64("0xe8d4a51000"), // 1_000_000_000_000
			},
			multiplier: 1.5,
			expected: rpc.ResourceBoundsMapping{
				L1Gas: rpc.ResourceBounds{
					// 1_000_000 * 1.5 = 1_500_000_
					MaxPricePerUnit: "0x124f80", // 1_200_000, capped by the limit
					// 500_000 * 1.5 = 750_000_
					MaxAmount: "0xb71b0", // 750_000, within the limit
				},
				L1DataGas: rpc.ResourceBounds{
					// 1_000_000_000 * 1.5 = 1_500_000_000_
					MaxPricePerUnit: "0xf4240", // 1_000_000, capped by the limit
					// 1_000_000 * 1.5 = 1_500_000_
					MaxAmount: "0x16e360", // 1_500_000, within the limit
				},
				L2Gas: rpc.ResourceBounds{
					// 800_000 * 1.5 = 1_200_000_
					MaxPricePerUnit: "0x7a120", // 500_000, capped by the limit
					// 200_000_000 * 1.5 = 300_000_000_
					MaxAmount: "0x11e1a300", // 300_000_000, within the limit
				},
			},
		},
		{
			name: "overflows, with only one limit set",
			feeEstimation: rpc.FeeEstimation{
				FeeEstimationCommon: rpc.FeeEstimationCommon{
					L1GasPrice:        BigIntToFelt(maxUint128BigInt),
					L1GasConsumed:     BigIntToFelt(maxUint128BigInt),
					L1DataGasPrice:    BigIntToFelt(maxUint128BigInt),
					L1DataGasConsumed: BigIntToFelt(maxUint128BigInt),
					L2GasPrice:        BigIntToFelt(maxUint128BigInt),
					L2GasConsumed:     BigIntToFelt(maxUint128BigInt),
				},
			},
			feeLimit: FeeLimits{
				L1GasPriceLimit: rpc.U128("0xf4240"), // 1_000_000
			},
			multiplier: 100,
			// All outputs but the L1Gas.MaxPricePerUnit should be the max U128 and U64 values.
			expected: rpc.ResourceBoundsMapping{
				L1Gas: rpc.ResourceBounds{
					MaxPricePerUnit: rpc.U128("0xf4240"), // The same as the limit.
					MaxAmount:       rpc.U64(maxUint64),
				},
				L1DataGas: rpc.ResourceBounds{
					MaxPricePerUnit: rpc.U128(maxUint128),
					MaxAmount:       rpc.U64(maxUint64),
				},
				L2Gas: rpc.ResourceBounds{
					MaxPricePerUnit: rpc.U128(maxUint128),
					MaxAmount:       rpc.U64(maxUint64),
				},
			},
		},
	}...)

	t.Run("Test CustomFeeEstToResBoundsMap", func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := CustomFeeEstToResBoundsMap(tt.feeEstimation, tt.multiplier, &tt.feeLimit)

				// Compare each field individually for better error messages
				assert.Equal(t, tt.expected.L1Gas.MaxAmount, got.L1Gas.MaxAmount,
					"L1Gas.MaxAmount mismatch")
				assert.Equal(t, tt.expected.L1Gas.MaxPricePerUnit, got.L1Gas.MaxPricePerUnit,
					"L1Gas.MaxPricePerUnit mismatch")

				assert.Equal(t, tt.expected.L1DataGas.MaxAmount, got.L1DataGas.MaxAmount,
					"L1DataGas.MaxAmount mismatch")
				assert.Equal(
					t,
					tt.expected.L1DataGas.MaxPricePerUnit,
					got.L1DataGas.MaxPricePerUnit,
					"L1DataGas.MaxPricePerUnit mismatch",
				)

				assert.Equal(t, tt.expected.L2Gas.MaxAmount, got.L2Gas.MaxAmount,
					"L2Gas.MaxAmount mismatch")
				assert.Equal(t, tt.expected.L2Gas.MaxPricePerUnit, got.L2Gas.MaxPricePerUnit,
					"L2Gas.MaxPricePerUnit mismatch")
			})
		}
	})
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

			// same behaviour as in Build...Txn functions
			if opts == nil {
				opts = new(TxnOptions)
			}

			assert.Equal(t, tt.expectedTip, opts.SafeTip())
			assert.Equal(t, tt.expectedVersion, opts.TxnVersion())
		})
	}
}

// TestBuildDeclareTxnWithBlake2sHash tests the BuildDeclareTxn function with the
// 'UseBlake2sHash' option. It checks whether the compiled class hash is calculated
// as expected when the option is set to true or false.
func TestBuildDeclareTxnWithBlake2sHash(t *testing.T) {
	t.Parallel()

	casmClass := *internalUtils.TestUnmarshalJSONFileToType[contracts.CasmClass](
		t,
		"../hash/testData/contracts_v2_HelloStarknet.compiled_contract_class.json",
		"",
	)

	testCases := []struct {
		name                      string
		opts                      *TxnOptions
		expectedCompiledClassHash string
	}{
		// Values taken from the hash/hash_test.go TestClassHashes test.
		{
			name: "UseBlake2sHash true",
			opts: &TxnOptions{
				UseBlake2sHash: true,
			},
			expectedCompiledClassHash: "0x23c2091df2547f77185ba592b06ee2e897b0c2a70f968521a6a24fc5bfc1b1e",
		},
		{
			name: "UseBlake2sHash false",
			opts: &TxnOptions{
				UseBlake2sHash: false,
			},
			expectedCompiledClassHash: "0x6ff9f7df06da94198ee535f41b214dce0b8bafbdb45e6c6b09d4b3b693b1f17",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			declareTxn, err := BuildDeclareTxn(
				nil,
				&casmClass,
				nil,
				nil,
				nil,
				test.opts,
			)
			require.NoError(t, err)
			require.NotNil(t, declareTxn)

			assert.Equal(
				t,
				test.expectedCompiledClassHash,
				declareTxn.CompiledClassHash.String(),
			)
		})
	}
}
