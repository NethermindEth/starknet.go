package rpc

import (
	"context"
	"strings"
)

// AddInvokeTransaction estimates the fee for a given StarkNet transaction.
func (provider *Provider) AddInvokeTransaction(ctx context.Context, invokeTxn AddInvokeTxnInput) (*AddInvokeTransactionResponse, error) {
	var output AddInvokeTransactionResponse
	if err := do(ctx, provider.c, "starknet_addInvokeTransaction", &output, invokeTxn); err != nil {
		return nil, err
	}
	return &output, nil
}

// AddDeclareTransaction submits a new class declaration transaction.
func (provider *Provider) AddDeclareTransaction(ctx context.Context, declareTransaction AddDeclareTxnInput) (*AddDeclareTransactionResponse, error) {
	var result AddDeclareTransactionResponse
	if err := do(ctx, provider.c, "starknet_addDeclareTransaction", &result, declareTransaction); err != nil {
		if strings.Contains(err.Error(), "Invalid contract class") {
			return nil, ErrInvalidContractClass
		}
		return nil, err
	}
	return &result, nil
}

// AddDeployAccountTransaction manages the DEPLOY_ACCOUNT syscall
func (provider *Provider) AddDeployAccountTransaction(ctx context.Context, deployAccountTransaction AddDeployAccountTxnInput) (*AddDeployTransactionResponse, error) {
	var result AddDeployTransactionResponse
	if err := do(ctx, provider.c, "starknet_addDeployAccountTransaction", &result, deployAccountTransaction); err != nil {
		if strings.Contains(err.Error(), "Class hash not found") {
			return nil, ErrClassHashNotFound
		}
		return nil, err
	}
	return &result, nil
}
