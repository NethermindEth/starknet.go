package rpc

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
)

// Transaction is an interface that represents a Starknet transaction.
// It is used to provide a common interface for all transaction types.
// The 'Type' and 'Version' fields are present in all transaction types.
type Transaction interface {
	GetType() TransactionType
	GetVersion() TransactionVersion
}

type tempTransaction struct {
	// Common fields for all transactions of all versions
	Type      TransactionType    `json:"type"`
	Signature []felt.Felt        `json:"signature"`
	Version   TransactionVersion `json:"version"`
	Nonce     *felt.Felt         `json:"nonce"`   // Only in v1 onwards
	MaxFee    *felt.Felt         `json:"max_fee"` // Only before v3

	// Common fields for all v3 transactions
	FeeMode        DataAvailabilityMode   `json:"fee_data_availability_mode"`
	NonceDataMode  DataAvailabilityMode   `json:"nonce_data_availability_mode"`
	PayMasterData  []felt.Felt            `json:"paymaster_data"`
	ResourceBounds *ResourceBoundsMapping `json:"resource_bounds"`
	Tip            U64                    `json:"tip"`

	// Common fields for Invoke transactions.
	// Also present in the L1Handler transaction type.
	Calldata           []felt.Felt `json:"calldata"`             // In all versions
	ContractAddress    *felt.Felt  `json:"contract_address"`     // Only in v0
	EntryPointSelector *felt.Felt  `json:"entry_point_selector"` // Only in v0

	// Common fields for Invoke and Declare transactions
	SenderAddress         *felt.Felt  `json:"sender_address"`          // v1 onwards
	AccountDeploymentData []felt.Felt `json:"account_deployment_data"` // Only in v3

	// Common field for Declare transactions
	CompiledClassHash *felt.Felt `json:"compiled_class_hash"` // v2 onwards

	// Common field for all Declare and DeployAccount transaction versions.
	// Also present in the deprecated Deploy transaction type.
	ClassHash *felt.Felt `json:"class_hash"` // In all versions

	// Common fields for all DeployAccount transaction versions.
	// Also present in the deprecated Deploy transaction type.
	ContractAddressSalt *felt.Felt  `json:"contract_address_salt"`
	ConstructorCalldata []felt.Felt `json:"constructor_calldata"`
}

var (
	_ Transaction = InvokeTxnV0{}
	_ Transaction = InvokeTxnV1{}
	_ Transaction = InvokeTxnV3{}
	_ Transaction = DeclareTxnV1{}
	_ Transaction = DeclareTxnV2{}
	_ Transaction = DeclareTxnV3{}
	_ Transaction = DeployTxn{}
	_ Transaction = DeployAccountTxnV1{}
	_ Transaction = DeployAccountTxnV3{}
	_ Transaction = L1HandlerTxn{}
)

// unmarshalTxn unmarshals a given txn as a byte slice and returns a concrete
// transaction type wrapped in the Transaction interface.
//
// Parameters:
//   - data: The transaction to be unmarshaled
//
// Returns:
//   - Transaction: a concrete transaction type wrapped in the Transaction interface
//   - error: an error if the unmarshaling process fails
//
//nolint:gocyclo // Inevitable due to many switch cases
func unmarshalTxn(data []byte) (Transaction, error) {
	var txnAsMap map[string]interface{}
	if err := json.Unmarshal(data, &txnAsMap); err != nil {
		return nil, err
	}

	switch TransactionType(txnAsMap["type"].(string)) {
	case TransactionTypeDeclare:
		switch TransactionVersion(txnAsMap["version"].(string)) {
		case TransactionV0:
			return unmarshalTxnToType[DeclareTxnV0](data)
		case TransactionV1:
			return unmarshalTxnToType[DeclareTxnV1](data)
		case TransactionV2:
			return unmarshalTxnToType[DeclareTxnV2](data)
		case TransactionV3:
			return unmarshalTxnToType[DeclareTxnV3](data)
		default:
			return nil, errors.New(
				"internal error with Declare transaction version and unmarshalTxn()",
			)
		}
	case TransactionTypeDeploy:
		return unmarshalTxnToType[DeployTxn](data)
	case TransactionTypeDeployAccount:
		switch TransactionVersion(txnAsMap["version"].(string)) {
		case TransactionV1:
			return unmarshalTxnToType[DeployAccountTxnV1](data)
		case TransactionV3:
			return unmarshalTxnToType[DeployAccountTxnV3](data)
		}
	case TransactionTypeInvoke:
		switch TransactionVersion(txnAsMap["version"].(string)) {
		case TransactionV0:
			return unmarshalTxnToType[InvokeTxnV0](data)
		case TransactionV1:
			return unmarshalTxnToType[InvokeTxnV1](data)
		case TransactionV3:
			return unmarshalTxnToType[InvokeTxnV3](data)
		}
	case TransactionTypeL1Handler:
		return unmarshalTxnToType[L1HandlerTxn](data)
	}

	return nil, fmt.Errorf("unknown transaction type: %v", txnAsMap["type"])
}

// unmarshalTxnToType is a generic function that takes in a byte slice 'data',
// unmarshals it to a concrete transaction of type T, and returns the concrete
// transaction wrapped in the Transaction interface.
func unmarshalTxnToType[T Transaction](data []byte) (T, error) {
	var resp T

	if err := json.Unmarshal(data, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}

// Invoke transactions
func (tx InvokeTxnV0) GetType() TransactionType {
	return tx.Type
}

func (tx InvokeTxnV0) GetVersion() TransactionVersion {
	return tx.Version
}

func (tx InvokeTxnV1) GetType() TransactionType {
	return tx.Type
}

func (tx InvokeTxnV1) GetVersion() TransactionVersion {
	return tx.Version
}

func (tx InvokeTxnV3) GetType() TransactionType {
	return tx.Type
}

func (tx InvokeTxnV3) GetVersion() TransactionVersion {
	return tx.Version
}

// Declare transactions
func (tx DeclareTxnV0) GetType() TransactionType {
	return tx.Type
}

func (tx DeclareTxnV0) GetVersion() TransactionVersion {
	return tx.Version
}

func (tx DeclareTxnV1) GetType() TransactionType {
	return tx.Type
}

func (tx DeclareTxnV1) GetVersion() TransactionVersion {
	return tx.Version
}

func (tx DeclareTxnV2) GetType() TransactionType {
	return tx.Type
}

func (tx DeclareTxnV2) GetVersion() TransactionVersion {
	return tx.Version
}

func (tx DeclareTxnV3) GetType() TransactionType {
	return tx.Type
}

func (tx DeclareTxnV3) GetVersion() TransactionVersion {
	return tx.Version
}

func (tx BroadcastDeclareTxnV3) GetType() TransactionType {
	return tx.Type
}

func (tx BroadcastDeclareTxnV3) GetVersion() TransactionVersion {
	return tx.Version
}

// Deploy transaction
func (tx DeployTxn) GetType() TransactionType {
	return tx.Type
}

func (tx DeployTxn) GetVersion() TransactionVersion {
	return tx.Version
}

// DeployAccount transactions
func (tx DeployAccountTxnV1) GetType() TransactionType {
	return tx.Type
}

func (tx DeployAccountTxnV1) GetVersion() TransactionVersion {
	return tx.Version
}

func (tx DeployAccountTxnV3) GetType() TransactionType {
	return tx.Type
}

func (tx DeployAccountTxnV3) GetVersion() TransactionVersion {
	return tx.Version
}

// L1Handler transaction
func (tx L1HandlerTxn) GetType() TransactionType {
	return tx.Type
}

func (tx L1HandlerTxn) GetVersion() TransactionVersion {
	return tx.Version
}

// InvokeTxnType is an interface that represents a Starknet invoke transaction.
// It is used to provide a common interface for all invoke transaction types.
// The 'Calldata' field is present in all invoke transaction types.
type InvokeTxnType interface {
	GetCalldata() []*felt.Felt
}

func (tx InvokeTxnV0) GetCalldata() []*felt.Felt {
	return tx.Calldata
}

func (tx InvokeTxnV1) GetCalldata() []*felt.Felt {
	return tx.Calldata
}

func (tx InvokeTxnV3) GetCalldata() []*felt.Felt {
	return tx.Calldata
}

var (
	_ InvokeTxnType = InvokeTxnV0{}
	_ InvokeTxnType = InvokeTxnV1{}
	_ InvokeTxnType = InvokeTxnV3{}
	_ InvokeTxnType = BroadcastInvokeTxnV3{}
)

// DeclareTxnType is an interface that represents a Starknet declare transaction.
// It is used to provide a common interface for all declare transaction types.
// The 'SenderAddress' field is present in all declare transaction types.
type DeclareTxnType interface {
	GetSenderAddress() *felt.Felt
}

func (tx DeclareTxnV0) GetSenderAddress() *felt.Felt {
	return tx.SenderAddress
}

func (tx DeclareTxnV1) GetSenderAddress() *felt.Felt {
	return tx.SenderAddress
}

func (tx DeclareTxnV2) GetSenderAddress() *felt.Felt {
	return tx.SenderAddress
}

func (tx DeclareTxnV3) GetSenderAddress() *felt.Felt {
	return tx.SenderAddress
}

func (tx BroadcastDeclareTxnV3) GetSenderAddress() *felt.Felt {
	return tx.SenderAddress
}

var (
	_ DeclareTxnType = DeclareTxnV0{}
	_ DeclareTxnType = DeclareTxnV1{}
	_ DeclareTxnType = DeclareTxnV2{}
	_ DeclareTxnType = DeclareTxnV3{}
	_ DeclareTxnType = BroadcastDeclareTxnV3{}
)

// DeployAccountType is an interface that represents a Starknet deploy account transaction.
// It is used to provide a common interface for all deploy account transaction types.
// The 'ConstructorCalldata' field is present in all deploy account transaction types.
type DeployAccountType interface {
	GetConstructorCalldata() []*felt.Felt
}

func (tx DeployAccountTxnV1) GetConstructorCalldata() []*felt.Felt {
	return tx.ConstructorCalldata
}

func (tx DeployAccountTxnV3) GetConstructorCalldata() []*felt.Felt {
	return tx.ConstructorCalldata
}

var (
	_ DeployAccountType = DeployAccountTxnV1{}
	_ DeployAccountType = DeployAccountTxnV3{}
	_ DeployAccountType = BroadcastDeployAccountTxnV3{}
)
