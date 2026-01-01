package rpc

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
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
func (provider *Provider) EstimateMessageFee(
	ctx context.Context,
	msg MsgFromL1,
	blockID BlockID,
) (MessageFeeEstimation, error) {
	var raw MessageFeeEstimation
	if err := do(
		ctx, provider.c, "starknet_estimateMessageFee", &raw, msg, blockID,
	); err != nil {
		return raw, rpcerr.UnwrapToRPCErr(err,
			ErrContractError,
			ErrContractNotFound,
			ErrBlockNotFound,
		)
	}

	return raw, nil
}
