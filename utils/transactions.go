package utils

import (
	"github.com/NethermindEth/juno/core/felt"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/rpc/types"
)

// @changed
// InvokeFuncCallsToFunctionCalls converts a slice of [types.InvokeFunctionCall] to a
// slice of [types.FunctionCall].
//
// Parameters:
//   - invokeFuncCalls: A slice of invoke function calls to convert
//
// Returns:
//   - []types.FunctionCall: a slice of function calls
func InvokeFuncCallsToFunctionCalls(
	invokeFuncCalls []types.InvokeFunctionCall,
) []types.FunctionCall {
	functionCalls := make([]types.FunctionCall, len(invokeFuncCalls))

	for i, call := range invokeFuncCalls {
		functionCalls[i] = types.FunctionCall{
			ContractAddress:    call.ContractAddress,
			EntryPointSelector: GetSelectorFromNameFelt(call.FunctionName),
			Calldata:           call.CallData,
		}
	}

	return functionCalls
}

// FeeLimits is a struct with custom limits for the fee values, used
// as a parameter for the `CustomFeeEstToResBoundsMap` function.
type FeeLimits struct {
	// Custom max value for L1 gas price
	L1GasPriceLimit types.U128
	// Custom max value for L1 gas amount
	L1GasAmountLimit types.U64

	// Custom max value for L2 gas price
	L2GasPriceLimit types.U128
	// Custom max value for L2 gas amount
	L2GasAmountLimit types.U64

	// Custom max value for L1 data gas price
	L1DataGasPriceLimit types.U128
	// Custom max value for L1 data gas amount
	L1DataGasAmountLimit types.U64
}

// FillHexWithZeroes normalises a hex string to have a '0x' prefix and pads it with leading zeros
// to a total length of 66 characters (including the '0x' prefix).
func FillHexWithZeroes(hex string) string {
	return internalUtils.FillHexWithZeroes(hex)
}

// WeiToETH converts a Wei amount to ETH
// Returns the ETH value as a float64
func WeiToETH(wei *felt.Felt) float64 {
	return internalUtils.WeiToETH(wei)
}

// ETHToWei converts an ETH amount to Wei
// Returns the Wei value as a *felt.Felt
func ETHToWei(eth float64) *felt.Felt {
	return internalUtils.ETHToWei(eth)
}

// FRIToSTRK converts a FRI amount to STRK
// Returns the STRK value as a float64
func FRIToSTRK(fri *felt.Felt) float64 {
	return internalUtils.WeiToETH(fri)
}

// STRKToFRI converts a STRK amount to FRI
// Returns the FRI value as a *felt.Felt
func STRKToFRI(strk float64) *felt.Felt {
	return internalUtils.ETHToWei(strk)
}
