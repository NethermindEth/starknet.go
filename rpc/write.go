package rpc

import (
	"context"
	"fmt"

	"github.com/dontpanicdao/caigo/rpc/types"
	ctypes "github.com/dontpanicdao/caigo/types"
)

// AddInvokeTransaction estimates the fee for a given StarkNet transaction.
func (provider *Provider) AddInvokeTransaction(ctx context.Context, call ctypes.FunctionCall, signature []string, maxFee string, version string) (*ctypes.AddInvokeTransactionOutput, error) {
	call.EntryPointSelector = fmt.Sprintf("0x%s", ctypes.GetSelectorFromName(call.EntryPointSelector).Text(16))
	var output ctypes.AddInvokeTransactionOutput
	if err := do(ctx, provider.c, "starknet_addInvokeTransaction", &output, call, signature, maxFee, version); err != nil {
		return nil, err
	}
	return &output, nil
}

// AddDeclareTransaction submits a new class declaration transaction.
func (provider *Provider) AddDeclareTransaction(ctx context.Context, contractClass ctypes.ContractClass, version string) (*types.AddDeclareTransactionOutput, error) {
	var result types.AddDeclareTransactionOutput
	if err := do(ctx, provider.c, "starknet_addDeclareTransaction", &result, contractClass, version); err != nil {
		return nil, err
	}
	return &result, nil
}

// AddDeployTransaction allows to declare a class and instantiate the
// associated contract in one command. This function will be deprecated and
// replaced by AddDeclareTransaction to declare a class, followed by
// AddInvokeTransaction to instantiate the contract. For now, it remains the only
// way to deploy an account without being charged for it.
func (provider *Provider) AddDeployTransaction(ctx context.Context, salt string, inputs []string, contractClass ctypes.ContractClass) (*types.AddDeployTransactionOutput, error) {
	var result types.AddDeployTransactionOutput
	if err := do(ctx, provider.c, "starknet_addDeployTransaction", &result, salt, inputs, contractClass); err != nil {
		return nil, err
	}
	return &result, nil
}
