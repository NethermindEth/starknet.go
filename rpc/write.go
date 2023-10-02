package rpc

import (
	"context"
)

// AddInvokeTransaction adds an invoke transaction to the provider.
// INVOKE syscall
//
// ctx: The context.Context used for the request.
// invokeTxn: The invoke transaction to be added.
// Returns the AddInvokeTransactionResponse and an error.
func (provider *Provider) AddInvokeTransaction(ctx context.Context, invokeTxn InvokeTxnV1) (*AddInvokeTransactionResponse, error) {
	var output AddInvokeTransactionResponse
	if err := do(ctx, provider.c, "starknet_addInvokeTransaction", &output, invokeTxn); err != nil {
		if unexpectedErr, ok := isErrUnexpectedError(err); ok {
			return nil, unexpectedErr
		}
		return nil, tryUnwrapToRPCErr(
			err,
			ErrInsufficientAccountBalance,
			ErrInsufficientMaxFee,
			ErrInvalidTransactionNonce,
			ErrValidationFailure,
			ErrNonAccount,
			ErrDuplicateTx,
			ErrUnsupportedTxVersion,
		)
	}
	return &output, nil
}

// AddDeclareTransaction adds a declare transaction to the Provider.
// DECLARE syscall
//
// ctx is the context.Context for the function.
// declareTransaction is the input for the transaction.
// It returns *AddDeclareTransactionResponse and an error.
func (provider *Provider) AddDeclareTransaction(ctx context.Context, declareTransaction AddDeclareTxnInput) (*AddDeclareTransactionResponse, error) {
	var result AddDeclareTransactionResponse
	if err := do(ctx, provider.c, "starknet_addDeclareTransaction", &result, declareTransaction); err != nil {
		if unexpectedErr, ok := isErrUnexpectedError(err); ok {
			return nil, unexpectedErr
		}
		return nil, tryUnwrapToRPCErr(
			err,
			ErrClassAlreadyDeclared,
			ErrCompilationFailed,
			ErrCompiledClassHashMismatch,
			ErrInsufficientAccountBalance,
			ErrInsufficientMaxFee,
			ErrInvalidTransactionNonce,
			ErrValidationFailure,
			ErrNonAccount,
			ErrDuplicateTx,
			ErrContractClassSizeTooLarge,
			ErrUnsupportedTxVersion,
			ErrUnsupportedContractClassVersion,
		)
	}
	return &result, nil
}

// AddDeployAccountTransaction adds a deploy account transaction to the Provider.
// DEPLOY_ACCOUNT syscall
//
// ctx: context.Context - The context for the function.
// deployAccountTransaction: DeployAccountTxn - The deploy account transaction to be added.
// Return type: *AddDeployAccountTransactionResponse - The response from adding the deploy account transaction.
// error - An error if any occurred.
func (provider *Provider) AddDeployAccountTransaction(ctx context.Context, deployAccountTransaction DeployAccountTxn) (*AddDeployAccountTransactionResponse, error) {
	var result AddDeployAccountTransactionResponse
	if err := do(ctx, provider.c, "starknet_addDeployAccountTransaction", &result, deployAccountTransaction); err != nil {
		if unexpectedErr, ok := isErrUnexpectedError(err); ok {
			return nil, unexpectedErr
		}

		return nil, tryUnwrapToRPCErr(
			err,
			ErrInsufficientAccountBalance,
			ErrInsufficientMaxFee,
			ErrInvalidTransactionNonce,
			ErrValidationFailure,
			ErrNonAccount,
			ErrClassHashNotFound,
			ErrDuplicateTx,
			ErrUnsupportedTxVersion,
		)
	}
	return &result, nil
}
