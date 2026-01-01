package rpc

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
)

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
	c callCloser,
	storageProofInput StorageProofInput,
) (*StorageProofResult, error) {
	err := checkForPreConfirmed(storageProofInput.BlockID)
	if err != nil {
		return nil, err
	}

	var raw StorageProofResult
	if err := doAsObject(
		ctx, c, "starknet_getStorageProof", &raw, storageProofInput,
	); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrBlockNotFound, ErrStorageProofNotSupported)
	}

	return &raw, nil
}
