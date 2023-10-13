package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
)

// Class gets the contract class definition associated with the given hash.
func (provider *Provider) Class(ctx context.Context, blockID BlockID, classHash *felt.Felt) (ClassOutput, error) {
	var rawClass map[string]any
	if err := do(ctx, provider.c, "starknet_getClass", &rawClass, blockID, classHash); err != nil {
		
		return nil, tryUnwrapToRPCErr(err, ErrClassHashNotFound, ErrBlockNotFound)
	}

	return typecastClassOutput(&rawClass)

}

// ClassAt get the contract class definition at the given address.
func (provider *Provider) ClassAt(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (ClassOutput, error) {
	var rawClass map[string]any
	if err := do(ctx, provider.c, "starknet_getClassAt", &rawClass, blockID, contractAddress); err != nil {
		
		return nil, tryUnwrapToRPCErr(err, ErrContractNotFound, ErrBlockNotFound)
	}
	return typecastClassOutput(&rawClass)
}

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

// ClassHashAt gets the contract class hash for the contract deployed at the given address.
func (provider *Provider) ClassHashAt(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (*felt.Felt, error) {
	var result *felt.Felt
	if err := do(ctx, provider.c, "starknet_getClassHashAt", &result, blockID, contractAddress); err != nil {
		
		return nil, tryUnwrapToRPCErr(err, ErrContractNotFound, ErrBlockNotFound)
	}
	return result, nil
}

// StorageAt gets the value of the storage at the given address and key.
func (provider *Provider) StorageAt(ctx context.Context, contractAddress *felt.Felt, key string, blockID BlockID) (string, error) {
	var value string
	hashKey := fmt.Sprintf("0x%x", utils.GetSelectorFromName(key))
	if err := do(ctx, provider.c, "starknet_getStorageAt", &value, contractAddress, hashKey, blockID); err != nil {
		
		return "", tryUnwrapToRPCErr(err, ErrContractNotFound, ErrBlockNotFound)
	}
	return value, nil
}

// Nonce returns the Nonce of a contract
func (provider *Provider) Nonce(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (*string, error) {
	nonce := ""
	if err := do(ctx, provider.c, "starknet_getNonce", &nonce, blockID, contractAddress); err != nil {
		
		return nil, tryUnwrapToRPCErr(err, ErrContractNotFound, ErrBlockNotFound)
	}
	return &nonce, nil
}

// EstimateFee estimates the fee for a given Starknet transaction.
func (provider *Provider) EstimateFee(ctx context.Context, requests []EstimateFeeInput, blockID BlockID) ([]FeeEstimate, error) {
	var raw []FeeEstimate
	if err := do(ctx, provider.c, "starknet_estimateFee", &raw, requests, blockID); err != nil {
		
		return nil, tryUnwrapToRPCErr(err, ErrContractNotFound,ErrContractError, ErrBlockNotFound)
	}
	return raw, nil
}

// EstimateMessageFee estimates the L2 fee of a message sent on L1
func (provider *Provider) EstimateMessageFee(ctx context.Context, msg MsgFromL1, blockID BlockID) (*FeeEstimate, error) {
	var raw FeeEstimate
	if err := do(ctx, provider.c, "starknet_estimateMessageFee", &raw, msg, blockID); err != nil {
		
		return nil, tryUnwrapToRPCErr(err, ErrContractNotFound,ErrContractError, ErrBlockNotFound)
	}
	return &raw, nil
}
