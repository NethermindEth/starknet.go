package rpc

import (
	"context"

	"github.com/dontpanicdao/caigo/rpc/types"
)

// BlockNumber gets the most recent accepted block number.
func (sc *Provider) BlockNumber(ctx context.Context) (uint64, error) {
	var blockNumber uint64
	if err := sc.c.CallContext(ctx, &blockNumber, "starknet_blockNumber"); err != nil {
		return 0, err
	}
	return blockNumber, nil
}

// BlockHashAndNumber gets block information given the block number or its hash.
func (sc *Provider) BlockHashAndNumber(ctx context.Context) (*types.BlockHashAndNumberOutput, error) {
	var block types.BlockHashAndNumberOutput
	if err := do(ctx, sc.c, "starknet_blockHashAndNumber", &block); err != nil {
		return nil, err
	}
	return &block, nil
}

func WithBlockNumber(n uint64) types.BlockID {
	return types.BlockID{
		Number: &n,
	}
}

func WithBlockHash(h types.Hash) types.BlockID {
	return types.BlockID{
		Hash: &h,
	}
}

func WithBlockTag(tag string) types.BlockID {
	return types.BlockID{
		Tag: tag,
	}
}

// BlockWithTxHashes gets block information given the block id.
func (sc *Provider) BlockWithTxHashes(ctx context.Context, blockID types.BlockID) (types.Block, error) {
	var result types.Block
	if err := do(ctx, sc.c, "starknet_getBlockWithTxHashes", &result, blockID); err != nil {
		return types.Block{}, err
	}
	return result, nil
}

// BlockTransactionCount gets the number of transactions in a block
func (sc *Provider) BlockTransactionCount(ctx context.Context, blockID types.BlockID) (uint64, error) {
	var result uint64
	if err := do(ctx, sc.c, "starknet_getBlockTransactionCount", &result, blockID); err != nil {
		return 0, err
	}
	return result, nil
}

// BlockWithTxs get block information with full transactions given the block id.
func (sc *Provider) BlockWithTxs(ctx context.Context, blockID types.BlockID) (interface{}, error) {
	var result types.Block
	if err := do(ctx, sc.c, "starknet_getBlockWithTxs", &result, blockID); err != nil {
		return nil, err
	}
	return &result, nil
}
