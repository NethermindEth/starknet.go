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

// TODO: remove this function before merge in case of no use
// ToResourceBounds converts a FeeEstimation to ResourceBoundsMapping with applied multipliers.
//
// Calculates max amount as max_amount = overall_fee / gas_price (unless gas_price is 0,
// then max_amount is 0). Calculates max price per unit as max_price_per_unit = gas_price.
//
// Then multiplies max_amount by amountMultiplier and max_price_per_unit by unitPriceMultiplier.
//
// Parameters:
//   - feeEstimation: The fee estimation to convert
//   - amountMultiplier: Multiplier for max amount, defaults to 1.5
//   - unitPriceMultiplier: Multiplier for max price per unit, defaults to 1.5
//
// Returns:
//   - rpc.ResourceBoundsMapping: Resource bounds with applied multipliers
// func ToResourceBounds(
// 	feeEstimation *rpc.FeeEstimation,
// 	amountMultiplier float64,
// 	unitPriceMultiplier float64,
// ) (rpc.ResourceBoundsMapping, error) {
// 	if amountMultiplier <= 0 || unitPriceMultiplier <= 0 {
// 		return rpc.ResourceBoundsMapping{}, fmt.Errorf("values of 'amountMultiplier' and 'unitPriceMultiplier' must be greater than 0")
// 	}

// 	// Convert felt.Felt values to big.Int for calculations
// 	overallFee := feeEstimation.OverallFee.Uint64()
// 	l1GasPrice := feeEstimation.L1GasPrice.Uint64()
// 	// l1DataGasPrice := feeEstimation.L1DataGasPrice.Uint64()
// 	// l2GasPrice := feeEstimation.L2GasPrice.Uint64()
// 	// l2DataGasPrice := feeEstimation.L2DataGasPrice.Uint64()

// 	// Calculate max amount
// 	var maxAmount float64
// 	if l1GasPrice != 0 {
// 		maxAmount = (float64(overallFee) / float64(l1GasPrice)) * amountMultiplier
// 	} else {
// 		maxAmount = 0
// 	}

// 	// Apply unit price multiplier to gas price
// 	maxPricePerUnit := float64(l1GasPrice) * unitPriceMultiplier

// 	// Convert big.Int values to U64 and U128 strings
// 	maxAmountHex := fmt.Sprintf("0x%x", maxAmount)
// 	maxPricePerUnitHex := fmt.Sprintf("0x%x", maxPricePerUnit)

// 	// Create L1 resource bounds
// 	l1ResourceBounds := rpc.ResourceBounds{
// 		MaxAmount:       rpc.U64(maxAmountHex),
// 		MaxPricePerUnit: rpc.U128(maxPricePerUnitHex),
// 	}

// 	// Create empty L2 resource bounds
// 	l2ResourceBounds := rpc.ResourceBounds{
// 		MaxAmount:       "0x0",
// 		MaxPricePerUnit: "0x0",
// 	}

// 	// Create resource bounds mapping
// 	resourceBounds := rpc.ResourceBoundsMapping{
// 		L1Gas: l1ResourceBounds,
// 		L1DataGas: rpc.ResourceBounds{
// 			MaxAmount:       "0x0",
// 			MaxPricePerUnit: "0x0",
// 		},
// 		L2Gas: l2ResourceBounds,
// 	}

// 	return resourceBounds, nil
// }

// func calcResourceBounds(
// 	overallFee uint64,
// 	gasPrice uint64,
// 	gasConsumed uint64,
// 	amountMultiplier float64,
// 	unitPriceMultiplier float64,
// ) rpc.ResourceBounds {

// 	var maxAmount float64
// 	if gasPrice != 0 {
// 		maxAmount = (float64(overallFee) / float64(gasPrice)) * amountMultiplier
// 	}

// 	maxPricePerUnit := float64(gasPrice) * unitPriceMultiplier

// 	return rpc.ResourceBounds{
// 		MaxAmount:       rpc.U64(fmt.Sprintf("0x%x", maxAmount)),
// 		MaxPricePerUnit: rpc.U128(fmt.Sprintf("0x%x", maxPricePerUnit)),
// 	}
// }

// func ToResourceBoundsC(
// 	feeEstimation *rpc.FeeEstimation,
// 	amountMultiplier float64,
// 	unitPriceMultiplier float64,
// ) (rpc.ResourceBoundsMapping, error) {
// 	if amountMultiplier <= 0 || unitPriceMultiplier <= 0 {
// 		return rpc.ResourceBoundsMapping{}, fmt.Errorf("values of 'amountMultiplier' and 'unitPriceMultiplier' must be greater than 0")
// 	}

// 	// Convert felt.Felt values to big.Int for calculations
// 	overallFee := feeEstimation.OverallFee.BigInt(nil)
// 	gasPrice := feeEstimation.L1GasPrice.BigInt(nil)

// 	// Calculate max amount
// 	var maxAmount *big.Int
// 	if gasPrice.Cmp(big.NewInt(0)) != 0 {
// 		maxAmount = new(big.Int).Div(overallFee, gasPrice)

// 		// Apply amount multiplier
// 		amountMultiplierBig := new(big.Float).SetFloat64(amountMultiplier)
// 		maxAmountFloat := new(big.Float).SetInt(maxAmount)
// 		maxAmountFloat.Mul(maxAmountFloat, amountMultiplierBig)

// 		maxAmount = new(big.Int)
// 		maxAmountFloat.Int(maxAmount)
// 	} else {
// 		maxAmount = big.NewInt(0)
// 	}

// 	// Apply unit price multiplier to gas price
// 	unitPriceMultiplierBig := new(big.Float).SetFloat64(unitPriceMultiplier)
// 	maxPricePerUnitFloat := new(big.Float).SetInt(gasPrice)
// 	maxPricePerUnitFloat.Mul(maxPricePerUnitFloat, unitPriceMultiplierBig)

// 	maxPricePerUnit := new(big.Int)
// 	maxPricePerUnitFloat.Int(maxPricePerUnit)

// 	// Convert big.Int values to U64 and U128 strings
// 	maxAmountHex := fmt.Sprintf("0x%x", maxAmount)
// 	maxPricePerUnitHex := fmt.Sprintf("0x%x", maxPricePerUnit)

// 	// Create L1 resource bounds
// 	l1ResourceBounds := rpc.ResourceBounds{
// 		MaxAmount:       rpc.U64(maxAmountHex),
// 		MaxPricePerUnit: rpc.U128(maxPricePerUnitHex),
// 	}

// 	// Create empty L2 resource bounds
// 	l2ResourceBounds := rpc.ResourceBounds{
// 		MaxAmount:       "0x0",
// 		MaxPricePerUnit: "0x0",
// 	}

// 	// Create resource bounds mapping
// 	resourceBounds := rpc.ResourceBoundsMapping{
// 		L1Gas: l1ResourceBounds,
// 		L1DataGas: rpc.ResourceBounds{
// 			MaxAmount:       "0x0",
// 			MaxPricePerUnit: "0x0",
// 		},
// 		L2Gas: l2ResourceBounds,
// 	}

// 	return resourceBounds, nil
// }
