package rpc

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
)

// AddDeclareTransaction submits a declare transaction to the StarkNet contract.
//
// Parameters:
//   - ctx: The context.Context object for the request.
//   - declareTransaction: The input for the declare transaction.
//
// Returns:
//   - AddDeclareTransactionResponse: The response of submitting the declare transaction
//   - error: an error if any
func AddDeclareTransaction(
	ctx context.Context,
	c callCloser,
	declareTransaction *BroadcastDeclareTxnV3,
) (AddDeclareTransactionResponse, error) {
	var result AddDeclareTransactionResponse
	if err := do(
		ctx, c, "starknet_addDeclareTransaction", &result, declareTransaction,
	); err != nil {
		return AddDeclareTransactionResponse{}, rpcerr.UnwrapToRPCErr(
			err,
			ErrClassAlreadyDeclared,
			ErrCompilationFailed,
			ErrCompiledClassHashMismatch,
			ErrInsufficientAccountBalance,
			ErrInsufficientResourcesForValidate,
			ErrInvalidTransactionNonce,
			ErrReplacementTransactionUnderpriced,
			ErrFeeBelowMinimum,
			ErrValidationFailure,
			ErrNonAccount,
			ErrDuplicateTx,
			ErrContractClassSizeTooLarge,
			ErrUnsupportedTxVersion,
			ErrUnsupportedContractClassVersion,
		)
	}

	return result, nil
}
