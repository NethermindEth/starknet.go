package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/types"
)

// Class returns the ClassOutput contract class definition for a given block ID and class hash.
//
// It takes the following parameters:
// - ctx: the context.Context object for cancellation and deadline propagation.
// - blockID: the ID of the block to retrieve the class from.
// - classHash: the hash of the class to retrieve.
//
// It returns a ClassOutput object and an error.
func (provider *Provider) Class(ctx context.Context, blockID BlockID, classHash *felt.Felt) (ClassOutput, error) {
	var rawClass map[string]any
	if err := do(ctx, provider.c, "starknet_getClass", &rawClass, blockID, classHash); err != nil {
		switch {
		case errors.Is(err, ErrClassHashNotFound):
			return nil, ErrClassHashNotFound
		case errors.Is(err, ErrBlockNotFound):
			return nil, ErrBlockNotFound
		}
		return nil, err
	}

	return typecastClassOutput(&rawClass)

}

// ClassAt retrieves the contract class at the specified block ID and contract address.
//
// ctx: The context.Context object for cancellation and timeouts.
// blockID: The ID of the block to retrieve the class from.
// contractAddress: The address of the contract.
//
// ClassOutput: The class output.
// error: An error if any occurred.
func (provider *Provider) ClassAt(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (ClassOutput, error) {
	var rawClass map[string]any
	if err := do(ctx, provider.c, "starknet_getClassAt", &rawClass, blockID, contractAddress); err != nil {
		switch {
		case errors.Is(err, ErrContractNotFound):
			return nil, ErrContractNotFound
		case errors.Is(err, ErrBlockNotFound):
			return nil, ErrBlockNotFound
		}
		return nil, err
	}
	return typecastClassOutput(&rawClass)
}

// typecastClassOutput typecasts the rawClass output into the appropriate ClassOutput type.
//
// rawClass: A pointer to a map[string]any containing the raw class data.
// Returns the typecasted ClassOutput and any error encountered during the typecasting process.
func typecastClassOutput(rawClass *map[string]any) (ClassOutput, error) {
	rawClassByte, err := json.Marshal(rawClass)
	if err != nil {
		return nil, err
	}

	// if contract_class_version exists, then it's a ContractClass type
	if _, exists := (*rawClass)["contract_class_version"]; exists {
		var contractClass ContractClass
		err = json.Unmarshal(rawClassByte, &contractClass)
		if err != nil {
			return nil, err
		}
		return &contractClass, nil
	}
	var depContractClass DeprecatedContractClass
	err = json.Unmarshal(rawClassByte, &depContractClass)
	if err != nil {
		return nil, err
	}
	return &depContractClass, nil
}

// ClassHashAt returns the contract class hash at the given block ID and contract address.
//
// It takes the following parameters:
// - ctx: the context.Context object for cancellation and timeouts.
// - blockID: the ID of the block.
// - contractAddress: the address of the contract.
//
// It returns a *felt.Felt object and an error.
func (provider *Provider) ClassHashAt(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (*felt.Felt, error) {
	var result *felt.Felt
	if err := do(ctx, provider.c, "starknet_getClassHashAt", &result, blockID, contractAddress); err != nil {
		switch {
		case errors.Is(err, ErrContractNotFound):
			return nil, ErrContractNotFound
		case errors.Is(err, ErrBlockNotFound):
			return nil, ErrBlockNotFound
		}
		return nil, err
	}
	return result, nil
}

// StorageAt gets the value of the storage at the given address and key.
func (provider *Provider) StorageAt(ctx context.Context, contractAddress *felt.Felt, key string, blockID BlockID) (string, error) {
	var value string
	hashKey := fmt.Sprintf("0x%x", types.GetSelectorFromName(key))
	if err := do(ctx, provider.c, "starknet_getStorageAt", &value, contractAddress, hashKey, blockID); err != nil {
		switch {
		case errors.Is(err, ErrContractNotFound):
			return "", ErrContractNotFound
		case errors.Is(err, ErrBlockNotFound):
			return "", ErrBlockNotFound
		}
		return "", err
	}
	return value, nil
}

// Nonce returns the Nonce of a contract
func (provider *Provider) Nonce(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (*string, error) {
	nonce := ""
	if err := do(ctx, provider.c, "starknet_getNonce", &nonce, blockID, contractAddress); err != nil {
		switch {
		case errors.Is(err, ErrContractNotFound):
			return nil, ErrContractNotFound
		case errors.Is(err, ErrBlockNotFound):
			return nil, ErrBlockNotFound
		}
		return nil, err
	}
	return &nonce, nil
}

// EstimateFee estimates the fee for a given Starknet transaction.
func (provider *Provider) EstimateFee(ctx context.Context, requests []EstimateFeeInput, blockID BlockID) ([]FeeEstimate, error) {
	var raw []FeeEstimate
	if err := do(ctx, provider.c, "starknet_estimateFee", &raw, requests, blockID); err != nil {
		switch {
		case errors.Is(err, ErrContractNotFound):
			return nil, ErrContractNotFound
		case errors.Is(err, ErrContractError):
			return nil, ErrContractError
		case errors.Is(err, ErrBlockNotFound):
			return nil, ErrBlockNotFound
		}
		return nil, err
	}
	return raw, nil
}

// EstimateMessageFee estimates the L2 fee of a message sent on L1
func (provider *Provider) EstimateMessageFee(ctx context.Context, msg MsgFromL1, blockID BlockID) (*FeeEstimate, error) {
	var raw FeeEstimate
	if err := do(ctx, provider.c, "starknet_estimateMessageFee", &raw, msg, blockID); err != nil {
		switch {
		case errors.Is(err, ErrContractNotFound):
			return nil, ErrContractNotFound
		case errors.Is(err, ErrContractError):
			return nil, ErrContractError
		case errors.Is(err, ErrBlockNotFound):
			return nil, ErrBlockNotFound
		}
		return nil, err
	}
	return &raw, nil
}
