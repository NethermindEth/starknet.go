package rpc

import (
	"context"
	"errors"

	"github.com/NethermindEth/juno/core/felt"
)

// BlockNumber gets the most recent accepted block number.
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

// BlockHashAndNumber gets block information given the block number or its hash.
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

func WithBlockNumber(n uint64) BlockID {
	return BlockID{
		Number: &n,
	}
}

func WithBlockHash(h *felt.Felt) BlockID {
	return BlockID{
		Hash: h,
	}
}

func WithBlockTag(tag string) BlockID {
	return BlockID{
		Tag: tag,
	}
}

// BlockWithTxHashes gets block information given the block id.
func (provider *Provider) BlockWithTxHashes(ctx context.Context, blockID BlockID) (interface{}, error) {
	var result BlockTxHashes
	if err := do(ctx, provider.c, "starknet_getBlockWithTxHashes", &result, blockID); err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrBlockNotFound)
	}

	// if header.Hash == nil it's a pending block
	if result.BlockHeader.BlockHash == nil {
		return &PendingBlockTxHashes{
			PendingBlockHeader{
				ParentHash:       result.ParentHash,
				Timestamp:        result.Timestamp,
				SequencerAddress: result.SequencerAddress},
			result.Transactions,
		}, nil
	}

	return &result, nil
}

// StateUpdate gets the information about the result of executing the requested block.
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

// BlockTransactionCount gets the number of transactions in a block
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

// BlockWithTxs get block information with full transactions given the block id.
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
			PendingBlockHeader{
				ParentHash:       result.ParentHash,
				Timestamp:        result.Timestamp,
				SequencerAddress: result.SequencerAddress},
			result.Transactions,
		}, nil
	}
	return &result, nil
}
