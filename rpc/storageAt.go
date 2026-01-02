package rpc

import (
	"context"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
)

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
	c callCloser,
	contractAddress *felt.Felt,
	key string,
	blockID BlockID,
) (string, error) {
	var value string
	hashKey := fmt.Sprintf("0x%x", internalUtils.GetSelectorFromName(key))
	if err := do(
		ctx, c, "starknet_getStorageAt", &value, contractAddress, hashKey, blockID,
	); err != nil {
		return "", rpcerr.UnwrapToRPCErr(err, ErrContractNotFound, ErrBlockNotFound)
	}

	return value, nil
}
