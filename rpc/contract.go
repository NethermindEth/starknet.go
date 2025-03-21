package rpc

import (
	"context"
	"encoding/json"

	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
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

	return typecastClassOutput(rawClass)

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
	return typecastClassOutput(rawClass)
}

// typecastClassOutput typecasts the rawClass output to the appropriate ClassOutput type.
//
// Parameters:
// rawClass - A pointer to a map[string]any containing the raw class data.
// Returns:
// - ClassOutput: a ClassOutput interface
// - error: an error if any
func typecastClassOutput(rawClass map[string]any) (ClassOutput, error) {
	rawClassByte, err := json.Marshal(rawClass)
	if err != nil {
		return nil, Err(InternalError, StringErrData(err.Error()))
	}

	// if contract_class_version exists, then it's a ContractClass type
	if _, exists := (rawClass)["contract_class_version"]; exists {
		var contractClass contracts.ContractClass
		err = json.Unmarshal(rawClassByte, &contractClass)
		if err != nil {
			return nil, Err(InternalError, StringErrData(err.Error()))
		}
		return &contractClass, nil
	}
	var depContractClass contracts.DeprecatedContractClass
	err = json.Unmarshal(rawClassByte, &depContractClass)
	if err != nil {
		return nil, Err(InternalError, StringErrData(err.Error()))
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
	hashKey := fmt.Sprintf("0x%x", internalUtils.GetSelectorFromName(key))
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

// Estimates the resources required by a given sequence of transactions when applied on a given state.
// If one of the transactions reverts or fails due to any reason (e.g. validation failure or an internal error),
// a TRANSACTION_EXECUTION_ERROR is returned. For v0-2 transactions the estimate is given in wei, and for v3 transactions it is given in fri.
//
// Parameters:
// - ctx: The context of the function call
// - requests: A sequence of transactions to estimate, running each transaction on the state resulting from applying all the previous ones
// - simulationFlags: Describes what parts of the transaction should be executed
// - blockID: The hash of the requested block, or number (height) of the requested block, or a block tag, for the block referencing the state or call the transaction on
// Returns:
// - []FeeEstimation: A sequence of fee estimation where the i'th estimate corresponds to the i'th transaction
// - error: An error if any occurred during the execution
func (provider *Provider) EstimateFee(ctx context.Context, requests []BroadcastTxn, simulationFlags []SimulationFlag, blockID BlockID) ([]FeeEstimation, error) {
	var raw []FeeEstimation
	if err := do(ctx, provider.c, "starknet_estimateFee", &raw, requests, simulationFlags, blockID); err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrTxnExec, ErrBlockNotFound)
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
// - *FeeEstimation: the fee estimated for the message
// - error: an error if any occurred during the execution
func (provider *Provider) EstimateMessageFee(ctx context.Context, msg MsgFromL1, blockID BlockID) (*FeeEstimation, error) {
	var raw FeeEstimation
	if err := do(ctx, provider.c, "starknet_estimateMessageFee", &raw, msg, blockID); err != nil {

		return nil, tryUnwrapToRPCErr(err, ErrContractError, ErrBlockNotFound)
	}
	return &raw, nil
}

// Get merkle paths in one of the state tries: global state, classes, individual contract.
// A single request can query for any mix of the three types of storage proofs (classes, contracts, and storage)
//
// Parameters:
// - ctx: The context of the function call
// - storageProofInput: an input containing optional and required fields for the request
// Returns:
// - *StorageProofResult: The requested storage proofs. Note that if a requested leaf has the default value,
// the path to it may end in an edge node whose path is not a prefix of the requested leaf, thus effectively proving non-membership
// - error: an error if any occurred during the execution
func (provider *Provider) GetStorageProof(ctx context.Context, storageProofInput StorageProofInput) (*StorageProofResult, error) {
	var raw StorageProofResult
	if err := do(ctx, provider.c, "starknet_getStorageProof", &raw, storageProofInput); err != nil {

		return nil, tryUnwrapToRPCErr(err, ErrBlockNotFound, ErrStorageProofNotSupported)
	}
	return &raw, nil
}
