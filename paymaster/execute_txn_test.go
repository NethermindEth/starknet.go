package paymaster

import (
	"encoding/json"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestExecuteTransaction(t *testing.T) {
	t.Parallel()

	t.Run("integration", func(t *testing.T) {
		t.Parallel()
		tests.RunTestOn(t, tests.IntegrationEnv)

		privKey, pubKey, _, err := curve.GetRandomKeys()
		require.NoError(t, err)

		pubKeyFelt := new(felt.Felt).SetBigInt(pubKey)
		privKeyFelt := new(felt.Felt).SetBigInt(privKey)

		t.Run("execute deploy transaction", func(t *testing.T) {
			t.Parallel()

			pm, spy := SetupPaymaster(t)
			t.Log("paymster successfully initialised")

			deployTxn := buildDeployTxn(t, pm, pubKeyFelt)
			assert.NotNil(t, deployTxn)

			t.Log("executing the deploy transaction in the paymaster")

			request := ExecuteTransactionRequest{
				Transaction: ExecutableUserTransaction{
					Type:       UserTxnDeploy,
					Deployment: deployTxn.Deployment,
				},
				Parameters: UserParameters{
					Version: UserParamV1,
					FeeMode: FeeMode{
						Mode: FeeModeSponsored,
						Tip: &TipPriority{
							Priority: TipPriorityNormal,
						},
					},
				},
			}

			resp, err := pm.ExecuteTransaction(t.Context(), &request)
			require.NoError(t, err)

			t.Log("transaction successfully executed")
			t.Logf("Tracking ID: %s", resp.TrackingID)
			t.Logf("Transaction Hash: %s", resp.TransactionHash)

			rawResp, err := json.Marshal(resp)
			require.NoError(t, err)
			assert.JSONEq(t, string(spy.LastResponse()), string(rawResp))
		})

		t.Run("execute invoke transaction", func(t *testing.T) {
			t.Parallel()

			pm, spy := SetupPaymaster(t)
			t.Log("paymaster successfully initialised")

			privK, _, accAdd := GetStrkAccountData(t)
			t.Log("account data fetched")

			invokeTxn := buildInvokeTxn(t, pm, &accAdd)
			assert.NotNil(t, invokeTxn)

			mshHash, err := invokeTxn.TypedData.GetMessageHash(accAdd.String())
			require.NoError(t, err)
			t.Log("message hash:", mshHash)

			r, s, err := curve.SignFelts(mshHash, &privK)
			require.NoError(t, err)
			t.Log("typed data signature:", r, s)

			t.Log("executing the invoke transaction in the paymaster")

			request := ExecuteTransactionRequest{
				Transaction: ExecutableUserTransaction{
					Type: UserTxnInvoke,
					Invoke: &ExecutableUserInvoke{
						UserAddress: &accAdd,
						TypedData:   invokeTxn.TypedData,
						Signature:   []*felt.Felt{r, s},
					},
				},
				Parameters: UserParameters{
					Version: UserParamV1,
					FeeMode: FeeMode{
						Mode: FeeModeSponsored,
						Tip: &TipPriority{
							Priority: TipPriorityNormal,
						},
					},
				},
			}

			resp, err := pm.ExecuteTransaction(t.Context(), &request)
			require.NoError(t, err)

			t.Log("transaction successfully executed")
			t.Logf("Tracking ID: %s", resp.TrackingID)
			t.Logf("Transaction Hash: %s", resp.TransactionHash)

			rawResp, err := json.Marshal(resp)
			require.NoError(t, err)
			assert.JSONEq(t, string(spy.LastResponse()), string(rawResp))
		})

		t.Run("execute deploy_and_invoke transaction", func(t *testing.T) {
			t.Parallel()

			pm, spy := SetupPaymaster(t)
			t.Log("paymster successfully initialised")

			builtTxn := buildDeployAndInvokeTxn(t, pm, pubKeyFelt)
			assert.NotNil(t, builtTxn)

			accAdd := builtTxn.Deployment.Address
			mshHash, err := builtTxn.TypedData.GetMessageHash(accAdd.String())
			require.NoError(t, err)
			t.Log("message hash:", mshHash)

			r, s, err := curve.SignFelts(mshHash, privKeyFelt)
			require.NoError(t, err)
			t.Log("typed data signature:", r, s)

			t.Log("executing the deploy_and_invoke transaction in the paymaster")

			request := ExecuteTransactionRequest{
				Transaction: ExecutableUserTransaction{
					Type:       UserTxnDeployAndInvoke,
					Deployment: builtTxn.Deployment,
					Invoke: &ExecutableUserInvoke{
						UserAddress: accAdd,
						TypedData:   builtTxn.TypedData,
						Signature:   []*felt.Felt{r, s},
					},
				},
				Parameters: UserParameters{
					Version: UserParamV1,
					FeeMode: FeeMode{
						Mode: FeeModeSponsored,
						Tip: &TipPriority{
							Priority: TipPriorityFast,
						},
					},
				},
			}

			resp, err := pm.ExecuteTransaction(t.Context(), &request)
			require.NoError(t, err)

			t.Log("transaction successfully executed")
			t.Logf("Tracking ID: %s", resp.TrackingID)
			t.Logf("Transaction Hash: %s", resp.TransactionHash)

			rawResp, err := json.Marshal(resp)
			require.NoError(t, err)
			assert.JSONEq(t, string(spy.LastResponse()), string(rawResp))
		})
	})

	t.Run("mock", func(t *testing.T) {
		t.Parallel()
		tests.RunTestOn(t, tests.MockEnv)

		pubKey := internalUtils.TestHexToFelt(
			t,
			"0x1cf6046c81f47d488c528e52066482f6756029bed10cf5df35608bb8eebac9",
		)

		t.Run("execute deploy transaction", func(t *testing.T) {
			t.Parallel()
			t.Log("building deploy request")

			deploymentData := createDeploymentData(t, pubKey)

			request := ExecuteTransactionRequest{
				Transaction: ExecutableUserTransaction{
					Type:       UserTxnDeploy,
					Deployment: deploymentData,
				},
				Parameters: UserParameters{
					Version: UserParamV1,
					FeeMode: FeeMode{
						Mode: FeeModeSponsored,
					},
				},
			}

			t.Log("asserting the request marshalled is equal to the expected request")
			expectedReqs := internalUtils.TestUnmarshalJSONFileToType[[]json.RawMessage](
				t,
				"testdata/execute_txn/deploy-request.json",
				"params",
			)
			expectedReq := expectedReqs[0]

			rawReq, err := json.Marshal(request)
			require.NoError(t, err)

			assert.JSONEq(t, string(expectedReq), string(rawReq))

			t.Log("asserting the response marshalled is equal to the expected response")
			expectedResp := internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
				t,
				"testdata/execute_txn/response.json",
				"result",
			)

			var response ExecuteTransactionResponse
			err = json.Unmarshal(expectedResp, &response)
			require.NoError(t, err)

			t.Log("setting up mock paymaster and mock call")
			pm := SetupMockPaymaster(t)
			pm.c.EXPECT().CallContextWithSliceArgs(
				t.Context(),
				gomock.AssignableToTypeOf(new(ExecuteTransactionResponse)),
				"paymaster_executeTransaction",
				&request,
			).Return(nil).
				SetArg(1, response)

			resp, err := pm.ExecuteTransaction(t.Context(), &request)
			require.NoError(t, err)

			rawResp, err := json.Marshal(resp)
			require.NoError(t, err)
			assert.JSONEq(t, string(expectedResp), string(rawResp))
		})

		t.Run("execute invoke transaction", func(t *testing.T) {
			t.Parallel()
			t.Log("building invoke request")

			t.Log("asserting the request marshalled is equal to the expected request")
			expectedReqs := internalUtils.TestUnmarshalJSONFileToType[[]json.RawMessage](
				t,
				"testdata/execute_txn/invoke-request.json",
				"params",
			)
			expectedReq := expectedReqs[0]

			// since the invoke request is more complex, let's take it from the file
			var request ExecuteTransactionRequest
			err := json.Unmarshal(expectedReq, &request)
			require.NoError(t, err)

			rawReq, err := json.Marshal(request)
			require.NoError(t, err)

			// assert if the MarshalJSON is correct
			assert.JSONEq(t, string(expectedReq), string(rawReq))

			t.Log("asserting the response marshalled is equal to the expected response")
			expectedResp := internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
				t,
				"testdata/execute_txn/response.json",
				"result",
			)

			var response ExecuteTransactionResponse
			err = json.Unmarshal(expectedResp, &response)
			require.NoError(t, err)

			t.Log("setting up mock paymaster and mock call")
			pm := SetupMockPaymaster(t)
			pm.c.EXPECT().CallContextWithSliceArgs(
				t.Context(),
				gomock.AssignableToTypeOf(new(ExecuteTransactionResponse)),
				"paymaster_executeTransaction",
				&request,
			).Return(nil).
				SetArg(1, response)

			t.Log("executing the invoke transaction in the mock paymaster")
			resp, err := pm.ExecuteTransaction(t.Context(), &request)
			require.NoError(t, err)

			rawResp, err := json.Marshal(resp)
			require.NoError(t, err)
			assert.JSONEq(t, string(expectedResp), string(rawResp))
		})

		t.Run("execute deploy_and_invoke transaction", func(t *testing.T) {
			t.Parallel()
			t.Log("building deploy_and_invoke request")

			t.Log("asserting the request marshalled is equal to the expected request")
			expectedReqs := internalUtils.TestUnmarshalJSONFileToType[[]json.RawMessage](
				t,
				"testdata/execute_txn/deploy_and_invoke-request.json",
				"params",
			)
			expectedReq := expectedReqs[0]

			// since the deploy_and_invoke request is more complex, let's take it from the file
			var request ExecuteTransactionRequest
			err := json.Unmarshal(expectedReq, &request)
			require.NoError(t, err)

			rawReq, err := json.Marshal(request)
			require.NoError(t, err)

			// assert if the MarshalJSON is correct
			assert.JSONEq(t, string(expectedReq), string(rawReq))

			t.Log("asserting the response marshalled is equal to the expected response")
			expectedResp := internalUtils.TestUnmarshalJSONFileToType[json.RawMessage](
				t,
				"testdata/execute_txn/response.json",
				"result",
			)

			var response ExecuteTransactionResponse
			err = json.Unmarshal(expectedResp, &response)
			require.NoError(t, err)

			t.Log("setting up mock paymaster and mock call")
			pm := SetupMockPaymaster(t)
			pm.c.EXPECT().CallContextWithSliceArgs(
				t.Context(),
				gomock.AssignableToTypeOf(new(ExecuteTransactionResponse)),
				"paymaster_executeTransaction",
				&request,
			).Return(nil).
				SetArg(1, response)

			t.Log("executing the deploy_and_invoke transaction in the mock paymaster")
			resp, err := pm.ExecuteTransaction(t.Context(), &request)
			require.NoError(t, err)

			rawResp, err := json.Marshal(resp)
			require.NoError(t, err)
			assert.JSONEq(t, string(expectedResp), string(rawResp))
		})
	})
}

// same as account.PrecomputeAccountAddress, but to avoid circular dependency
func precomputeAccountAddress(
	salt,
	classHash *felt.Felt,
	constructorCalldata []*felt.Felt,
) *felt.Felt {
	return contracts.PrecomputeAddress(&felt.Zero, salt, classHash, constructorCalldata)
}

// createDeploymentData creates the deployment data for a deploy transaction
func createDeploymentData(t *testing.T, pubKey *felt.Felt) *AccountDeploymentData {
	t.Helper()

	t.Log("creating deployment data")

	// Argent account class hash that supports outside executions
	classHash := internalUtils.TestHexToFelt(
		t,
		"0x036078334509b514626504edc9fb252328d1a240e4e948bef8d0c08dff45927f",
	)
	constructorCalldata := []*felt.Felt{&felt.Zero, pubKey, new(felt.Felt).SetUint64(1)}
	precAddress := precomputeAccountAddress(internalUtils.DeadBeef, classHash, constructorCalldata)
	t.Log("precomputed address:", precAddress)

	deploymentData := &AccountDeploymentData{
		Address:       precAddress,
		ClassHash:     classHash,
		Salt:          internalUtils.DeadBeef,
		Calldata:      constructorCalldata,
		SignatureData: []*felt.Felt{},
		Version:       Cairo1,
	}
	t.Logf("deployment data: %+v", deploymentData)

	return deploymentData
}

// buildDeployTxn builds a deploy transaction calling the paymaster_buildTransaction method
//

func buildDeployTxn(
	t *testing.T,
	pm *Paymaster,
	pubKey *felt.Felt,
) *BuildTransactionResponse {
	t.Helper()

	t.Log("building deploy transaction")
	t.Log("public key:", pubKey)

	deploymentData := createDeploymentData(t, pubKey)

	t.Log("calling paymaster_buildTransaction method")

	resp, err := pm.BuildTransaction(t.Context(), &BuildTransactionRequest{
		Transaction: UserTransaction{
			Type:       UserTxnDeploy,
			Deployment: deploymentData,
		},
		Parameters: UserParameters{
			Version: UserParamV1,
			FeeMode: FeeMode{
				Mode: FeeModeSponsored,
			},
		},
	})
	require.NoError(t, err)
	t.Log("deploy transaction successfully built by the paymaster")

	return &resp
}

// createInvokeData creates the invoke data for an invoke transaction
func createInvokeData(t *testing.T, accAdd *felt.Felt) *UserInvoke {
	t.Helper()

	t.Log("creating invoke data")

	invokeData := &UserInvoke{
		UserAddress: accAdd,
		Calls: []Call{
			{
				// same ERC20 contract as in examples/simpleInvoke
				To: internalUtils.TestHexToFelt(
					t,
					"0x0669e24364ce0ae7ec2864fb03eedbe60cfbc9d1c74438d10fa4b86552907d54",
				),
				Selector: internalUtils.GetSelectorFromNameFelt("mint"),
				Calldata: []*felt.Felt{new(felt.Felt).SetUint64(10000), &felt.Zero},
			},
		},
	}
	t.Logf("invoke data: %+v", invokeData)

	return invokeData
}

// buildInvokeTxn builds an invoke transaction calling the paymaster_buildTransaction method
//

func buildInvokeTxn(
	t *testing.T,
	pm *Paymaster,
	accAdd *felt.Felt,
) *BuildTransactionResponse {
	t.Helper()

	t.Log("building deploy transaction")
	t.Log("account address:", accAdd)

	invokeData := createInvokeData(t, accAdd)

	t.Log("calling paymaster_buildTransaction method")

	resp, err := pm.BuildTransaction(t.Context(), &BuildTransactionRequest{
		Transaction: UserTransaction{
			Type:   UserTxnInvoke,
			Invoke: invokeData,
		},
		Parameters: UserParameters{
			Version: UserParamV1,
			FeeMode: FeeMode{
				Mode: FeeModeSponsored,
			},
		},
	})
	require.NoError(t, err)
	t.Log("invoke transaction successfully built by the paymaster")

	return &resp
}

// buildDeployAndInvokeTxn builds a deploy and invoke transaction calling the paymaster_buildTransaction method
func buildDeployAndInvokeTxn(
	t *testing.T,
	pm *Paymaster,
	pubKey *felt.Felt,
) *BuildTransactionResponse {
	t.Helper()

	t.Log("building deploy_and_invoke transaction")
	t.Log("public key:", pubKey)

	deploymentData := createDeploymentData(t, pubKey)
	invokeData := createInvokeData(t, deploymentData.Address)

	t.Log("calling paymaster_buildTransaction method")

	resp, err := pm.BuildTransaction(t.Context(), &BuildTransactionRequest{
		Transaction: UserTransaction{
			Type:       UserTxnDeployAndInvoke,
			Deployment: deploymentData,
			Invoke:     invokeData,
		},
		Parameters: UserParameters{
			Version: UserParamV1,
			FeeMode: FeeMode{
				Mode: FeeModeSponsored,
				Tip: &TipPriority{
					Priority: TipPriorityFast,
				},
			},
		},
	})
	require.NoError(t, err)
	t.Log("deploy_and_invoke transaction successfully built by the paymaster")

	return &resp
}
