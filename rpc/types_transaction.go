package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
)

// https://github.com/starkware-libs/starknet-specs/blob/a789ccc3432c57777beceaa53a34a7ae2f25fda0/api/starknet_api_openrpc.json#L1252
type TXN struct {
	Hash                *felt.Felt      `json:"transaction_hash,omitempty"`
	Type                TransactionType `json:"type"`
	Version             *felt.Felt      `json:"version,omitempty"`
	Nonce               *felt.Felt      `json:"nonce,omitempty"`
	MaxFee              *felt.Felt      `json:"max_fee,omitempty"`
	ContractAddress     *felt.Felt      `json:"contract_address,omitempty"`
	ContractAddressSalt *felt.Felt      `json:"contract_address_salt,omitempty"`
	ClassHash           *felt.Felt      `json:"class_hash,omitempty"`
	ConstructorCalldata []*felt.Felt    `json:"constructor_calldata,omitempty"`
	SenderAddress       *felt.Felt      `json:"sender_address,omitempty"`
	Signature           *[]*felt.Felt   `json:"signature,omitempty"`
	Calldata            *[]*felt.Felt   `json:"calldata,omitempty"`
	EntryPointSelector  *felt.Felt      `json:"entry_point_selector,omitempty"`
	CompiledClassHash   *felt.Felt      `json:"compiled_class_hash,omitempty"`
}

type InvokeTxnV0 struct {
	MaxFee    *felt.Felt         `json:"max_fee"`
	Version   TransactionVersion `json:"version"`
	Signature []*felt.Felt       `json:"signature"`
	Nonce     *felt.Felt         `json:"nonce"`
	Type      TransactionType    `json:"type"`
	FunctionCall
}

type InvokeTxnV1 struct {
	MaxFee        *felt.Felt         `json:"max_fee"`
	Version       TransactionVersion `json:"version"`
	Signature     []*felt.Felt       `json:"signature"`
	Nonce         *felt.Felt         `json:"nonce"`
	Type          TransactionType    `json:"type"`
	SenderAddress *felt.Felt         `json:"sender_address"`
	// Calldata The parameters passed to the function
	Calldata []*felt.Felt `json:"calldata"`
}

type L1HandlerTxn struct {
	Type TransactionType `json:"type,omitempty"`
	// Version of the transaction scheme
	Version NumAsHex `json:"version"`
	// Nonce
	Nonce string `json:"nonce,omitempty"`
	FunctionCall
}

type DeclareTxnV0 struct {
	MaxFee    *felt.Felt         `json:"max_fee"`
	Version   TransactionVersion `json:"version"`
	Signature []*felt.Felt       `json:"signature"`
	Nonce     *felt.Felt         `json:"nonce"`
	Type      TransactionType    `json:"type"`

	// SenderAddress the address of the account contract sending the declaration transaction
	SenderAddress *felt.Felt `json:"sender_address"`

	DeprecatedContractClass `json:"contract_class,omitempty"`
	ClassHash               *felt.Felt `json:"class_hash,omitempty"`
}

type DeclareTxnV1 struct {
	MaxFee    *felt.Felt         `json:"max_fee"`
	Version   TransactionVersion `json:"version"`
	Signature []*felt.Felt       `json:"signature"`
	Nonce     *felt.Felt         `json:"nonce"`
	Type      TransactionType    `json:"type"`

	// ClassHash the hash of the declared class
	ClassHash *felt.Felt `json:"class_hash,omitempty"`

	DeprecatedContractClass `json:"contract_class,omitempty"`

	// SenderAddress the address of the account contract sending the declaration transaction
	SenderAddress *felt.Felt `json:"sender_address"`
}

type DeclareTxnV2 struct {
	MaxFee    *felt.Felt         `json:"max_fee"`
	Version   TransactionVersion `json:"version"`
	Signature []*felt.Felt       `json:"signature"`
	Nonce     *felt.Felt         `json:"nonce"`
	Type      TransactionType    `json:"type"`

	// SenderAddress the address of the account contract sending the declaration transaction
	SenderAddress *felt.Felt `json:"sender_address"`

	CompiledClassHash *felt.Felt `json:"compiled_class_hash"`

	ContractClass `json:"contract_class,omitempty"`
	ClassHash     *felt.Felt `json:"class_hash,omitempty"`
}

type Transaction interface {
	GetType() TransactionType
}

var _ Transaction = InvokeTxnV0{}
var _ Transaction = InvokeTxnV1{}
var _ Transaction = DeclareTxnV1{}
var _ Transaction = DeclareTxnV2{}
var _ Transaction = DeployTxn{}
var _ Transaction = DeployAccountTxn{}
var _ Transaction = L1HandlerTxn{}

func (tx InvokeTxnV0) GetType() TransactionType {
	return tx.Type
}

func (tx InvokeTxnV1) GetType() TransactionType {
	return tx.Type
}
func (tx DeclareTxnV0) GetType() TransactionType {
	return tx.Type
}
func (tx DeclareTxnV1) GetType() TransactionType {
	return tx.Type
}
func (tx DeclareTxnV2) GetType() TransactionType {
	return tx.Type
}
func (tx DeployTxn) GetType() TransactionType {
	return tx.Type
}
func (tx DeployAccountTxn) GetType() TransactionType {
	return tx.Type
}
func (tx L1HandlerTxn) GetType() TransactionType {
	return tx.Type
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

// DeployAccountTxn The structure of a deployAccount transaction.
type DeployAccountTxn struct {
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

type UnknownTransaction struct{ Transaction }

func (txn *UnknownTransaction) UnmarshalJSON(data []byte) error {
	var dec map[string]interface{}
	if err := json.Unmarshal(data, &dec); err != nil {
		return err
	}

	t, err := unmarshalTxn(dec)
	if err != nil {
		return err
	}

	*txn = UnknownTransaction{t}
	return nil
}

func unmarshalTxn(t interface{}) (Transaction, error) {
	switch casted := t.(type) {
	case map[string]interface{}:
		switch TransactionType(casted["type"].(string)) {
		case TransactionType_Declare:

			switch TransactionType(casted["version"].(string)) {
			case "0x0":
				var txn DeclareTxnV0
				remarshal(casted, &txn)
				return txn, nil
			case "0x1":
				var txn DeclareTxnV1
				remarshal(casted, &txn)
				return txn, nil
			case "0x2":
				var txn DeclareTxnV2
				remarshal(casted, &txn)
				return txn, nil
			default:
				return nil, errors.New("Internal error with Declare transaction version and unmarshalTxn()")
			}
		case TransactionType_Deploy:
			var txn DeployTxn
			remarshal(casted, &txn)
			return txn, nil
		case TransactionType_DeployAccount:
			var txn DeployAccountTxn
			remarshal(casted, &txn)
			return txn, nil
		case TransactionType_Invoke:
			if casted["version"].(string) == "0x0" {
				var txn InvokeTxnV0
				remarshal(casted, &txn)
				return txn, nil
			} else {
				var txn InvokeTxnV1
				remarshal(casted, &txn)
				return txn, nil
			}
		case TransactionType_L1Handler:
			var txn L1HandlerTxn
			remarshal(casted, &txn)
			return txn, nil
		}
	}

	return nil, fmt.Errorf("unknown transaction type: %v", t)
}

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
	TransactionV0 TransactionVersion = "0x0"
	TransactionV1 TransactionVersion = "0x1"
	TransactionV2 TransactionVersion = "0x2"
)

func (v *TransactionVersion) BigInt() (*big.Int, error) {
	switch *v {
	case TransactionV0:
		return big.NewInt(0), nil
	case TransactionV1:
		return big.NewInt(1), nil
	default:
		return big.NewInt(-1), errors.New(fmt.Sprint("TransactionVersion %i not supported", *v))
	}
}

type AddDeclareTxnInput interface{}

var _ AddDeclareTxnInput = DeclareTxnV1{}
var _ AddDeclareTxnInput = DeclareTxnV2{}

type EstimateFeeInput interface{}

var _ EstimateFeeInput = InvokeTxnV0{}
var _ EstimateFeeInput = InvokeTxnV1{}
var _ EstimateFeeInput = DeployAccountTxn{}
var _ EstimateFeeInput = DeclareTxnV1{}
var _ EstimateFeeInput = DeclareTxnV2{}
