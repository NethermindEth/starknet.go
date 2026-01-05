package methods

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
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
func ClassHashAt(
	ctx context.Context,
	c rpc.Caller,
	blockID rpcv10.BlockID,
	contractAddress *felt.Felt,
) (*felt.Felt, error) {
	var result *felt.Felt
	if err := internal.Do(
		ctx, c, "starknet_getClassHashAt", &result, blockID, contractAddress,
	); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrContractNotFound, ErrBlockNotFound)
	}

	return result, nil
}
