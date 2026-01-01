package rpc

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
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
func (provider *Provider) Nonce(
	ctx context.Context,
	blockID BlockID,
	contractAddress *felt.Felt,
) (*felt.Felt, error) {
	var nonce *felt.Felt
	if err := do(
		ctx, provider.c, "starknet_getNonce", &nonce, blockID, contractAddress,
	); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrContractNotFound, ErrBlockNotFound)
	}

	return nonce, nil
}
