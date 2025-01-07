package rpc

import "context"

// New block headers subscription.
// Creates a WebSocket stream which will fire events for new block headers
//
// Parameters:
// - ctx: The context.Context object for controlling the function call
// - blockID: The ID of the block to retrieve the transactions from
// Returns:
// - subscriptionId: The subscription ID
// - error: An error, if any
func (provider *WsProvider) SubscribeNewHeads(ctx context.Context, blockID BlockID) (subscriptionId int, err error) {
	if err = do(ctx, provider.c, "starknet_subscribeNewHeads", &subscriptionId, blockID); err != nil {
		return 0, tryUnwrapToRPCErr(err, ErrTooManyBlocksBack, ErrBlockNotFound, ErrCallOnPending)
	}
	return subscriptionId, nil
}
