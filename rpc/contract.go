package rpc

import (
	"context"
	"fmt"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/rpc/types"
)

// Class gets the contract class definition associated with the given hash.
func (sc *Provider) Class(ctx context.Context, classHash string) (*types.ContractClass, error) {
	var rawClass types.ContractClass
	if err := do(ctx, sc.c, "starknet_getClass", &rawClass, classHash); err != nil {
		return nil, err
	}
	return &rawClass, nil
}

// ClassAt get the contract class definition at the given address.
func (sc *Provider) ClassAt(ctx context.Context, blockID types.BlockID, contractAddress types.Hash) (*types.ContractClass, error) {
	var rawClass types.ContractClass
	if err := do(ctx, sc.c, "starknet_getClassAt", &rawClass, blockID, contractAddress); err != nil {
		return nil, err
	}
	return &rawClass, nil
}

// ClassHashAt gets the contract class hash for the contract deployed at the given address.
func (sc *Provider) ClassHashAt(ctx context.Context, blockID types.BlockID, contractAddress types.Hash) (*string, error) {
	var result string
	if err := do(ctx, sc.c, "starknet_getClassHashAt", &result, blockID, contractAddress); err != nil {
		return nil, err
	}
	return &result, nil
}

// StorageAt gets the value of the storage at the given address and key.
func (sc *Provider) StorageAt(ctx context.Context, contractAddress types.Hash, key string, blockID types.BlockID) (string, error) {
	var value string
	hashKey := fmt.Sprintf("0x%s", caigo.GetSelectorFromName(key).Text(16))
	if err := do(ctx, sc.c, "starknet_getStorageAt", &value, contractAddress, hashKey, blockID); err != nil {
		return "", err
	}
	return value, nil
}

// Nonce returns the Nonce of a contract
func (sc *Provider) Nonce(ctx context.Context, contractAddress types.Hash) (*string, error) {
	nonce := ""
	if err := do(ctx, sc.c, "starknet_getNonce", &nonce, contractAddress); err != nil {
		return nil, err
	}
	return &nonce, nil
}

// EstimateFee estimates the fee for a given StarkNet transaction.
func (sc *Provider) EstimateFee(ctx context.Context, request types.Call, blockID types.BlockID) (*types.FeeEstimate, error) {
	if request.EntryPointSelector != nil {
		entrypointSelector := fmt.Sprintf("0x%s", caigo.GetSelectorFromName(*request.EntryPointSelector).Text(16))
		request.EntryPointSelector = &entrypointSelector
	}
	var raw types.FeeEstimate
	if err := do(ctx, sc.c, "starknet_estimateFee", &raw, request, blockID); err != nil {
		return nil, err
	}
	return &raw, nil
}
