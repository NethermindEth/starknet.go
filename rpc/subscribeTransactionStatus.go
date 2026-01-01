package rpc

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
)

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
