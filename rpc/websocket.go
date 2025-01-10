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
// - headers: The channel to send the new block headers to
// - blockID (optional): The block to get notifications from, default is latest, limited to 1024 blocks back
// Returns:
// - subscriptionId: The subscription ID
// - error: An error, if any
func (provider *WsProvider) SubscribeNewHeads(ctx context.Context, headers chan<- *BlockHeader, blockID ...BlockID) (*client.ClientSubscription, error) {
	// Convert blockID to []any
	params := make([]any, len(blockID))
	for i, v := range blockID {
		params[i] = v
	}

	sub, err := provider.c.Subscribe(ctx, "starknet", "_subscribeNewHeads", headers, params...)
	if err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrTooManyBlocksBack, ErrBlockNotFound, ErrCallOnPending)
	}
	return sub, nil
}
