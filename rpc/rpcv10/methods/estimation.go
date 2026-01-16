package methods

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
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
	c callers.Caller,
	requests []rpcv10.BroadcastTxn,
	simulationFlags []rpcv10.SimulationFlag,
	blockID rpcv10.BlockID,
) ([]rpcv10.FeeEstimation, error) {
	var raw []rpcv10.FeeEstimation
	if err := internal.Do(
		ctx, c, "starknet_estimateFee", &raw, requests, simulationFlags, blockID,
	); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(
			err,
			rpcv10.ErrBlockNotFound,
			rpcv10.ErrContractNotFound,
			rpcv10.ErrTxnExec,
		)
	}

	return raw, nil
}

// EstimateMessageFee estimates the L2 fee of a message sent on L1 (Provider struct).
//
// Parameters:
//   - ctx: The context of the function call
//   - msg: The message to estimate the fee for
//   - blockID: The ID of the block to estimate the fee in
//
// Returns:
//   - MessageFeeEstimation: the fee estimated for the message
//   - error: an error if any occurred during the execution
func EstimateMessageFee(
	ctx context.Context,
	c callers.Caller,
	msg rpcv10.MsgFromL1,
	blockID rpcv10.BlockID,
) (rpcv10.MessageFeeEstimation, error) {
	var raw rpcv10.MessageFeeEstimation
	if err := internal.Do(
		ctx, c, "starknet_estimateMessageFee", &raw, msg, blockID,
	); err != nil {
		return raw, rpcerr.UnwrapToRPCErr(err,
			rpcv10.ErrContractError,
			rpcv10.ErrContractNotFound,
			rpcv10.ErrBlockNotFound,
		)
	}

	return raw, nil
}
