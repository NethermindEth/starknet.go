package paymaster

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
)

// BuildTransaction receives the transaction the user wants to execute. Returns the typed
// data along with the estimated gas cost and the maximum gas cost suggested to ensure execution
//
// Parameters:
//   - ctx: The context.Context object for controlling the function call
//   - request: The BuildTransactionRequest containing the transaction and parameters
//
// Returns:
//   - *BuildTransactionResponse: The response containing typed data and fee estimate
//   - error: An error if the request fails
func (p *Paymaster) BuildTransaction(ctx context.Context, request *BuildTransactionRequest) (*BuildTransactionResponse, error) {
	var response BuildTransactionResponse
	if err := p.c.CallContextWithSliceArgs(ctx, &response, "paymaster_buildTransaction", request); err != nil {
		return nil, err
	}

	return &response, nil
}

// BuildTransactionRequest is the request to build a transaction for the paymaster (transaction + parameters).
type BuildTransactionRequest struct {
	// The transaction to be executed by the paymaster
	Transaction UserTransaction `json:"transaction"`
	// Execution parameters to be used when executing the transaction
	Parameters UserParameters `json:"parameters"`
}

// UserParameters are execution parameters for the transaction (version, fee mode, time bounds).
type UserParameters struct {
	Version    string      `json:"version"` // "0x1"
	FeeMode    FeeMode     `json:"fee_mode"`
	TimeBounds interface{} `json:"time_bounds,omitempty"`
}

// UserTransaction represents a user transaction (deploy, invoke, or deploy_and_invoke).
type UserTransaction struct {
	Type   UserTxnType `json:"type"`
	Deploy interface{} `json:"deployment,omitempty"`
	Invoke UserInvoke  `json:"invoke,omitempty"`
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
