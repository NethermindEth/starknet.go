package methods

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
)

// Given an L1 tx hash, returns the associated l1_handler tx hashes and statuses
// for all L1 -> L2 messages sent by the l1 transaction, ordered by the L1 tx sending order
//
// Parameters:
//   - ctx: the context.Context object for cancellation and timeouts.
//   - transactionHash: The hash of the L1 transaction that sent L1->L2 messages
//
// Returns:
//   - [] MessageStatusResp: An array containing the status of the messages sent
//     by the L1 transaction
//   - error, if one arose.
func MessagesStatus(
	ctx context.Context,
	c callers.Caller,
	transactionHash rpcv10.NumAsHex,
) ([]rpcv10.MessageStatus, error) {
	var response []rpcv10.MessageStatus
	err := internal.Do(ctx, c, "starknet_getMessagesStatus", &response, transactionHash)
	if err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, rpcv10.ErrHashNotFound)
	}

	return response, nil
}
