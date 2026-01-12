//nolint:dupl // Similar to BlockWithTxHashes, but it's a different method.
package methods

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/NethermindEth/starknet.go/rpc/types"
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
func BlockWithTxs(
	ctx context.Context,
	c callers.Caller,
	blockID types.BlockID,
) (interface{}, error) {
	var result rpcv10.Block
	if err := internal.Do(ctx, c, "starknet_getBlockWithTxs", &result, blockID); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, rpcv10.ErrBlockNotFound)
	}
	// if header.Hash == nil it's a pre_confirmed block
	if result.Hash == nil {
		return &rpcv10.PreConfirmedBlock{
			PreConfirmedBlockHeader: rpcv10.PreConfirmedBlockHeader{
				Number:           result.Number,
				Timestamp:        result.Timestamp,
				SequencerAddress: result.SequencerAddress,
				L1GasPrice:       result.L1GasPrice,
				L2GasPrice:       result.L2GasPrice,
				StarknetVersion:  result.StarknetVersion,
				L1DataGasPrice:   result.L1DataGasPrice,
				L1DAMode:         result.L1DAMode,
			},
			Transactions: result.Transactions,
		}, nil
	}

	return &result, nil
}
