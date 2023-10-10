package rpc

import (
	"context"
	"errors"
)

// Events retrieves events from the provider matching the given filter.
//
// ctx: The context to use for the request.
// input: The input parameters for retrieving events.
// Returns the chunk of events and an error if any.
func (provider *Provider) Events(ctx context.Context, input EventsInput) (*EventChunk, error) {
	var result EventChunk
	if err := do(ctx, provider.c, "starknet_getEvents", &result, input); err != nil {
		switch {
		case errors.Is(err, ErrPageSizeTooBig):
			return nil, ErrPageSizeTooBig
		case errors.Is(err, ErrInvalidContinuationToken):
			return nil, ErrInvalidContinuationToken
		case errors.Is(err, ErrBlockNotFound):
			return nil, ErrBlockNotFound
		case errors.Is(err, ErrTooManyKeysInFilter):
			return nil, ErrTooManyKeysInFilter
		}
		return nil, err
	}
	return &result, nil
}
