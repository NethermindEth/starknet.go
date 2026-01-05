package methods

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
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
	c callCloser,
	invokeTxn *rpcv10.BroadcastInvokeTxnV3,
) (rpcv10.AddInvokeTransactionResponse, error) {
	var output rpcv10.AddInvokeTransactionResponse
	if err := do(ctx, c, "starknet_addInvokeTransaction", &output, invokeTxn); err != nil {
		return rpcv10.AddInvokeTransactionResponse{}, rpcerr.UnwrapToRPCErr(
			err,
			ErrInsufficientAccountBalance,
			ErrInsufficientResourcesForValidate,
			ErrInvalidTransactionNonce,
			ErrReplacementTransactionUnderpriced,
			ErrFeeBelowMinimum,
			ErrValidationFailure,
			ErrNonAccount,
			ErrDuplicateTx,
			ErrUnsupportedTxVersion,
			ErrUnexpectedError,
		)
	}

	return output, nil
}
