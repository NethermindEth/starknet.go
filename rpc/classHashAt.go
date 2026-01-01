package rpc

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
)

// ClassHashAt retrieves the class hash at the given block ID and contract address.
//
// Parameters:
//   - ctx: The context.Context used for the request
//   - blockID: The ID of the block
//   - contractAddress: The address of the contract
//
// Returns:
//   - *felt.Felt: The class hash
//   - error: An error if any occurred during the execution
func (provider *Provider) ClassHashAt(
	ctx context.Context,
	blockID BlockID,
	contractAddress *felt.Felt,
) (*felt.Felt, error) {
	var result *felt.Felt
	if err := do(
		ctx, provider.c, "starknet_getClassHashAt", &result, blockID, contractAddress,
	); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrContractNotFound, ErrBlockNotFound)
	}

	return result, nil
}
