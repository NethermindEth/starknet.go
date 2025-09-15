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

// OutsideExecutionTypedData represents the EIP-712 typed data structure for outside execution (used for signing and validation).
type OutsideExecutionTypedData struct {
	Types       map[string][]TypedDataField `json:"types"`
	PrimaryType string                      `json:"primaryType"`
	Domain      TypedDataDomain             `json:"domain"`
	Message     interface{}                 `json:"message"`
}

// OutsideCallV1 represents a single contract call within a V1 outside execution message.
type OutsideCallV1 struct {
	To          *felt.Felt   `json:"to"`
	Selector    *felt.Felt   `json:"selector"`
	CalldataLen *felt.Felt   `json:"calldata_len"`
	Calldata    []*felt.Felt `json:"calldata"`
}

// OutsideExecutionMessageV1 is the message payload for a V1 outside execution.
type OutsideExecutionMessageV1 struct {
	Caller        *felt.Felt       `json:"caller"`
	Nonce         *felt.Felt       `json:"nonce"`
	ExecuteAfter  *felt.Felt       `json:"execute_after"`
	ExecuteBefore *felt.Felt       `json:"execute_before"`
	CallsLen      *felt.Felt       `json:"calls_len"`
	Calls         []*OutsideCallV1 `json:"calls"`
}

// OutsideExecutionMessageV2 is the message payload for a V2 outside execution.
type OutsideExecutionMessageV2 struct {
	Caller        *felt.Felt `json:"Caller"`
	Nonce         *felt.Felt `json:"Nonce"`
	ExecuteAfter  string     `json:"Execute After"`  // u128
	ExecuteBefore string     `json:"Execute Before"` // u128
	Calls         []Call     `json:"Calls"`
}

// OutsideExecutionMessageV3 is the message payload for a V3-rc outside execution.
// Note: The 'Fee' field is represented as an interface{} to accommodate different fee structures.
type OutsideExecutionMessageV3 struct {
	Caller        *felt.Felt  `json:"Caller"`
	Nonce         *felt.Felt  `json:"Nonce"`
	ExecuteAfter  string      `json:"Execute After"`  // u128
	ExecuteBefore string      `json:"Execute Before"` // u128
	Calls         []Call      `json:"Calls"`
	Fee           interface{} `json:"Fee"`
}

// TypedDataField describes a single field in a typed data struct (name and type).
type TypedDataField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// TypedDataDomain is the domain separator for EIP-712 typed data (name, version, chainId).
type TypedDataDomain struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	ChainID string `json:"chainId"`
}

// ===== SNIP-X TYPES (Official Specification) =====

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

// ExecutableUserInvoke is an invoke transaction ready for execution (user address, typed data, signature).
type ExecutableUserInvoke struct {
	UserAddress *felt.Felt   `json:"user_address"`
	TypedData   interface{}  `json:"typed_data"`
	Signature   []*felt.Felt `json:"signature"`
}

// ExecutableUserTransaction is a user transaction ready for execution (deploy, invoke, or both).
type ExecutableUserTransaction struct {
	Type   string               `json:"type"` // "deploy", "invoke", "deploy_and_invoke"
	Deploy interface{}          `json:"deployment,omitempty"`
	Invoke ExecutableUserInvoke `json:"invoke,omitempty"`
}

// ExecuteTransactionRequest is the request to execute a transaction via the paymaster (transaction + parameters).
type ExecuteTransactionRequest struct {
	Transaction ExecutableUserTransaction `json:"transaction"`
	Parameters  UserParameters            `json:"parameters"`
}

// ExecuteTransactionResponse is the response from executing a transaction (tracking ID and transaction hash).
type ExecuteTransactionResponse struct {
	TrackingId      *felt.Felt `json:"tracking_id"`
	TransactionHash *felt.Felt `json:"transaction_hash"`
}
