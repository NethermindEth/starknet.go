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
		expectedError         error
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
			expectedError:       errClassHashNotProvided,
		},
		{
			name:                "Invalid UDC version",
			classHash:           internalUtils.RANDOM_FELT,
			constructorCalldata: []*felt.Felt{new(felt.Felt).SetUint64(100)},
			opts:                &UDCOptions{UDCVersion: 999}, // Invalid version
			expectedError:       errInvalidUDCVersion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := BuildUDCCalldata(tt.classHash, tt.constructorCalldata, tt.opts)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
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

func checkCallDataContents(t *testing.T, result rpc.InvokeFunctionCall, classHash *felt.Felt, constructorCalldata []*felt.Felt, opts *UDCOptions) {
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
	if opts != nil && opts.UDCVersion == UDCCairoV2 {
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

	} else {
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

func TestPrecomputeAddressForUDC(t *testing.T) {
	t.Parallel()

	type TestSet struct {
		name                string
		classHash           *felt.Felt
		salt                *felt.Felt
		originAccAddress    *felt.Felt
		udcVersion          UDCVersion
		constructorCalldata []*felt.Felt
		expectedAddress     *felt.Felt
	}

	tests := []TestSet{
		{ //https://sepolia.voyager.online/tx/0x36f1ae1379dcea5e14938f6e9b5189d5ddb34e8b2eff404d5886e7a7e1ebb48
			name:             "UDCCairoV0: Origin-independent 0x486002",
			originAccAddress: nil,
			classHash:        internalUtils.TestHexToFelt(t, "0x486002e9e1d1fd7f07852663c0e476853e68c1aaa6afaeaa58cce954e3cf7cf"),
			salt:             internalUtils.TestHexToFelt(t, "0x23fcbf2e9a467d3081182b138fc605f8743610c4327edf4f24104f478143f66"),
			constructorCalldata: []*felt.Felt{
				internalUtils.TestHexToFelt(t, "0x23fcbf2e9a467d3081182b138fc605f8743610c4327edf4f24104f478143f66"),
			},
			udcVersion:      UDCCairoV0,
			expectedAddress: internalUtils.TestHexToFelt(t, "0x77a249aefd797d3ce41ec9d434c26ca7fe929a6ec3458bcec97cc456741fa3f"),
		},
		{ //https://sepolia.voyager.online/tx/0x3c79058b1a7a8904d4a5ade66ae87a1dcec6a568a5193a9f95f5c9646d383f3
			name:             "UDCCairoV0: Origin-independent 0x54328a",
			originAccAddress: nil,
			classHash:        internalUtils.TestHexToFelt(t, "0x54328a1075b8820eb43caf0caa233923148c983742402dcfc38541dd843d01a"),
			salt:             internalUtils.TestHexToFelt(t, "0x52a06f2789d5229e87d5d0dc0d933d0e1e83366306ae069192ca9a3e4e6881f"),
			constructorCalldata: []*felt.Felt{
				internalUtils.TestHexToFelt(t, "0x546f6b656e"),
				internalUtils.TestHexToFelt(t, "0x4552433230"),
				internalUtils.TestHexToFelt(t, "0xf9e998b2853e6d01f3ae3c598c754c1b9a7bd398fec7657de022f3b778679"),
			},
			udcVersion:      UDCCairoV0,
			expectedAddress: internalUtils.TestHexToFelt(t, "0x48634f9843983eeb06b47bf7f5d156a55a1d297e958da1c86427f9ce077425b"),
		},
		{ //https://sepolia.voyager.online/tx/0x4dadb8b32b286acd11ee6698a71206f274f1cf93fdf602f6fcf2d376c197ca4
			name:             "UDCCairoV0: Origin-dependent 0x54328a",
			originAccAddress: internalUtils.TestHexToFelt(t, "0x000f9e998b2853e6d01f3ae3c598c754c1b9a7bd398fec7657de022f3b778679"),
			classHash:        internalUtils.TestHexToFelt(t, "0x54328a1075b8820eb43caf0caa233923148c983742402dcfc38541dd843d01a"),
			salt:             internalUtils.TestHexToFelt(t, "0xb2334ec640982eb272cbfa72f4f9d32769e77166460c620ac950c8a4d94606"),
			constructorCalldata: []*felt.Felt{
				internalUtils.TestHexToFelt(t, "0x546f6b656e"),
				internalUtils.TestHexToFelt(t, "0x4552433230"),
				internalUtils.TestHexToFelt(t, "0x7d1f349d4d1c93d3e95bf584fd3f806fa61f4c72aa9b42ae624ef25470da0c6"),
			},
			udcVersion:      UDCCairoV0,
			expectedAddress: internalUtils.TestHexToFelt(t, "0x770f3b98c5a23250bc237dbeeb0a9385621f99c04a1fe3842d899264f4a1268"),
		},
		{ //https://sepolia.voyager.online/tx/0x17ffdf141abc5e5db10b7ec0f69a5f099e70390c0844e281cdd7877c7c98a54
			name:                "UDCCairoV0: Origin-dependent 0x387edd",
			originAccAddress:    internalUtils.TestHexToFelt(t, "0x02d54b7dc47eafa80f8e451cf39e7601f51fef6f1bfe5cea44ff12fa563e5457"),
			classHash:           internalUtils.TestHexToFelt(t, "0x387edd4804deba7af741953fdf64189468f37593a66b618d00d2476be3168f8"),
			salt:                internalUtils.TestHexToFelt(t, "0x18809706355007129cd38d1719447f454b7ff081c38817c53d9e2951d185243"),
			constructorCalldata: []*felt.Felt{},
			udcVersion:          UDCCairoV0,
			expectedAddress:     internalUtils.TestHexToFelt(t, "0x43267890ad2798db4a6a5374ee361cb6f3669facad58e9e58ea88c078c567bb"),
		},
		{ //https://sepolia.voyager.online/tx/0x1683dd3917260a8b1572cf703b1525ea9bd6501cd76a4c20d6e4aa316ddff86
			name:                "UDCCairoV2: Origin-independent 0x387edd",
			originAccAddress:    nil,
			classHash:           internalUtils.TestHexToFelt(t, "0x387edd4804deba7af741953fdf64189468f37593a66b618d00d2476be3168f8"),
			salt:                internalUtils.TestHexToFelt(t, "0x357f0d6bf3f6931e6a0d579ec8936510dad2dbabb8593a3dda16c4cd8fe3dc6"),
			constructorCalldata: []*felt.Felt{},
			udcVersion:          UDCCairoV2,
			expectedAddress:     internalUtils.TestHexToFelt(t, "0x2543d193907205cff1ba8818002fb0b7477493e093c7ea605516a21e838648e"),
		},
		{ //https://sepolia.voyager.online/tx/0x5c76c4f41451f64276571d5f225ab203e9de19d1749372e548a6c5e5df775af
			name:                "UDCCairoV2: Origin-dependent 0x387edd",
			originAccAddress:    internalUtils.TestHexToFelt(t, "0x02d54b7dc47eafa80f8e451cf39e7601f51fef6f1bfe5cea44ff12fa563e5457"),
			classHash:           internalUtils.TestHexToFelt(t, "0x387edd4804deba7af741953fdf64189468f37593a66b618d00d2476be3168f8"),
			salt:                internalUtils.TestHexToFelt(t, "0x35fd26c7a05d64130bb27636931282cd1fbc19930933c97974d15028771ab72"),
			constructorCalldata: []*felt.Felt{},
			udcVersion:          UDCCairoV2,
			expectedAddress:     internalUtils.TestHexToFelt(t, "0x6493b4567665272ef04b30d57abe07a51b5f2c7862eb1f02a84e633ae2f12c7"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			address := PrecomputeAddressForUDC(test.classHash, test.salt, test.constructorCalldata, test.udcVersion, test.originAccAddress)
			assert.Equal(t, test.expectedAddress, address)
		})
	}
}
