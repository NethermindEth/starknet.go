//nolint:dupl // Similar to BlockWithTxHashes, but it's a different method.
package rpc

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
)

// BlockWithTxs retrieves a block with its transactions given the block id.
//
// Parameters:
//   - ctx: The context.Context object for the request
//   - blockID: The ID of the block to retrieve
//
// Returns:
//   - interface{}: The retrieved block
//   - error: An error, if any
func BlockWithTxs(ctx context.Context, c callCloser, blockID BlockID) (interface{}, error) {
	var result Block
	if err := do(ctx, c, "starknet_getBlockWithTxs", &result, blockID); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrBlockNotFound)
	}
	// if header.Hash == nil it's a pre_confirmed block
	if result.Hash == nil {
		return &PreConfirmedBlock{
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
