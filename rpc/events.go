package rpc

import (
	"context"
)

// Events returns all events matching the given filter
func (provider *Provider) Events(ctx context.Context, input EventsInput) (*EventsOutput, error) {
	var result EventsOutput
	if err := do(ctx, provider.c, "starknet_getEvents", &result, input); err != nil {
		// TODO: Check with Pathfinder/Devnet for errors
		// PAGE_SIZE_TOO_BIG, INVALID_CONTINUATION_TOKEN, BLOCK_NOT_FOUND or TOO_MANY_KEYS_IN_FILTER
		return nil, err
	}

	return &result, nil
}
