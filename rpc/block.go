package rpc

import (
	"context"
	"errors"

	"github.com/NethermindEth/juno/core/felt"
)

// BlockNumber returns the block number of the current block.
//
// It takes a context as input and returns a uint64 and an error.
func (provider *Provider) BlockNumber(ctx context.Context) (uint64, error) {
	var blockNumber uint64
	if err := provider.c.CallContext(ctx, &blockNumber, "starknet_blockNumber"); err != nil {
		if errors.Is(err, errNotFound) {
			return 0, ErrNoBlocks
		}
		return 0, err
	}
	return blockNumber, nil
}

// BlockHashAndNumber retrieves the hash and number of the current block.
//
// ctx: The context to use for the request.
// Returns a pointer to a BlockHashAndNumberOutput struct and an error.
func (provider *Provider) BlockHashAndNumber(ctx context.Context) (*BlockHashAndNumberOutput, error) {
	var block BlockHashAndNumberOutput
	if err := do(ctx, provider.c, "starknet_blockHashAndNumber", &block); err != nil {
		if errors.Is(err, errNotFound) {
			return nil, ErrNoBlocks
		}
		return nil, err
	}
	return &block, nil
}

// WithBlockNumber returns a BlockID with the given block number.
//
// Parameters:
//   - n: The block number.
//
// Returns:
//   - BlockID: A BlockID struct with the specified block number.
func WithBlockNumber(n uint64) BlockID {
	return BlockID{
		Number: &n,
	}
}

// WithBlockHash returns a BlockID with the given hash.
//
// It takes a pointer to a felt.Felt object as its parameter, and returns a BlockID.
func WithBlockHash(h *felt.Felt) BlockID {
	return BlockID{
		Hash: h,
	}
}

// WithBlockTag creates a new BlockID with the specified tag.
//
// tag: The tag for the BlockID.
// returns: A new BlockID with the specified tag.
func WithBlockTag(tag string) BlockID {
	return BlockID{
		Tag: tag,
	}
}

// BlockWithTxHashes retrieves the block with transaction hashes for the given block ID.
//
// ctx: The context.Context object for controlling the function call.
// blockID: The ID of the block to retrieve.
// Returns: The retrieved block with transaction hashes or an error if it fails.
func (provider *Provider) BlockWithTxHashes(ctx context.Context, blockID BlockID) (interface{}, error) {
	var result BlockTxHashes
	if err := do(ctx, provider.c, "starknet_getBlockWithTxHashes", &result, blockID); err != nil {
		if errors.Is(err, errNotFound) {
			return nil, ErrBlockNotFound
		}
		return nil, err
	}

	// if header.Hash == nil it's a pending block
	if result.BlockHeader.BlockHash == nil {
		return &PendingBlockTxHashes{
			ParentHash:       result.ParentHash,
			Timestamp:        result.Timestamp,
			SequencerAddress: result.SequencerAddress,
			Transactions:     result.Transactions,
		}, nil
	}

	return &result, nil
}

// StateUpdate is a function that performs a state update operation
// (gets the information about the result of executing the requested block).
//
// It takes a context and a blockID as parameters.
// It returns a pointer to a StateUpdateOutput and an error.
func (provider *Provider) StateUpdate(ctx context.Context, blockID BlockID) (*StateUpdateOutput, error) {
	var state StateUpdateOutput
	if err := do(ctx, provider.c, "starknet_getStateUpdate", &state, blockID); err != nil {
		if errors.Is(err, errNotFound) {
			return nil, ErrBlockNotFound
		}
		return nil, err
	}
	return &state, nil
}

// BlockTransactionCount returns the number of transactions in a specific block.
//
// It takes a context.Context object to handle cancellation signals and timeouts,
// and a BlockID object representing the ID of the block.
// The function returns a uint64 value representing the number of transactions in the block,
// and an error in case of any failures.
func (provider *Provider) BlockTransactionCount(ctx context.Context, blockID BlockID) (uint64, error) {
	var result uint64
	if err := do(ctx, provider.c, "starknet_getBlockTransactionCount", &result, blockID); err != nil {
		if errors.Is(err, errNotFound) {
			return 0, ErrBlockNotFound
		}
		return 0, err
	}
	return result, nil
}

// BlockWithTxs retrieves a block with its transactions given the block id..
//
// It takes a context.Context and a BlockID as parameters.
// It returns an interface{} and an error.
func (provider *Provider) BlockWithTxs(ctx context.Context, blockID BlockID) (interface{}, error) {
	var result Block
	if err := do(ctx, provider.c, "starknet_getBlockWithTxs", &result, blockID); err != nil {
		if errors.Is(err, errNotFound) {
			return nil, ErrBlockNotFound
		}
		return nil, err
	}
	// if header.Hash == nil it's a pending block
	if result.BlockHeader.BlockHash == nil {
		return &PendingBlock{
			ParentHash:       result.ParentHash,
			Timestamp:        result.Timestamp,
			SequencerAddress: result.SequencerAddress,
			Transactions:     result.Transactions,
		}, nil
	}
	return &result, nil
}
