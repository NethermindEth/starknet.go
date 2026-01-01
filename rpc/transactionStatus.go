package rpc

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
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
func (provider *Provider) TransactionStatus(
	ctx context.Context,
	transactionHash *felt.Felt,
) (*TxnStatusResult, error) {
	var receipt TxnStatusResult
	err := do(ctx, provider.c, "starknet_getTransactionStatus", &receipt, transactionHash)
	if err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrHashNotFound)
	}

	return &receipt, nil
}
