package methods

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/NethermindEth/starknet.go/rpc/types"
)

// TransactionReceipt fetches the transaction receipt for a given transaction hash.
//
// Parameters:
//   - ctx: the context.Context object for the request
//   - transactionHash: the hash of the transaction as a Felt
//
// Returns:
//   - TransactionReceipt: the transaction receipt
//   - error: an error if any
func GetTransactionReceipt(
	ctx context.Context,
	c callers.Caller,
	transactionHash *felt.Felt,
) (*types.TransactionReceiptWithBlockInfo, error) {
	var receipt types.TransactionReceiptWithBlockInfo
	err := internal.Do(ctx, c, "starknet_getTransactionReceipt", &receipt, transactionHash)
	if err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, rpcv10.ErrHashNotFound)
	}

	return &receipt, nil
}
