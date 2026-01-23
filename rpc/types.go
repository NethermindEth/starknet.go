package rpc

import "github.com/NethermindEth/juno/core/felt"

// InvokeFunctionCall represents a function call to be invoked on a contract.
// It's a helper type used to build a rpcvX.FunctionCall.
type InvokeFunctionCall struct {
	// The address of the contract to invoke
	ContractAddress *felt.Felt
	// The name of the function to invoke
	FunctionName string
	// The parameters passed to the function
	CallData []*felt.Felt
}
