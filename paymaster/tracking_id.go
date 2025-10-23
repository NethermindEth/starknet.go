package paymaster

import (
	"context"
	"fmt"
	"strconv"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
)

// TrackingIDToLatestHash gets the latest transaction hash and status for a given tracking ID.
// Returns a TrackingIdResponse.
//
// Parameters:
//   - ctx: The context.Context object for controlling the function call
//   - trackingID: A unique identifier used to track an execution request of a user.
//     This identitifier is returned by the paymaster after a successful call to `execute`.
//     Its purpose is to track the possibly different transaction hashes in the mempool which
//     are associated with a same user request.
//
// Returns:
//   - *TrackingIDResponse: The hash of the latest transaction broadcasted by the paymaster
//     corresponding to the requested ID and the status of the ID.
//   - error: An error if any
func (p *Paymaster) TrackingIDToLatestHash(
	ctx context.Context,
	trackingID *felt.Felt,
) (TrackingIDResponse, error) {
	var response TrackingIDResponse
	if err := p.c.CallContextWithSliceArgs(
		ctx,
		&response,
		"paymaster_trackingIdToLatestHash",
		trackingID,
	); err != nil {
		return TrackingIDResponse{}, rpcerr.UnwrapToRPCErr(err, ErrInvalidID)
	}

	return response, nil
}

// TrackingIDResponse is the response for the `paymaster_trackingIdToLatestHash` method.
type TrackingIDResponse struct {
	// The hash of the most recent tx sent by the paymaster and corresponding to the ID
	TransactionHash *felt.Felt `json:"transaction_hash"`
	// The status of the transaction associated with the ID
	Status TxnStatus `json:"status"`
}

// An enum representing the status of the transaction associated with a tracking ID
type TxnStatus int

const (
	// Indicates that the latest transaction associated with the ID is not yet
	// included in a block but is still being handled and monitored by the paymaster.
	// Represents the "active" string value.
	TxnStatusActive TxnStatus = iota + 1
	// Indicates that a transaction associated with the ID has been accepted on L2.
	// Represents the "accepted" string value.
	TxnStatusAccepted
	// Indicates that no transaction associated with the ID managed to enter a block
	// and the request has been dropped by the paymaster.
	// Represents the "dropped" string value.
	TxnStatusDropped
)

// String returns the string representation of the TxnStatus.
func (t TxnStatus) String() string {
	return []string{"active", "accepted", "dropped"}[t-1]
}

// MarshalJSON marshals the TxnStatus to JSON.
func (t TxnStatus) MarshalJSON() ([]byte, error) {
	return strconv.AppendQuote(nil, t.String()), nil
}

// UnmarshalJSON unmarshals the JSON data into a TxnStatus.
func (t *TxnStatus) UnmarshalJSON(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}

	switch s {
	case "active":
		*t = TxnStatusActive
	case "accepted":
		*t = TxnStatusAccepted
	case "dropped":
		*t = TxnStatusDropped
	default:
		return fmt.Errorf("invalid transaction status: %s", s)
	}

	return nil
}
