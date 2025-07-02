package rpc

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClassAt tests the ClassAt function.
//
// The function tests the ClassAt function by creating different test sets for different environments
// (mock, testnet, and mainnet). It then iterates over each test set and performs the following steps:
//   - Calls the ClassAt function with the specified block tag and contract address.
//   - Checks the response type and performs the following actions based on the type:
//   - If the response type is DeprecatedContractClass or ContractClass:
//   - Checks if the program code exists in the response object.
//   - Checks if the expected operation exist in the provided contract address
//   - If the response type is of unknown type: log and fail the test
//
// Parameters:
//   - t: the testing object for running the test cases
//
// Returns:
//
//	none
func TestClassAt(t *testing.T) {
	testConfig := beforeEach(t, false)

	type testSetType struct {
		ContractAddress   *felt.Felt
		ExpectedOperation string
		Block             BlockID
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				ContractAddress:   internalUtils.TestHexToFelt(t, "0xdeadbeef"),
				ExpectedOperation: "0xdeadbeef",
				Block:             WithBlockNumber(58344),
			},
		},
		tests.TestnetEnv: {
			// v0 contract
			{
				ContractAddress:   internalUtils.TestHexToFelt(t, "0x073ad76dCF68168cBF68EA3EC0382a3605F3dEAf24dc076C355e275769b3c561"),
				ExpectedOperation: internalUtils.GetSelectorFromNameFelt("getPublicKey").String(),
				Block:             WithBlockNumber(58344),
			},
			// v2 contract
			{
				ContractAddress:   internalUtils.TestHexToFelt(t, "0x04dAadB9d30c887E1ab2cf7D78DFE444A77AAB5a49C3353d6d9977e7eD669902"),
				ExpectedOperation: internalUtils.GetSelectorFromNameFelt("name_get").String(),
				Block:             WithBlockNumber(65168),
			},
		},
		tests.MainnetEnv: {
			{
				ContractAddress:   internalUtils.TestHexToFelt(t, "0x004b3d247e79c58e77c93e2c52025d0bb1727957cc9c33b33f7216f369c77be5"),
				ExpectedOperation: internalUtils.GetSelectorFromNameFelt("get_name").String(),
				Block:             WithBlockNumber(643360),
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		resp, err := testConfig.provider.ClassAt(context.Background(), test.Block, test.ContractAddress)
		require.NoError(t, err)

		switch class := resp.(type) {
		case *contracts.DeprecatedContractClass:
			require.NotEmpty(t, class.Program, "code should exist")

			require.Condition(t, func() bool {
				for _, deprecatedCairoEntryPoint := range class.DeprecatedEntryPointsByType.External {
					if test.ExpectedOperation == deprecatedCairoEntryPoint.Selector.String() {
						return true
					}
				}

				return false
			}, "operation not found in the class")
		case *contracts.ContractClass:
			require.NotEmpty(t, class.SierraProgram, "code should exist")

			require.Condition(t, func() bool {
				for _, entryPointsByType := range class.EntryPointsByType.External {
					if test.ExpectedOperation == entryPointsByType.Selector.String() {
						return true
					}
				}

				return false
			}, "operation not found in the class")
		default:
			t.Fatalf("Received unknown response type: %v", reflect.TypeOf(resp))
		}
	}
}

// TestClassHashAt tests the ClassHashAt function.
//
// This function tests the behaviour of the ClassHashAt function by providing
// different test cases for the contract hash and expected class hash. It
// verifies if the returned class hash matches the expected class hash and
// if there are any differences between the two. It also checks if the
// returned class hash is not nil. The function takes in a testing.T
// parameter and does not return anything.
//
// Parameters:
//   - t: the testing object for running the test cases
//
// Returns:
//
//	none
func TestClassHashAt(t *testing.T) {
	testConfig := beforeEach(t, false)

	type testSetType struct {
		ContractHash      *felt.Felt
		ExpectedClassHash *felt.Felt
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				ContractHash:      internalUtils.TestHexToFelt(t, "0xdeadbeef"),
				ExpectedClassHash: internalUtils.TestHexToFelt(t, "0xdeadbeef"),
			},
		},
		tests.DevnetEnv: {
			{
				ContractHash:      internalUtils.TestHexToFelt(t, "0x41A78E741E5AF2FEC34B695679BC6891742439F7AFB8484ECD7766661AD02BF"),
				ExpectedClassHash: internalUtils.TestHexToFelt(t, "0x7B3E05F48F0C69E4A65CE5E076A66271A527AFF2C34CE1083EC6E1526997A69"),
			},
		},
		tests.TestnetEnv: {
			// v0 contracts
			{
				ContractHash:      internalUtils.TestHexToFelt(t, "0x05C0f2F029693e7E3A5500710F740f59C5462bd617A48F0Ed14b6e2d57adC2E9"),
				ExpectedClassHash: internalUtils.TestHexToFelt(t, "0x054328a1075b8820eb43caf0caa233923148c983742402dcfc38541dd843d01a"),
			},
			{
				ContractHash:      internalUtils.TestHexToFelt(t, "0x073ad76dcf68168cbf68ea3ec0382a3605f3deaf24dc076c355e275769b3c561"),
				ExpectedClassHash: internalUtils.TestHexToFelt(t, "0x036c7e49a16f8fc760a6fbdf71dde543d98be1fee2eda5daff59a0eeae066ed9"),
			},
			// v2 contract
			{
				ContractHash:      internalUtils.TestHexToFelt(t, "0x04dAadB9d30c887E1ab2cf7D78DFE444A77AAB5a49C3353d6d9977e7eD669902"),
				ExpectedClassHash: internalUtils.TestHexToFelt(t, "0x01f372292df22d28f2d4c5798734421afe9596e6a566b8bc9b7b50e26521b855"),
			},
		},
		tests.MainnetEnv: {
			{
				ContractHash:      internalUtils.TestHexToFelt(t, "0x3b4be7def2fc08589348966255e101824928659ebb724855223ff3a8c831efa"),
				ExpectedClassHash: internalUtils.TestHexToFelt(t, "0x4c53698c9a42341e4123632e87b752d6ae470ddedeb8b0063eaa2deea387eeb"),
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		classhash, err := testConfig.provider.ClassHashAt(context.Background(), WithBlockTag("latest"), test.ContractHash)
		require.NoError(t, err)
		require.NotEmpty(t, classhash, "should return a class")
		require.Equal(t, test.ExpectedClassHash, classhash)
	}
}

// TestClass is a test function that tests the behaviour of the Class function.
//
// It creates a test configuration and defines a testSet containing different scenarios
// for testing the Class function. The testSet is a map where the keys represent the
// test environment and the values are an array of testSetType structs. Each struct in
// the array represents a specific test case with properties such as BlockID,
// ClassHash, ExpectedProgram, and ExpectedEntryPointConstructor.
//
// The function iterates over each test case in the testSet and performs the following steps:
//   - Calls the Class function with the appropriate parameters.
//   - Handles the response based on its type:
//   - If the response is of type DeprecatedContractClass:
//   - Checks if the class program starts with the expected program.
//   - If not, it reports an error.
//   - If the response is of type ContractClass:
//   - Checks if the class program ends with the expected program.
//   - Compares the constructor entry point with the expected entry point constructor.
//   - If they are not equal, it reports an error.
//
// The function is used for testing the behaviour of the Class function in different scenarios.
//
// Parameters:
//   - t: A *testing.T object used for reporting test failures and logging
//
// Returns:
//
//	none
func TestClass(t *testing.T) {
	testConfig := beforeEach(t, false)

	type testSetType struct {
		BlockID                       BlockID
		ClassHash                     *felt.Felt
		ExpectedProgram               string
		ExpectedEntryPointConstructor contracts.SierraEntryPoint
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				BlockID:         WithBlockTag("pending"),
				ClassHash:       internalUtils.TestHexToFelt(t, "0xdeadbeef"),
				ExpectedProgram: "H4sIAAAAAAAA",
			},
		},
		tests.TestnetEnv: {
			// v0 class
			{
				BlockID:         WithBlockTag("latest"),
				ClassHash:       internalUtils.TestHexToFelt(t, "0x036c7e49a16f8fc760a6fbdf71dde543d98be1fee2eda5daff59a0eeae066ed9"),
				ExpectedProgram: "H4sIAAAAAAAA",
			},
			// v2 classes
			{
				BlockID:                       WithBlockTag("latest"),
				ClassHash:                     internalUtils.TestHexToFelt(t, "0x00816dd0297efc55dc1e7559020a3a825e81ef734b558f03c83325d4da7e6253"),
				ExpectedProgram:               internalUtils.TestHexToFelt(t, "0x576402000a0028a9c00a010").String(),
				ExpectedEntryPointConstructor: contracts.SierraEntryPoint{FunctionIdx: 34, Selector: internalUtils.TestHexToFelt(t, "0x28ffe4ff0f226a9107253e17a904099aa4f63a02a5621de0576e5aa71bc5194")},
			},
			{
				BlockID:                       WithBlockTag("latest"),
				ClassHash:                     internalUtils.TestHexToFelt(t, "0x01f372292df22d28f2d4c5798734421afe9596e6a566b8bc9b7b50e26521b855"),
				ExpectedProgram:               internalUtils.TestHexToFelt(t, "0xe70d09071117174f17170d4fe60d09071117").String(),
				ExpectedEntryPointConstructor: contracts.SierraEntryPoint{FunctionIdx: 2, Selector: internalUtils.TestHexToFelt(t, "0x28ffe4ff0f226a9107253e17a904099aa4f63a02a5621de0576e5aa71bc5194")},
			},
		},
		tests.MainnetEnv: {
			// v2 class
			{
				BlockID:                       WithBlockTag("latest"),
				ClassHash:                     internalUtils.TestHexToFelt(t, "0x029927c8af6bccf3f6fda035981e765a7bdbf18a2dc0d630494f8758aa908e2b"),
				ExpectedProgram:               internalUtils.TestHexToFelt(t, "0x9fa00900700e00712e12500712e").String(),
				ExpectedEntryPointConstructor: contracts.SierraEntryPoint{FunctionIdx: 32, Selector: internalUtils.TestHexToFelt(t, "0x28ffe4ff0f226a9107253e17a904099aa4f63a02a5621de0576e5aa71bc5194")},
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		resp, err := testConfig.provider.Class(context.Background(), test.BlockID, test.ClassHash)
		require.NoError(t, err)

		switch class := resp.(type) {
		case *contracts.DeprecatedContractClass:
			assert.Contains(t, class.Program, test.ExpectedProgram)
		case *contracts.ContractClass:
			assert.Equal(t, class.SierraProgram[len(class.SierraProgram)-1].String(), test.ExpectedProgram)
			assert.Equal(t, class.EntryPointsByType.Constructor[0], test.ExpectedEntryPointConstructor)
		default:
			t.Fatalf("Received unknown response type: %v", reflect.TypeOf(resp))
		}
	}
}

// TestStorageAt tests the StorageAt function.
//
// It tests the StorageAt function by creating a test environment and executing
// a series of test cases. Each test case represents a different scenario, such
// as different contract hashes, storage keys, and expected values. The function
// checks if the actual value returned by the StorageAt function matches the
// expected value. If there is a mismatch, the test fails and an error is
// reported.
//
// Parameters:
//   - t: The testing.T instance used for reporting test failures and logging
//
// Returns:
//
//	none
func TestStorageAt(t *testing.T) {
	testConfig := beforeEach(t, false)

	type testSetType struct {
		ContractHash  *felt.Felt
		StorageKey    string
		Block         BlockID
		ExpectedValue string
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				ContractHash:  internalUtils.TestHexToFelt(t, "0xdeadbeef"),
				StorageKey:    "_signer",
				Block:         WithBlockTag("latest"),
				ExpectedValue: "0xdeadbeef",
			},
		},
		tests.DevnetEnv: {
			{
				ContractHash:  internalUtils.TestHexToFelt(t, "0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d"),
				StorageKey:    "ERC20_name",
				Block:         WithBlockTag("latest"),
				ExpectedValue: "0x537461726b4e657420546f6b656e",
			},
		},
		tests.TestnetEnv: {
			{
				ContractHash:  internalUtils.TestHexToFelt(t, "0x0200AB5CE3D7aDE524335Dc57CaF4F821A0578BBb2eFc2166cb079a3D29cAF9A"),
				StorageKey:    "_signer",
				Block:         WithBlockNumber(69399),
				ExpectedValue: "0x38bd4cad8706e3a5d167ef7af12e28268c6122df3e0e909839a103039871b9e",
			},
		},
		tests.MainnetEnv: {
			{
				ContractHash:  internalUtils.TestHexToFelt(t, "0x8d17e6a3B92a2b5Fa21B8e7B5a3A794B05e06C5FD6C6451C6F2695Ba77101"),
				StorageKey:    "_signer",
				Block:         WithBlockTag("latest"),
				ExpectedValue: "0x7f72660ca40b8ca85f9c0dd38db773f17da7a52f5fc0521cb8b8d8d44e224b8",
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		value, err := testConfig.provider.StorageAt(context.Background(), test.ContractHash, test.StorageKey, test.Block)
		require.NoError(t, err)
		require.EqualValues(t, test.ExpectedValue, value)
	}
}

// TestNonce is a test function for testing the Nonce functionality.
//
// It initialises a test configuration, sets up a test data set, and then performs a series of tests.
// The tests involve calling the Nonce function.
// The expected result is a successful response from the Nonce function and a matching value with the expected nonce.
// If any errors occur during the tests, the function will fail and display an error message.
//
// Parameters:
//   - t: the testing object for running the test cases
//
// Returns:
//
//	none
func TestNonce(t *testing.T) {
	testConfig := beforeEach(t, false)

	type testSetType struct {
		ContractAddress *felt.Felt
		Block           BlockID
		ExpectedNonce   *felt.Felt
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				ContractAddress: internalUtils.TestHexToFelt(t, "0x0207acc15dc241e7d167e67e30e769719a727d3e0fa47f9e187707289885dfde"),
				Block:           WithBlockTag("latest"),
				ExpectedNonce:   internalUtils.TestHexToFelt(t, "0xdeadbeef"),
			},
		},
		tests.DevnetEnv: {
			{
				ContractAddress: internalUtils.TestHexToFelt(t, "0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d"),
				Block:           WithBlockTag("latest"),
				ExpectedNonce:   internalUtils.TestHexToFelt(t, "0x0"),
			},
		},
		tests.TestnetEnv: {
			{
				ContractAddress: internalUtils.TestHexToFelt(t, "0x0200AB5CE3D7aDE524335Dc57CaF4F821A0578BBb2eFc2166cb079a3D29cAF9A"),
				Block:           WithBlockNumber(69399),
				ExpectedNonce:   internalUtils.TestHexToFelt(t, "0x1"),
			},
		},
		tests.MainnetEnv: {
			{
				ContractAddress: internalUtils.TestHexToFelt(t, "0x00bE9AeF00Ec751Ba252A595A473315FBB8DA629850e13b8dB83d0fACC44E4f2"),
				Block:           WithBlockNumber(644060),
				ExpectedNonce:   internalUtils.TestHexToFelt(t, "0x2"),
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		nonce, err := testConfig.provider.Nonce(context.Background(), test.Block, test.ContractAddress)
		require.NoError(t, err)
		require.NotNil(t, nonce, "should return a nonce")
		require.Equal(t, test.ExpectedNonce, nonce)
	}
}

// TestEstimateMessageFee is a test function to test the EstimateMessageFee function.
//
// Parameters:
//   - t: the testing object for running the test cases
//
// Returns:
//
//	none
func TestEstimateMessageFee(t *testing.T) {
	testConfig := beforeEach(t, false)

	type testSetType struct {
		MsgFromL1
		BlockID
		ExpectedFeeEst *FeeEstimation
		ExpectedError  *RPCError
	}

	// https://sepolia.voyager.online/message/0x273f4e20fc522098a60099e5872ab3deeb7fb8321a03dadbd866ac90b7268361
	l1Handler := MsgFromL1{
		FromAddress: "0x8453fc6cd1bcfe8d4dfc069c400b433054d47bdc",
		ToAddress:   internalUtils.TestHexToFelt(t, "0x04c5772d1914fe6ce891b64eb35bf3522aeae1315647314aac58b01137607f3f"),
		Selector:    internalUtils.TestHexToFelt(t, "0x1b64b1b3b690b43b9b514fb81377518f4039cd3e4f4914d8a6bdf01d679fb19"),
		Payload: internalUtils.TestHexArrToFelt(t, []string{
			"0x455448",
			"0x2f14d277fc49e0e2d2967d019aea8d6bd9cb3998",
			"0x02000e6213e24b84012b1f4b1cbd2d7a723fb06950aeab37bedb6f098c7e051a",
			"0x01a055690d9db80000",
			"0x00",
		}),
	}

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				MsgFromL1: MsgFromL1{FromAddress: "0x0", ToAddress: &felt.Zero, Selector: &felt.Zero, Payload: []*felt.Felt{&felt.Zero}},
				BlockID:   BlockID{Tag: "latest"},
				ExpectedFeeEst: &FeeEstimation{
					L1GasConsumed: internalUtils.RANDOM_FELT,
					L1GasPrice:    internalUtils.RANDOM_FELT,
					L2GasConsumed: internalUtils.RANDOM_FELT,
					L2GasPrice:    internalUtils.RANDOM_FELT,
					OverallFee:    internalUtils.RANDOM_FELT,
				},
			},
		},
		tests.TestnetEnv: {
			{
				MsgFromL1: l1Handler,
				BlockID:   WithBlockNumber(523066),
				ExpectedFeeEst: &FeeEstimation{
					L1GasConsumed:     internalUtils.TestHexToFelt(t, "0x4ed1"),
					L1GasPrice:        internalUtils.TestHexToFelt(t, "0x7e15d2b5"),
					L2GasConsumed:     internalUtils.TestHexToFelt(t, "0x0"),
					L2GasPrice:        internalUtils.TestHexToFelt(t, "0x0"),
					L1DataGasConsumed: internalUtils.TestHexToFelt(t, "0x80"),
					L1DataGasPrice:    internalUtils.TestHexToFelt(t, "0x1"),
					OverallFee:        internalUtils.TestHexToFelt(t, "0x26d196042c45"),
					FeeUnit:           UnitWei,
				},
			},
			{ // invalid msg data
				MsgFromL1: MsgFromL1{
					FromAddress: "0x8453fc6cd1bcfe8d4dfc069c400b433054d47bdc",
					ToAddress:   internalUtils.RANDOM_FELT,
					Selector:    internalUtils.RANDOM_FELT,
					Payload:     []*felt.Felt{},
				},
				BlockID:       WithBlockNumber(523066),
				ExpectedError: ErrContractError,
			},
			{ // invalid block number
				MsgFromL1:     l1Handler,
				BlockID:       WithBlockNumber(9999999999999999999),
				ExpectedError: ErrBlockNotFound,
			},
		},
		"mainnet": {},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		resp, err := testConfig.provider.EstimateMessageFee(context.Background(), test.MsgFromL1, test.BlockID)
		if err != nil {
			rpcErr, ok := err.(*RPCError)
			require.True(t, ok)
			require.Equal(t, test.ExpectedError.Code, rpcErr.Code)
			require.Equal(t, test.ExpectedError.Message, rpcErr.Message)
		} else {
			require.Exactly(t, test.ExpectedFeeEst, resp)
		}
	}
}

func TestEstimateFee(t *testing.T) {
	testConfig := beforeEach(t, false)

	type testSetType struct {
		description   string
		txs           []BroadcastTxn
		simFlags      []SimulationFlag
		blockID       BlockID
		expectedResp  []FeeEstimation
		expectedError *RPCError
	}

	bradcastInvokeV3 := *internalUtils.TestUnmarshalJSONFileToType[BroadcastInvokeTxnV3](t, "./testData/transactions/sepoliaInvokeV3_0x6035477af07a1b0a0186bec85287a6f629791b2f34b6e90eec9815c7a964f64.json", "")

	testSet, ok := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				description: "without flag",
				txs: []BroadcastTxn{
					bradcastInvokeV3,
				},
				simFlags:      []SimulationFlag{},
				blockID:       WithBlockTag("latest"),
				expectedError: nil,
				expectedResp: []FeeEstimation{
					{
						L1GasConsumed:     internalUtils.RANDOM_FELT,
						L1GasPrice:        internalUtils.RANDOM_FELT,
						L2GasConsumed:     internalUtils.RANDOM_FELT,
						L2GasPrice:        internalUtils.RANDOM_FELT,
						L1DataGasConsumed: internalUtils.RANDOM_FELT,
						L1DataGasPrice:    internalUtils.RANDOM_FELT,
						OverallFee:        internalUtils.RANDOM_FELT,
						FeeUnit:           UnitWei,
					},
				},
			},
			{
				description: "with flag",
				txs: []BroadcastTxn{
					bradcastInvokeV3,
				},
				simFlags:      []SimulationFlag{SKIP_VALIDATE},
				blockID:       WithBlockTag("latest"),
				expectedError: nil,
				expectedResp: []FeeEstimation{
					{
						L1GasConsumed:     new(felt.Felt).SetUint64(1234),
						L1GasPrice:        new(felt.Felt).SetUint64(1234),
						L2GasConsumed:     new(felt.Felt).SetUint64(1234),
						L2GasPrice:        new(felt.Felt).SetUint64(1234),
						L1DataGasConsumed: new(felt.Felt).SetUint64(1234),
						L1DataGasPrice:    new(felt.Felt).SetUint64(1234),
						OverallFee:        new(felt.Felt).SetUint64(1234),
						FeeUnit:           UnitWei,
					},
				},
			},
		},
		tests.TestnetEnv: {
			{
				description: "without flag",
				txs: []BroadcastTxn{
					bradcastInvokeV3,
				},
				simFlags:      []SimulationFlag{},
				blockID:       WithBlockNumber(574447),
				expectedError: nil,
				expectedResp: []FeeEstimation{
					{
						L1GasConsumed:     internalUtils.TestHexToFelt(t, "0x0"),
						L1GasPrice:        internalUtils.TestHexToFelt(t, "0xa7fe9fec104"),
						L2GasConsumed:     internalUtils.TestHexToFelt(t, "0xf49c0"),
						L2GasPrice:        internalUtils.TestHexToFelt(t, "0x1020990a5"),
						L1DataGasConsumed: internalUtils.TestHexToFelt(t, "0x140"),
						L1DataGasPrice:    internalUtils.TestHexToFelt(t, "0x617"),
						OverallFee:        internalUtils.TestHexToFelt(t, "0xf68e5bb1e2580"),
						FeeUnit:           UnitStrk,
					},
				},
			},
			{
				description: "with flag",
				txs: []BroadcastTxn{
					bradcastInvokeV3,
				},
				simFlags:      []SimulationFlag{SKIP_VALIDATE},
				blockID:       WithBlockNumber(574447),
				expectedError: nil,
				expectedResp: []FeeEstimation{
					{
						L1GasConsumed:     internalUtils.TestHexToFelt(t, "0x0"),
						L1GasPrice:        internalUtils.TestHexToFelt(t, "0xa7fe9fec104"),
						L2GasConsumed:     internalUtils.TestHexToFelt(t, "0xe1140"),
						L2GasPrice:        internalUtils.TestHexToFelt(t, "0x1020990a5"),
						L1DataGasConsumed: internalUtils.TestHexToFelt(t, "0x140"),
						L1DataGasPrice:    internalUtils.TestHexToFelt(t, "0x617"),
						OverallFee:        internalUtils.TestHexToFelt(t, "0xe2de90e0cbb00"),
						FeeUnit:           UnitStrk,
					},
				},
			},
			{
				description: "invalid transaction",
				txs: []BroadcastTxn{
					InvokeTxnV3{
						ResourceBounds: &ResourceBoundsMapping{
							L1Gas: ResourceBounds{
								MaxAmount:       "0x0",
								MaxPricePerUnit: "0x4305031628668",
							},
							L1DataGas: ResourceBounds{
								MaxAmount:       "0x210",
								MaxPricePerUnit: "0x948",
							},
							L2Gas: ResourceBounds{
								MaxAmount:       "0x15cde0",
								MaxPricePerUnit: "0x18955dc56",
							},
						},
						Type:                  TransactionType_Invoke,
						Version:               TransactionV3,
						SenderAddress:         internalUtils.RANDOM_FELT,
						Nonce:                 &felt.Zero,
						Calldata:              []*felt.Felt{},
						Signature:             []*felt.Felt{},
						Tip:                   "0x0",
						PayMasterData:         []*felt.Felt{},
						AccountDeploymentData: []*felt.Felt{},
						NonceDataMode:         DAModeL1,
						FeeMode:               DAModeL1,
					},
				},
				simFlags:      []SimulationFlag{},
				blockID:       WithBlockNumber(100000),
				expectedError: ErrTxnExec,
			},
			{
				description: "invalid block",
				txs: []BroadcastTxn{
					bradcastInvokeV3,
				},
				simFlags:      []SimulationFlag{},
				blockID:       WithBlockNumber(9999999999999999999),
				expectedError: ErrBlockNotFound,
			},
		},
	}[tests.TEST_ENV]

	// TODO: implement this pattern in all tests to know which test cases are being skipped
	if !ok {
		t.Skipf("'%s' environment testset not implemented by this test", tests.TEST_ENV)
	}

	for _, test := range testSet {
		t.Run(test.description, func(t *testing.T) {
			resp, err := testConfig.provider.EstimateFee(context.Background(), test.txs, test.simFlags, test.blockID)
			if test.expectedError != nil {
				require.Error(t, err)
				rpcErr, ok := err.(*RPCError)
				require.True(t, ok)
				assert.Equal(t, test.expectedError.Code, rpcErr.Code)
				assert.Equal(t, test.expectedError.Message, rpcErr.Message)
				assert.IsType(t, rpcErr.Data, rpcErr.Data)
			} else {
				require.NoError(t, err)
			}

			assert.Exactly(t, test.expectedResp, resp)
		})
	}
}

func TestGetStorageProof(t *testing.T) {
	testConfig := beforeEach(t, false)

	type testSetType struct {
		Description       string
		StorageProofInput StorageProofInput
		ExpectedError     error
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv:   {},
		tests.DevnetEnv: {},
		tests.TestnetEnv: {
			{
				Description: "normal call, only required field block_id with 'latest' tag",
				StorageProofInput: StorageProofInput{
					BlockID: BlockID{Tag: "latest"},
				},
				ExpectedError: nil,
			},
			{
				Description: "block_id + class_hashes parameter",
				StorageProofInput: StorageProofInput{
					BlockID: BlockID{Tag: "latest"},
					ClassHashes: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x076791ef97c042f81fbf352ad95f39a22554ee8d7927b2ce3c681f3418b5206a"),
					},
				},
				ExpectedError: nil,
			},
			{
				Description: "block_id + contract_addresses parameter",
				StorageProofInput: StorageProofInput{
					BlockID: BlockID{Tag: "latest"},
					ContractAddresses: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
					},
				},
				ExpectedError: nil,
			},
			{
				Description: "block_id + contracts_storage_keys parameter",
				StorageProofInput: StorageProofInput{
					BlockID: BlockID{Tag: "latest"},
					ContractsStorageKeys: []ContractStorageKeys{
						{
							ContractAddress: internalUtils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
							StorageKeys: []*felt.Felt{
								internalUtils.TestHexToFelt(t, "0x0341c1bdfd89f69748aa00b5742b03adbffd79b8e80cab5c50d91cd8c2a79be1"),
							},
						},
					},
				},
				ExpectedError: nil,
			},
			{
				Description: "block_id + class_hashes + contract_addresses + contracts_storage_keys parameter",
				StorageProofInput: StorageProofInput{
					BlockID: BlockID{Tag: "latest"},
					ClassHashes: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x076791ef97c042f81fbf352ad95f39a22554ee8d7927b2ce3c681f3418b5206a"),
						internalUtils.TestHexToFelt(t, "0x009524a94b41c4440a16fd96d7c1ef6ad6f44c1c013e96662734502cd4ee9b1f"),
					},
					ContractAddresses: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x04718f5a0Fc34cC1AF16A1cdee98fFB20C31f5cD61D6Ab07201858f4287c938D"),
						internalUtils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
					},
					ContractsStorageKeys: []ContractStorageKeys{
						{
							ContractAddress: internalUtils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
							StorageKeys: []*felt.Felt{
								internalUtils.TestHexToFelt(t, "0x0341c1bdfd89f69748aa00b5742b03adbffd79b8e80cab5c50d91cd8c2a79be1"),
								internalUtils.TestHexToFelt(t, "0x00b6ce5410fca59d078ee9b2a4371a9d684c530d697c64fbef0ae6d5e8f0ac72"),
							},
						},
						{
							ContractAddress: internalUtils.TestHexToFelt(t, "0x04718f5a0Fc34cC1AF16A1cdee98fFB20C31f5cD61D6Ab07201858f4287c938D"),
							StorageKeys: []*felt.Felt{
								internalUtils.TestHexToFelt(t, "0x0341c1bdfd89f69748aa00b5742b03adbffd79b8e80cab5c50d91cd8c2a79be1"),
								internalUtils.TestHexToFelt(t, "0x00b6ce5410fca59d078ee9b2a4371a9d684c530d697c64fbef0ae6d5e8f0ac72"),
							},
						},
					},
				},
				ExpectedError: nil,
			},
			{
				Description: "error: using pending tag in block_id",
				StorageProofInput: StorageProofInput{
					BlockID: BlockID{Tag: "pending"},
				},
				ExpectedError: ErrInvalidBlockID,
			},
			{
				Description: "error: invalid block number",
				StorageProofInput: StorageProofInput{
					BlockID: func() BlockID {
						num := uint64(999999999)

						return BlockID{Number: &num}
					}(),
				},
				ExpectedError: ErrBlockNotFound,
			},
			{
				Description: "error: storage proof not supported",
				StorageProofInput: StorageProofInput{
					BlockID: func() BlockID {
						num := uint64(123456)

						return BlockID{Number: &num}
					}(),
				},
				ExpectedError: ErrStorageProofNotSupported,
			},
		},
		tests.MainnetEnv: {},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		t.Run(test.Description, func(t *testing.T) {
			result, err := testConfig.provider.GetStorageProof(context.Background(), test.StorageProofInput)
			if test.ExpectedError != nil {
				require.Error(t, err)
				require.ErrorContains(t, err, test.ExpectedError.Error())

				return
			}

			require.NoError(t, err)
			require.NotNil(t, result, "empty result from starknet_getStorageProof")

			// verify JSON equality
			var rawResult any

			// call the RPC method directly to get the raw result
			input := test.StorageProofInput
			input.BlockID = WithBlockHash(
				result.GlobalRoots.BlockHash,
			) // using the same block returned by GetStorageProof to avoid temporal coupling
			err = testConfig.provider.c.CallContext(
				context.Background(),
				&rawResult,
				"starknet_getStorageProof",
				input,
			)
			require.NoError(t, err)
			// marshal the results to JSON
			rawResultJSON, err := json.Marshal(rawResult)
			require.NoError(t, err)
			resultJSON, err := json.Marshal(result)
			require.NoError(t, err)

			assertStorageProofJSONEquality(t, rawResultJSON, resultJSON)
		})
	}
}

func assertStorageProofJSONEquality(t *testing.T, expectedResult, result []byte) {
	// unmarshal to map[string]any
	var expectedResultMap, resultMap map[string]any
	require.NoError(t, json.Unmarshal(expectedResult, &expectedResultMap))
	require.NoError(t, json.Unmarshal(result, &resultMap))

	// compare 'classes_proof'
	expectedClassesProof, ok := expectedResultMap["classes_proof"].([]any)
	require.True(t, ok)
	resultClassesProof, ok := resultMap["classes_proof"].([]any)
	require.True(t, ok)
	assert.ElementsMatch(t, expectedClassesProof, resultClassesProof)

	// compare 'contracts_proof'
	expectedContractsProof, ok := expectedResultMap["contracts_proof"].(map[string]any)
	require.True(t, ok)
	resultContractsProof, ok := resultMap["contracts_proof"].(map[string]any)
	require.True(t, ok)
	// compare 'contracts_proof.nodes'
	expectedContractsProofNodes, ok := expectedContractsProof["nodes"].([]any)
	require.True(t, ok)
	resultContractsProofNodes, ok := resultContractsProof["nodes"].([]any)
	require.True(t, ok)
	assert.ElementsMatch(t, expectedContractsProofNodes, resultContractsProofNodes)
	// compare 'contracts_proof.contract_leaves_data'
	expectedContractsProofContractLeavesData, ok := expectedContractsProof["contract_leaves_data"].([]any)
	require.True(t, ok)
	resultContractsProofContractLeavesData, ok := resultContractsProof["contract_leaves_data"].([]any)
	require.True(t, ok)
	assert.ElementsMatch(t, expectedContractsProofContractLeavesData, resultContractsProofContractLeavesData)

	// compare 'contracts_storage_proofs'
	expectedContractsStorageProofs, ok := expectedResultMap["contracts_storage_proofs"].([]any)
	require.True(t, ok)
	expectedGeneralSlice := make([]any, 0)
	resultContractsStorageProofs, ok := resultMap["contracts_storage_proofs"].([]any)
	require.True(t, ok)
	resultGeneralSlice := make([]any, 0)
	for i, expectedContractStorageProof := range expectedContractsStorageProofs {
		expectedContractStorageProofArray, ok := expectedContractStorageProof.([]any)
		require.True(t, ok)
		expectedGeneralSlice = append(expectedGeneralSlice, expectedContractStorageProofArray...)

		resultContractStorageProofArray, ok := resultContractsStorageProofs[i].([]any)
		require.True(t, ok)
		resultGeneralSlice = append(resultGeneralSlice, resultContractStorageProofArray...)
	}
	assert.ElementsMatch(t, expectedGeneralSlice, resultGeneralSlice)

	// compare 'global_roots'
	assert.Equal(t, expectedResultMap["global_roots"], resultMap["global_roots"])
}
