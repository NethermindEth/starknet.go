package rpc

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/test-go/testify/require"
)

// TestDeclareTransaction is a test function for declaring a transaction.
//
// It sets up different test sets based on the test environment and runs the tests accordingly.
// Each test set contains a transaction hash, class hash, and an expected error.
// The function reads a JSON file, unmarshals it into an `AddDeclareTxnInput` struct,
// and performs various tests using the `AddDeclareTransaction` function.
// It checks if the error matches the expected error and if the transaction hash matches the expected hash.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestDeclareTransaction(t *testing.T) {

	testConfig := beforeEach(t)

	type testSetType struct {
		TransactionHash *felt.Felt
		ClassHash       *felt.Felt
		ExpectedError   string
	}
	testSet := map[string][]testSetType{
		"devnet":  {},
		"mainnet": {},
		"mock":    {},
		"testnet": {{
			TransactionHash: utils.TestHexToFelt(t, "0x55b094dc5c84c2042e067824f82da90988674314d37e45cb0032aca33d6e0b9"),
			ClassHash:       utils.TestHexToFelt(t, "0xdeadbeef"),
			ExpectedError:   "Invalid Params",
		}},
	}[testEnv]

	for _, test := range testSet {

		declareTxJSON, err := os.ReadFile("./tests/write/declareTx.json")
		if err != nil {
			t.Fatal("should be able to read file", err)
		}

		var declareTx AddDeclareTxnInput
		err = json.Unmarshal(declareTxJSON, &declareTx)
		require.Nil(t, err, "Error unmarshalling decalreTx")

		spy := NewSpy(testConfig.provider.c)
		testConfig.provider.c = spy

		// TODO: test transaction against client that supports RPC method (currently Sequencer uses
		// "sierra_program" instead of "program"
		dec, err := testConfig.provider.AddDeclareTransaction(context.Background(), declareTx)
		if err != nil {
			require.Equal(t, err.Error(), test.ExpectedError)
			continue
		}
		if dec.TransactionHash != test.TransactionHash {
			t.Fatalf("classHash does not match expected, current: %s", dec.ClassHash)
		}

	}
}

// TestAddInvokeTransaction is a test function that checks the AddInvokeTransaction functionality.
//
// It initializes a test configuration and defines sets of test cases for different environments.
// Each test case includes an InvokeTxnV1 object, an expected AddInvokeTransactionResponse, and an expected RPCError.
// The function iterates through the test cases, invokes AddInvokeTransaction, and compares the response and error with the expected values.
//
// Parameters:
// - t: the testing object for running the test cases
// Returns:
//
//	none
func TestAddInvokeTransaction(t *testing.T) {

	testConfig := beforeEach(t)

	type testSetType struct {
		InvokeTx      InvokeTxnType
		ExpectedResp  AddInvokeTransactionResponse
		ExpectedError RPCError
	}
	testSet := map[string][]testSetType{
		"devnet":  {},
		"mainnet": {},
		"mock": {
			{
				InvokeTx:     InvokeTxnV1{SenderAddress: new(felt.Felt).SetUint64(123)},
				ExpectedResp: AddInvokeTransactionResponse{&felt.Zero},
				ExpectedError: RPCError{
					code:    ErrUnexpectedError.code,
					message: ErrUnexpectedError.message,
					data:    "Something crazy happened"},
			},
			{
				InvokeTx:      InvokeTxnV1{},
				ExpectedResp:  AddInvokeTransactionResponse{utils.TestHexToFelt(t, "0xdeadbeef")},
				ExpectedError: RPCError{},
			},
			{
				InvokeTx: InvokeTxnV3{
					Type:    TransactionType_Invoke,
					Version: TransactionV3,
					Signature: []*felt.Felt{
						utils.TestHexToFelt(t, "0x71a9b2cd8a8a6a4ca284dcddcdefc6c4fd20b92c1b201bd9836e4ce376fad16"),
						utils.TestHexToFelt(t, "0x6bef4745194c9447fdc8dd3aec4fc738ab0a560b0d2c7bf62fbf58aef3abfc5"),
					},
					Nonce:         utils.TestHexToFelt(t, "0xe97"),
					NonceDataMode: DAModeL1,
					FeeMode:       DAModeL1,
					ResourceBounds: ResourceBoundsMapping{
						L1Gas: ResourceBounds{
							MaxAmount:       utils.TestHexToFelt(t, "0x186a0"),
							MaxPricePerUnit: utils.TestHexToFelt(t, "0x5af3107a4000"),
						},
						L2Gas: ResourceBounds{
							MaxAmount:       utils.TestHexToFelt(t, "0x0"),
							MaxPricePerUnit: utils.TestHexToFelt(t, "0x0"),
						},
					},
					Tip:           new(felt.Felt),
					PayMasterData: []*felt.Felt{},
					SenderAddress: utils.TestHexToFelt(t, "0x3f6f3bc663aedc5285d6013cc3ffcbc4341d86ab488b8b68d297f8258793c41"),
					Calldata: []*felt.Felt{
						utils.TestHexToFelt(t, "0x2"),
						utils.TestHexToFelt(t, "0x450703c32370cf7ffff540b9352e7ee4ad583af143a361155f2b485c0c39684"),
						utils.TestHexToFelt(t, "0x27c3334165536f239cfd400ed956eabff55fc60de4fb56728b6a4f6b87db01c"),
						utils.TestHexToFelt(t, "0x0"),
						utils.TestHexToFelt(t, "0x4"),
						utils.TestHexToFelt(t, "0x4c312760dfd17a954cdd09e76aa9f149f806d88ec3e402ffaf5c4926f568a42"),
						utils.TestHexToFelt(t, "0x5df99ae77df976b4f0e5cf28c7dcfe09bd6e81aab787b19ac0c08e03d928cf"),
						utils.TestHexToFelt(t, "0x4"),
						utils.TestHexToFelt(t, "0x1"),
						utils.TestHexToFelt(t, "0x5"),
						utils.TestHexToFelt(t, "0x450703c32370cf7ffff540b9352e7ee4ad583af143a361155f2b485c0c39684"),
						utils.TestHexToFelt(t, "0x5df99ae77df976b4f0e5cf28c7dcfe09bd6e81aab787b19ac0c08e03d928cf"),
						utils.TestHexToFelt(t, "0x1"),
						utils.TestHexToFelt(t, "0x7fe4fd616c7fece1244b3616bb516562e230be8c9f29668b46ce0369d5ca829"),
						utils.TestHexToFelt(t, "0x287acddb27a2f9ba7f2612d72788dc96a5b30e401fc1e8072250940e024a587"),
					},
					AccountDeploymentData: []*felt.Felt{},
				},
				ExpectedResp:  AddInvokeTransactionResponse{utils.TestHexToFelt(t, "0x49728601e0bb2f48ce506b0cbd9c0e2a9e50d95858aa41463f46386dca489fd")},
				ExpectedError: RPCError{},
			},
		},
		"testnet": {},
	}[testEnv]

	for _, test := range testSet {
		resp, err := testConfig.provider.AddInvokeTransaction(context.Background(), test.InvokeTx)
		if err != nil {
			require.Equal(t, err, &test.ExpectedError, "AddInvokeTransaction did not give expected error")
		} else {
			require.Equal(t, *resp, test.ExpectedResp)
		}

	}
}
