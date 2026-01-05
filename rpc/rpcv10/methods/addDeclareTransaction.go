package methods

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc"
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
	c rpc.Caller,
	declareTransaction *rpcv10.BroadcastDeclareTxnV3,
) (rpcv10.AddDeclareTransactionResponse, error) {
	var result rpcv10.AddDeclareTransactionResponse
	if err := internal.Do(
		ctx, c, "starknet_addDeclareTransaction", &result, declareTransaction,
	); err != nil {
		return rpcv10.AddDeclareTransactionResponse{}, rpcerr.UnwrapToRPCErr(
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
