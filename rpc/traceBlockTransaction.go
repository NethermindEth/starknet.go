package rpc

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
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
	c callCloser,
	blockID BlockID,
) ([]Trace, error) {
	err := checkForPreConfirmed(blockID)
	if err != nil {
		return nil, err
	}

	var output []Trace
	if err := do(
		ctx, c, "starknet_traceBlockTransactions", &output, blockID,
	); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrBlockNotFound)
	}

	return output, nil
}
