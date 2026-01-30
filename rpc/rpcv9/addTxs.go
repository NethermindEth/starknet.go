package rpcv9

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
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
	c callers.Caller,
	declareTransaction *BroadcastDeclareTxnV3,
) (AddDeclareTransactionResponse, error) {
	var result AddDeclareTransactionResponse
	if err := internal.Do(
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

// AddDeployAccountTransaction adds a DEPLOY_ACCOUNT transaction to the provider.
//
// Parameters:
//   - ctx: The context of the function
//   - deployAccountTransaction: The deploy account transaction to be added
//
// Returns:
//   - AddDeployAccountTransactionResponse: the response of adding the deploy
//     account transaction or an error
func AddDeployAccountTransaction(
	ctx context.Context,
	c callers.Caller,
	deployAccountTransaction *BroadcastDeployAccountTxnV3,
) (AddDeployAccountTransactionResponse, error) {
	var result AddDeployAccountTransactionResponse
	if err := internal.Do(
		ctx, c, "starknet_addDeployAccountTransaction", &result, deployAccountTransaction,
	); err != nil {
		return AddDeployAccountTransactionResponse{}, rpcerr.UnwrapToRPCErr(
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
			ErrClassHashNotFound,
		)
	}

	return result, nil
}

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
	invokeTxn *BroadcastInvokeTxnV3,
) (AddInvokeTransactionResponse, error) {
	var output AddInvokeTransactionResponse
	if err := internal.Do(ctx, c, "starknet_addInvokeTransaction", &output, invokeTxn); err != nil {
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
