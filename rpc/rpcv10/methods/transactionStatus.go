package methods

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
)

// TransactionStatus gets the transaction status (possibly reflecting that
// the tx is still in the mempool, or dropped from it)
// Parameters:
//   - ctx: the context.Context object for cancellation and timeouts.
//   - transactionHash: The hash of the requested transaction
//
// Returns:
//   - *TxnStatusResult: Transaction status result, including finality status
//     and execution status
//   - error, if one arose.
func TransactionStatus(
	ctx context.Context,
	c callers.Caller,
	transactionHash *felt.Felt,
) (*rpcv10.TxnStatusResult, error) {
	var receipt rpcv10.TxnStatusResult
	err := internal.Do(ctx, c, "starknet_getTransactionStatus", &receipt, transactionHash)
	if err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, rpcv10.ErrHashNotFound)
	}

	return &receipt, nil
}
