package rpc

import (
	"context"
	"errors"
)

// Events returns all events matching the given filter
func (provider *Provider) Events(ctx context.Context, input EventsInput) (*EventsOutput, error) {
	var result EventsOutput
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
