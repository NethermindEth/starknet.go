package methods

import (
	"context"

	"github.com/NethermindEth/starknet.go/client"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
)

// New transactions subscription
// Creates a WebSocket stream which will fire events when new transaction are created.
// The endpoint receives a vector of finality statuses. An event is fired for each finality
// status update. It is possible for events for pre-confirmed and candidate transactions
// to be received multiple times, or not at all.
//
// Parameters:
//   - ctx: The context.Context object for controlling the function call
//   - newTxns: The channel to send the new transactions to
//   - options: The optional input struct containing the optional filters. Set to nil if
//     no filters are needed.
//
// Returns:
//   - clientSubscription: The client subscription object, used to unsubscribe from
//     the stream and to get errors
//   - error: An error, if any
func SubscribeNewTransactions(
	ctx context.Context,
	ws Subscriber,
	newTxns chan<- *rpcv10.TxnWithHashAndStatus,
	options *rpcv10.SubNewTxnsInput,
) (*client.ClientSubscription, error) {
	sub, err := ws.Subscribe(ctx, "starknet", "_subscribeNewTransactions", newTxns, options)
	if err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, rpcv10.ErrTooManyAddressesInFilter)
	}

	return sub, nil
}
