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
// - clientSubscription: The client subscription object, used to unsubscribe from the stream and to get errors
// - error: An error, if any
func (provider *WsProvider) SubscribeNewHeads(ctx context.Context, headers chan<- *BlockHeader, blockID ...BlockID) (*client.ClientSubscription, error) {
	// Convert blockID to []any
	params := make([]any, len(blockID))
	for i, v := range blockID {
		params[i] = v
	}

	sub, err := provider.c.SubscribeWithSliceArgs(ctx, "starknet", "_subscribeNewHeads", headers, params...)
	if err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrTooManyBlocksBack, ErrBlockNotFound, ErrCallOnPending)
	}
	return sub, nil
}

// Events subscription.
// Creates a WebSocket stream which will fire events for new Starknet events with applied filters
//
// Parameters:
// - ctx: The context.Context object for controlling the function call
// - events: The channel to send the new events to
// - input: The input struct containing the optional filters
// Returns:
// - clientSubscription: The client subscription object, used to unsubscribe from the stream and to get errors
// - error: An error, if any
func (provider *WsProvider) SubscribeEvents(ctx context.Context, events chan<- *EmittedEvent, input EventSubscriptionInput) (*client.ClientSubscription, error) {
	sub, err := provider.c.Subscribe(ctx, "starknet", "_subscribeEvents", events, input)
	if err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrTooManyKeysInFilter, ErrTooManyBlocksBack, ErrBlockNotFound, ErrCallOnPending)
	}
	return sub, nil
}
