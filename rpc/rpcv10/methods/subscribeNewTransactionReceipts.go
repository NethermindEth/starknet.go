package methods

import (
	"context"

	"github.com/NethermindEth/starknet.go/client"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
)

// New transactions receipts subscription
// Creates a WebSocket stream which will fire events when new transaction
// receipts are created. The endpoint receives a vector of finality statuses.
// An event is fired for each finality status update.
// It is possible for receipts for pre-confirmed transactions to be received
// multiple times, or not at all.
//
// Parameters:
//   - ctx: The context.Context object for controlling the function call
//   - txnReceipts: The channel to send the new transaction receipts to
//   - options: The optional input struct containing the optional filters. Set
//     to nil if no filters are needed.
//
// Returns:
//   - clientSubscription: The client subscription object, used to unsubscribe
//     from the stream and to get errors
//   - error: An error, if any
func SubscribeNewTransactionReceipts(
	ctx context.Context,
	ws callers.Subscriber,
	txnReceipts chan<- *rpcv10.TransactionReceiptWithBlockInfo,
	options *rpcv10.SubNewTxnReceiptsInput,
) (*client.ClientSubscription, error) {
	sub, err := ws.Subscribe(
		ctx,
		"starknet",
		"_subscribeNewTransactionReceipts",
		txnReceipts,
		options,
	)
	if err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, rpcv10.ErrTooManyAddressesInFilter)
	}

	return sub, nil
}
