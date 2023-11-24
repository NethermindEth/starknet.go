package rpc

import (
	"context"
	"encoding/json"
	
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
)

// Class retrieves the class information from the Provider with the given hash.
//
// Parameters:
// - ctx: The context.Context object
// - blockID: The BlockID object
// - classHash: The *felt.Felt object
// Returns:
// - ClassOutput: The output of the class.
// - error: An error if any occurred during the execution.
func (provider *Provider) Class(ctx context.Context, blockID BlockID, classHash *felt.Felt) (ClassOutput, error) {
	var rawClass map[string]any
	if err := do(ctx, provider.c, "starknet_getClass", &rawClass, blockID, classHash); err != nil {
		
		return nil, tryUnwrapToRPCErr(err, ErrClassHashNotFound, ErrBlockNotFound)
	}

	return typecastClassOutput(&rawClass)

}

// ClassAt returns the class at the specified blockID and contractAddress.
//
// Parameters:
// - ctx: The context.Context object for the function
// - blockID: The BlockID of the class
// - contractAddress: The address of the contract
// Returns:
// - ClassOutput: The output of the class
// - error: An error if any occurred during the execution
func (provider *Provider) ClassAt(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (ClassOutput, error) {
	var rawClass map[string]any
	if err := do(ctx, provider.c, "starknet_getClassAt", &rawClass, blockID, contractAddress); err != nil {
		
		return nil, tryUnwrapToRPCErr(err, ErrContractNotFound, ErrBlockNotFound)
	}
	return typecastClassOutput(&rawClass)
}

// typecastClassOutput typecasts the rawClass output to the appropriate ClassOutput type.
//
// Parameters:
// rawClass - A pointer to a map[string]any containing the raw class data.
// Returns:
// - ClassOutput: a ClassOutput interface
// - error: an error if any
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

// ClassHashAt retrieves the class hash at the given block ID and contract address.
//
// Parameters:
// - ctx: The context.Context used for the request
// - blockID: The ID of the block
// - contractAddress: The address of the contract
// Returns:
// - *felt.Felt: The class hash
// - error: An error if any occurred during the execution
func (provider *Provider) ClassHashAt(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (*felt.Felt, error) {
	var result *felt.Felt
	if err := do(ctx, provider.c, "starknet_getClassHashAt", &result, blockID, contractAddress); err != nil {
		
		return nil, tryUnwrapToRPCErr(err, ErrContractNotFound, ErrBlockNotFound)
	}
	return result, nil
}

// StorageAt retrieves the storage value of a given contract at a specific key and block ID.
//
// Parameters:
// - ctx: The context.Context for the function
// - contractAddress: The address of the contract
// - key: The key for which to retrieve the storage value
// - blockID: The ID of the block at which to retrieve the storage value
// Returns:
// - string: The value of the storage
// - error: An error if any occurred during the execution
func (provider *Provider) StorageAt(ctx context.Context, contractAddress *felt.Felt, key string, blockID BlockID) (string, error) {
	var value string
	hashKey := fmt.Sprintf("0x%x", utils.GetSelectorFromName(key))
	if err := do(ctx, provider.c, "starknet_getStorageAt", &value, contractAddress, hashKey, blockID); err != nil {
		
		return "", tryUnwrapToRPCErr(err, ErrContractNotFound, ErrBlockNotFound)
	}
	return value, nil
}

// Nonce retrieves the nonce for a given block ID and contract address.
//
// Parameters:
// - ctx: is the context.Context for the function call
// - blockID: is the ID of the block
// - contractAddress: is the address of the contract
// Returns:
// - *felt.Felt: the contract's nonce at the requested state
// - error: an error if any
func (provider *Provider) Nonce(ctx context.Context, blockID BlockID, contractAddress *felt.Felt) (*felt.Felt, error) {
	var nonce *felt.Felt
	if err := do(ctx, provider.c, "starknet_getNonce", &nonce, blockID, contractAddress); err != nil {
		
		return nil, tryUnwrapToRPCErr(err, ErrContractNotFound, ErrBlockNotFound)
	}
	return nonce, nil
}

// EstimateFee estimates the fee for executing a set of requests on the StarkNet blockchain.
// The function returns a slice of FeeEstimate objects and an error. The FeeEstimate objects
// represent the estimated fees for executing the requests. The error is nil if the fee estimation
// is successful. If there is an error, it can be one of the following:
//   - ErrContractNotFound: when the contract is not found.
//   - ErrContractError: when there is an error with the contract.
//   - ErrBlockNotFound: when the block is not found.
//   - any other error that occurred during the fee estimation.
//
// Parameters:
// - ctx: the context.Context object for cancellation and timeouts
// - requests: a slice of BroadcastTxn objects representing the requests to be executed
// - blockID: the ID of the block on which the requests should be executed
// Returns:
// - []FeeEstimate: a slice of FeeEstimate objects representing the estimated fees for executing the requests
// - error: an error if any occurred during the execution
func (provider *Provider) EstimateFee(ctx context.Context, requests []BroadcastTxn, blockID BlockID) ([]FeeEstimate, error) {
	var raw []FeeEstimate
	if err := do(ctx, provider.c, "starknet_estimateFee", &raw, requests, blockID); err != nil {

		return nil, tryUnwrapToRPCErr(err, ErrContractNotFound, ErrBlockNotFound)

	}
	return raw, nil
}

// EstimateMessageFee estimates the L2 fee of a message sent on L1 (Provider struct).
//
// Parameters:
// - ctx: The context of the function call
// - msg: The message to estimate the fee for
// - blockID: The ID of the block to estimate the fee in
// Returns:
// - *FeeEstimate: the fee estimated for the message
// - error: an error if any occurred during the execution
func (provider *Provider) EstimateMessageFee(ctx context.Context, msg MsgFromL1, blockID BlockID) (*FeeEstimate, error) {
	var raw FeeEstimate
	if err := do(ctx, provider.c, "starknet_estimateMessageFee", &raw, msg, blockID); err != nil {
		

		return nil, tryUnwrapToRPCErr(err, ErrContractNotFound, ErrBlockNotFound)

	}
	return &raw, nil
}
