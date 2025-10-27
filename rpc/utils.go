package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

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
	rawLatestBlock, err := provider.BlockWithTxs(ctx, WithBlockTag(BlockTagLatest))
	if err != nil {
		return tip, fmt.Errorf("failed to get latest block: %w", err)
	}

	latestBlock, ok := rawLatestBlock.(*Block)
	if !ok {
		return tip, fmt.Errorf("unexpected block type: %T", rawLatestBlock)
	}

	var tipStruct struct {
		Tip U64 `json:"tip"`
	}

	var tipCounter uint64
	// sum up the tips from all transactions
	for _, transaction := range latestBlock.Transactions {
		// TODO: replace this with the future V3 txn interface
		rawTxn, err := json.Marshal(transaction.Transaction)
		if err != nil {
			return tip, fmt.Errorf("failed to marshal transaction: %w", err)
		}
		err = json.Unmarshal(rawTxn, &tipStruct)
		if err != nil {
			return tip, fmt.Errorf("failed to get tip from transaction: %w", err)
		}

		uintTip, err := tipStruct.Tip.ToUint64()
		if err != nil {
			return tip, fmt.Errorf("failed to convert tip to uint64: %w", err)
		}
		tipCounter += uintTip
	}

	// take the average of the tips
	tip = U64(strconv.FormatUint(tipCounter/uint64(len(latestBlock.Transactions)), 16))

	return tip, nil
}
