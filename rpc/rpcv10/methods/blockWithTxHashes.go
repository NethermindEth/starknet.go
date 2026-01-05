//nolint:dupl // Similar to BlockWithTxs, but it's a different method.
package methods

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
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
func BlockWithTxHashes(
	ctx context.Context,
	c callCloser,
	blockID rpcv10.BlockID,
) (interface{}, error) {
	var result rpcv10.BlockTxHashes
	if err := do(ctx, c, "starknet_getBlockWithTxHashes", &result, blockID); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrBlockNotFound)
	}

	// if header.Hash == nil it's a pre_confirmed block
	if result.Hash == nil {
		return &rpcv10.PreConfirmedBlockTxHashes{
			rpcv10.PreConfirmedBlockHeader{
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
