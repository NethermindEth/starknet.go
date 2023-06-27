package rpcv02

import (
	"context"
	"fmt"

	"github.com/NethermindEth/caigo/types"
	"github.com/NethermindEth/juno/core/felt"
)

// Class gets the contract class definition associated with the given hash.
func (provider *Provider) Class(ctx context.Context, blockID BlockID, classHash string) (*types.ContractClass, error) {
	var rawClass types.ContractClass
	if err := do(ctx, provider.c, "starknet_getClass", &rawClass, blockID, classHash); err != nil {
		// TODO: bind pathfinder/devnet errors with the correct errors;
		// it should return CLASS_HASH_NOT_FOUND and BLOCK_NOT_FOUND
		return nil, err
	}
	return &rawClass, nil
}

// ClassAt get the contract class definition at the given address.
func (provider *Provider) ClassAt(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (*types.ContractClass, error) {
	var rawClass types.ContractClass
	if err := do(ctx, provider.c, "starknet_getClassAt", &rawClass, blockID, contractAddress); err != nil {
		// TODO: bind pathfinder/devnet errors with the correct errors;
		// it should return CONTRACT_NOT_FOUND and BLOCK_NOT_FOUND
		return nil, err
	}
	return &rawClass, nil
}

// ClassHashAt gets the contract class hash for the contract deployed at the given address.
func (provider *Provider) ClassHashAt(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (*string, error) {
	var result string
	if err := do(ctx, provider.c, "starknet_getClassHashAt", &result, blockID, contractAddress); err != nil {
		// TODO: bind pathfinder/devnet errors with the correct errors;
		// it should return CONTRACT_NOT_FOUND and BLOCK_NOT_FOUND
		return nil, err
	}
	return &result, nil
}

// StorageAt gets the value of the storage at the given address and key.
func (provider *Provider) StorageAt(ctx context.Context, contractAddress *felt.Felt, key string, blockID BlockID) (string, error) {
	var value string
	hashKey := fmt.Sprintf("0x%x", types.GetSelectorFromName(key))
	if err := do(ctx, provider.c, "starknet_getStorageAt", &value, contractAddress, hashKey, blockID); err != nil {
		// TODO: bind pathfinder/devnet errors with the correct errors;
		// it should return CONTRACT_NOT_FOUND and BLOCK_NOT_FOUND
		return "", err
	}
	return value, nil
}

// Nonce returns the Nonce of a contract
func (provider *Provider) Nonce(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (*string, error) {
	nonce := ""
	if err := do(ctx, provider.c, "starknet_getNonce", &nonce, blockID, contractAddress); err != nil {
		// TODO: bind pathfinder/devnet errors with the correct errors;
		// it should return CONTRACT_NOT_FOUND and BLOCK_NOT_FOUND
		return nil, err
	}
	return &nonce, nil
}

// EstimateFee estimates the fee for a given StarkNet transaction.
func (provider *Provider) EstimateFee(ctx context.Context, requests []BroadcastedTransaction, blockID BlockID) ([]types.FeeEstimate, error) {
	// tx, ok := request.(*BroadcastedInvokeV1Transaction)
	// TODO:
	// NOTE: EntryPointSelector is now just part of Calldata
	// if ok {
	// 	tx.EntryPointSelector = fmt.Sprintf("0x%x", types.GetSelectorFromName(tx.EntryPointSelector))
	// 	request = tx
	// }
	var raw []types.FeeEstimate
	if err := do(ctx, provider.c, "starknet_estimateFee", &raw, requests, blockID); err != nil {
		// TODO: Bind Pathfinder/Devnet errors to
		// CONTRACT_NOT_FOUND, INVALID_MESSAGE_SELECTOR, INVALID_CALL_DATA, CONTRACT_ERROR and BLOCK_NOT_FOUND
		return nil, err
	}
	return raw, nil
}
