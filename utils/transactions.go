package utils

import (
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/hash"
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
			NonceDataMode:         rpc.DAModeL1,
			FeeMode:               rpc.DAModeL1,
		},
	}

	return invokeTxn
}

// BuildDeclareTxn creates a new declare transaction (v3) for the StarkNet network.
// A declare transaction is used to declare a new contract class on the network.
//
// Parameters:
//   - senderAddress: The address of the account sending the transaction
//   - casmClass: The casm class of the contract to be declared
//   - contractClass: The contract class to be declared
//   - nonce: The account's nonce
//   - resourceBounds: Resource bounds for the transaction execution
//
// Returns:
//   - rpc.BroadcastDeclareTxnV3: A broadcast declare transaction with default values
//     for signature, tip, paymaster data, etc. Need to be signed before being sent.
func BuildDeclareTxn(
	senderAddress *felt.Felt,
	casmClass contracts.CasmClass,
	contractClass *rpc.ContractClass,
	nonce *felt.Felt,
	resourceBounds rpc.ResourceBoundsMapping,
) (rpc.BroadcastDeclareTxnV3, error) {
	compiledClassHash, err := hash.CompiledClassHash(casmClass)
	if err != nil {
		return rpc.BroadcastDeclareTxnV3{}, err
	}

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
		NonceDataMode:         rpc.DAModeL1,
		FeeMode:               rpc.DAModeL1,
	}

	return declareTxn, nil
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
			NonceDataMode:       rpc.DAModeL1,
			FeeMode:             rpc.DAModeL1,
		},
	}

	return deployAccountTxn
}

// InvokeFuncCallsToFunctionCalls converts a slice of InvokeFunctionCall to a slice of FunctionCall.
//
// Parameters:
//   - invokeFuncCalls: The invoke function calls to convert
//
// Returns:
//   - []*rpc.FunctionCall: A new function calls
func InvokeFuncCallsToFunctionCalls(invokeFuncCalls []rpc.InvokeFunctionCall) []rpc.FunctionCall {
	functionCalls := make([]rpc.FunctionCall, len(invokeFuncCalls))
	for i, call := range invokeFuncCalls {
		functionCalls[i] = rpc.FunctionCall{
			ContractAddress:    call.ContractAddress,
			EntryPointSelector: GetSelectorFromNameFelt(call.FunctionName),
			Calldata:           call.CallData,
		}
	}

	return functionCalls
}

// FeeEstToResBoundsMap converts a FeeEstimation to ResourceBoundsMapping with applied multipliers.
// Parameters:
//   - feeEstimation: The fee estimation to convert
//   - multiplier: Multiplier for max amount and max price per unit. Recommended to be 1.5, but at least 1
//
// Returns:
//   - rpc.ResourceBoundsMapping: Resource bounds with applied multipliers
func FeeEstToResBoundsMap(
	feeEstimation rpc.FeeEstimation,
	multiplier float64,
) rpc.ResourceBoundsMapping {

	// Create L1 resources bounds
	l1Gas := toResourceBounds(feeEstimation.L1GasPrice.Uint64(), feeEstimation.L1GasConsumed.Uint64(), multiplier)
	l1DataGas := toResourceBounds(feeEstimation.L1DataGasPrice.Uint64(), feeEstimation.L1DataGasConsumed.Uint64(), multiplier)
	// Create L2 resource bounds
	l2Gas := toResourceBounds(feeEstimation.L2GasPrice.Uint64(), feeEstimation.L2GasConsumed.Uint64(), multiplier)

	return rpc.ResourceBoundsMapping{
		L1Gas:     l1Gas,
		L1DataGas: l1DataGas,
		L2Gas:     l2Gas,
	}
}

// toResourceBounds converts a gas price and gas consumed to a ResourceBounds with applied multiplier.
//
// Parameters:
//   - gasPrice: The gas price
//   - gasConsumed: The gas consumed
//   - multiplier: Multiplier for max amount and max price per unit
//
// Returns:
//   - rpc.ResourceBounds: Resource bounds with applied multiplier
func toResourceBounds(
	gasPrice uint64,
	gasConsumed uint64,
	multiplier float64,
) rpc.ResourceBounds {
	maxAmount := float64(gasConsumed) * multiplier
	maxPricePerUnit := float64(gasPrice) * multiplier

	return rpc.ResourceBounds{
		MaxAmount:       rpc.U64(fmt.Sprintf("0x%x", uint64(maxAmount))),
		MaxPricePerUnit: rpc.U128(fmt.Sprintf("0x%x", uint64(maxPricePerUnit))),
	}
}
