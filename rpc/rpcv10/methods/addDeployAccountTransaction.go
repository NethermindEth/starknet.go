package methods

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/NethermindEth/starknet.go/types"
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
	c callers.Caller,
	deployAccountTransaction *types.BroadcastDeployAccountTxnV3,
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
