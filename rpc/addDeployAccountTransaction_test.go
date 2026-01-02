package rpc

import (
	"encoding/json"
	"testing"

	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestAddDeployAccountTransaction tests the AddDeployAccountTransaction function.
func TestAddDeployAccountTransaction(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv)

	type testSetType struct {
		Description   string
		DeployTxn     *BroadcastDeployAccountTxnV3
		ExpectedError *RPCError

		// there are multiple errors that could be returned by the function, and
		// this is a way to specify which error we want to test in the Mock environment.
		// It's better than modifying the `DeclareTxn`, having to create a new variable
		// for each error variant, and compare inside the mock `DoAndReturn` function.
		ErrorIndex int
	}

	deployTxn := internalUtils.TestUnmarshalJSONFileToType[*BroadcastDeployAccountTxnV3](
		t,
		"./testData/addTxn/sepoliaDeployAccount.json",
		"result",
	)

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				Description: "normal call",
				DeployTxn:   deployTxn,
			},
			{
				Description:   "error - insufficient account balance",
				DeployTxn:     deployTxn,
				ExpectedError: ErrInsufficientAccountBalance,
				ErrorIndex:    1,
			},
			{
				Description:   "error - insufficient resources for validate",
				DeployTxn:     deployTxn,
				ExpectedError: ErrInsufficientResourcesForValidate,
				ErrorIndex:    2,
			},
			{
				Description:   "error - invalid transaction nonce",
				DeployTxn:     deployTxn,
				ExpectedError: ErrInvalidTransactionNonce,
				ErrorIndex:    3,
			},
			{
				Description:   "error - replacement transaction underpriced",
				DeployTxn:     deployTxn,
				ExpectedError: ErrReplacementTransactionUnderpriced,
				ErrorIndex:    4,
			},
			{
				Description:   "error - fee below minimum",
				DeployTxn:     deployTxn,
				ExpectedError: ErrFeeBelowMinimum,
				ErrorIndex:    5,
			},
			{
				Description:   "error - validation failure",
				DeployTxn:     deployTxn,
				ExpectedError: ErrValidationFailure,
				ErrorIndex:    6,
			},
			{
				Description:   "error - non account",
				DeployTxn:     deployTxn,
				ExpectedError: ErrNonAccount,
				ErrorIndex:    7,
			},
			{
				Description:   "error - duplicate tx",
				DeployTxn:     deployTxn,
				ExpectedError: ErrDuplicateTx,
				ErrorIndex:    8,
			},
			{
				Description:   "error - unsupported tx version",
				DeployTxn:     deployTxn,
				ExpectedError: ErrUnsupportedTxVersion,
				ErrorIndex:    9,
			},
			{
				Description:   "error - class hash not found",
				DeployTxn:     deployTxn,
				ExpectedError: ErrClassHashNotFound,
				ErrorIndex:    10,
			},
		},
		tests.TestnetEnv: {
			{
				Description: "normal call - with error",
				DeployTxn:   deployTxn,
				// this test sends an already sent transaction, and this is the error
				// returned by the node for this case.
				// We do this because it's not feasible to create a new transaction each time.
				// But with this test, we can assure our txn is correctly received by the node.
				ExpectedError: ErrInvalidTransactionNonce,
			},
		},
	}[tests.TEST_ENV]
	for _, test := range testSet {
		t.Run(test.Description, func(t *testing.T) {
			testConfig := BeforeEach(t, false)

			if tests.TEST_ENV == tests.MockEnv {
				testConfig.MockClient.EXPECT().
					CallContextWithSliceArgs(
						t.Context(),
						gomock.Any(),
						"starknet_addDeployAccountTransaction",
						test.DeployTxn,
					).
					DoAndReturn(func(_, result, _ any, _ ...any) error {
						rawResp := result.(*json.RawMessage)

						switch test.ErrorIndex {
						case 1:
							return RPCError{
								Code:    54,
								Message: "Account balance is smaller than the transaction's maximal fee (calculated as the sum of each resource's limit x max price)",
							}
						case 2:
							return RPCError{
								Code:    53,
								Message: "The transaction's resources don't cover validation or the minimal transaction fee",
							}
						case 3:
							return RPCError{
								Code:    52,
								Message: "Invalid transaction nonce",
								Data:    StringErrData(""),
							}
						case 4:
							return RPCError{
								Code:    64,
								Message: "Replacement transaction is underpriced",
							}
						case 5:
							return RPCError{
								Code:    65,
								Message: "Transaction fee below minimum",
							}
						case 6:
							return RPCError{
								Code:    55,
								Message: "Account validation failed",
								Data:    StringErrData(""),
							}
						case 7:
							return RPCError{
								Code:    58,
								Message: "Sender address is not an account contract",
							}
						case 8:
							return RPCError{
								Code:    59,
								Message: "A transaction with the same hash already exists in the mempool",
							}
						case 9:
							return RPCError{
								Code:    61,
								Message: "The transaction version is not supported",
							}
						case 10:
							return RPCError{
								Code:    28,
								Message: "Class hash not found",
							}
						}

						*rawResp = json.RawMessage(`
							{
								"transaction_hash": "0x32b272b6d0d584305a460197aa849b5c7a9a85903b66e9d3e1afa2427ef093e",
								"contract_address": "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"
							}
						`)

						return nil
					}).
					Times(1)
			}

			resp, err := testConfig.Provider.AddDeployAccountTransaction(
				t.Context(),
				test.DeployTxn,
			)
			if test.ExpectedError != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.ExpectedError.Message)

				return
			}
			require.NoError(t, err)

			rawExpectedResp := testConfig.RPCSpy.LastResponse()
			rawResp, err := json.Marshal(resp)
			require.NoError(t, err)
			assert.JSONEq(t, string(rawExpectedResp), string(rawResp))
		})
	}
}
