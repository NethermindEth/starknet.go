package rpc

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
)

// BlockNumber returns the block number of the current block.
//
// Parameters:
//   - ctx: The context to use for the request
//
// Returns:
//   - uint64: The block number
//   - error: An error if any
func (provider *Provider) BlockNumber(ctx context.Context) (uint64, error) {
	var blockNumber uint64
	if err := do(ctx, provider.c, "starknet_blockNumber", &blockNumber); err != nil {
		if errors.Is(err, errNotFound) {
			return 0, ErrNoBlocks
		}

		return 0, rpcerr.UnwrapToRPCErr(err)
	}

	return blockNumber, nil
}

// BlockHashAndNumber retrieves the hash and number of the current block.
//
// Parameters:
//   - ctx: The context to use for the request.
//
// Returns:
//   - *BlockHashAndNumberOutput: The hash and number of the current block
//   - error: An error if any
func (provider *Provider) BlockHashAndNumber(
	ctx context.Context,
) (*BlockHashAndNumberOutput, error) {
	var block BlockHashAndNumberOutput
	if err := do(ctx, provider.c, "starknet_blockHashAndNumber", &block); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrNoBlocks)
	}

	return &block, nil
}

// WithBlockNumber returns a BlockID with the given block number.
//
// Parameters:
//   - n: The block number to use for the BlockID.
//
// Returns:
//   - BlockID: A BlockID struct with the specified block number
func WithBlockNumber(n uint64) BlockID {
	var blockID BlockID
	blockID.Number = &n

	return blockID
}

// WithBlockHash returns a BlockID with the given hash.
//
// Parameters:
//   - h: The hash to use for the BlockID.
//
// Returns:
//   - BlockID: A BlockID struct with the specified hash
func WithBlockHash(h *felt.Felt) BlockID {
	var blockID BlockID
	blockID.Hash = h

	return blockID
}

// WithBlockTag creates a new BlockID with the specified tag.
//
// Parameters:
//   - tag: The tag for the BlockID
//
// Returns:
//   - BlockID: A BlockID struct with the specified tag
func WithBlockTag(tag BlockTag) BlockID {
	var blockID BlockID
	blockID.Tag = tag

	return blockID
}

// BlockWithTxHashes retrieves the block with transaction hashes for the given block ID.
//
// Parameters:
//   - ctx: The context.Context object for controlling the function call
//   - blockID: The ID of the block to retrieve the transactions from
//
// Returns:
//   - *BlockWithTxHashesOutput: The retrieved block (confirmed or pre-confirmed)
//   - error: An error, if any
//
//nolint:dupl // Similar to BlockWithTxs, but it's a different method.
func (provider *Provider) BlockWithTxHashes(
	ctx context.Context,
	blockID BlockID,
) (*BlockWithTxHashesOutput, error) {
	var result BlockTxHashes
	if err := do(ctx, provider.c, "starknet_getBlockWithTxHashes", &result, blockID); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrBlockNotFound)
	}

	// if header.Hash == nil it's a pre_confirmed block
	if result.Hash == nil {
		return &BlockWithTxHashesOutput{
			PreConfirmed: &PreConfirmedBlockTxHashes{
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
			},
		}, nil
	}
	return &BlockWithTxHashesOutput{Block: &result}, nil
}


// StateUpdate is a function that performs a state update operation
// (gets the information about the result of executing the requested block).
//
// Parameters:
//   - ctx: The context.Context object for controlling the function call
//   - blockID: The ID of the block to retrieve the transactions from
//
// Returns:
//   - *StateUpdateOutput: The retrieved state update
//   - error: An error, if any
func (provider *Provider) StateUpdate(
	ctx context.Context,
	blockID BlockID,
) (*StateUpdateOutput, error) {
	var state StateUpdateOutput
	if err := do(ctx, provider.c, "starknet_getStateUpdate", &state, blockID); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrBlockNotFound)
	}

	return &state, nil
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
func (provider *Provider) BlockTransactionCount(
	ctx context.Context,
	blockID BlockID,
) (uint64, error) {
	var result uint64
	if err := do(ctx, provider.c, "starknet_getBlockTransactionCount", &result, blockID); err != nil {
		return 0, rpcerr.UnwrapToRPCErr(err, ErrBlockNotFound)
	}

	return result, nil
}

// BlockWithTxs retrieves a block with its transactions given the block id.
//
// Parameters:
//   - ctx: The context.Context object for the request
//   - blockID: The ID of the block to retrieve
//
// Returns:
//   - *BlockWithTxsOutput: The retrieved block (confirmed or pre-confirmed)
//   - error: An error, if any
//
//nolint:dupl // Similar to BlockWithTxHashes, but it's a different method.
func (provider *Provider) BlockWithTxs(ctx context.Context, blockID BlockID) (*BlockWithTxsOutput, error) {
	var result Block
	if err := do(ctx, provider.c, "starknet_getBlockWithTxs", &result, blockID); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrBlockNotFound)
	}
	// if header.Hash == nil it's a pre_confirmed block
	if result.Hash == nil {
		return &BlockWithTxsOutput{
			PreConfirmed: &PreConfirmedBlock{
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
			},
		}, nil
	}

	return &BlockWithTxsOutput{Block: &result}, nil
}

// Get block information with full transactions and receipts given the block id
func (provider *Provider) BlockWithReceipts(
	ctx context.Context,
	blockID BlockID,
) (interface{}, error) {
	var result json.RawMessage
	if err := do(ctx, provider.c, "starknet_getBlockWithReceipts", &result, blockID); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrBlockNotFound)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(result, &m); err != nil {
		return nil, rpcerr.Err(rpcerr.InternalError, StringErrData(err.Error()))
	}

	// Pre_confirmedBlockWithReceipts doesn't contain a "status" field
	if _, ok := m["status"]; ok {
		var block BlockWithReceipts
		if err := json.Unmarshal(result, &block); err != nil {
			return nil, rpcerr.Err(rpcerr.InternalError, StringErrData(err.Error()))
		}

		return &block, nil
	} else {
		var preConfirmedBlock PreConfirmedBlockWithReceipts
		if err := json.Unmarshal(result, &preConfirmedBlock); err != nil {
			return nil, rpcerr.Err(rpcerr.InternalError, StringErrData(err.Error()))
		}

		return &preConfirmedBlock, nil
	}
}
