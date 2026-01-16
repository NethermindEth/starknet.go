package methods

import (
	"context"
	"errors"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
)

// StateUpdate is a function that performs a state update operation
// (gets the information about the result of executing the requested block).
//
// Parameters:
//   - ctx: The context.Context object for controlling the function call
//   - blockID: The ID of the block to retrieve the transactions from
//
// Returns:
//   - *StateUpdateOutput: The retrieved state update
//   - error: An error, if any
func GetStateUpdate(
	ctx context.Context,
	c callers.Caller,
	blockID rpcv10.BlockID,
) (*rpcv10.StateUpdateOutput, error) {
	var state rpcv10.StateUpdateOutput
	if err := internal.Do(ctx, c, "starknet_getStateUpdate", &state, blockID); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, rpcv10.ErrBlockNotFound)
	}

	return &state, nil
}

// StorageAt retrieves the storage value of a given contract at a specific key and block ID.
//
// Parameters:
//   - ctx: The context.Context for the function
//   - contractAddress: The address of the contract
//   - key: The key for which to retrieve the storage value
//   - blockID: The ID of the block at which to retrieve the storage value
//
// Returns:
//   - string: The value of the storage
//   - error: An error if any occurred during the execution
func StorageAt(
	ctx context.Context,
	c callers.Caller,
	contractAddress *felt.Felt,
	key string,
	blockID rpcv10.BlockID,
) (string, error) {
	var value string
	hashKey := fmt.Sprintf("0x%x", internalUtils.GetSelectorFromName(key))
	if err := internal.Do(
		ctx, c, "starknet_getStorageAt", &value, contractAddress, hashKey, blockID,
	); err != nil {
		return "", rpcerr.UnwrapToRPCErr(err, rpcv10.ErrContractNotFound, rpcv10.ErrBlockNotFound)
	}

	return value, nil
}

// Get merkle paths in one of the state tries: global state, classes, individual contract.
// A single request can query for any mix of the three types of storage proofs (classes,
// contracts, and storage)
//
// Parameters:
//   - ctx: The context of the function call
//   - storageProofInput: an input containing optional and required fields for the request
//
// Returns:
//   - *StorageProofResult: The requested storage proofs. Note that if a requested leaf
//     has the default value, the path to it may end in an edge node whose path is not a
//     prefix of the requested leaf, thus effectively proving non-membership
//   - error: an error if any occurred during the execution
func StorageProof(
	ctx context.Context,
	c callers.Caller,
	storageProofInput rpcv10.StorageProofInput,
) (*rpcv10.StorageProofResult, error) {
	err := checkForPreConfirmed(storageProofInput.BlockID)
	if err != nil {
		return nil, err
	}

	var raw rpcv10.StorageProofResult
	if err := internal.DoAsObject(
		ctx, c, "starknet_getStorageProof", &raw, storageProofInput,
	); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(
			err,
			rpcv10.ErrBlockNotFound,
			rpcv10.ErrStorageProofNotSupported,
		)
	}

	return &raw, nil
}

// checkForPreConfirmed checks if the block ID has the 'pre_confirmed' tag. If it
// does, it returns an error. This is used to prevent the user from using the
// 'pre_confirmed' tag on methods that do not support it.
func checkForPreConfirmed(b rpcv10.BlockID) error {
	if b.Tag == rpcv10.BlockTagPreConfirmed {
		return errors.Join(
			rpcv10.ErrInvalidBlockID,
			errors.New("'pre_confirmed' tag is not supported on this method"),
		)
	}

	return nil
}
