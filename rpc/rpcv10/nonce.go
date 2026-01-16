package rpcv10

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
)

// Nonce retrieves the nonce for a given block ID and contract address.
//
// Parameters:
//   - ctx: is the context.Context for the function call
//   - blockID: is the ID of the block
//   - contractAddress: is the address of the contract
//
// Returns:
//   - *felt.Felt: the contract's nonce at the requested state
//   - error: an error if any
func Nonce(
	ctx context.Context,
	c callers.Caller,
	blockID BlockID,
	contractAddress *felt.Felt,
) (*felt.Felt, error) {
	var nonce *felt.Felt
	if err := internal.Do(
		ctx, c, "starknet_getNonce", &nonce, blockID, contractAddress,
	); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrContractNotFound, ErrBlockNotFound)
	}

	return nonce, nil
}
