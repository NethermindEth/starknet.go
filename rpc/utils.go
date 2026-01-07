package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	"github.com/Masterminds/semver/v3"
)

// IsCompatible compares the version of the Starknet JSON-RPC Specification
// implemented by the node with the version implemented by the Provider type,
// and returns whether they are the same or not.
//
// Parameters:
//   - ctx: The context for the function.
//   - provider: The provider to use.
//
// Returns:
//   - isCompatible: True if the node version is compatible with the SDK version, false otherwise.
//   - nodeVersion: The version of the Starknet JSON-RPC Specification implemented by the node.
//   - err: An error if any.
func IsCompatible(ctx context.Context, provider RPCProvider) (
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

// EstimateTip returns the estimated tip to be used in a transaction
// based on the average tip of all transactions in the latest block.
//
// Parameters:
//   - ctx: The context for the function.
//   - provider: The provider to use.
//   - multiplier: The multiplier to be used against the estimated tip
//     (E.g: 1.5 means estimated tip + 50% of it).
//     If multiplier <= 0, it'll be set to 1.0 (no multiplier, just the estimated tip).
//
// Returns:
//   - tip: The estimated tip to be used in a transaction (the average of
//     the tips of all transactions in the latest block) multiplied by the multiplier.
//   - err: An error if any.
func EstimateTip(
	ctx context.Context,
	provider RPCProvider,
	multiplier float64,
) (
	tip U64,
	err error,
) {
	rawLatestBlock, err := provider.BlockWithTxs(ctx, WithBlockTag(BlockTagLatest))
	if err != nil {
		return tip, fmt.Errorf("failed to get latest block: %w", err)
	}

	var transactions []BlockTransaction
	switch {
	case rawLatestBlock.Block != nil:
		transactions = rawLatestBlock.Block.Transactions
	case rawLatestBlock.PreConfirmed != nil:
		transactions = rawLatestBlock.PreConfirmed.Transactions
	default:
		return tip, fmt.Errorf("unexpected block output: %+v", rawLatestBlock)
	}
	var tipStruct struct {
		Tip U64 `json:"tip"`
	}

	tipCounter := new(big.Int)
	// sum up the tips from all transactions
	for _, transaction := range transactions {
		// L1Handler transactions don't have a tip
		if transaction.GetType() == TransactionTypeL1Handler {
			continue
		}

		// TODO: refactor this in the RPC package refactoring
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
		bigTip := new(big.Int).SetUint64(uintTip)
		tipCounter.Add(tipCounter, bigTip)
	}

	// No transactions in the block OR all transactions have a tip of 0
	if tipCounter.Cmp(new(big.Int)) == 0 {
		return U64("0x0"), nil
	}

	bigLength := new(big.Int).SetUint64(uint64(len(transactions)))

	averageTip := tipCounter.Div(tipCounter, bigLength).Uint64()

	if multiplier <= 0 || averageTip == 0 {
		return U64("0x" + strconv.FormatUint(averageTip, 16)), nil
	}

	multipliedAverageTip := float64(averageTip) * multiplier
	tip = U64("0x" + strconv.FormatUint(uint64(multipliedAverageTip), 16))

	return tip, nil
}
