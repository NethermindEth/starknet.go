package rpc

import (
	"context"
	"log"
	"reflect"
	"strings"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/test-go/testify/require"
)

// TestClassAt tests the ClassAt function.
//
// The function tests the ClassAt function by creating different test sets for different environments
// (mock, testnet, and mainnet). It then iterates over each test set and performs the following steps:
//   - Creates a spy object to intercept calls to the provider.
//   - Sets the provider of the test configuration to the spy object.
//   - Calls the ClassAt function with the specified block tag and contract address.
//   - Checks the response type and performs the following actions based on the type:
//   - If the response type is DeprecatedContractClass:
//   - Compares the response object with the spy object and checks for a full match.
//   - If the objects do not match, compares them again and logs an error if the match is still not achieved.
//   - Checks if the program code exists in the response object.
//   - If the response type is ContractClass:
//   - Throws an error indicating that the case is not covered.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
// Todo: Is not yet implemented completely
//
//	none
func TestClassAt(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		ContractAddress   *felt.Felt
		ExpectedOperation string
		BlockHash         string
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				ContractAddress:   utils.TestHexToFelt(t, "0xdeadbeef"),
				ExpectedOperation: "0xdeadbeef",
				BlockHash:         "0x561eeb100ad42aedc8810cce883caccc77eda75a9af58b24aabb770c027d249",
			},
		},
		"testnet": {
			// v0 contract
			{
				ContractAddress:   utils.TestHexToFelt(t, "0x073ad76dCF68168cBF68EA3EC0382a3605F3dEAf24dc076C355e275769b3c561"),
				ExpectedOperation: utils.GetSelectorFromNameFelt("getPublicKey").String(),
				BlockHash:         "0x561eeb100ad42aedc8810cce883caccc77eda75a9af58b24aabb770c027d249",
			},
			// v2 contract
			{
				ContractAddress:   utils.TestHexToFelt(t, "0x04dAadB9d30c887E1ab2cf7D78DFE444A77AAB5a49C3353d6d9977e7eD669902"),
				ExpectedOperation: utils.GetSelectorFromNameFelt("name_get").String(),
				BlockHash:         "0x6d49f7047818b6e002ab2ae7ee0376fe1632fb4fe4c80775ec7ed728fa99ecc",
			},
		},
		"mainnet": {
			{
				ContractAddress:   utils.TestHexToFelt(t, "0x004b3d247e79c58e77c93e2c52025d0bb1727957cc9c33b33f7216f369c77be5"),
				ExpectedOperation: utils.GetSelectorFromNameFelt("get_name").String(),
				BlockHash:         "0x05b277fbda1ca1a24dcfe7d9b45e3083d44dd1bb873349b7183dbbf63db74acf",
			},
		},
	}[testEnv]

	for _, test := range testSet {
		require := require.New(t)
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		resp, err := testConfig.provider.ClassAt(context.Background(), WithBlockHash(utils.TestHexToFelt(t, test.BlockHash)), test.ContractAddress)
		require.NoError(err)

		switch class := resp.(type) {
		case *DeprecatedContractClass:
			diff, err := spy.Compare(class, false)
			require.NoError(err, "expecting to match")
			require.Equal(diff, "FullMatch", "structure expecting to be FullMatch")
			require.NotEmpty(class.Program, "code should exist")

			require.Condition(func() bool {
				for _, deprecatedCairoEntryPoint := range class.DeprecatedEntryPointsByType.External {
					t.Log(deprecatedCairoEntryPoint)
					if test.ExpectedOperation == deprecatedCairoEntryPoint.Selector.String() {
						return true
					}
				}
				return false
			}, "operation not found in the class")
		case *ContractClass:
			diff, err := spy.Compare(class, false)
			require.NoError(err, "expecting to match")
			require.Equal(diff, "FullMatch", "structure expecting to be FullMatch")
			require.NotEmpty(class.SierraProgram, "code should exist")

			require.Condition(func() bool {
				for _, entryPointsByType := range class.EntryPointsByType.External {
					t.Log(entryPointsByType)
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
// This function tests the behavior of the ClassHashAt function by providing
// different test cases for the contract hash and expected class hash. It
// verifies if the returned class hash matches the expected class hash and
// if there are any differences between the two. It also checks if the
// returned class hash is not nil. The function takes in a testing.T
// parameter and does not return anything.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestClassHashAt(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		ContractHash      *felt.Felt
		ExpectedClassHash *felt.Felt
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				ContractHash:      utils.TestHexToFelt(t, "0xdeadbeef"),
				ExpectedClassHash: utils.TestHexToFelt(t, "0xdeadbeef"),
			},
		},
		"testnet": {
			// v0 contracts
			{
				ContractHash:      utils.TestHexToFelt(t, "0x05C0f2F029693e7E3A5500710F740f59C5462bd617A48F0Ed14b6e2d57adC2E9"),
				ExpectedClassHash: utils.TestHexToFelt(t, "0x054328a1075b8820eb43caf0caa233923148c983742402dcfc38541dd843d01a"),
			},
			{
				ContractHash:      utils.TestHexToFelt(t, "0x073ad76dcf68168cbf68ea3ec0382a3605f3deaf24dc076c355e275769b3c561"),
				ExpectedClassHash: utils.TestHexToFelt(t, "0x036c7e49a16f8fc760a6fbdf71dde543d98be1fee2eda5daff59a0eeae066ed9"),
			},
			// v2 contract
			{
				ContractHash:      utils.TestHexToFelt(t, "0x04dAadB9d30c887E1ab2cf7D78DFE444A77AAB5a49C3353d6d9977e7eD669902"),
				ExpectedClassHash: utils.TestHexToFelt(t, "0x01f372292df22d28f2d4c5798734421afe9596e6a566b8bc9b7b50e26521b855"),
			},
		},
		"mainnet": {
			{
				ContractHash:      utils.TestHexToFelt(t, "0x3b4be7def2fc08589348966255e101824928659ebb724855223ff3a8c831efa"),
				ExpectedClassHash: utils.TestHexToFelt(t, "0x4c53698c9a42341e4123632e87b752d6ae470ddedeb8b0063eaa2deea387eeb"),
			},
		},
	}[testEnv]

	for _, test := range testSet {

		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		classhash, err := testConfig.provider.ClassHashAt(context.Background(), WithBlockTag("latest"), test.ContractHash)
		if err != nil {
			t.Fatal(err)
		}
		diff, err := spy.Compare(classhash, false)
		if err != nil {
			t.Fatal("expecting to match", err)
		}
		if diff != "FullMatch" {
			if _, err := spy.Compare(classhash, true); err != nil {
				log.Fatal(err)
			}
			t.Fatal("structure expecting to be FullMatch, instead", diff)
		}
		if classhash == nil {
			t.Fatalf("should return a class, instead %v", classhash)
		}
		require.Equal(t, test.ExpectedClassHash, classhash)
	}
}

// TestClass is a test function that tests the behavior of the Class function.
//
// It creates a test configuration and defines a testSet containing different scenarios
// for testing the Class function. The testSet is a map where the keys represent the
// test environment and the values are an array of testSetType structs. Each struct in
// the array represents a specific test case with properties such as BlockID,
// ClassHash, ExpectedProgram, and ExpectedEntryPointConstructor.
//
// The function iterates over each test case in the testSet and performs the following steps:
// - Creates a new spy object to spy on the provider.
// - Sets the provider of the test configuration to the spy object.
// - Calls the Class function with the appropriate parameters.
// - Handles the response based on its type:
//   - If the response is of type DeprecatedContractClass:
//   - Compares the response with the spy object to check for any differences.
//   - If there is a difference, it reports an error and prints the difference.
//   - Checks if the class program starts with the expected program.
//   - If not, it reports an error.
//   - If the response is of type ContractClass:
//   - Compares the constructor entry point with the expected entry point constructor.
//   - If they are not equal, it reports an error.
//
// The function is used for testing the behavior of the Class function in different scenarios.
//
// Parameters:
// - t: A *testing.T object used for reporting test failures and logging
// Returns:
//
//	none
func TestClass(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		BlockID                       BlockID
		ClassHash                     *felt.Felt
		ExpectedProgram               string
		ExpectedEntryPointConstructor SierraEntryPoint
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				BlockID:         WithBlockTag("pending"),
				ClassHash:       utils.TestHexToFelt(t, "0xdeadbeef"),
				ExpectedProgram: "H4sIAAAAAAAA",
			},
		},
		"testnet": {
			// v0 class
			{
				BlockID:         WithBlockNumber(15329),
				ClassHash:       utils.TestHexToFelt(t, "0x036c7e49a16f8fc760a6fbdf71dde543d98be1fee2eda5daff59a0eeae066ed9"),
				ExpectedProgram: "H4sIAAAAAAAA",
			},
			{
				BlockID:                       WithBlockHash(utils.TestHexToFelt(t, "0x258dc3bf21fbefb29b5dfd782c9d9472f73075213e9b63a0421ff7d2d3106d2")),
				ClassHash:                     utils.TestHexToFelt(t, "0x036c7e49a16f8fc760a6fbdf71dde543d98be1fee2eda5daff59a0eeae066ed9"),
				ExpectedEntryPointConstructor: SierraEntryPoint{FunctionIdx: 16, Selector: utils.TestHexToFelt(t, "0x28ffe4ff0f226a9107253e17a904099aa4f63a02a5621de0576e5aa71bc5194")},
			},
			// v2 class
			{
				BlockID:         WithBlockNumber(15329),
				ClassHash:       utils.TestHexToFelt(t, "0x079b7ec8fdf40a4ff6ed47123049dfe36b5c02db93aa77832682344775ef70c6"),
				ExpectedProgram: "H4sIAAAAAAAA",
			},
			{
				BlockID:   WithBlockHash(utils.TestHexToFelt(t, "0x7b7f2d9b2e4502326eac1615e754d414df22b8266e7206f0cc90380e8052ee3")),
				ClassHash: utils.TestHexToFelt(t, "0x01f372292df22d28f2d4c5798734421afe9596e6a566b8bc9b7b50e26521b855"),
				// ExpectedEntryPointConstructor: SierraEntryPoint{FunctionIdx: 16, Selector: utils.TestHexToFelt(t, "0x28ffe4ff0f226a9107253e17a904099aa4f63a02a5621de0576e5aa71bc5194")},
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		resp, err := testConfig.provider.Class(context.Background(), test.BlockID, test.ClassHash)
		if err != nil {
			t.Fatal(err)
		}

		switch class := resp.(type) {
		case DeprecatedContractClass:
			diff, err := spy.Compare(class, false)
			if err != nil {
				t.Fatal("expecting to match", err)
			}
			if diff != "FullMatch" {
				if _, err := spy.Compare(class, true); err != nil {
					log.Fatal(err)
				}
				t.Fatal("structure expecting to be FullMatch, instead", diff)
			}

			if !strings.HasPrefix(class.Program, test.ExpectedProgram) {
				t.Fatal("code should exist")
			}
		case ContractClass:
			require.Equal(t, class.EntryPointsByType.Constructor, test.ExpectedEntryPointConstructor)
		default:
			log.Fatalln("Received unknown response type:", reflect.TypeOf(resp))
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
// - t: The testing.T instance used for reporting test failures and logging
// Returns:
//
//	none
func TestStorageAt(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		ContractHash  *felt.Felt
		StorageKey    string
		Block         BlockID
		ExpectedValue string
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				ContractHash:  utils.TestHexToFelt(t, "0xdeadbeef"),
				StorageKey:    "_signer",
				Block:         WithBlockTag("latest"),
				ExpectedValue: "0xdeadbeef",
			},
		},
		"testnet": {
			{
				ContractHash:  utils.TestHexToFelt(t, "0x073ad76dCF68168cBF68EA3EC0382a3605F3dEAf24dc076C355e275769b3c561"),
				StorageKey:    "balance",
				Block:         WithBlockHash(utils.TestHexToFelt(t, "0x561eeb100ad42aedc8810cce883caccc77eda75a9af58b24aabb770c027d249")),
				ExpectedValue: "0x0",
			},
		},
		"mainnet": {
			{
				ContractHash:  utils.TestHexToFelt(t, "0x8d17e6a3B92a2b5Fa21B8e7B5a3A794B05e06C5FD6C6451C6F2695Ba77101"),
				StorageKey:    "_signer",
				Block:         WithBlockTag("latest"),
				ExpectedValue: "0x7f72660ca40b8ca85f9c0dd38db773f17da7a52f5fc0521cb8b8d8d44e224b8",
			},
		},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		value, err := testConfig.provider.StorageAt(context.Background(), test.ContractHash, test.StorageKey, test.Block)
		if err != nil {
			t.Fatal(err)
		}
		diff, err := spy.Compare(value, false)
		if err != nil {
			t.Fatal("expecting to match", err)
		}
		if diff != "FullMatch" {
			if _, err := spy.Compare(value, true); err != nil {
				log.Fatal(err)
			}
			t.Fatal("structure expecting to be FullMatch, instead", diff)
		}
		if value != test.ExpectedValue {
			t.Fatalf("expecting value %s, got %s", test.ExpectedValue, value)
		}
	}
}

// TestNonce is a test function for testing the Nonce functionality.
//
// It initializes a test configuration, sets up a test data set, and then performs a series of tests.
// The tests involve creating a spy object, modifying the test configuration provider, and calling the Nonce function.
// The expected result is a successful response from the Nonce function and a matching value from the spy object.
// If any errors occur during the tests, the function will fail and display an error message.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestNonce(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		ContractAddress *felt.Felt
		ExpectedNonce   *felt.Felt
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				ContractAddress: utils.TestHexToFelt(t, "0x0207acc15dc241e7d167e67e30e769719a727d3e0fa47f9e187707289885dfde"),
				ExpectedNonce:   utils.TestHexToFelt(t, "0x0"),
			},
		},
		"testnet": {
			{
				ContractAddress: utils.TestHexToFelt(t, "0x00a3d19d9e80d74dd6140fed379e2c10a21609374811b244cc9d7d1f6d9e0037"),
				ExpectedNonce:   utils.TestHexToFelt(t, "0x0"),
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		nonce, err := testConfig.provider.Nonce(context.Background(), WithBlockTag("latest"), test.ContractAddress)
		if err != nil {
			t.Fatal(err)
		}
		diff, err := spy.Compare(nonce, false)
		if err != nil {
			t.Fatal("expecting to match", err)
		}
		if diff != "FullMatch" {
			if _, err := spy.Compare(nonce, true); err != nil {
				log.Fatal(err)
			}
			t.Fatal("structure expecting to be FullMatch, instead", diff)
		}

		if nonce == nil {
			t.Fatalf("should return a nonce, instead %v", nonce)
		}
		require.Equal(t, test.ExpectedNonce, nonce)
	}
}

// TestEstimateMessageFee is a test function to test the EstimateMessageFee function.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestEstimateMessageFee(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		MsgFromL1
		BlockID
		ExpectedFeeEst FeeEstimate
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				MsgFromL1: MsgFromL1{FromAddress: "0x0", ToAddress: &felt.Zero, Selector: &felt.Zero, Payload: []*felt.Felt{&felt.Zero}},
				BlockID:   BlockID{Tag: "latest"},
				ExpectedFeeEst: FeeEstimate{
					GasConsumed: new(felt.Felt).SetUint64(1),
					GasPrice:    new(felt.Felt).SetUint64(2),
					OverallFee:  new(felt.Felt).SetUint64(3),
				},
			},
		},
		"testnet": {},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		value, err := testConfig.provider.EstimateMessageFee(context.Background(), test.MsgFromL1, test.BlockID)
		if err != nil {
			t.Fatal(err)
		}
		require.Equal(t, *value, test.ExpectedFeeEst)

	}
}

func TestEstimateFee(t *testing.T) {
	testConfig := beforeEach(t)

	testBlockNumber := uint64(15643)
	type testSetType struct {
		txs           []BroadcastTxn
		simFlags      []SimulationFlag
		blockID       BlockID
		expectedResp  []FeeEstimate
		expectedError error
	}
	testSet := map[string][]testSetType{
		"mainnet": {
			{
				txs: []BroadcastTxn{
					InvokeTxnV0{
						Type:    TransactionType_Invoke,
						Version: TransactionV0,
						MaxFee:  utils.TestHexToFelt(t, "0x95e566845d000"),
						FunctionCall: FunctionCall{
							ContractAddress:    utils.TestHexToFelt(t, "0x45e92c365ba0908382bc346159f896e528214470c60ae2cd4038a0fff747b1e"),
							EntryPointSelector: utils.TestHexToFelt(t, "0x15d40a3d6ca2ac30f4031e42be28da9b056fef9bb7357ac5e85627ee876e5ad"),
							Calldata: utils.TestHexArrToFelt(t, []string{
								"0x1",
								"0x4a3621276a83251b557a8140e915599ae8e7b6207b067ea701635c0d509801e",
								"0x2f0b3c5710379609eb5495f1ecd348cb28167711b73609fe565a72734550354",
								"0x0",
								"0x3",
								"0x3",
								"0x697066733a2f2f516d57554c7a475135556a52616953514776717765347931",
								"0x4731796f4757324e6a5a76564e77776a66514577756a",
								"0x0",
								"0x2"}),
						},
						Signature: []*felt.Felt{
							utils.TestHexToFelt(t, "0x63e4618ca2e323a45b9f860f12a4f5c4984648f1d110aa393e79d596d82abcc"),
							utils.TestHexToFelt(t, "0x2844257b088ad4f49e2fe3df1ea6a8530aa2d21d8990112b7e88c4bd0ce9d50"),
						},
					},
				},
				simFlags:      []SimulationFlag{},
				blockID:       BlockID{Number: &testBlockNumber},
				expectedError: nil,
				expectedResp: []FeeEstimate{
					{
						GasConsumed: utils.TestHexToFelt(t, "0x39b8"),
						GasPrice:    utils.TestHexToFelt(t, "0x350da9915"),
						OverallFee:  utils.TestHexToFelt(t, "0xbf62c933b418"),
						FeeUnit:     UnitWei,
					},
				},
			},
			{

				txs: []BroadcastTxn{
					DeployAccountTxn{

						Type:    TransactionType_DeployAccount,
						Version: TransactionV1,
						MaxFee:  utils.TestHexToFelt(t, "0xdec823b1380c"),
						Nonce:   utils.TestHexToFelt(t, "0x0"),
						Signature: []*felt.Felt{
							utils.TestHexToFelt(t, "0x41dbc4b41f6506502a09eb7aea85759de02e91f49d0565776125946e54a2ec6"),
							utils.TestHexToFelt(t, "0x85dcf2bc8e3543071a6657947cc9c157a9f6ad7844a686a975b588199634a9"),
						},
						ContractAddressSalt: utils.TestHexToFelt(t, "0x74ddc51af144d1bd805eb4184d07453d7c4388660270a7851fec387e654a50e"),
						ClassHash:           utils.TestHexToFelt(t, "0x25ec026985a3bf9d0cc1fe17326b245dfdc3ff89b8fde106542a3ea56c5a918"),
						ConstructorCalldata: utils.TestHexArrToFelt(t, []string{
							"0x33434ad846cdd5f23eb73ff09fe6fddd568284a0fb7d1be20ee482f044dabe2",
							"0x79dc0da7c54b95f10aa182ad0a46400db63156920adb65eca2654c0945a463",
							"0x2",
							"0x74ddc51af144d1bd805eb4184d07453d7c4388660270a7851fec387e654a50e",
							"0x0",
						}),
					},
				},
				simFlags:      []SimulationFlag{},
				blockID:       BlockID{Hash: utils.TestHexToFelt(t, "0x1b0df1bafcb826b1fc053495aef5cdc24d0345cbfa1259b15939d01b89dc6d9")},
				expectedError: nil,
				expectedResp: []FeeEstimate{
					{
						GasConsumed: utils.TestHexToFelt(t, "0x15be"),
						GasPrice:    utils.TestHexToFelt(t, "0x378f962c4"),
						OverallFee:  utils.TestHexToFelt(t, "0x4b803e316178"),
						FeeUnit:     UnitWei,
					},
				},
			},
		},
		"mock":    {},
		"testnet": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		resp, err := testConfig.provider.EstimateFee(context.Background(), test.txs, test.simFlags, test.blockID)
		require.Equal(t, test.expectedError, err)
		require.Equal(t, test.expectedResp, resp)
	}
}
