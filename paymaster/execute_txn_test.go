package paymaster

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecuteTransaction(t *testing.T) {
	t.Parallel()

	t.Run("integration", func(t *testing.T) {
		tests.RunTestOn(t, tests.IntegrationEnv)

		_, pubKey, _, _ := curve.GetRandomKeys()
		pubKeyFelt := new(felt.Felt).SetBigInt(pubKey)

		t.Run("execute deploy transaction", func(t *testing.T) {
			t.Parallel()

			pm, spy := SetupPaymaster(t)
			t.Log("paymster successfully initialized")

			deployTxn := buildDeployTxn(t, pm, pubKeyFelt)
			assert.NotNil(t, deployTxn)

			t.Log("executing the deploy transaction in the paymaster")

			request := ExecuteTransactionRequest{
				Transaction: &ExecutableUserTransaction{
					Type:       UserTxnDeploy,
					Deployment: deployTxn.Deployment,
				},
				Parameters: &UserParameters{
					Version: UserParamV1,
					FeeMode: FeeMode{
						Mode: FeeModeSponsored,
						Tip: &TipPriority{
							Priority: TipPriorityNormal,
						},
					},
				},
			}

			resp, err := pm.ExecuteTransaction(context.Background(), &request)
			require.NoError(t, err)

			t.Log("transaction successfully executed")
			t.Logf("Tracking ID: %s", resp.TrackingId)
			t.Logf("Transaction Hash: %s", resp.TransactionHash)

			rawResp, err := json.Marshal(resp)
			require.NoError(t, err)
			assert.JSONEq(t, string(spy.LastResponse()), string(rawResp))
		})

		t.Run("execute invoke transaction", func(t *testing.T) {
			t.Parallel()

			pm, spy := SetupPaymaster(t)
			t.Log("paymster successfully initialized")

			privK, _, accAdd := GetStrkAccountData(t)
			t.Log("account data fetched")

			invokeTxn := buildInvokeTxn(t, pm, accAdd)
			assert.NotNil(t, invokeTxn)

			mshHash, err := invokeTxn.TypedData.GetMessageHash(accAdd.String())
			require.NoError(t, err)

			r, s, err := curve.SignFelts(mshHash, privK)
			require.NoError(t, err)
			t.Log("typed data signed")

			t.Log("executing the invoke transaction in the paymaster")

			request := ExecuteTransactionRequest{
				Transaction: &ExecutableUserTransaction{
					Type: UserTxnInvoke,
					Invoke: &ExecutableUserInvoke{
						UserAddress: accAdd,
						TypedData:   invokeTxn.TypedData,
						Signature:   []*felt.Felt{r, s},
					},
				},
				Parameters: &UserParameters{
					Version: UserParamV1,
					FeeMode: FeeMode{
						Mode: FeeModeSponsored,
						Tip: &TipPriority{
							Priority: TipPriorityNormal,
						},
					},
				},
			}

			resp, err := pm.ExecuteTransaction(context.Background(), &request)
			require.NoError(t, err)

			t.Log("transaction successfully executed")
			t.Logf("Tracking ID: %s", resp.TrackingId)
			t.Logf("Transaction Hash: %s", resp.TransactionHash)

			rawResp, err := json.Marshal(resp)
			require.NoError(t, err)
			assert.JSONEq(t, string(spy.LastResponse()), string(rawResp))
		})
	})

}

// same as account.PrecomputeAccountAddress, but to avoid circular dependency
func precomputeAccountAddress(salt, classHash *felt.Felt, constructorCalldata []*felt.Felt) *felt.Felt {
	return contracts.PrecomputeAddress(&felt.Zero, salt, classHash, constructorCalldata)
}

// buildDeployTxn builds a deploy transaction calling the paymaster_buildTransaction method
func buildDeployTxn(t *testing.T, pm *Paymaster, pubKey *felt.Felt) (resp *BuildTransactionResponse) {
	t.Helper()

	t.Log("building deploy transaction")
	t.Log("public key:", pubKey)

	// OZ account class hash
	classHash := internalUtils.TestHexToFelt(t, "0x61dac032f228abef9c6626f995015233097ae253a7f72d68552db02f2971b8f")
	constructorCalldata := []*felt.Felt{pubKey}
	precAddress := precomputeAccountAddress(internalUtils.RANDOM_FELT, classHash, constructorCalldata)
	t.Log("precomputed address:", precAddress)

	deploymentData := &AccDeploymentData{
		Address:             precAddress,
		ClassHash:           classHash,
		Salt:                internalUtils.RANDOM_FELT,
		ConstructorCalldata: constructorCalldata,
		SignatureData:       []*felt.Felt{},
		Version:             2,
	}
	t.Logf("deployment data: %+v", deploymentData)

	t.Log("calling paymaster_buildTransaction method")
	var err error
	resp, err = pm.BuildTransaction(context.Background(), &BuildTransactionRequest{
		Transaction: &UserTransaction{
			Type:       UserTxnDeploy,
			Deployment: deploymentData,
		},
		Parameters: &UserParameters{
			Version: UserParamV1,
			FeeMode: FeeMode{
				Mode: FeeModeSponsored,
			},
		},
	})
	require.NoError(t, err)
	t.Log("deploy transaction successfully built by the paymaster")

	return resp
}

func buildInvokeTxn(t *testing.T, pm *Paymaster, pubKey *felt.Felt) (resp *BuildTransactionResponse) {
	t.Helper()

	return resp
}
