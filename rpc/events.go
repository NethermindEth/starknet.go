package rpc

import (
	"context"
	"errors"
)

// Events returns all events based on the given input filter parameters.
//
// ctx - The context in which the function is called.
// input - The input parameters for retrieving events.
// Returns a pointer to the EventChunk struct and an error.
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
