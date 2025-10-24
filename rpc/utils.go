package rpc

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
)

// IsCompatible compares the version of the Starknet JSON-RPC Specification
// implemented by the node with the version implemented by the Provider type,
// and returns whether they are the same or not.
//
// Parameters:
//   - ctx: The context for the function.
//
// Returns:
//   - isCompatible: True if the node version is compatible with the SDK version, false otherwise.
//   - nodeVersion: The version of the Starknet JSON-RPC Specification implemented by the node.
//   - err: An error if any.
func (provider *Provider) IsCompatible(ctx context.Context) (
	isCompatible bool,
	nodeVersion string,
	err error,
) {
	rawNodeVersion, err := provider.SpecVersion(ctx)
	if err != nil {
		return false, "", fmt.Errorf("failed to get the node's RPC spec version: %w", err)
	}

	nodeVersionParsed, err := semver.NewVersion(rawNodeVersion)
	if err != nil {
		return false, "", fmt.Errorf("failed to parse node version: %w", err)
	}

	return rpcVersion.Compare(nodeVersionParsed) == 0, rawNodeVersion, nil
}

func (provider *Provider) EstimateTip(ctx context.Context) (
	tip U64,
	err error,
) {
	// rawLatestBlock, err := provider.BlockWithTxs(ctx, WithBlockTag(BlockTagLatest))
	// if err != nil {
	// 	return tip, fmt.Errorf("failed to get latest block: %w", err)
	// }

	// latestBlock, ok := rawLatestBlock.(*Block)
	// if !ok {
	// 	return tip, fmt.Errorf("unexpected block type: %T", rawLatestBlock)
	// }

	// var tempTip uint64
	// for _, transaction := range latestBlock.Transactions {
	// 	tempTip += transaction.Transaction
	// }

	// tip = rpc.U64(tempTip)

	return tip, nil
}
