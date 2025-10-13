package rpc

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
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
func (ws *WsProvider) SubscribeNewHeads(
	ctx context.Context,
	headers chan<- *BlockHeader,
	subBlockID SubscriptionBlockID,
) (*client.ClientSubscription, error) {
	var sub *client.ClientSubscription
	var err error

	// if subBlockID is empty, don't send it to the server to avoid it being marshalled as 'null'
	if subBlockID == (SubscriptionBlockID{}) { //nolint:exhaustruct // Asserting the type
		sub, err = ws.c.SubscribeWithSliceArgs(ctx, "starknet", "_subscribeNewHeads", headers)
	} else {
		sub, err = ws.c.SubscribeWithSliceArgs(
			ctx, "starknet", "_subscribeNewHeads", headers, subBlockID,
		)
	}

	if err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrTooManyBlocksBack, ErrBlockNotFound)
	}

	return sub, nil
}

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
func (ws *WsProvider) SubscribeNewTransactionReceipts(
	ctx context.Context,
	txnReceipts chan<- *TransactionReceiptWithBlockInfo,
	options *SubNewTxnReceiptsInput,
) (*client.ClientSubscription, error) {
	sub, err := ws.c.Subscribe(
		ctx,
		"starknet",
		"_subscribeNewTransactionReceipts",
		txnReceipts,
		options,
	)
	if err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrTooManyAddressesInFilter)
	}

	return sub, nil
}

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
func (ws *WsProvider) SubscribeNewTransactions(
	ctx context.Context,
	newTxns chan<- *TxnWithHashAndStatus,
	options *SubNewTxnsInput,
) (*client.ClientSubscription, error) {
	sub, err := ws.c.Subscribe(ctx, "starknet", "_subscribeNewTransactions", newTxns, options)
	if err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrTooManyAddressesInFilter)
	}

	return sub, nil
}

// Transaction Status subscription.
// Creates a WebSocket stream which at first fires an event with the current known
// transaction status,
// followed by events for every transaction status update
//
// Parameters:
//   - ctx: The context.Context object for controlling the function call
//   - newStatus: The channel to send the new transaction status to
//   - transactionHash: The transaction hash to fetch status updates for
//
// Returns:
//   - clientSubscription: The client subscription object, used to unsubscribe from
//     the stream and to get errors
//   - error: An error, if any
func (ws *WsProvider) SubscribeTransactionStatus(
	ctx context.Context,
	newStatus chan<- *NewTxnStatus,
	transactionHash *felt.Felt,
) (*client.ClientSubscription, error) {
	sub, err := ws.c.SubscribeWithSliceArgs(
		ctx,
		"starknet",
		"_subscribeTransactionStatus",
		newStatus,
		transactionHash,
	)
	if err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err)
	}

	return sub, nil
}
