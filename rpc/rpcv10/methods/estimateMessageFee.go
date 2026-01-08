package methods

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/NethermindEth/starknet.go/rpc/types"
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
	c callers.Caller,
	msg rpcv10.MsgFromL1,
	blockID rpcv10.BlockID,
) (types.MessageFeeEstimation, error) {
	var raw types.MessageFeeEstimation
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
