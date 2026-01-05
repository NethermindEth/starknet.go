package methods

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
)

// TraceBlockTransactions retrieves the traces of transactions in a given block.
//
// Parameters:
//   - ctx: the context.Context object for controlling the request
//   - blockID: the block to retrieve the traces from. `pre_confirmed` tag is not allowed
//
// Returns:
//   - []Trace: a slice of Trace objects representing the traces of transactions in the block
//   - error: an error if there was a problem retrieving the traces.
func TraceBlockTransactions(
	ctx context.Context,
	c rpc.Caller,
	blockID rpcv10.BlockID,
) ([]rpcv10.Trace, error) {
	err := checkForPreConfirmed(blockID)
	if err != nil {
		return nil, err
	}

	var output []rpcv10.Trace
	if err := internal.Do(
		ctx, c, "starknet_traceBlockTransactions", &output, blockID,
	); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrBlockNotFound)
	}

	return output, nil
}
