package rpc

import (
	"context"

	"github.com/NethermindEth/starknet.go/client"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
)

// Events subscription.
// Creates a WebSocket stream which will fire events for new Starknet events
// with applied filters. Events are emitted for all events from the specified
// block_id, up to the latest block. If PRE_CONFIRMED finality status is provided,
// events might appear multiple times, for each finality status update. If a single
// event is required, ACCEPTED_ON_L2 must be selected (the default).
//
// Parameters:
//
//   - ctx: The context.Context object for controlling the function call
//   - events: The channel to send the new events to
//   - options: The optional input struct containing the optional filters. Set to
//     nil if no filters are needed.
//
// Returns:
//   - clientSubscription: The client subscription object, used to unsubscribe from
//     the stream and to get errors
//   - error: An error, if any
func (ws *WsProvider) SubscribeEvents(
	ctx context.Context,
	events chan<- *EmittedEventWithFinalityStatus,
	options *EventSubscriptionInput,
) (*client.ClientSubscription, error) {
	sub, err := ws.c.Subscribe(ctx, "starknet", "_subscribeEvents", events, options)
	if err != nil {
		return nil, rpcerr.UnwrapToRPCErr(
			err,
			ErrTooManyKeysInFilter,
			ErrTooManyBlocksBack,
			ErrBlockNotFound,
		)
	}

	return sub, nil
}
