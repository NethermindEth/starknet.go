package methods

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
)

// Events retrieves events from the provider matching the given filter.
//
// Parameters:
//   - ctx: The context to use for the request
//   - input: The input parameters for retrieving events
//
// Returns
//   - eventChunk: The retrieved events
//   - error: An error if any
func Events(
	ctx context.Context,
	c rpc.Caller,
	input rpcv10.EventsInput,
) (*rpcv10.EventChunk, error) {
	var result rpcv10.EventChunk
	if err := internal.Do(ctx, c, "starknet_getEvents", &result, input); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(
			err,
			ErrPageSizeTooBig,
			ErrInvalidContinuationToken,
			ErrBlockNotFound,
			ErrTooManyKeysInFilter,
		)
	}

	return &result, nil
}
