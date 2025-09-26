package paymaster

import (
	"encoding/json"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
)

// aliases to facilitate usage

type (
	RPCError      = rpcerr.RPCError
	StringErrData = rpcerr.StringErrData
)

// Paymaster-specific errors based on SNIP-29 specification
//
//nolint:exhaustruct
var (
	ErrInvalidAddress = &RPCError{
		Code:    150,
		Message: "An error occurred (INVALID_ADDRESS)",
	}

	ErrTokenNotSupported = &RPCError{
		Code:    151,
		Message: "An error occurred (TOKEN_NOT_SUPPORTED)",
	}

	ErrInvalidSignature = &RPCError{
		Code:    153,
		Message: "An error occurred (INVALID_SIGNATURE)",
	}

	ErrMaxAmountTooLow = &RPCError{
		Code:    154,
		Message: "An error occurred (MAX_AMOUNT_TOO_LOW)",
	}

	ErrClassHashNotSupported = &RPCError{
		Code:    155,
		Message: "An error occurred (CLASS_HASH_NOT_SUPPORTED)",
	}

	ErrTransactionExecutionError = &RPCError{
		Code:    156,
		Message: "An error occurred (TRANSACTION_EXECUTION_ERROR)",
		Data:    &TxnExecutionErrData{},
	}

	ErrInvalidTimeBounds = &RPCError{
		Code:    157,
		Message: "An error occurred (INVALID_TIME_BOUNDS)",
	}

	ErrInvalidDeploymentData = &RPCError{
		Code:    158,
		Message: "An error occurred (INVALID_DEPLOYMENT_DATA)",
	}

	ErrInvalidClassHash = &RPCError{
		Code:    159,
		Message: "An error occurred (INVALID_ADDRESS)",
	}

	ErrInvalidID = &RPCError{
		Code:    160,
		Message: "An error occurred (INVALID_ID)",
	}

	ErrUnknownError = &RPCError{
		Code:    163,
		Message: "An error occurred (UNKNOWN_ERROR)",
		Data:    StringErrData(""),
	}
)

// TxnExecutionErrData represents the structured data for TRANSACTION_EXECUTION_ERROR
type TxnExecutionErrData struct {
	ExecutionError string `json:"execution_error"`
}

// ErrorMessage implements the RPCData interface
func (t TxnExecutionErrData) ErrorMessage() string {
	return t.ExecutionError
}

// Object containing data about the token: contract address, number of decimals and current price in STRK
type TokenData struct {
	// Token contract address
	TokenAddress *felt.Felt `json:"token_address"`
	// The number of decimals of the token
	Decimals uint8 `json:"decimals"`
	// Price in STRK (in FRI units)
	PriceInStrk string `json:"price_in_strk"` // u256 as a hex string
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

// TrackingIdResponse is the response for the `paymaster_trackingIdToLatestHash` method.
type TrackingIdResponse struct {
	// The hash of the most recent tx sent by the paymaster and corresponding to the ID
	TransactionHash *felt.Felt `json:"transaction_hash"`
	// The status of the transaction associated with the ID
	Status TxnStatus `json:"status"`
}
