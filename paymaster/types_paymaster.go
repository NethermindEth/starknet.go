package paymaster

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
)

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

// TokenData contains information about a supported token (address, decimals, price in STRK).
type TokenData struct {
	TokenAddress *felt.Felt `json:"token_address"`
	Decimals     int        `json:"decimals"`
	PriceInStrk  string     `json:"price_in_strk"` // u256 as string
}

// TrackingIdResponse is the response for tracking a transaction by ID (latest tx hash and status).
type TrackingIdResponse struct {
	TransactionHash *felt.Felt `json:"transaction_hash"`
	Status          string     `json:"status"` // "active", "accepted", "dropped"
}

// Call represents a single contract call (to, selector, calldata).
type Call struct {
	To       *felt.Felt   `json:"to"`
	Selector *felt.Felt   `json:"selector"`
	Calldata []*felt.Felt `json:"calldata"`
}

// UserInvoke represents an invoke transaction from a user (user address and calls).
type UserInvoke struct {
	UserAddress *felt.Felt `json:"user_address"`
	Calls       []Call     `json:"calls"`
}

// UserTransaction represents a user transaction (deploy, invoke, or deploy_and_invoke).
type UserTransaction struct {
	Type   string      `json:"type"` // "deploy", "invoke", "deploy_and_invoke"
	Deploy interface{} `json:"deployment,omitempty"`
	Invoke UserInvoke  `json:"invoke,omitempty"`
}

// FeeMode specifies how the transaction fee should be paid (mode, gas token, tip).
type FeeMode struct {
	Mode      string     `json:"mode"` // "sponsored", "default", "priority"
	GasToken  *felt.Felt `json:"gas_token,omitempty"`
	TipInStrk *felt.Felt `json:"tip_in_strk,omitempty"`
}

// UserParameters are execution parameters for the transaction (version, fee mode, time bounds).
type UserParameters struct {
	Version    string      `json:"version"` // "0x1"
	FeeMode    FeeMode     `json:"fee_mode"`
	TimeBounds interface{} `json:"time_bounds,omitempty"`
}

// BuildTransactionRequest is the request to build a transaction for the paymaster (transaction + parameters).
type BuildTransactionRequest struct {
	Transaction UserTransaction `json:"transaction"`
	Parameters  UserParameters  `json:"parameters"`
}

// FeeEstimateResponse is a detailed fee estimation (in STRK and gas token, with suggested max).
type FeeEstimateResponse struct {
	GasTokenPriceInStrk       *felt.Felt `json:"gas_token_price_in_strk"`
	EstimatedFeeInStrk        *felt.Felt `json:"estimated_fee_in_strk"`
	EstimatedFeeInGasToken    *felt.Felt `json:"estimated_fee_in_gas_token"`
	SuggestedMaxFeeInStrk     *felt.Felt `json:"suggested_max_fee_in_strk"`
	SuggestedMaxFeeInGasToken *felt.Felt `json:"suggested_max_fee_in_gas_token"`
}

// BuildTransactionResponse is the response from building a transaction (typed data, fee, parameters, etc.).
type BuildTransactionResponse struct {
	Type       string              `json:"type"` // "deploy", "invoke", "deploy_and_invoke"
	Deployment interface{}         `json:"deployment,omitempty"`
	TypedData  interface{}         `json:"typed_data,omitempty"`
	Parameters UserParameters      `json:"parameters"`
	Fee        FeeEstimateResponse `json:"fee"`
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

// PaymasterProvider is the interface for paymaster operations (SNIP-X methods only).
//
// Methods:
//   - IsAvailable: Check if the paymaster service is available.
//   - GetSupportedTokens: List supported tokens.
//   - TrackingIdToLatestHash: Get latest tx hash and status for a tracking ID.
//   - BuildTransaction: Build a transaction (get typed data and fee).
//   - ExecuteTransaction: Execute a signed transaction.
//
//go:generate go run go.uber.org/mock/mockgen -destination=../mocks/mock_paymaster_provider.go -package=mocks -source=types_paymaster.go PaymasterProvider
type PaymasterProvider interface {
	IsAvailable(ctx context.Context) (bool, error)
	GetSupportedTokens(ctx context.Context) ([]TokenData, error)
	TrackingIdToLatestHash(ctx context.Context, trackingId *felt.Felt) (*TrackingIdResponse, error)
	BuildTransaction(ctx context.Context, request BuildTransactionRequest) (*BuildTransactionResponse, error)
	ExecuteTransaction(ctx context.Context, request ExecuteTransactionRequest) (*ExecuteTransactionResponse, error)
}
