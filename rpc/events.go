package rpc

import (
	"context"

	"github.com/dontpanicdao/caigo/rpc/types"
)

// Events returns all events matching the given filter
func (sc *Provider) Events(ctx context.Context, filter types.EventFilter) (*types.EventsOutput, error) {
	var result types.EventsOutput
	if err := do(ctx, sc.c, "starknet_getEvents", &result, filter); err != nil {
		return nil, err
	}

	return &result, nil
}
