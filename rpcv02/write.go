package rpcv02

import (
	"context"
	"errors"
	"fmt"

	"github.com/dontpanicdao/caigo/types"
)

type BroadcastedInvokeTransaction interface{}

// AddInvokeTransaction estimates the fee for a given StarkNet transaction.
func (provider *Provider) AddInvokeTransaction(ctx context.Context, broadcastedInvoke BroadcastedInvokeTransaction) (*types.AddInvokeTransactionOutput, error) {
	tx, ok := broadcastedInvoke.(BroadcastedInvokeV0Transaction)
	if ok {
		tx.EntryPointSelector = fmt.Sprintf("0x0%s", types.GetSelectorFromName(tx.EntryPointSelector).Text(16))
		broadcastedInvoke = tx
	}
	var output types.AddInvokeTransactionOutput
	switch invoke := broadcastedInvoke.(type) {
	case BroadcastedInvokeV0Transaction:
		if err := do(ctx, provider.c, "starknet_addInvokeTransaction", &output, invoke); err != nil {
			return nil, err
		}
		return &output, nil
	case BroadcastedInvokeV1Transaction:
		if err := do(ctx, provider.c, "starknet_addInvokeTransaction", &output, invoke); err != nil {
			return nil, err
		}
		return &output, nil
	}
	return nil, errors.New("invalid invoke type")
}

// AddDeclareTransaction submits a new class declaration transaction.
func (provider *Provider) AddDeclareTransaction(ctx context.Context, declareTransaction BroadcastedDeclareTransaction) (*AddDeclareTransactionOutput, error) {
	var result AddDeclareTransactionOutput
	if err := do(ctx, provider.c, "starknet_addDeclareTransaction", &result, declareTransaction); err != nil {
		// TODO: check Pathfinder/Devnet error and return
		// INVALID_CONTRACT_CLASS
		return nil, err
	}
	return &result, nil
}

// AddDeployTransaction allows to declare a class and instantiate the
// associated contract in one command. This function will be deprecated and
// replaced by AddDeclareTransaction to declare a class, followed by
// AddInvokeTransaction to instantiate the contract. For now, it remains the only
// way to deploy an account without being charged for it.
func (provider *Provider) AddDeployTransaction(ctx context.Context, deployTransaction BroadcastedDeployTransaction) (*AddDeployTransactionOutput, error) {
	var result AddDeployTransactionOutput
	if err := do(ctx, provider.c, "starknet_addDeployTransaction", &result, deployTransaction); err != nil {
		// TODO: check Pathfinder/Devnet error and return
		// INVALID_CONTRACT_CLASS
		return nil, err
	}
	return &result, nil
}

// AddDeployAccountTransaction manages the DEPLOY_ACCOUNT syscall
func (provider *Provider) AddDeployAccountTransaction(ctx context.Context, deployAccountTransaction BroadcastedDeployAccountTransaction) (*AddDeployTransactionOutput, error) {
	var result AddDeployTransactionOutput
	if err := do(ctx, provider.c, "starknet_addDeployAccountTransaction", &result, deployAccountTransaction); err != nil {
		// TODO: check Pathfinder/Devnet error and return
		// CLASS_HASH_NOT_FOUND
		return nil, err
	}
	return &result, nil
}
