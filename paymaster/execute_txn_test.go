package paymaster

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecuteTransaction(t *testing.T) {
	t.Parallel()

	t.Run("integration", func(t *testing.T) {
		tests.RunTestOn(t, tests.IntegrationEnv)

		t.Run("execute deploy transaction", func(t *testing.T) {
			t.Parallel()

			pm, spy := SetupPaymaster(t, true)

			deployTxn := buildDeployTxn(t, pm)
			assert.NotNil(t, deployTxn)

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

			rawResp, err := json.Marshal(resp)
			require.NoError(t, err)
			assert.JSONEq(t, string(spy.LastResponse()), string(rawResp))
		})
	})

}

// buildDeployTxn builds a deploy transaction calling the paymaster_buildTransaction method
func buildDeployTxn(t *testing.T, pm *Paymaster) (resp *BuildTransactionResponse) {
	t.Helper()

	_, pub, _ := account.GetRandomKeys()

	// OZ account class hash
	classHash := internalUtils.TestHexToFelt(t, "0x61dac032f228abef9c6626f995015233097ae253a7f72d68552db02f2971b8f")
	constructorCalldata := []*felt.Felt{pub}
	precAddress := account.PrecomputeAccountAddress(internalUtils.RANDOM_FELT, classHash, constructorCalldata)

	deploymentData := &AccDeploymentData{
		Address:             precAddress,
		ClassHash:           classHash,
		Salt:                internalUtils.RANDOM_FELT,
		ConstructorCalldata: constructorCalldata,
		SignatureData:       []*felt.Felt{},
		Version:             2,
	}

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

	return resp
}
