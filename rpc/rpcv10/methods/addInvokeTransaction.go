package methods

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
)

// AddInvokeTransaction adds an invoke transaction to the provider.
//
// Parameters:
//   - ctx: The context for the function.
//   - invokeTxn: The invoke transaction to be added.
//
// Returns:
//   - AddInvokeTransactionResponse: the response of adding the invoke transaction
//   - error: an error if any
func AddInvokeTransaction(
	ctx context.Context,
	c callers.Caller,
	invokeTxn *rpcv10.BroadcastInvokeTxnV3,
) (rpcv10.AddInvokeTransactionResponse, error) {
	var output rpcv10.AddInvokeTransactionResponse
	if err := internal.Do(ctx, c, "starknet_addInvokeTransaction", &output, invokeTxn); err != nil {
		return rpcv10.AddInvokeTransactionResponse{}, rpcerr.UnwrapToRPCErr(
			err,
			rpcv10.ErrInsufficientAccountBalance,
			rpcv10.ErrInsufficientResourcesForValidate,
			rpcv10.ErrInvalidTransactionNonce,
			rpcv10.ErrReplacementTransactionUnderpriced,
			rpcv10.ErrFeeBelowMinimum,
			rpcv10.ErrValidationFailure,
			rpcv10.ErrNonAccount,
			rpcv10.ErrDuplicateTx,
			rpcv10.ErrUnsupportedTxVersion,
			rpcv10.ErrUnexpectedError,
		)
	}

	return output, nil
}
