package rpc

import (
	"context"
	"errors"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/stretchr/testify/require"
)

func TestDeclareTransaction(t *testing.T) {

	testConfig := beforeEach(t)

	type testSetType struct {
		DeclareTx     BroadcastDeclareTxnType
		ExpectedResp  AddDeclareTransactionResponse
		ExpectedError error
	}
	testSet := map[string][]testSetType{
		"devnet":  {},
		"mainnet": {},
		"mock": {
			{
				DeclareTx: BroadcastDeclareTxnV2{},
				ExpectedResp: AddDeclareTransactionResponse{
					TransactionHash: utils.TestHexToFelt(t, "0x41d1f5206ef58a443e7d3d1ca073171ec25fa75313394318fc83a074a6631c3")},
				ExpectedError: nil,
			},
			{
				DeclareTx: BroadcastDeclareTxnV3{},
				ExpectedResp: AddDeclareTransactionResponse{
					TransactionHash: utils.TestHexToFelt(t, "0x41d1f5206ef58a443e7d3d1ca073171ec25fa75313394318fc83a074a6631c3")},
				ExpectedError: nil,
			},
		},
		"testnet": {{
			DeclareTx: BroadcastDeclareTxnV1{},
			ExpectedResp: AddDeclareTransactionResponse{
				TransactionHash: utils.TestHexToFelt(t, "0x55b094dc5c84c2042e067824f82da90988674314d37e45cb0032aca33d6e0b9")},
			ExpectedError: errors.New("Invalid Params"),
		},
		},
	}[testEnv]

	for _, test := range testSet {
		resp, err := testConfig.provider.AddDeclareTransaction(context.Background(), test.DeclareTx)
		if err != nil {
			require.Equal(t, test.ExpectedError.Error(), err.Error())
		} else {
			require.Equal(t, (*resp.TransactionHash).String(), (*test.ExpectedResp.TransactionHash).String())
		}

	}
}

func TestAddInvokeTransaction(t *testing.T) {

	testConfig := beforeEach(t)

	type testSetType struct {
		InvokeTx      BroadcastInvokeTxnType
		ExpectedResp  AddInvokeTransactionResponse
		ExpectedError *RPCError
	}
	testSet := map[string][]testSetType{
		"devnet":  {},
		"mainnet": {},
		"mock": {
			{
				InvokeTx:     BroadcastInvokev1Txn{InvokeTxnV1{SenderAddress: new(felt.Felt).SetUint64(123)}},
				ExpectedResp: AddInvokeTransactionResponse{&felt.Zero},
				ExpectedError: &RPCError{
					Code:    ErrUnexpectedError.Code,
					Message: ErrUnexpectedError.Message,
					Data:    "Something crazy happened"},
			},
			{
				InvokeTx:      BroadcastInvokev1Txn{InvokeTxnV1{}},
				ExpectedResp:  AddInvokeTransactionResponse{utils.TestHexToFelt(t, "0xdeadbeef")},
				ExpectedError: nil,
			},
			{
				InvokeTx: BroadcastInvokev3Txn{
					InvokeTxnV3{
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
					}},
				ExpectedResp:  AddInvokeTransactionResponse{utils.TestHexToFelt(t, "0x49728601e0bb2f48ce506b0cbd9c0e2a9e50d95858aa41463f46386dca489fd")},
				ExpectedError: nil,
			},
		},
		"testnet": {},
	}[testEnv]

	for _, test := range testSet {
		resp, err := testConfig.provider.AddInvokeTransaction(context.Background(), test.InvokeTx)
		if test.ExpectedError != nil {
			require.Equal(t, test.ExpectedError, err)
		} else {
			require.Equal(t, *resp, test.ExpectedResp)
		}

	}
}

func TestAddDeployAccountTansaction(t *testing.T) {

	testConfig := beforeEach(t)

	type testSetType struct {
		DeployTx      BroadcastAddDeployTxnType
		ExpectedResp  AddDeployAccountTransactionResponse
		ExpectedError error
	}
	testSet := map[string][]testSetType{
		"devnet":  {},
		"mainnet": {},
		"mock": {
			{
				DeployTx: BroadcastDeployAccountTxn{DeployAccountTxn{}},
				ExpectedResp: AddDeployAccountTransactionResponse{
					TransactionHash: utils.TestHexToFelt(t, "0x32b272b6d0d584305a460197aa849b5c7a9a85903b66e9d3e1afa2427ef093e"),
					ContractAddress: utils.TestHexToFelt(t, "0x0"),
				},
				ExpectedError: nil,
			},
			{
				DeployTx: BroadcastDeployAccountTxnV3{
					DeployAccountTxnV3{
						Type:      TransactionType_DeployAccount,
						Version:   TransactionV3,
						ClassHash: utils.TestHexToFelt(t, "0x2338634f11772ea342365abd5be9d9dc8a6f44f159ad782fdebd3db5d969738"),
						Signature: []*felt.Felt{
							utils.TestHexToFelt(t, "0x6d756e754793d828c6c1a89c13f7ec70dbd8837dfeea5028a673b80e0d6b4ec"),
							utils.TestHexToFelt(t, "0x4daebba599f860daee8f6e100601d98873052e1c61530c630cc4375c6bd48e3"),
						},
						Nonce:         new(felt.Felt),
						NonceDataMode: DAModeL1,
						FeeMode:       DAModeL1,
						ResourceBounds: ResourceBoundsMapping{
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
							utils.TestHexToFelt(t, "0x5cd65f3d7daea6c63939d659b8473ea0c5cd81576035a4d34e52fb06840196c"),
						},
					}},
				ExpectedResp: AddDeployAccountTransactionResponse{
					TransactionHash: utils.TestHexToFelt(t, "0x32b272b6d0d584305a460197aa849b5c7a9a85903b66e9d3e1afa2427ef093e"),
					ContractAddress: utils.TestHexToFelt(t, "0x0")},
				ExpectedError: nil,
			},
		},
	}[testEnv]

	for _, test := range testSet {

		resp, err := testConfig.provider.AddDeployAccountTransaction(context.Background(), test.DeployTx)
		if err != nil {
			require.Equal(t, err.Error(), test.ExpectedError)
		} else {
			require.Equal(t, (*resp.TransactionHash).String(), (*test.ExpectedResp.TransactionHash).String())
		}

	}
}
