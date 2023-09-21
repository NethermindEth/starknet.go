package rpc

import (
	"context"
	"errors"
)

// Events returns all events matching the given filter
func (provider *Provider) Events(ctx context.Context, input EventsInput) (*EventChunk, error) {
	var result EventChunk
	if err := do(ctx, provider.c, "starknet_getEvents", &result, input); err != nil {

		return nil, tryUnwrapToRPCErr(err,  ErrPageSizeTooBig , ErrInvalidContinuationToken , ErrBlockNotFound ,ErrTooManyKeysInFilter)
	}
	return &result, nil
}
