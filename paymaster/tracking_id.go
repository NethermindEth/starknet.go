package paymaster

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
)

// TrackingIdToLatestHash gets the latest transaction hash and status for a given tracking ID.
// Returns a TrackingIdResponse.
//
// Parameters:
//   - ctx: The context.Context object for controlling the function call
//   - trackingId: A unique identifier used to track an execution request of a user.
//     This identitifier is returned by the paymaster after a successful call to `execute`.
//     Its purpose is to track the possibly different transaction hashes in the mempool which
//     are associated with a same user request.
//
// Returns:
//   - *TrackingIdResponse: The hash of the latest transaction broadcasted by the paymaster
//     corresponding to the requested ID and the status of the ID.
//   - error: An error if any
func (p *Paymaster) TrackingIdToLatestHash(ctx context.Context, trackingId *felt.Felt) (TrackingIdResponse, error) {
	var response TrackingIdResponse
	if err := p.c.CallContextWithSliceArgs(ctx, &response, "paymaster_trackingIdToLatestHash", trackingId); err != nil {
		return TrackingIdResponse{}, err
	}

	return response, nil
}
