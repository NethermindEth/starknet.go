package methods

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
)

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
	c callCloser,
	deployAccountTransaction *rpcv10.BroadcastDeployAccountTxnV3,
) (rpcv10.AddDeployAccountTransactionResponse, error) {
	var result rpcv10.AddDeployAccountTransactionResponse
	if err := do(
		ctx, c, "starknet_addDeployAccountTransaction", &result, deployAccountTransaction,
	); err != nil {
		return rpcv10.AddDeployAccountTransactionResponse{}, rpcerr.UnwrapToRPCErr(
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
