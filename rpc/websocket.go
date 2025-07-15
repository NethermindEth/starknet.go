package rpc

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client"
)

// Events subscription.
// Creates a WebSocket stream which will fire events for new Starknet events with applied filters
//
// Parameters:
//
// - ctx: The context.Context object for controlling the function call
//
// - events: The channel to send the new events to
//
// - options: The optional input struct containing the optional filters. Set to nil if no filters are needed.
//
//   - fromAddress: Filter events by from_address which emitted the event
//   - keys: Per key (by position), designate the possible values to be matched for events to be returned.
//     Empty array designates 'any' value
//   - blockID: The block to get notifications from, limited to 1024 blocks back. If empty, the latest block will be used
//
// Returns:
//   - clientSubscription: The client subscription object, used to unsubscribe from the stream and to get errors
//   - error: An error, if any
func (provider *WsProvider) SubscribeEvents(
	ctx context.Context,
	events chan<- *EmittedEvent,
	options *EventSubscriptionInput,
) (*client.ClientSubscription, error) {
	if options == nil {
		options = &EventSubscriptionInput{}
	}

	err := checkForPre_confirmed(options.BlockID)
	if err != nil {
		return nil, err
	}

	sub, err := provider.c.Subscribe(ctx, "starknet", "_subscribeEvents", events, options)
	if err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrTooManyKeysInFilter, ErrTooManyBlocksBack, ErrBlockNotFound)
	}

	return sub, nil
}

// New block headers subscription.
// Creates a WebSocket stream which will fire events for new block headers
//
// Parameters:
//   - ctx: The context.Context object for controlling the function call
//   - headers: The channel to send the new block headers to
//   - blockID (optional): The block to get notifications from, limited to 1024 blocks back. If empty, the latest block will be used
//
// Returns:
//   - clientSubscription: The client subscription object, used to unsubscribe from the stream and to get errors
//   - error: An error, if any
func (provider *WsProvider) SubscribeNewHeads(
	ctx context.Context,
	headers chan<- *BlockHeader,
	blockID BlockID,
) (*client.ClientSubscription, error) {
	err := checkForPre_confirmed(blockID)
	if err != nil {
		return nil, err
	}

	sub, err := provider.c.SubscribeWithSliceArgs(ctx, "starknet", "_subscribeNewHeads", headers, blockID)
	if err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrTooManyBlocksBack, ErrBlockNotFound)
	}

	return sub, nil
}

// New Pending Transactions subscription
// Creates a WebSocket stream which will fire events when a new pending transaction is added.
// While there is no mempool, this notifies of transactions in the pending block.
//
// Parameters:
//   - ctx: The context.Context object for controlling the function call
//   - pendingTxns: The channel to send the new pending transactions to
//   - options: The optional input struct containing the optional filters. Set to nil if no filters are needed.
//
// Returns:
//   - clientSubscription: The client subscription object, used to unsubscribe from the stream and to get errors
//   - error: An error, if any
func (provider *WsProvider) SubscribePendingTransactions(
	ctx context.Context,
	pendingTxns chan<- *PendingTxn,
	options *SubPendingTxnsInput,
) (*client.ClientSubscription, error) {
	sub, err := provider.c.Subscribe(ctx, "starknet", "_subscribePendingTransactions", pendingTxns, options)
	if err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrTooManyAddressesInFilter)
	}

	return sub, nil
}

// Transaction Status subscription.
// Creates a WebSocket stream which at first fires an event with the current known transaction status,
// followed by events for every transaction status update
//
// Parameters:
//   - ctx: The context.Context object for controlling the function call
//   - newStatus: The channel to send the new transaction status to
//   - transactionHash: The transaction hash to fetch status updates for
//
// Returns:
//   - clientSubscription: The client subscription object, used to unsubscribe from the stream and to get errors
//   - error: An error, if any
func (provider *WsProvider) SubscribeTransactionStatus(
	ctx context.Context,
	newStatus chan<- *NewTxnStatus,
	transactionHash *felt.Felt,
) (*client.ClientSubscription, error) {
	sub, err := provider.c.SubscribeWithSliceArgs(ctx, "starknet", "_subscribeTransactionStatus", newStatus, transactionHash)
	if err != nil {
		return nil, tryUnwrapToRPCErr(err)
	}

	return sub, nil
}
