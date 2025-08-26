package rpc

import (
	"context"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeclareTransaction(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		DeclareTx     BroadcastDeclareTxnV3
		ExpectedResp  AddDeclareTransactionResponse
		ExpectedError *RPCError
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				DeclareTx: BroadcastDeclareTxnV3{},
				ExpectedResp: AddDeclareTransactionResponse{
					Hash: internalUtils.TestHexToFelt(t, "0x41d1f5206ef58a443e7d3d1ca073171ec25fa75313394318fc83a074a6631c3"),
				},
				ExpectedError: nil,
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		resp, err := testConfig.Provider.AddDeclareTransaction(context.Background(), &test.DeclareTx)
		if test.ExpectedError != nil {
			require.Error(t, err)
			rpcErr, ok := err.(*RPCError)
			require.True(t, ok)
			assert.Equal(t, test.ExpectedError.Code, rpcErr.Code)
			assert.Equal(t, test.ExpectedError.Message, rpcErr.Message)

			continue
		}
		require.NoError(t, err)
		assert.Equal(t, resp.Hash.String(), test.ExpectedResp.Hash.String())
	}
}

func TestAddInvokeTransaction(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		InvokeTx      BroadcastInvokeTxnV3
		ExpectedResp  AddInvokeTransactionResponse
		ExpectedError *RPCError
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				InvokeTx: BroadcastInvokeTxnV3{
					Type:    TransactionType_Invoke,
					Version: TransactionV3,
					Signature: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x71a9b2cd8a8a6a4ca284dcddcdefc6c4fd20b92c1b201bd9836e4ce376fad16"),
						internalUtils.TestHexToFelt(t, "0x6bef4745194c9447fdc8dd3aec4fc738ab0a560b0d2c7bf62fbf58aef3abfc5"),
					},
					Nonce:         internalUtils.TestHexToFelt(t, "0xe97"),
					NonceDataMode: DAModeL1,
					FeeMode:       DAModeL1,
					ResourceBounds: &ResourceBoundsMapping{
						L1Gas: ResourceBounds{
							MaxAmount:       "0x186a0",
							MaxPricePerUnit: "0x5af3107a4000",
						},
						L2Gas: ResourceBounds{
							MaxAmount:       "0x0",
							MaxPricePerUnit: "0x0",
						},
					},
					Tip:           "",
					PayMasterData: []*felt.Felt{},
					SenderAddress: internalUtils.TestHexToFelt(t, "0x3f6f3bc663aedc5285d6013cc3ffcbc4341d86ab488b8b68d297f8258793c41"),
					Calldata: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x2"),
						internalUtils.TestHexToFelt(t, "0x450703c32370cf7ffff540b9352e7ee4ad583af143a361155f2b485c0c39684"),
						internalUtils.TestHexToFelt(t, "0x27c3334165536f239cfd400ed956eabff55fc60de4fb56728b6a4f6b87db01c"),
						internalUtils.TestHexToFelt(t, "0x0"),
						internalUtils.TestHexToFelt(t, "0x4"),
						internalUtils.TestHexToFelt(t, "0x4c312760dfd17a954cdd09e76aa9f149f806d88ec3e402ffaf5c4926f568a42"),
						internalUtils.TestHexToFelt(t, "0x5df99ae77df976b4f0e5cf28c7dcfe09bd6e81aab787b19ac0c08e03d928cf"),
						internalUtils.TestHexToFelt(t, "0x4"),
						internalUtils.TestHexToFelt(t, "0x1"),
						internalUtils.TestHexToFelt(t, "0x5"),
						internalUtils.TestHexToFelt(t, "0x450703c32370cf7ffff540b9352e7ee4ad583af143a361155f2b485c0c39684"),
						internalUtils.TestHexToFelt(t, "0x5df99ae77df976b4f0e5cf28c7dcfe09bd6e81aab787b19ac0c08e03d928cf"),
						internalUtils.TestHexToFelt(t, "0x1"),
						internalUtils.TestHexToFelt(t, "0x7fe4fd616c7fece1244b3616bb516562e230be8c9f29668b46ce0369d5ca829"),
						internalUtils.TestHexToFelt(t, "0x287acddb27a2f9ba7f2612d72788dc96a5b30e401fc1e8072250940e024a587"),
					},
					AccountDeploymentData: []*felt.Felt{},
				},
				ExpectedResp:  AddInvokeTransactionResponse{internalUtils.TestHexToFelt(t, "0x49728601e0bb2f48ce506b0cbd9c0e2a9e50d95858aa41463f46386dca489fd")},
				ExpectedError: nil,
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		resp, err := testConfig.Provider.AddInvokeTransaction(context.Background(), &test.InvokeTx)
		if test.ExpectedError != nil {
			require.Equal(t, test.ExpectedError, err)
		} else {
			require.Equal(t, resp, test.ExpectedResp)
		}
	}
}

func TestAddDeployAccountTransaction(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv)

	testConfig := BeforeEach(t, false)

	type testSetType struct {
		DeployTx      BroadcastDeployAccountTxnV3
		ExpectedResp  AddDeployAccountTransactionResponse
		ExpectedError error
	}
	testSet := map[tests.TestEnv][]testSetType{
		tests.MockEnv: {
			{
				DeployTx: BroadcastDeployAccountTxnV3{
					Type:      TransactionType_DeployAccount,
					Version:   TransactionV3,
					ClassHash: internalUtils.TestHexToFelt(t, "0x2338634f11772ea342365abd5be9d9dc8a6f44f159ad782fdebd3db5d969738"),
					Signature: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x6d756e754793d828c6c1a89c13f7ec70dbd8837dfeea5028a673b80e0d6b4ec"),
						internalUtils.TestHexToFelt(t, "0x4daebba599f860daee8f6e100601d98873052e1c61530c630cc4375c6bd48e3"),
					},
					Nonce:         new(felt.Felt),
					NonceDataMode: DAModeL1,
					FeeMode:       DAModeL1,
					ResourceBounds: &ResourceBoundsMapping{
						L1Gas: ResourceBounds{
							MaxAmount:       "0x186a0",
							MaxPricePerUnit: "0x5af3107a4000",
						},
						L2Gas: ResourceBounds{
							MaxAmount:       "",
							MaxPricePerUnit: "",
						},
					},
					Tip:                 "",
					PayMasterData:       []*felt.Felt{},
					ContractAddressSalt: new(felt.Felt),
					ConstructorCalldata: []*felt.Felt{
						internalUtils.TestHexToFelt(t, "0x5cd65f3d7daea6c63939d659b8473ea0c5cd81576035a4d34e52fb06840196c"),
					},
				},
				ExpectedResp: AddDeployAccountTransactionResponse{
					Hash:            internalUtils.TestHexToFelt(t, "0x32b272b6d0d584305a460197aa849b5c7a9a85903b66e9d3e1afa2427ef093e"),
					ContractAddress: internalUtils.TestHexToFelt(t, "0x0"),
				},
				ExpectedError: nil,
			},
		},
	}[tests.TEST_ENV]

	for _, test := range testSet {
		resp, err := testConfig.Provider.AddDeployAccountTransaction(context.Background(), &test.DeployTx)
		if err != nil {
			require.Equal(t, err.Error(), test.ExpectedError)
		} else {
			require.Equal(t, resp.Hash.String(), test.ExpectedResp.Hash.String())
		}
	}
}
