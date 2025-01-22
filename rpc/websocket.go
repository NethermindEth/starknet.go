package rpc

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client"
)

// New block headers subscription.
// Creates a WebSocket stream which will fire events for new block headers
//
// Parameters:
// - ctx: The context.Context object for controlling the function call
// - headers: The channel to send the new block headers to
// - blockID (optional): The block to get notifications from, limited to 1024 blocks back. If set to nil, the latest block will be used
// Returns:
// - clientSubscription: The client subscription object, used to unsubscribe from the stream and to get errors
// - error: An error, if any
func (provider *WsProvider) SubscribeNewHeads(ctx context.Context, headers chan<- *BlockHeader, blockID *BlockID) (*client.ClientSubscription, error) {
	if blockID == nil {
		blockID = &BlockID{Tag: "latest"}
	}

	sub, err := provider.c.SubscribeWithSliceArgs(ctx, "starknet", "_subscribeNewHeads", headers, blockID)
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
	var sub *client.ClientSubscription
	var err error

	var emptyBlockID BlockID
	if input.BlockID == emptyBlockID {
		// BlockID has a custom MarshalJSON that doesn't allow zero values.
		// Create a temporary struct without BlockID field to properly handle the optional parameter.
		tempInput := struct {
			FromAddress *felt.Felt     `json:"from_address,omitempty"`
			Keys        [][]*felt.Felt `json:"keys,omitempty"`
		}{
			FromAddress: input.FromAddress,
			Keys:        input.Keys,
		}

		sub, err = provider.c.Subscribe(ctx, "starknet", "_subscribeEvents", events, tempInput)
	} else {
		sub, err = provider.c.Subscribe(ctx, "starknet", "_subscribeEvents", events, input)
	}

	if err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrTooManyKeysInFilter, ErrTooManyBlocksBack, ErrBlockNotFound, ErrCallOnPending)
	}
	return sub, nil
}

// Transaction Status subscription.
// Creates a WebSocket stream which at first fires an event with the current known transaction status,
// followed by events for every transaction status update
//
// Parameters:
// - ctx: The context.Context object for controlling the function call
// - newStatus: The channel to send the new transaction status to
// - transactionHash: The transaction hash to fetch status updates for
// Returns:
// - clientSubscription: The client subscription object, used to unsubscribe from the stream and to get errors
// - error: An error, if any
func (provider *WsProvider) SubscribeTransactionStatus(ctx context.Context, newStatus chan<- *NewTxnStatusResp, transactionHash *felt.Felt) (*client.ClientSubscription, error) {
	sub, err := provider.c.SubscribeWithSliceArgs(ctx, "starknet", "_subscribeTransactionStatus", newStatus, transactionHash, WithBlockTag("latest"))
	if err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrTooManyBlocksBack, ErrBlockNotFound)
	}
	// TODO: wait for Juno to implement this. This is the correct implementation by the spec
	// 	sub, err := provider.c.SubscribeWithSliceArgs(ctx, "starknet", "_subscribeTransactionStatus", newStatus, transactionHash)
	// 	if err != nil {
	// 		return nil, tryUnwrapToRPCErr(err)
	// 	}
	return sub, nil
}

// New Pending Transactions subscription
// Creates a WebSocket stream which will fire events when a new pending transaction is added.
// While there is no mempool, this notifies of transactions in the pending block.
//
// Parameters:
// - ctx: The context.Context object for controlling the function call
// - pendingTxns: The channel to send the new pending transactions to
// - options: The optional input struct containing the optional filters. Set to nil if no filters are needed.
// Returns:
// - clientSubscription: The client subscription object, used to unsubscribe from the stream and to get errors
// - error: An error, if any
func (provider *WsProvider) SubscribePendingTransactions(ctx context.Context, pendingTxns chan<- *SubPendingTxns, options *SubPendingTxnsInput) (*client.ClientSubscription, error) {
	if options == nil {
		options = &SubPendingTxnsInput{}
	}

	sub, err := provider.c.Subscribe(ctx, "starknet", "_subscribePendingTransactions", pendingTxns, options)
	if err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrTooManyAddressesInFilter)
	}
	return sub, nil
}
