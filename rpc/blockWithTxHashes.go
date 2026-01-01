//nolint:dupl // Similar to BlockWithTxs, but it's a different method.
package rpc

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
)

// BlockWithTxHashes retrieves the block with transaction hashes for the given block ID.
//
// Parameters:
//   - ctx: The context.Context object for controlling the function call
//   - blockID: The ID of the block to retrieve the transactions from
//
// Returns:
//   - interface{}: The retrieved block
//   - error: An error, if any
func (provider *Provider) BlockWithTxHashes(
	ctx context.Context,
	blockID BlockID,
) (interface{}, error) {
	var result BlockTxHashes
	if err := do(ctx, provider.c, "starknet_getBlockWithTxHashes", &result, blockID); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrBlockNotFound)
	}

	// if header.Hash == nil it's a pre_confirmed block
	if result.Hash == nil {
		return &PreConfirmedBlockTxHashes{
			PreConfirmedBlockHeader{
				Number:           result.Number,
				Timestamp:        result.Timestamp,
				SequencerAddress: result.SequencerAddress,
				L1GasPrice:       result.L1GasPrice,
				L2GasPrice:       result.L2GasPrice,
				StarknetVersion:  result.StarknetVersion,
				L1DataGasPrice:   result.L1DataGasPrice,
				L1DAMode:         result.L1DAMode,
			},
			result.Transactions,
		}, nil
	}

	return &result, nil
}
