package rpc

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
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
func (provider *Provider) AddInvokeTransaction(
	ctx context.Context,
	invokeTxn *BroadcastInvokeTxnV3,
) (AddInvokeTransactionResponse, error) {
	var output AddInvokeTransactionResponse
	if err := do(ctx, provider.c, "starknet_addInvokeTransaction", &output, invokeTxn); err != nil {
		return AddInvokeTransactionResponse{}, rpcerr.UnwrapToRPCErr(
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
