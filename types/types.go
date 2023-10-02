package types

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/NethermindEth/juno/core/felt"
)

type NumAsHex string

type AddInvokeTransactionOutput struct {
	TransactionHash *felt.Felt `json:"transaction_hash"`
}

type AddDeclareResponse struct {
	Code            string `json:"code"`
	TransactionHash string `json:"transaction_hash"`
	ClassHash       string `json:"class_hash"`
}

type AddDeployResponse struct {
	Code            string `json:"code"`
	TransactionHash string `json:"transaction_hash"`
	ContractAddress string `json:"address"`
}

// // TODO: remove
// type DeployRequest struct {
// 	Type                string               `json:"type"`
// 	ContractAddressSalt string               `json:"contract_address_salt"`
// 	ConstructorCalldata []string             `json:"constructor_calldata"`
// 	ContractDefinition  rpc.ContractClass `json:"contract_definition"`
// }

type DeployAccountRequest struct {
	MaxFee *big.Int `json:"max_fee"`
	// Version of the transaction scheme, should be set to 0 or 1
	Version *big.Int `json:"version"`
	// Signature
	Signature Signature `json:"signature"`
	// Nonce should only be set with Transaction V1
	Nonce *big.Int `json:"nonce,omitempty"`

	Type                string   `json:"type"`
	ContractAddressSalt string   `json:"contract_address_salt"`
	ConstructorCalldata []string `json:"constructor_calldata"`
	ClassHash           string   `json:"class_hash"`
}

// FunctionCall function call information
type FunctionCall struct {
	ContractAddress    *felt.Felt `json:"contract_address"`
	EntryPointSelector *felt.Felt `json:"entry_point_selector,omitempty"`

	// Calldata The parameters passed to the function
	Calldata []*felt.Felt `json:"calldata"`
}

type Signature []*big.Int

// todo(): what is this used for?
type FunctionInvoke struct {
	MaxFee *big.Int `json:"max_fee"`
	// Version of the transaction scheme, should be set to 0 or 1
	Version *big.Int `json:"version"`
	// Signature
	Signature Signature `json:"signature"`
	// Nonce should only be set with Transaction V1
	Nonce *big.Int `json:"nonce,omitempty"`
	// Defines the transaction type to invoke
	Type string `json:"type,omitempty"`

	SenderAddress      *felt.Felt `json:"sender_address"`
	EntryPointSelector string     `json:"entry_point_selector,omitempty"`

	// Calldata The parameters passed to the function
	Calldata []string `json:"calldata"`
}

type FeeEstimate struct {
	// GasConsumed the Ethereum gas cost of the transaction (see https://docs.starknet.io/docs/Fees/fee-mechanism for more info)
	GasConsumed NumAsHex `json:"gas_consumed"`

	// GasPrice the gas price (in gwei) that was used in the cost estimation
	GasPrice NumAsHex `json:"gas_price"`

	// OverallFee the estimated fee for the transaction (in gwei), product of gas_consumed and gas_price
	OverallFee NumAsHex `json:"overall_fee"`
}

// ExecuteDetails provides some details about the execution.
type ExecuteDetails struct {
	MaxFee *big.Int
	Nonce  *big.Int
}

type TransactionState string

const (
	TransactionAcceptedOnL1 TransactionState = "ACCEPTED_ON_L1"
	TransactionAcceptedOnL2 TransactionState = "ACCEPTED_ON_L2"
	TransactionNotReceived  TransactionState = "NOT_RECEIVED"
	TransactionPending      TransactionState = "PENDING"
	TransactionReceived     TransactionState = "RECEIVED"
	TransactionRejected     TransactionState = "REJECTED"
)

// UnmarshalJSON unmarshals the JSON data into a TransactionState object.
//
// The function takes a byte slice `data` as its parameter, which represents the JSON data to be unmarshaled.
// It returns an error if there is an issue unmarshaling the data.
func (ts *TransactionState) UnmarshalJSON(data []byte) error {
	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	switch unquoted {
	case "ACCEPTED_ON_L2":
		*ts = TransactionAcceptedOnL2
	case "ACCEPTED_ON_L1":
		*ts = TransactionAcceptedOnL1
	case "NOT_RECEIVED":
		*ts = TransactionNotReceived
	case "PENDING":
		*ts = TransactionPending
	case "RECEIVED":
		*ts = TransactionReceived
	case "REJECTED":
		*ts = TransactionRejected
	default:
		return fmt.Errorf("unsupported status: %s", data)
	}
	return nil
}

// MarshalJSON returns the JSON encoding of the TransactionState.
//
// It returns a byte slice containing the JSON-encoded string representation of the TransactionState and
// an error if there was any error during the encoding process.
func (ts TransactionState) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(string(ts))), nil
}

// String returns the string representation of the TransactionState.
//
// It does not take any parameters.
// It returns a string.
func (s TransactionState) String() string {
	return string(s)
}

// IsTransactionFinal checks if the transaction state is final.
//
// This function takes no parameters.
// It returns a boolean value indicating whether the transaction state is final or not.
func (s TransactionState) IsTransactionFinal() bool {
	if s == TransactionAcceptedOnL2 ||
		s == TransactionAcceptedOnL1 ||
		s == TransactionRejected {
		return true
	}
	return false
}
