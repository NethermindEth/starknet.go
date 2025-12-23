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

// TestDeclareTransaction tests the AddDeclareTransaction function.
func TestDeclareTransaction(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		Description   string
		DeclareTxn    *BroadcastDeclareTxnV3
		ExpectedError *RPCError

		// there are multiple errors that could be returned by the function, and
		// this is a way to specify which error we want to test in the Mock environment.
		// It's better than modifying the `DeclareTxn`, having to create a new variable
		// for each error variant, and compare inside the mock `DoAndReturn` function.
		ErrorIndex int
	}

	declareTxn := internalUtils.TestUnmarshalJSONFileToType[BroadcastDeclareTxnV3](
		t,
		"./testData/addTxn/sepoliaDeclare.json",
		"",
	)

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				Description: "normal call",
				DeclareTxn:  declareTxn,
			},
			{
				Description:   "error - class already declared",
				DeclareTxn:    declareTxn,
				ExpectedError: ErrClassAlreadyDeclared,
				ErrorIndex:    1,
			},
			{
				Description:   "error - compilation failed",
				DeclareTxn:    declareTxn,
				ExpectedError: ErrCompilationFailed,
				ErrorIndex:    2,
			},
			{
				Description:   "error - compiled class hash mismatch",
				DeclareTxn:    declareTxn,
				ExpectedError: ErrCompiledClassHashMismatch,
				ErrorIndex:    3,
			},
			{
				Description:   "error - insufficient account balance",
				DeclareTxn:    declareTxn,
				ExpectedError: ErrInsufficientAccountBalance,
				ErrorIndex:    4,
			},
			{
				Description:   "error - insufficient resources for validate",
				DeclareTxn:    declareTxn,
				ExpectedError: ErrInsufficientResourcesForValidate,
				ErrorIndex:    5,
			},
			{
				Description:   "error - invalid transaction nonce",
				DeclareTxn:    declareTxn,
				ExpectedError: ErrInvalidTransactionNonce,
				ErrorIndex:    6,
			},
			{
				Description:   "error - replacement transaction underpriced",
				DeclareTxn:    declareTxn,
				ExpectedError: ErrReplacementTransactionUnderpriced,
				ErrorIndex:    7,
			},
			{
				Description:   "error - fee below minimum",
				DeclareTxn:    declareTxn,
				ExpectedError: ErrFeeBelowMinimum,
				ErrorIndex:    8,
			},
			{
				Description:   "error - validation failure",
				DeclareTxn:    declareTxn,
				ExpectedError: ErrValidationFailure,
				ErrorIndex:    9,
			},
			{
				Description:   "error - non account",
				DeclareTxn:    declareTxn,
				ExpectedError: ErrNonAccount,
				ErrorIndex:    10,
			},
			{
				Description:   "error - duplicate tx",
				DeclareTxn:    declareTxn,
				ExpectedError: ErrDuplicateTx,
				ErrorIndex:    11,
			},
			{
				Description:   "error - contract class size too large",
				DeclareTxn:    declareTxn,
				ExpectedError: ErrContractClassSizeTooLarge,
				ErrorIndex:    12,
			},
			{
				Description:   "error - unsupported tx version",
				DeclareTxn:    declareTxn,
				ExpectedError: ErrUnsupportedTxVersion,
				ErrorIndex:    13,
			},
			{
				Description:   "error - unsupported contract class version",
				DeclareTxn:    declareTxn,
				ExpectedError: ErrUnsupportedContractClassVersion,
				ErrorIndex:    14,
			},
		},
		tests.TestnetEnv: {
			{
				Description: "normal call - with error",
				DeclareTxn:  declareTxn,
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
			testConfig.MockClient.EXPECT().
				CallContextWithSliceArgs(
					t.Context(),
					gomock.Any(),
					"starknet_addDeclareTransaction",
					test.DeclareTxn,
				).
				DoAndReturn(func(_, result, _ any, _ ...any) error {
					rawResp := result.(*json.RawMessage)

					switch test.ErrorIndex {
					case 1:
						return RPCError{
							Code:    51,
							Message: "Class already declared",
						}
					case 2:
						return RPCError{
							Code:    56,
							Message: "Compilation failed",
							Data:    StringErrData(""),
						}
					case 3:
						return RPCError{
							Code:    60,
							Message: "The compiled class hash did not match the one supplied in the transaction",
						}
					case 4:
						return RPCError{
							Code:    54,
							Message: "Account balance is smaller than the transaction's maximal fee (calculated as the sum of each resource's limit x max price)",
						}
					case 5:
						return RPCError{
							Code:    53,
							Message: "The transaction's resources don't cover validation or the minimal transaction fee",
						}
					case 6:
						return RPCError{
							Code:    52,
							Message: "Invalid transaction nonce",
							Data:    StringErrData(""),
						}
					case 7:
						return RPCError{
							Code:    64,
							Message: "Replacement transaction is underpriced",
						}
					case 8:
						return RPCError{
							Code:    65,
							Message: "Transaction fee below minimum",
						}
					case 9:
						return RPCError{
							Code:    55,
							Message: "Account validation failed",
							Data:    StringErrData(""),
						}
					case 10:
						return RPCError{
							Code:    58,
							Message: "Sender address is not an account contract",
						}
					case 11:
						return RPCError{
							Code:    59,
							Message: "A transaction with the same hash already exists in the mempool",
						}
					case 12:
						return RPCError{
							Code:    57,
							Message: "Contract class size is too large",
						}
					case 13:
						return RPCError{
							Code:    61,
							Message: "The transaction version is not supported",
						}
					case 14:
						return RPCError{
							Code:    62,
							Message: "The contract class version is not supported",
						}
					}

					*rawResp = json.RawMessage(`
						{
							"transaction_hash": "0x41d1f5206ef58a443e7d3d1ca073171ec25fa75313394318fc83a074a6631c3",
							"class_hash": "0x5d68906f23c7e96713002a9ef6a7b1b6ec19e18c31a32710446d87b2aca762d"
						}
					`)

					return nil
				}).
				Times(1)

			resp, err := testConfig.Provider.AddDeclareTransaction(
				t.Context(),
				test.DeclareTxn,
			)
			if test.ExpectedError != nil {
				require.Error(t, err)
				rpcErr, ok := err.(*RPCError)
				require.True(t, ok)
				assert.Equal(t, test.ExpectedError.Code, rpcErr.Code)
				assert.Equal(t, test.ExpectedError.Message, rpcErr.Message)

				return
			}
			require.NoError(t, err)

			rawExpectedResp := testConfig.Spy.LastResponse()
			rawResp, err := json.Marshal(resp)
			require.NoError(t, err)
			assert.JSONEq(t, string(rawExpectedResp), string(rawResp))
		})
	}
}

// TestAddInvokeTransaction tests the AddInvokeTransaction function.
func TestAddInvokeTransaction(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		Description   string
		InvokeTxn     *BroadcastInvokeTxnV3
		ExpectedError *RPCError

		// there are multiple errors that could be returned by the function, and
		// this is a way to specify which error we want to test in the Mock environment.
		// It's better than modifying the `DeclareTxn`, having to create a new variable
		// for each error variant, and compare inside the mock `DoAndReturn` function.
		ErrorIndex int
	}

	invokeTxn := internalUtils.TestUnmarshalJSONFileToType[BroadcastInvokeTxnV3](
		t,
		"./testData/addTxn/sepoliaInvoke.json",
		"",
	)

	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				Description: "normal call",
				InvokeTxn:   invokeTxn,
			},
			{
				Description:   "error - insufficient account balance",
				InvokeTxn:     invokeTxn,
				ExpectedError: ErrInsufficientAccountBalance,
				ErrorIndex:    1,
			},
			{
				Description:   "error - insufficient resources for validate",
				InvokeTxn:     invokeTxn,
				ExpectedError: ErrInsufficientResourcesForValidate,
				ErrorIndex:    2,
			},
			{
				Description:   "error - invalid transaction nonce",
				InvokeTxn:     invokeTxn,
				ExpectedError: ErrInvalidTransactionNonce,
				ErrorIndex:    3,
			},
			{
				Description:   "error - replacement transaction underpriced",
				InvokeTxn:     invokeTxn,
				ExpectedError: ErrReplacementTransactionUnderpriced,
				ErrorIndex:    4,
			},
			{
				Description:   "error - fee below minimum",
				InvokeTxn:     invokeTxn,
				ExpectedError: ErrFeeBelowMinimum,
				ErrorIndex:    5,
			},
			{
				Description:   "error - validation failure",
				InvokeTxn:     invokeTxn,
				ExpectedError: ErrValidationFailure,
				ErrorIndex:    6,
			},
			{
				Description:   "error - non account",
				InvokeTxn:     invokeTxn,
				ExpectedError: ErrNonAccount,
				ErrorIndex:    7,
			},
			{
				Description:   "error - duplicate tx",
				InvokeTxn:     invokeTxn,
				ExpectedError: ErrDuplicateTx,
				ErrorIndex:    8,
			},
			{
				Description:   "error - unsupported tx version",
				InvokeTxn:     invokeTxn,
				ExpectedError: ErrUnsupportedTxVersion,
				ErrorIndex:    9,
			},
			{
				Description:   "error - unexpected error",
				InvokeTxn:     invokeTxn,
				ExpectedError: ErrUnexpectedError,
				ErrorIndex:    10,
			},
		},
		tests.TestnetEnv: {
			{
				Description: "normal call - with error",
				InvokeTxn:   invokeTxn,
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
			testConfig.MockClient.EXPECT().
				CallContextWithSliceArgs(
					t.Context(),
					gomock.Any(),
					"starknet_addInvokeTransaction",
					test.InvokeTxn,
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
							Code:    63,
							Message: "An unexpected error occurred",
							Data:    StringErrData(""),
						}
					}

					*rawResp = json.RawMessage(`
						{
							"transaction_hash": "0x49728601e0bb2f48ce506b0cbd9c0e2a9e50d95858aa41463f46386dca489fd"
						}
					`)

					return nil
				}).
				Times(1)

			resp, err := testConfig.Provider.AddInvokeTransaction(
				t.Context(),
				test.InvokeTxn,
			)
			if test.ExpectedError != nil {
				require.Error(t, err)
				rpcErr, ok := err.(*RPCError)
				require.True(t, ok)
				assert.Equal(t, test.ExpectedError.Code, rpcErr.Code)
				assert.Equal(t, test.ExpectedError.Message, rpcErr.Message)

				return
			}
			require.NoError(t, err)

			rawExpectedResp := testConfig.Spy.LastResponse()
			rawResp, err := json.Marshal(resp)
			require.NoError(t, err)
			assert.JSONEq(t, string(rawExpectedResp), string(rawResp))
		})
	}
}

// TestAddDeployAccountTransaction tests the AddDeployAccountTransaction function.
func TestAddDeployAccountTransaction(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv, tests.TestnetEnv)

	testConfig := BeforeEach(t, false)

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

	deployTxn := internalUtils.TestUnmarshalJSONFileToType[BroadcastDeployAccountTxnV3](
		t,
		"./testData/addTxn/sepoliaDeployAccount.json",
		"",
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

			resp, err := testConfig.Provider.AddDeployAccountTransaction(
				t.Context(),
				test.DeployTxn,
			)
			if test.ExpectedError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.ExpectedError.Error())

				return
			}
			require.NoError(t, err)

			rawExpectedResp := testConfig.Spy.LastResponse()
			rawResp, err := json.Marshal(resp)
			require.NoError(t, err)
			assert.JSONEq(t, string(rawExpectedResp), string(rawResp))
		})
	}
}
