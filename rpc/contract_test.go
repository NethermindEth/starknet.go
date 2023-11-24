package rpc

import (
	"context"
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
//	none
func TestClassAt(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		ContractAddress   *felt.Felt
		ExpectedOperation string
	}
	testSet := map[string][]testSetType{
		"mock": {
			{
				ContractAddress:   utils.TestHexToFelt(t, "0xdeadbeef"),
				ExpectedOperation: "0xdeadbeef",
			},
		},
		"testnet": {
			{
				ContractAddress:   utils.TestHexToFelt(t, "0x6fbd460228d843b7fbef670ff15607bf72e19fa94de21e29811ada167b4ca39"),
				ExpectedOperation: "0x480680017fff8000",
			},
		},
		"mainnet": {
			{
				ContractAddress:   utils.TestHexToFelt(t, "0x028105caf03e1c4eb96b1c18d39d9f03bd53e5d2affd0874792e5bf05f3e529f"),
				ExpectedOperation: "0x20780017fff7ffd",
			},
		},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		resp, err := testConfig.provider.ClassAt(context.Background(), WithBlockTag("latest"), test.ContractAddress)
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
				spy.Compare(class, true)
				t.Fatal("structure expecting to be FullMatch, instead", diff)
			}
			if class.Program == "" {
				t.Fatal("code should exist")
			}
		case ContractClass:
			panic("Not covered")
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
			{
				ContractHash:      utils.TestHexToFelt(t, "0x315e364b162653e5c7b23efd34f8da27ba9c069b68e3042b7d76ce1df890313"),
				ExpectedClassHash: utils.TestHexToFelt(t, "0x493af3546940eb96471cf95ae3a5aa1286217b07edd1e12d00143010ca904b1"),
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
			spy.Compare(classhash, true)
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
			{
				BlockID:         WithBlockTag("pending"),
				ClassHash:       utils.TestHexToFelt(t, "0x493af3546940eb96471cf95ae3a5aa1286217b07edd1e12d00143010ca904b1"),
				ExpectedProgram: "H4sIAAAAAAAA",
			},
			{
				BlockID:                       WithBlockHash(utils.TestHexToFelt(t, "0x464fc8c86a6452536a2e27cb301815e5f8b16a2f6872ba4f3d83701fbe99fb3")),
				ClassHash:                     utils.TestHexToFelt(t, "0x011fbe1adeb2afdf5b545f583f8b5a64fb35905f987d249193ad8185f6fcf571"),
				ExpectedEntryPointConstructor: SierraEntryPoint{FunctionIdx: 16, Selector: utils.TestHexToFelt(t, "0x28ffe4ff0f226a9107253e17a904099aa4f63a02a5621de0576e5aa71bc5194")},
			},
		},
		"mainnet": {},
	}[testEnv]

	for _, test := range testSet {
		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy
		resp, err := testConfig.provider.Class(context.Background(), WithBlockTag("latest"), test.ClassHash)
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
				spy.Compare(class, true)
				t.Fatal("structure expecting to be FullMatch, instead", diff)
			}

			if !strings.HasPrefix(class.Program, test.ExpectedProgram) {
				t.Fatal("code should exist")
			}
		case ContractClass:
			require.Equal(t, class.EntryPointsByType.Constructor, test.ExpectedEntryPointConstructor)
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
				ContractHash:  utils.TestHexToFelt(t, "0x6fbd460228d843b7fbef670ff15607bf72e19fa94de21e29811ada167b4ca39"),
				StorageKey:    "balance",
				Block:         WithBlockTag("latest"),
				ExpectedValue: "0x1e240",
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
			spy.Compare(value, true)
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
				ContractAddress: utils.TestHexToFelt(t, "0x0207acc15dc241e7d167e67e30e769719a727d3e0fa47f9e187707289885dfde"),
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
			spy.Compare(nonce, true)
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
