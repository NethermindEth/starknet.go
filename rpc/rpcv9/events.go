package rpcv9

import (
	"context"

	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
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
	c callers.Caller,
	input EventsInput,
) (*EventChunk, error) {
	var result EventChunk
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
