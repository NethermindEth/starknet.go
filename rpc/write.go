package rpc

import (
	"context"
	"errors"
	"strings"
)

type BroadcastedInvokeTransaction interface{}

// AddInvokeTransaction estimates the fee for a given StarkNet transaction.
func (provider *Provider) AddInvokeTransaction(ctx context.Context, broadcastedInvoke BroadcastedInvokeTransaction) (*AddInvokeTransactionResponse, error) {
	// TODO: EntryPointSelector now part of calldata
	// tx, ok := broadcastedInvoke.(BroadcastedInvokeV0Transaction)
	// if ok {
	// 	tx.EntryPointSelector = fmt.Sprintf("0x%x", types.GetSelectorFromName(tx.EntryPointSelector))
	// 	broadcastedInvoke = tx
	// }
	var output AddInvokeTransactionResponse
	switch invoke := broadcastedInvoke.(type) {
	case BroadcastedInvokeV1Transaction:
		if err := do(ctx, provider.c, "starknet_addInvokeTransaction", &output, invoke); err != nil {
			return nil, err
		}
		return &output, nil
	}
	return nil, errors.New("invalid invoke type")
}

// AddDeclareTransaction submits a new class declaration transaction.
func (provider *Provider) AddDeclareTransaction(ctx context.Context, declareTransaction BroadcastedDeclareTransaction) (*AddDeclareTransactionResponse, error) {
	var result AddDeclareTransactionResponse
	if err := do(ctx, provider.c, "starknet_addDeclareTransaction", &result, declareTransaction); err != nil {
		if strings.Contains(err.Error(), "Invalid contract class") {
			return nil, ErrInvalidContractClass
		}
		return nil, err
	}
	return &result, nil
}

// AddDeployTransaction allows to declare a class and instantiate the
// associated contract in one command. This function will be deprecated and
// replaced by AddDeclareTransaction to declare a class, followed by
// AddInvokeTransaction to instantiate the contract. For now, it remains the only
// way to deploy an account without being charged for it.
func (provider *Provider) AddDeployTransaction(ctx context.Context, deployTransaction BroadcastedDeployTxn) (*AddDeployTransactionResponse, error) {
	var result AddDeployTransactionResponse
	return &result, errors.New("AddDeployTransaction was removed, UDC should be used instead")
}

// AddDeployAccountTransaction manages the DEPLOY_ACCOUNT syscall
func (provider *Provider) AddDeployAccountTransaction(ctx context.Context, deployAccountTransaction BroadcastedDeployAccountTransaction) (*AddDeployTransactionResponse, error) {
	var result AddDeployTransactionResponse
	if err := do(ctx, provider.c, "starknet_addDeployAccountTransaction", &result, deployAccountTransaction); err != nil {
		if strings.Contains(err.Error(), "Class hash not found") {
			return nil, ErrClassHashNotFound
		}
		return nil, err
	}
	return &result, nil
}
