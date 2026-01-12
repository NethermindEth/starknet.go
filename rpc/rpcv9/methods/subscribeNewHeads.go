package methods

import (
	"context"

	"github.com/NethermindEth/starknet.go/client"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/NethermindEth/starknet.go/rpc/types"
)

// New block headers subscription.
// Creates a WebSocket stream which will fire events for new block headers
//
// Parameters:
//   - ctx: The context.Context object for controlling the function call
//   - headers: The channel to send the new block headers to
//   - blockID (optional): The block to get notifications from, limited to 1024
//     blocks back. If empty, the latest block will be used
//
// Returns:
//   - clientSubscription: The client subscription object, used to unsubscribe from
//     the stream and to get errors
//   - error: An error, if any
func SubscribeNewHeads(
	ctx context.Context,
	ws callers.Subscriber,
	headers chan<- *rpcv10.BlockHeader,
	subBlockID types.SubscriptionBlockID,
) (*client.ClientSubscription, error) {
	var sub *client.ClientSubscription
	var err error

	// @todo see why not accept subBlockID as a pointer
	// if subBlockID is empty, don't send it to the server to avoid it being marshalled as 'null'
	if subBlockID == (types.SubscriptionBlockID{}) { //nolint:exhaustruct // Asserting the type
		sub, err = ws.SubscribeWithSliceArgs(ctx, "starknet", "_subscribeNewHeads", headers)
	} else {
		sub, err = ws.SubscribeWithSliceArgs(
			ctx, "starknet", "_subscribeNewHeads", headers, subBlockID,
		)
	}

	if err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, rpcv10.ErrTooManyBlocksBack, rpcv10.ErrBlockNotFound)
	}

	return sub, nil
}
