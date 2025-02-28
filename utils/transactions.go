package utils

import (
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
)

// BuildInvokeTxn creates a new invoke transaction (v3) for the StarkNet network.
//
// Parameters:
//   - senderAddress: The address of the account sending the transaction
//   - nonce: The account's nonce
//   - calldata: The data expected by the account's `execute` function (in most usecases,
//     this includes the called contract address and a function selector)
//   - resourceBounds: Resource bounds for the transaction execution
//
// Returns:
//   - rpc.BroadcastInvokev3Txn: A broadcast invoke transaction with default values
//     for signature, tip, paymaster data, etc. Need to be signed before being sent.
func BuildInvokeTxn(
	senderAddress *felt.Felt,
	nonce *felt.Felt,
	calldata []*felt.Felt,
	resourceBounds rpc.ResourceBoundsMapping,
) rpc.BroadcastInvokev3Txn {
	invokeTxn := rpc.BroadcastInvokev3Txn{
		InvokeTxnV3: rpc.InvokeTxnV3{
			Type:                  rpc.TransactionType_Invoke,
			SenderAddress:         senderAddress,
			Calldata:              calldata,
			Version:               rpc.TransactionV3,
			Signature:             []*felt.Felt{},
			Nonce:                 nonce,
			ResourceBounds:        resourceBounds,
			Tip:                   "0x0",
			PayMasterData:         []*felt.Felt{},
			AccountDeploymentData: []*felt.Felt{},
			NonceDataMode:         "0x0",
			FeeMode:               "0x0",
		},
	}

	return invokeTxn
}

// BuildDeclareTxn creates a new declare transaction (v3) for the StarkNet network.
// A declare transaction is used to declare a new contract class on the network.
//
// Parameters:
//   - senderAddress: The address of the account sending the transaction
//   - compiledClassHash: The hash of the casm contract class
//   - nonce: The account's nonce
//   - contractClass: The contract class to be declared
//   - resourceBounds: Resource bounds for the transaction execution
//
// Returns:
//   - rpc.BroadcastDeclareTxnV3: A broadcast declare transaction with default values
//     for signature, tip, paymaster data, etc. Need to be signed before being sent.
func BuildDeclareTxn(
	senderAddress *felt.Felt,
	compiledClassHash *felt.Felt,
	nonce *felt.Felt,
	contractClass *rpc.ContractClass,
	resourceBounds rpc.ResourceBoundsMapping,
) rpc.BroadcastDeclareTxnV3 {

	declareTxn := rpc.BroadcastDeclareTxnV3{
		Type:                  rpc.TransactionType_Declare,
		SenderAddress:         senderAddress,
		CompiledClassHash:     compiledClassHash,
		Version:               rpc.TransactionV3,
		Signature:             []*felt.Felt{},
		Nonce:                 nonce,
		ContractClass:         contractClass,
		ResourceBounds:        resourceBounds,
		Tip:                   "0x0",
		PayMasterData:         []*felt.Felt{},
		AccountDeploymentData: []*felt.Felt{},
		NonceDataMode:         "0x0",
		FeeMode:               "0x0",
	}

	return declareTxn
}

// BuildDeployAccountTxn creates a new deploy account transaction (v3) for the StarkNet network.
// A deploy account transaction is used to deploy a new account contract on the network.
//
// Parameters:
//   - nonce: The account's nonce
//   - contractAddressSalt: A value used to randomize the deployed contract address
//   - constructorCalldata: The parameters for the constructor function
//   - classHash: The hash of the contract class to deploy
//   - resourceBounds: Resource bounds for the transaction execution
//
// Returns:
//   - rpc.BroadcastDeployAccountTxnV3: A broadcast deploy account transaction with default values
//     for signature, tip, paymaster data, etc. Need to be signed before being sent.
func BuildDeployAccountTxn(
	nonce *felt.Felt,
	contractAddressSalt *felt.Felt,
	constructorCalldata []*felt.Felt,
	classHash *felt.Felt,
	resourceBounds rpc.ResourceBoundsMapping,
) rpc.BroadcastDeployAccountTxnV3 {
	deployAccountTxn := rpc.BroadcastDeployAccountTxnV3{
		DeployAccountTxnV3: rpc.DeployAccountTxnV3{
			Type:                rpc.TransactionType_DeployAccount,
			Version:             rpc.TransactionV3,
			Signature:           []*felt.Felt{},
			Nonce:               nonce,
			ContractAddressSalt: contractAddressSalt,
			ConstructorCalldata: constructorCalldata,
			ClassHash:           classHash,
			ResourceBounds:      resourceBounds,
			Tip:                 "0x0",
			PayMasterData:       []*felt.Felt{},
			NonceDataMode:       "0x0",
			FeeMode:             "0x0",
		},
	}

	return deployAccountTxn
}
