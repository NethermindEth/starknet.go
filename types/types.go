package types

import (
	"math/big"
)

type NumAsHex string

type AddInvokeTransactionOutput struct {
	TransactionHash string `json:"transaction_hash"`
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

type DeployRequest struct {
	Type                string        `json:"type"`
	ContractAddressSalt string        `json:"contract_address_salt"`
	ConstructorCalldata []string      `json:"constructor_calldata"`
	ContractDefinition  ContractClass `json:"contract_definition"`
}

// FunctionCall function call information
type FunctionCall struct {
	ContractAddress    Hash   `json:"contract_address"`
	EntryPointSelector string `json:"entry_point_selector,omitempty"`

	// Calldata The parameters passed to the function
	Calldata []string `json:"calldata"`
}

type Signature []*big.Int

type FunctionInvoke struct {
	MaxFee *big.Int `json:"max_fee"`
	// Version of the transaction scheme, should be set to 0 or 1
	Version uint64 `json:"version"`
	// Signature
	Signature Signature `json:"signature"`
	// Nonce should only be set with Transaction V1
	Nonce *big.Int `json:"nonce,omitempty"`

	FunctionCall
}

type FeeEstimate struct {
	GasConsumed NumAsHex `json:"gas_consumed"`
	GasPrice    NumAsHex `json:"gas_price"`
	OverallFee  NumAsHex `json:"overall_fee"`
}

// ExecuteDetails provides some details about the execution.
type ExecuteDetails struct {
	MaxFee *big.Int
	Nonce  *big.Int
}
