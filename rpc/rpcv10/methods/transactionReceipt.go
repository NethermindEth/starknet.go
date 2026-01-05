package methods

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
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
	c rpc.Caller,
	transactionHash *felt.Felt,
) (*rpcv10.TransactionReceiptWithBlockInfo, error) {
	var receipt rpcv10.TransactionReceiptWithBlockInfo
	err := internal.Do(ctx, c, "starknet_getTransactionReceipt", &receipt, transactionHash)
	if err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, rpcv10.ErrHashNotFound)
	}

	return &receipt, nil
}
