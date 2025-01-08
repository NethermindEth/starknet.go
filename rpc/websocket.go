package rpc

import (
	"context"

	"github.com/NethermindEth/starknet.go/client"
)

// New block headers subscription.
// Creates a WebSocket stream which will fire events for new block headers
//
// Parameters:
// - ctx: The context.Context object for controlling the function call
// - blockID: The ID of the block to retrieve the transactions from
// Returns:
// - subscriptionId: The subscription ID
// - error: An error, if any

func (provider *WsProvider) SubscribeNewHeads(ctx context.Context, ch chan<- *BlockHeader) (*client.ClientSubscription, error) {
	sub, err := provider.c.Subscribe(ctx, "starknet", "_subscribeNewHeads", ch)
	if err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrTooManyBlocksBack, ErrBlockNotFound, ErrCallOnPending)
	}
	return sub, nil
}
