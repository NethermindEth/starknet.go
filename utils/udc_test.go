package utils

import (
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildUDCCalldata(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                  string
		classHash             *felt.Felt
		constructorCalldata   []*felt.Felt
		opts                  *UDCOptions
		expectedUDCAddress    *felt.Felt
		expectedFunctionName  string
		expectedCallDataLen   int
		expectedError         string
		checkCallDataContents bool
	}{
		{
			name:                  "UDC Cairo V0 with default options",
			classHash:             internalUtils.RANDOM_FELT,
			constructorCalldata:   []*felt.Felt{new(felt.Felt).SetUint64(100), new(felt.Felt).SetUint64(200)},
			opts:                  nil,
			expectedUDCAddress:    udcAddressCairoV0,
			expectedFunctionName:  "deployContract",
			expectedCallDataLen:   6, // classHash + salt + originInd + calldataLen + 2 constructor args
			checkCallDataContents: true,
		},
		{
			name:                  "UDC Cairo V0 with custom salt",
			classHash:             internalUtils.RANDOM_FELT,
			constructorCalldata:   []*felt.Felt{new(felt.Felt).SetUint64(100)},
			opts:                  &UDCOptions{Salt: new(felt.Felt).SetUint64(999)},
			expectedUDCAddress:    udcAddressCairoV0,
			expectedFunctionName:  "deployContract",
			expectedCallDataLen:   5, // classHash + salt + originInd + calldataLen + 1 constructor arg
			checkCallDataContents: true,
		},
		{
			name:                  "UDC Cairo V0 with origin independent",
			classHash:             internalUtils.RANDOM_FELT,
			constructorCalldata:   []*felt.Felt{},
			opts:                  &UDCOptions{OriginIndependent: true},
			expectedUDCAddress:    udcAddressCairoV0,
			expectedFunctionName:  "deployContract",
			expectedCallDataLen:   4, // classHash + salt + originInd + calldataLen
			checkCallDataContents: true,
		},
		{
			name:                  "UDC Cairo V2 with default options",
			classHash:             internalUtils.RANDOM_FELT,
			constructorCalldata:   []*felt.Felt{new(felt.Felt).SetUint64(100), new(felt.Felt).SetUint64(200)},
			opts:                  &UDCOptions{UDCVersion: UDCCairoV2},
			expectedUDCAddress:    udcAddressCairoV2,
			expectedFunctionName:  "deploy_contract",
			expectedCallDataLen:   5, // classHash + salt + originInd + 2 constructor args
			checkCallDataContents: true,
		},
		{
			name:                  "UDC Cairo V2 with origin independent",
			classHash:             internalUtils.RANDOM_FELT,
			constructorCalldata:   []*felt.Felt{new(felt.Felt).SetUint64(100)},
			opts:                  &UDCOptions{UDCVersion: UDCCairoV2, OriginIndependent: true},
			expectedUDCAddress:    udcAddressCairoV2,
			expectedFunctionName:  "deploy_contract",
			expectedCallDataLen:   4, // classHash + salt + originInd + 1 constructor arg
			checkCallDataContents: true,
		},
		{
			name:                  "UDC Cairo V2 with custom salt and origin independent",
			classHash:             internalUtils.RANDOM_FELT,
			constructorCalldata:   []*felt.Felt{},
			opts:                  &UDCOptions{UDCVersion: UDCCairoV2, Salt: new(felt.Felt).SetUint64(888), OriginIndependent: true},
			expectedUDCAddress:    udcAddressCairoV2,
			expectedFunctionName:  "deploy_contract",
			expectedCallDataLen:   3, // classHash + salt + originInd
			checkCallDataContents: true,
		},
		{
			name:                  "Empty constructor calldata",
			classHash:             internalUtils.RANDOM_FELT,
			constructorCalldata:   []*felt.Felt{},
			opts:                  nil,
			expectedUDCAddress:    udcAddressCairoV0,
			expectedFunctionName:  "deployContract",
			expectedCallDataLen:   4, // classHash + salt + originInd + calldataLen
			checkCallDataContents: true,
		},
		{
			name:                "Nil classHash",
			classHash:           nil,
			constructorCalldata: []*felt.Felt{new(felt.Felt).SetUint64(100)},
			opts:                nil,
			expectedError:       "classHash not provided",
		},
		{
			name:                "Invalid UDC version",
			classHash:           internalUtils.RANDOM_FELT,
			constructorCalldata: []*felt.Felt{new(felt.Felt).SetUint64(100)},
			opts:                &UDCOptions{UDCVersion: 999}, // Invalid version
			expectedError:       "invalid UDC version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := BuildUDCCalldata(tt.classHash, tt.constructorCalldata, tt.opts)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)

			// Check UDC address
			assert.Equal(t, tt.expectedUDCAddress, result.ContractAddress)

			// Check function name
			assert.Equal(t, tt.expectedFunctionName, result.FunctionName)

			// Check call data length
			assert.Equal(t, tt.expectedCallDataLen, len(result.CallData))

			// Check call data contents if requested
			if tt.checkCallDataContents {
				checkCallDataContents(t, result, tt.classHash, tt.constructorCalldata, tt.opts)
			}
		})
	}
}

func checkCallDataContents(t *testing.T, result *rpc.InvokeFunctionCall, classHash *felt.Felt, constructorCalldata []*felt.Felt, opts *UDCOptions) {
	t.Helper()
	callData := result.CallData
	require.NotEmpty(t, callData)

	// First element should always be the class hash
	assert.Equal(t, classHash, callData[0])

	// Second element should be the salt (either provided or random)
	require.NotNil(t, callData[1])
	if opts != nil && opts.Salt != nil {
		assert.Equal(t, opts.Salt, callData[1])
	} else {
		// If no salt provided, it should be a random value (not zero)
		assert.NotEqual(t, &felt.Zero, callData[1])
	}

	// Check the rest based on UDC version
	if opts != nil && opts.UDCVersion == UDCCairoV0 {

		// Cairo V0: [classHash, salt, originInd, calldataLen, ...constructorCalldata]
		originInd := callData[2]
		expectedOriginInd := new(felt.Felt).SetUint64(1)
		if opts != nil && opts.OriginIndependent {
			expectedOriginInd.SetUint64(0)
		}
		assert.Equal(t, expectedOriginInd, originInd)

		// Check calldata length
		calldataLen := callData[3]
		expectedCalldataLen := new(felt.Felt).SetUint64(uint64(len(constructorCalldata)))
		assert.Equal(t, expectedCalldataLen, calldataLen)

		// Check constructor calldata
		constructorStart := 4
		for i, expected := range constructorCalldata {
			assert.Equal(t, expected, callData[constructorStart+i])
		}
	} else {
		// Cairo V2: [classHash, salt, originInd, ...constructorCalldata]
		originInd := callData[2]
		expectedOriginInd := new(felt.Felt).SetUint64(0)
		if opts.OriginIndependent {
			expectedOriginInd.SetUint64(1)
		}
		assert.Equal(t, expectedOriginInd, originInd)

		// Check constructor calldata
		constructorStart := 3
		for i, expected := range constructorCalldata {
			assert.Equal(t, expected, callData[constructorStart+i])
		}
	}
}

func TestBuildUDCCalldata_UDCAddresses(t *testing.T) {
	t.Parallel()
	// Test that the UDC addresses are correctly set
	classHash := new(felt.Felt).SetUint64(12345)
	constructorCalldata := []*felt.Felt{new(felt.Felt).SetUint64(100)}

	// Test Cairo V0 address
	result, err := BuildUDCCalldata(classHash, constructorCalldata, &UDCOptions{UDCVersion: UDCCairoV0})
	require.NoError(t, err)
	assert.Equal(t, udcAddressCairoV0, result.ContractAddress)

	// Test Cairo V2 address
	result, err = BuildUDCCalldata(classHash, constructorCalldata, &UDCOptions{UDCVersion: UDCCairoV2})
	require.NoError(t, err)
	assert.Equal(t, udcAddressCairoV2, result.ContractAddress)
}

func TestBuildUDCCalldata_OriginIndependent(t *testing.T) {
	t.Parallel()
	classHash := new(felt.Felt).SetUint64(12345)
	constructorCalldata := []*felt.Felt{new(felt.Felt).SetUint64(100)}

	// **** Test Cairo V0 ****
	// origin independent
	result, err := BuildUDCCalldata(classHash, constructorCalldata, &UDCOptions{UDCVersion: UDCCairoV0, OriginIndependent: true})
	require.NoError(t, err)
	// Cairo V0: `unique` should be 0 (false) when OriginIndependent is true
	assert.Equal(t, new(felt.Felt).SetUint64(0), result.CallData[2])

	// not origin independent
	result, err = BuildUDCCalldata(classHash, constructorCalldata, &UDCOptions{UDCVersion: UDCCairoV0, OriginIndependent: false})
	require.NoError(t, err)
	// Cairo V0: `unique` should be 1 (true) when OriginIndependent is false
	assert.Equal(t, new(felt.Felt).SetUint64(1), result.CallData[2])

	// **** Test Cairo V2 ****
	// origin independent
	result, err = BuildUDCCalldata(classHash, constructorCalldata, &UDCOptions{UDCVersion: UDCCairoV2, OriginIndependent: true})
	require.NoError(t, err)
	// Cairo V2: `from_zero` should be 1 (true) when OriginIndependent is true
	assert.Equal(t, new(felt.Felt).SetUint64(1), result.CallData[2])

	// not origin independent
	result, err = BuildUDCCalldata(classHash, constructorCalldata, &UDCOptions{UDCVersion: UDCCairoV2, OriginIndependent: false})
	require.NoError(t, err)
	// Cairo V2: `from_zero` should be 0 (false) when OriginIndependent is false
	assert.Equal(t, new(felt.Felt).SetUint64(0), result.CallData[2])
}

func TestBuildUDCCalldata_Salt(t *testing.T) {
	t.Parallel()
	t.Run("Random salt", func(t *testing.T) {
		t.Parallel()
		classHash := new(felt.Felt).SetUint64(12345)
		constructorCalldata := []*felt.Felt{new(felt.Felt).SetUint64(100)}

		// Test that when no salt is provided, a random salt is generated
		result1, err := BuildUDCCalldata(classHash, constructorCalldata, nil)
		require.NoError(t, err)

		result2, err := BuildUDCCalldata(classHash, constructorCalldata, nil)
		require.NoError(t, err)

		// The salts should be different (random)
		assert.NotEqual(t, result1.CallData[1], result2.CallData[1])

		// The salts should not be zero
		assert.NotEqual(t, &felt.Zero, result1.CallData[1])
		assert.NotEqual(t, &felt.Zero, result2.CallData[1])
	})

	t.Run("Custom salt", func(t *testing.T) {
		t.Parallel()
		classHash := new(felt.Felt).SetUint64(12345)
		constructorCalldata := []*felt.Felt{new(felt.Felt).SetUint64(100)}
		customSalt := new(felt.Felt).SetUint64(999)

		// Test that when a custom salt is provided, it's used
		result, err := BuildUDCCalldata(classHash, constructorCalldata, &UDCOptions{Salt: customSalt})
		require.NoError(t, err)

		assert.Equal(t, customSalt, result.CallData[1])
	})
}

func TestBuildUDCCalldata_LargeConstructorCalldata(t *testing.T) {
	t.Parallel()
	classHash := new(felt.Felt).SetUint64(12345)

	// Create a large constructor calldata
	largeCalldata := make([]*felt.Felt, 100)
	for i := range largeCalldata {
		largeCalldata[i] = new(felt.Felt).SetUint64(uint64(i))
	}

	result, err := BuildUDCCalldata(classHash, largeCalldata, &UDCOptions{UDCVersion: UDCCairoV0})
	require.NoError(t, err)

	// Check that all constructor calldata is included
	expectedLen := 4 + len(largeCalldata) // classHash + salt + originInd + calldataLen + constructor args
	assert.Equal(t, expectedLen, len(result.CallData))

	// Check calldata length field
	calldataLen := result.CallData[3]
	expectedCalldataLen := new(felt.Felt).SetUint64(uint64(len(largeCalldata)))
	assert.Equal(t, expectedCalldataLen, calldataLen)

	// Check that all constructor arguments are present
	for i, expected := range largeCalldata {
		assert.Equal(t, expected, result.CallData[4+i])
	}
}
