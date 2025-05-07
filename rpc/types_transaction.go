package rpc

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
)

type BlockTransaction struct {
	Hash *felt.Felt `json:"transaction_hash"`
	Transaction
}

type InvokeTxnV0 struct {
	Type      TransactionType    `json:"type"`
	MaxFee    *felt.Felt         `json:"max_fee"`
	Version   TransactionVersion `json:"version"`
	Signature []*felt.Felt       `json:"signature"`
	FunctionCall
}

type InvokeTxnV1 struct {
	MaxFee        *felt.Felt         `json:"max_fee"`
	Version       TransactionVersion `json:"version"`
	Signature     []*felt.Felt       `json:"signature"`
	Nonce         *felt.Felt         `json:"nonce"`
	Type          TransactionType    `json:"type"`
	SenderAddress *felt.Felt         `json:"sender_address"`
	// The data expected by the account's `execute` function (in most usecases, this includes the called contract address and a function selector)
	Calldata []*felt.Felt `json:"calldata"`
}

type InvokeTxnV3 struct {
	Type           TransactionType        `json:"type"`
	SenderAddress  *felt.Felt             `json:"sender_address"`
	Calldata       []*felt.Felt           `json:"calldata"`
	Version        TransactionVersion     `json:"version"`
	Signature      []*felt.Felt           `json:"signature"`
	Nonce          *felt.Felt             `json:"nonce"`
	ResourceBounds *ResourceBoundsMapping `json:"resource_bounds"`
	Tip            U64                    `json:"tip"`
	// The data needed to allow the paymaster to pay for the transaction in native tokens
	PayMasterData []*felt.Felt `json:"paymaster_data"`
	// The data needed to deploy the account contract from which this tx will be initiated
	AccountDeploymentData []*felt.Felt `json:"account_deployment_data"`
	// The storage domain of the account's nonce (an account has a nonce per DA mode)
	NonceDataMode DataAvailabilityMode `json:"nonce_data_availability_mode"`
	// The storage domain of the account's balance from which fee will be charged
	FeeMode DataAvailabilityMode `json:"fee_data_availability_mode"`
}

type L1HandlerTxn struct {
	Type TransactionType `json:"type"`
	// Version of the transaction scheme
	Version TransactionVersion `json:"version"`
	// Nonce
	Nonce string `json:"nonce"`
	FunctionCall
}

type DeclareTxnV0 struct {
	Type TransactionType `json:"type"`
	// SenderAddress the address of the account contract sending the declaration transaction
	SenderAddress *felt.Felt         `json:"sender_address"`
	MaxFee        *felt.Felt         `json:"max_fee"`
	Version       TransactionVersion `json:"version"`
	Signature     []*felt.Felt       `json:"signature"`
	ClassHash     *felt.Felt         `json:"class_hash"`
}

type DeclareTxnV1 struct {
	Type TransactionType `json:"type"`
	// SenderAddress the address of the account contract sending the declaration transaction
	SenderAddress *felt.Felt         `json:"sender_address"`
	MaxFee        *felt.Felt         `json:"max_fee"`
	Version       TransactionVersion `json:"version"`
	Signature     []*felt.Felt       `json:"signature"`
	Nonce         *felt.Felt         `json:"nonce"`
	// ClassHash the hash of the declared class
	ClassHash *felt.Felt `json:"class_hash"`
}

type DeclareTxnV2 struct {
	Type TransactionType `json:"type"`
	// SenderAddress the address of the account contract sending the declaration transaction
	SenderAddress     *felt.Felt         `json:"sender_address"`
	CompiledClassHash *felt.Felt         `json:"compiled_class_hash"`
	MaxFee            *felt.Felt         `json:"max_fee"`
	Version           TransactionVersion `json:"version"`
	Signature         []*felt.Felt       `json:"signature"`
	Nonce             *felt.Felt         `json:"nonce"`
	ClassHash         *felt.Felt         `json:"class_hash"`
}

type DeclareTxnV3 struct {
	Type              TransactionType        `json:"type"`
	SenderAddress     *felt.Felt             `json:"sender_address"`
	CompiledClassHash *felt.Felt             `json:"compiled_class_hash"`
	Version           TransactionVersion     `json:"version"`
	Signature         []*felt.Felt           `json:"signature"`
	Nonce             *felt.Felt             `json:"nonce"`
	ClassHash         *felt.Felt             `json:"class_hash"`
	ResourceBounds    *ResourceBoundsMapping `json:"resource_bounds"`
	Tip               U64                    `json:"tip"`
	// The data needed to allow the paymaster to pay for the transaction in native tokens
	PayMasterData []*felt.Felt `json:"paymaster_data"`
	// The data needed to deploy the account contract from which this tx will be initiated
	AccountDeploymentData []*felt.Felt `json:"account_deployment_data"`
	// The storage domain of the account's nonce (an account has a nonce per DA mode)
	NonceDataMode DataAvailabilityMode `json:"nonce_data_availability_mode"`
	// The storage domain of the account's balance from which fee will be charged
	FeeMode DataAvailabilityMode `json:"fee_data_availability_mode"`
}

type ResourceBoundsMapping struct {
	// The max amount and max price per unit of L1 gas used in this tx
	L1Gas ResourceBounds `json:"l1_gas"`
	// The max amount and max price per unit of L1 blob gas used in this tx
	L1DataGas ResourceBounds `json:"l1_data_gas"`
	// The max amount and max price per unit of L2 gas used in this tx
	L2Gas ResourceBounds `json:"l2_gas"`
}

// DA_MODE: Specifies a storage domain in Starknet. Each domain has different guarantees regarding availability
type DataAvailabilityMode string

const (
	DAModeL1 DataAvailabilityMode = "L1"
	DAModeL2 DataAvailabilityMode = "L2"
)

// MarshalJSON implements the json.Marshaler interface.
// It validates that the DataAvailabilityMode is either L1 or L2 before marshalling.
func (da DataAvailabilityMode) MarshalJSON() ([]byte, error) {
	switch da {
	case DAModeL1, DAModeL2:
		return json.Marshal(string(da))
	default:
		return nil, fmt.Errorf("invalid DataAvailabilityMode: %s, must be either L1 or L2", string(da))
	}
}

func (da *DataAvailabilityMode) UInt64() (uint64, error) {
	switch *da {
	case DAModeL1:
		return uint64(0), nil
	case DAModeL2:
		return uint64(1), nil
	}

	return 0, errors.New("unknown DAMode")
}

type Resource string

// Values used in the Resource Bounds hash calculation
// Ref: https://docs.starknet.io/architecture-and-concepts/network-architecture/transactions/#v3_hash_calculation
const (
	ResourceL1Gas     Resource = "L1_GAS"
	ResourceL2Gas     Resource = "L2_GAS"
	ResourceL1DataGas Resource = "L1_DATA"
)

type ResourceBounds struct {
	// The max amount of the resource that can be used in the tx
	MaxAmount U64 `json:"max_amount"`
	// The max price per unit of this resource for this tx
	MaxPricePerUnit U128 `json:"max_price_per_unit"`
}

func (rb ResourceBounds) Bytes(resource Resource) ([]byte, error) {
	const eight = 8
	maxAmountBytes := make([]byte, eight)
	maxAmountUint64, err := rb.MaxAmount.ToUint64()
	if err != nil {
		return nil, err
	}
	binary.BigEndian.PutUint64(maxAmountBytes, maxAmountUint64)
	maxPricePerUnitFelt, err := new(felt.Felt).SetString(string(rb.MaxPricePerUnit))
	if err != nil {
		return nil, err
	}
	maxPriceBytes := maxPricePerUnitFelt.Bytes()

	return internalUtils.Flatten(
		[]byte{0},
		[]byte(resource),
		maxAmountBytes,
		maxPriceBytes[16:], // uint128.
	), nil
}

// DeployTxn The structure of a deploy transaction. Note that this transaction type is deprecated and will no longer be supported in future versions
type DeployTxn struct {
	// ClassHash The hash of the deployed contract's class
	ClassHash *felt.Felt `json:"class_hash"`

	Version             TransactionVersion `json:"version"`
	Type                TransactionType    `json:"type"`
	ContractAddressSalt *felt.Felt         `json:"contract_address_salt"`
	ConstructorCalldata []*felt.Felt       `json:"constructor_calldata"`
}

// DeployAccountTxnV1 The structure of a deployAccount transaction.
type DeployAccountTxnV1 struct {
	MaxFee    *felt.Felt         `json:"max_fee"`
	Version   TransactionVersion `json:"version"`
	Signature []*felt.Felt       `json:"signature"`
	Nonce     *felt.Felt         `json:"nonce"`
	Type      TransactionType    `json:"type"`
	// ClassHash The hash of the deployed contract's class
	ClassHash *felt.Felt `json:"class_hash"`

	// ContractAddressSalt The salt for the address of the deployed contract
	ContractAddressSalt *felt.Felt `json:"contract_address_salt"`

	// ConstructorCalldata The parameters passed to the constructor
	ConstructorCalldata []*felt.Felt `json:"constructor_calldata"`
}

type DeployAccountTxnV3 struct {
	Type                TransactionType        `json:"type"`
	Version             TransactionVersion     `json:"version"`
	Signature           []*felt.Felt           `json:"signature"`
	Nonce               *felt.Felt             `json:"nonce"`
	ContractAddressSalt *felt.Felt             `json:"contract_address_salt"`
	ConstructorCalldata []*felt.Felt           `json:"constructor_calldata"`
	ClassHash           *felt.Felt             `json:"class_hash"`
	ResourceBounds      *ResourceBoundsMapping `json:"resource_bounds"`
	Tip                 U64                    `json:"tip"`
	// The data needed to allow the paymaster to pay for the transaction in native tokens
	PayMasterData []*felt.Felt `json:"paymaster_data"`
	// The storage domain of the account's nonce (an account has a nonce per DA mode)
	NonceDataMode DataAvailabilityMode `json:"nonce_data_availability_mode"`
	// The storage domain of the account's balance from which fee will be charged
	FeeMode DataAvailabilityMode `json:"fee_data_availability_mode"`
}

// remarshal is a function that takes in an interface{} value 'v' and an interface{} value 'dst'.
// It marshals the 'v' value to JSON using the json.Marshal function and then unmarshals the JSON data to 'dst' using the json.Unmarshal function.
//
// Parameters:
//   - v: The interface{} value to be marshalled
//   - dst: The interface{} value to be unmarshaled
//
// Returns:
//   - error: An error if the marshalling or unmarshaling process fails
func remarshal(v interface{}, dst interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, dst); err != nil {
		return err
	}

	return nil
}

// string must be NUM_AS_HEX
type TransactionVersion string

const (
	TransactionV0             TransactionVersion = "0x0"
	TransactionV0WithQueryBit TransactionVersion = "0x100000000000000000000000000000000"
	TransactionV1             TransactionVersion = "0x1"
	TransactionV1WithQueryBit TransactionVersion = "0x100000000000000000000000000000001"
	TransactionV2             TransactionVersion = "0x2"
	TransactionV2WithQueryBit TransactionVersion = "0x100000000000000000000000000000002"
	TransactionV3             TransactionVersion = "0x3"
	TransactionV3WithQueryBit TransactionVersion = "0x100000000000000000000000000000003"
)

// BigInt returns a big integer corresponding to the transaction version.
//
// Returns:
//   - *big.Int: a pointer to a big.Int
//   - error: an error if the conversion fails
func (v *TransactionVersion) BigInt() (*big.Int, error) {
	switch *v {
	case TransactionV0:
		return big.NewInt(0), nil
	case TransactionV1:
		return big.NewInt(1), nil
	case TransactionV2:
		return big.NewInt(2), nil
	case TransactionV3:
		return big.NewInt(3), nil
	}

	// Handle versions with query bit.
	// Remove the 0x prefix and convert to big.Int
	version, ok := new(big.Int).SetString(string(*v)[2:], 16)
	if !ok {
		return big.NewInt(-1), errors.New(fmt.Sprint("TransactionVersion %i not supported", *v))
	}

	return version, nil
}

// Int returns an integer corresponding to the transaction version.
// For versions with query bit, it returns the base version number (e.g. TransactionV2WithQueryBit returns 2).
// Returns -1 for invalid versions.
//
// Returns:
//   - int: the integer version, or -1 for invalid versions
func (v *TransactionVersion) Int() int {
	switch *v {
	case TransactionV0, TransactionV0WithQueryBit:
		return 0
	case TransactionV1, TransactionV1WithQueryBit:
		return 1
	case TransactionV2, TransactionV2WithQueryBit:
		return 2
	case TransactionV3, TransactionV3WithQueryBit:
		return 3
	}

	// Handle invalid versions
	return -1
}

// SubPendingTxnsInput is the optional input of the starknet_subscribePendingTransactions subscription.
type SubPendingTxnsInput struct {
	// Optional: Get all transaction details, and not only the hash. If not provided, only hash is returned. Default is false
	TransactionDetails bool `json:"transaction_details,omitempty"`
	// Optional: Filter transactions to only receive notification from address list
	SenderAddress []*felt.Felt `json:"sender_address,omitempty"`
}

// PendingTxn is the response of the starknet_subscribePendingTransactions subscription.
type PendingTxn struct {
	// The hash of the pending transaction. Always present.
	Hash *felt.Felt
	// The full transaction details. Only present if transactionDetails is true.
	Transaction *BlockTransaction
}

// UnmarshalJSON unmarshals the JSON data into a PendingTxn object.
//
// Parameters:
//   - data: The JSON data to be unmarshalled
//
// Returns:
//   - error: An error if the unmarshalling process fails
func (s *PendingTxn) UnmarshalJSON(data []byte) error {
	var txn *BlockTransaction
	if err := json.Unmarshal(data, &txn); err == nil {
		s.Transaction = txn
		s.Hash = txn.Hash

		return nil
	}
	var txnHash *felt.Felt
	if err := json.Unmarshal(data, &txnHash); err == nil {
		s.Hash = txnHash

		return nil
	}

	return errors.New("failed to unmarshal PendingTxn")
}

// UnmarshalJSON unmarshals the data into a BlockTransaction object.
//
// It takes a byte slice as the parameter, representing the JSON data to be unmarshalled.
// The function returns an error if the unmarshalling process fails.
//
// Parameters:
//   - data: The JSON data to be unmarshalled
//
// Returns:
//   - error: An error if the unmarshalling process fails
func (blockTxn *BlockTransaction) UnmarshalJSON(data []byte) error {
	type alias BlockTransaction
	var aux alias

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	txn, err := unmarshalTxn(data)
	if err != nil {
		return err
	}

	blockTxn.Hash = aux.Hash
	blockTxn.Transaction = txn

	return nil
}

// MarshalJSON marshals the BlockTransaction object into a JSON byte slice.
//
// It takes a pointer to a BlockTransaction object as the parameter.
// The function returns a byte slice representing the JSON data and an error if the marshalling process fails.
func (blockTxn *BlockTransaction) MarshalJSON() ([]byte, error) {
	// First marshal the transaction to get all its fields
	txnData, err := json.Marshal(blockTxn.Transaction)
	if err != nil {
		return nil, err
	}

	// Unmarshal into a map to add the hash field
	var result map[string]interface{}
	if err := json.Unmarshal(txnData, &result); err != nil {
		return nil, err
	}

	// Add the hash field
	result["transaction_hash"] = blockTxn.Hash

	return json.Marshal(result)
}
