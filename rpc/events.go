package rpc

import (
	"context"

	internalutils "github.com/NethermindEth/starknet.go/internal/utils"
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
//
//nolint:gocritic
func (provider *Provider) Events(ctx context.Context, input EventsInput) (*EventChunk, error) {
	var result EventChunk
	if err := do(ctx, provider.c, "starknet_getEvents", &result, input); err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrPageSizeTooBig, ErrInvalidContinuationToken, ErrBlockNotFound, ErrTooManyKeysInFilter)
	}

	return &result, nil
}

func EventWith(events []Event, key string) *Event {
	feltKey, err := internalutils.HexToFelt(key)
	if err != nil {
		return nil
	}

	for i := range events {
		for _, k := range events[i].Keys {
			if k.Equal(feltKey) {
				return &events[i]
			}
		}
	}
	return nil
}
