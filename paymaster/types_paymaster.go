package paymaster

import (
	"encoding/json"
	"fmt"

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
func (t TxnStatus) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	switch s {
	case "active":
		t = TxnActive
	case "accepted":
		t = TxnAccepted
	case "dropped":
		t = TxnDropped
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

// An enum representing the type of the transaction to be executed by the paymaster
type UserTxnType string

const (
	// Represents a deploy transaction
	UserTxnDeploy UserTxnType = "deploy"
	// Represents an invoke transaction
	UserTxnInvoke UserTxnType = "invoke"
	// Represents a deploy and invoke transaction
	UserTxnDeployAndInvoke UserTxnType = "deploy_and_invoke"
)

// MarshalJSON marshals the UserTxnType to JSON.
func (u UserTxnType) MarshalJSON() ([]byte, error) {
	switch u {
	case UserTxnDeploy, UserTxnInvoke, UserTxnDeployAndInvoke:
		return json.Marshal(string(u))
	}

	return nil, fmt.Errorf("invalid user transaction type: %s", u)
}

// UnmarshalJSON unmarshals the JSON data into a UserTxnType.
func (u UserTxnType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch s {
	case "deploy":
		u = UserTxnDeploy
	case "invoke":
		u = UserTxnInvoke
	case "deploy_and_invoke":
		u = UserTxnDeployAndInvoke
	default:
		return fmt.Errorf("invalid user transaction type: %s", s)
	}
	return nil
}

// UserTransaction represents a user transaction (deploy, invoke, or deploy_and_invoke).
type UserTransaction struct {
	Type   UserTxnType `json:"type"`
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
	// The transaction to be executed by the paymaster
	Transaction UserTransaction `json:"transaction"`
	// Execution parameters to be used when executing the transaction
	Parameters UserParameters `json:"parameters"`
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
