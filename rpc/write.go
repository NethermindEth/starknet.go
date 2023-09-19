package rpc

import (
	"context"
)

type BroadcastedInvokeTransaction interface{}

func (provider *Provider) AddInvokeTransaction(ctx context.Context, broadcastedInvoke BroadcastedInvokeTransaction) (*AddInvokeTransactionResponse, error) {
	var output AddInvokeTransactionResponse
	switch invoke := broadcastedInvoke.(type) {
	case BroadcastedInvokeV1Transaction:
		if err := do(ctx, provider.c, "starknet_addInvokeTransaction", &output, invoke); err != nil {
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
	return nil, Err(InvalidParams, "invalid method parameter(s)")
}

func (provider *Provider) AddDeclareTransaction(ctx context.Context, declareTransaction BroadcastedDeclareTransaction) (*AddDeclareTransactionResponse, error) {
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

func (provider *Provider) AddDeployAccountTransaction(ctx context.Context, deployAccountTransaction BroadcastedDeployAccountTransaction) (*AddDeployAccountTransactionResponse, error) {
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
