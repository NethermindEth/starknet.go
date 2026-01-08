package methods

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/NethermindEth/starknet.go/rpc/types"
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
	declareTransaction *types.BroadcastDeclareTxnV3,
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
