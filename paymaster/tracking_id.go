package paymaster

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
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
func (p *Paymaster) TrackingIdToLatestHash(
	ctx context.Context,
	trackingId *felt.Felt,
) (TrackingIdResponse, error) {
	var response TrackingIdResponse
	if err := p.c.CallContextWithSliceArgs(ctx, &response, "paymaster_trackingIdToLatestHash", trackingId); err != nil {
		return TrackingIdResponse{}, rpcerr.UnwrapToRPCErr(err, ErrInvalidID)
	}

	return response, nil
}

// TrackingIdResponse is the response for the `paymaster_trackingIdToLatestHash` method.
type TrackingIdResponse struct {
	// The hash of the most recent tx sent by the paymaster and corresponding to the ID
	TransactionHash *felt.Felt `json:"transaction_hash"`
	// The status of the transaction associated with the ID
	Status TxnStatus `json:"status"`
}

// An enum representing the status of the transaction associated with a tracking ID
type TxnStatus string

const (
	// Indicates that the latest transaction associated with the ID is not yet
	// included in a block but is still being handled and monitored by the paymaster
	TxnActive TxnStatus = "active"
	// Indicates that a transaction associated with the ID has been accepted on L2
	TxnAccepted TxnStatus = "accepted"
	// Indicates that no transaction associated with the ID managed to enter a block
	// and the request has been dropped by the paymaster
	TxnDropped TxnStatus = "dropped"
)

// MarshalJSON marshals the TxnStatus to JSON.
func (t TxnStatus) MarshalJSON() ([]byte, error) {
	switch t {
	case TxnActive, TxnAccepted, TxnDropped:
		return json.Marshal(string(t))
	}

	return nil, fmt.Errorf("invalid transaction status: %s", t)
}

// UnmarshalJSON unmarshals the JSON data into a TxnStatus.
func (t *TxnStatus) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	switch s {
	case "active":
		*t = TxnActive
	case "accepted":
		*t = TxnAccepted
	case "dropped":
		*t = TxnDropped
	default:
		return fmt.Errorf("invalid transaction status: %s", s)
	}

	return nil
}
