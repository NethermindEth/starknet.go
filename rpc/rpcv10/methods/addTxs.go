package methods

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
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
	declareTransaction *rpcv10.BroadcastDeclareTxnV3,
) (rpcv10.AddDeclareTransactionResponse, error) {
	var result rpcv10.AddDeclareTransactionResponse
	if err := internal.Do(
		ctx, c, "starknet_addDeclareTransaction", &result, declareTransaction,
	); err != nil {
		return rpcv10.AddDeclareTransactionResponse{}, rpcerr.UnwrapToRPCErr(
			err,
			rpcv10.ErrClassAlreadyDeclared,
			rpcv10.ErrCompilationFailed,
			rpcv10.ErrCompiledClassHashMismatch,
			rpcv10.ErrInsufficientAccountBalance,
			rpcv10.ErrInsufficientResourcesForValidate,
			rpcv10.ErrInvalidTransactionNonce,
			rpcv10.ErrReplacementTransactionUnderpriced,
			rpcv10.ErrFeeBelowMinimum,
			rpcv10.ErrValidationFailure,
			rpcv10.ErrNonAccount,
			rpcv10.ErrDuplicateTx,
			rpcv10.ErrContractClassSizeTooLarge,
			rpcv10.ErrUnsupportedTxVersion,
			rpcv10.ErrUnsupportedContractClassVersion,
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
	deployAccountTransaction *rpcv10.BroadcastDeployAccountTxnV3,
) (rpcv10.AddDeployAccountTransactionResponse, error) {
	var result rpcv10.AddDeployAccountTransactionResponse
	if err := internal.Do(
		ctx, c, "starknet_addDeployAccountTransaction", &result, deployAccountTransaction,
	); err != nil {
		return rpcv10.AddDeployAccountTransactionResponse{}, rpcerr.UnwrapToRPCErr(
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
			rpcv10.ErrClassHashNotFound,
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
