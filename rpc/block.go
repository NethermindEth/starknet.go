package rpc

import (
	"context"
	"errors"

	"github.com/NethermindEth/juno/core/felt"
)

// BlockNumber returns the most recent accepted block number of the blockchain.
//
// It takes a context.Context parameter and returns a uint64 and an error.
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

// BlockHashAndNumber retrieves the block hash information and number from the provider (node).
//
// ctx - The context to use for the request.
// Returns a pointer to BlockHashAndNumberOutput and an error.
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

// WithBlockNumber sets the block number of the BlockID.
//
// n: The block number to set.
// Returns a BlockID with the specified block number set.
func WithBlockNumber(n uint64) BlockID {
	return BlockID{
		Number: &n,
	}
}

// WithBlockHash generates a BlockID with the given Felt hash.
//
// Parameters:
// - h: A pointer to a Felt object representing the hash.
//
// Return type: BlockID.
func WithBlockHash(h *felt.Felt) BlockID {
	return BlockID{
		Hash: h,
	}
}

// WithBlockTag returns a new BlockID with the specified tag.
//
// It takes a string parameter `tag` and returns a BlockID.
func WithBlockTag(tag string) BlockID {
	return BlockID{
		Tag: tag,
	}
}

// BlockWithTxHashes returns the block with transaction hashes for the given block ID.
//
// It takes a context.Context and a BlockID as parameters.
// It returns an interface{} and an error.
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

// StateUpdate is a function that gets the information about the result of executing the requested block.
//
// It takes a context.Context object and a BlockID object as parameters and returns a
// pointer to a StateUpdateOutput object and an error object.
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

// BlockTransactionCount returns the number of transactions in a block.
//
// It takes a context and a block ID as parameters.
// It returns a uint64 representing the number of transactions in the block, and an error if any.
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

// BlockWithTxs retrieves a block information with transactions given the block id.
//
// It takes a context.Context and a BlockID as input parameters.
// The function returns an interface{} and an error.
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
