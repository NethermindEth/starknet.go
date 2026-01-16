package methods

import (
	"context"
	"encoding/json"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
)

// BlockHashAndNumber retrieves the hash and number of the current block.
//
// Parameters:
//   - ctx: The context to use for the request.
//
// Returns:
//   - *BlockHashAndNumberOutput: The hash and number of the current block
//   - error: An error if any
func BlockHashAndNumber(
	ctx context.Context,
	c callers.Caller,
) (*rpcv10.BlockHashAndNumberOutput, error) {
	var block rpcv10.BlockHashAndNumberOutput
	if err := internal.Do(ctx, c, "starknet_blockHashAndNumber", &block); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, rpcv10.ErrNoBlocks)
	}

	return &block, nil
}

// BlockNumber returns the block number of the current block.
//
// Parameters:
//   - ctx: The context to use for the request
//
// Returns:
//   - uint64: The block number
//   - error: An error if any
func BlockNumber(ctx context.Context, c callers.Caller) (uint64, error) {
	var blockNumber uint64
	if err := internal.Do(ctx, c, "starknet_blockNumber", &blockNumber); err != nil {
		return 0, rpcerr.UnwrapToRPCErr(err, rpcv10.ErrNoBlocks)
	}

	return blockNumber, nil
}

// BlockTransactionCount returns the number of transactions in a specific block.
//
// Parameters:
//   - ctx: The context.Context object to handle cancellation signals and timeouts
//   - blockID: The ID of the block to retrieve the number of transactions from
//
// Returns:
//   - uint64: The number of transactions in the block
//   - error: An error, if any
func BlockTransactionCount(
	ctx context.Context,
	c callers.Caller,
	blockID rpcv10.BlockID,
) (uint64, error) {
	var result uint64
	if err := internal.Do(ctx, c, "starknet_getBlockTransactionCount", &result, blockID); err != nil {
		return 0, rpcerr.UnwrapToRPCErr(err, rpcv10.ErrBlockNotFound)
	}

	return result, nil
}

// Get block information with full transactions and receipts given the block id
func GetBlockWithReceipts(
	ctx context.Context,
	c callers.Caller,
	blockID rpcv10.BlockID,
) (interface{}, error) {
	var result json.RawMessage
	if err := internal.Do(ctx, c, "starknet_getBlockWithReceipts", &result, blockID); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, rpcv10.ErrBlockNotFound)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(result, &m); err != nil {
		return nil, rpcerr.Err(rpcerr.InternalError, rpcv10.StringErrData(err.Error()))
	}

	// Pre_confirmedBlockWithReceipts doesn't contain a "status" field
	if _, ok := m["status"]; ok {
		var block rpcv10.BlockWithReceipts
		if err := json.Unmarshal(result, &block); err != nil {
			return nil, rpcerr.Err(rpcerr.InternalError, rpcv10.StringErrData(err.Error()))
		}

		return &block, nil
	} else {
		var preConfirmedBlock rpcv10.PreConfirmedBlockWithReceipts
		if err := json.Unmarshal(result, &preConfirmedBlock); err != nil {
			return nil, rpcerr.Err(rpcerr.InternalError, rpcv10.StringErrData(err.Error()))
		}

		return &preConfirmedBlock, nil
	}
}

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
	c callers.Caller,
	blockID rpcv10.BlockID,
) (interface{}, error) {
	var result rpcv10.BlockTxHashes
	if err := internal.Do(ctx, c, "starknet_getBlockWithTxHashes", &result, blockID); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, rpcv10.ErrBlockNotFound)
	}

	// if header.Hash == nil it's a pre_confirmed block
	if result.Hash == nil {
		return &rpcv10.PreConfirmedBlockTxHashes{
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
	blockID rpcv10.BlockID,
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
