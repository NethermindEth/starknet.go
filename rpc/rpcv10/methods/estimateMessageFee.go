package methods

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
)

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
	c callCloser,
	msg rpcv10.MsgFromL1,
	blockID rpcv10.BlockID,
) (rpcv10.MessageFeeEstimation, error) {
	var raw rpcv10.MessageFeeEstimation
	if err := do(
		ctx, c, "starknet_estimateMessageFee", &raw, msg, blockID,
	); err != nil {
		return raw, rpcerr.UnwrapToRPCErr(err,
			ErrContractError,
			ErrContractNotFound,
			ErrBlockNotFound,
		)
	}

	return raw, nil
}
