package rpc

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
)

// Estimates the resources required by a given sequence of transactions when applied
// on a given state. If one of the transactions reverts or fails due to any reason
// (e.g. validation failure or an internal error), a TRANSACTION_EXECUTION_ERROR is returned.
// The estimate is given in fri.
//
// Parameters:
//   - ctx: The context of the function call
//   - requests: A sequence of transactions to estimate, running each transaction on the
//     state resulting from applying all the previous ones
//   - simulationFlags: Describes what parts of the transaction should be executed
//   - blockID: The hash of the requested block, or number (height) of the requested block,
//     or a block tag, for the block referencing the state or call the transaction on.
//
// Returns:
//   - []FeeEstimation: A sequence of fee estimation where the i'th estimate corresponds
//     to the i'th transaction
//   - error: An error if any occurred during the execution
func EstimateFee(
	ctx context.Context,
	c callCloser,
	requests []BroadcastTxn,
	simulationFlags []SimulationFlag,
	blockID BlockID,
) ([]FeeEstimation, error) {
	var raw []FeeEstimation
	if err := do(
		ctx, c, "starknet_estimateFee", &raw, requests, simulationFlags, blockID,
	); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(
			err,
			ErrBlockNotFound,
			ErrContractNotFound,
			ErrTxnExec,
		)
	}

	return raw, nil
}
