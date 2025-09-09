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
	Transaction *UserTransaction `json:"transaction"`
	// Execution parameters to be used when executing the transaction
	Parameters *UserParameters `json:"parameters"`
}

// UserTransaction represents a user transaction (deploy, invoke, or deploy_and_invoke).
type UserTransaction struct {
	// The type of the transaction to be executed by the paymaster
	Type UserTxnType `json:"type"`
	// The deployment data for the transaction, used for `deploy` and `deploy_and_invoke` transaction types.
	// Should be omitted for `invoke` transaction types.
	Deployment *AccDeploymentData `json:"deployment,omitempty"`
	// The invoke data for the transaction, used for `invoke` and `deploy_and_invoke` transaction types.
	// Should be omitted for `deploy` transaction types.
	Invoke *UserInvoke `json:"invoke,omitempty"`
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
func (u *UserTxnType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch s {
	case "deploy":
		*u = UserTxnDeploy
	case "invoke":
		*u = UserTxnInvoke
	case "deploy_and_invoke":
		*u = UserTxnDeployAndInvoke
	default:
		return fmt.Errorf("invalid user transaction type: %s", s)
	}

	return nil
}

// Data required to deploy an account at an address.
type AccDeploymentData struct {
	// The expected address to be deployed, used to double check
	Address *felt.Felt `json:"address"`
	// The hash of the deployed contract's class
	ClassHash *felt.Felt `json:"class_hash"`
	// The salt used for the contract address calculation
	Salt *felt.Felt `json:"salt"`
	// The parameters passed to the constructor
	ConstructorCalldata []*felt.Felt `json:"calldata"`
	// Optional array of felts to be added to the signature
	SignatureData []*felt.Felt `json:"sigdata,omitempty"`
	// The Cairo version (CairoZero is not supported)
	Version uint8 `json:"version"`
}

// Execution parameters to be used when executing the transaction through the paymaster
type UserParameters struct {
	// Version of the execution parameters which is not tied to the 'execute from outside' version.
	Version UserParamVersion `json:"version"`
	// Fee mode to use for the execution
	FeeMode FeeMode `json:"fee_mode"`
	// Optional. Time constraint on the execution
	TimeBounds *TimeBounds `json:"time_bounds,omitempty"`
}

// An enum representing the version of the execution parameters
type UserParamVersion string

const (
	// Represents the v1 of the execution parameters
	UserParamV1 UserParamVersion = "0x1"
)

// MarshalJSON marshals the UserParamVersion to JSON.
func (v UserParamVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(v))
}

// UnmarshalJSON unmarshals the JSON data into a UserParamVersion.
func (v *UserParamVersion) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch s {
	case "0x1":
		*v = UserParamV1
	default:
		return fmt.Errorf("invalid user parameter version: %s", s)
	}

	return nil
}

// An enum representing the fee mode to use for the transaction
type FeeModeType string

const (
	// Specify that the transaction should be sponsored. This argument does not
	// guaranteed sponsorship and will depend on the paymaster provider
	FeeModeSponsored FeeModeType = "sponsored"
	// Default fee mode where the transaction is paid by the user in the given gas token
	FeeModeDefault FeeModeType = "default"
	// Fee mode where the transaction is paid by the user in the given gas token and
	// the user can specify a tip
	FeeModePriority FeeModeType = "priority"
)

// MarshalJSON marshals the FeeModeType to JSON.
func (feeMode FeeModeType) MarshalJSON() ([]byte, error) {
	switch feeMode {
	case FeeModeSponsored, FeeModeDefault, FeeModePriority:
		return json.Marshal(string(feeMode))
	}

	return nil, fmt.Errorf("invalid fee mode: %s", feeMode)
}

// UnmarshalJSON unmarshals the JSON data into a FeeModeType.
func (feeMode *FeeModeType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	switch s {
	case "sponsored":
		*feeMode = FeeModeSponsored
	case "default":
		*feeMode = FeeModeDefault
	case "priority":
		*feeMode = FeeModePriority
	default:
		return fmt.Errorf("invalid fee mode: %s", s)
	}

	return nil
}

// Specify how the transaction should be paid. Either by the user specifying a gas token or through sponsorship
type FeeMode struct {
	// The fee mode type to use for the transaction
	Mode FeeModeType `json:"mode"`
	// The gas token to use for the transaction. Should be omitted for `sponsored` fee mode
	GasToken *felt.Felt `json:"gas_token,omitempty"`
	// The tip to use for the transaction. Only used for `priority` fee mode
	TipInStrk *felt.Felt `json:"tip_in_strk,omitempty"`
}

// Object containing timestamps corresponding to `Execute After` and `Execute Before`
type TimeBounds struct {
	// A lower bound after which an outside call is valid in UNIX timestamp format
	ExecuteAfter string `json:"execute_after"`
	// An upper bound before which an outside call is valid in UNIX timestamp format
	ExecuteBefore string `json:"execute_before"`
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
