package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
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

type TransactionHash struct {
	TransactionHash *felt.Felt `json:"transaction_hash"`
}

func (t TransactionHash) Hash() *felt.Felt {
	return t.TransactionHash
}

func (t *TransactionHash) UnmarshalJSON(input []byte) error {
	return t.TransactionHash.UnmarshalJSON(input)
}

func (t TransactionHash) MarshalJSON() ([]byte, error) {
	return t.TransactionHash.MarshalJSON()
}

func (t TransactionHash) MarshalText() ([]byte, error) {
	return t.TransactionHash.MarshalJSON()
}

func (t *TransactionHash) UnmarshalText(input []byte) error {
	return t.TransactionHash.UnmarshalJSON(input)
}

type CommonTransaction struct {
	TransactionHash *felt.Felt `json:"transaction_hash,omitempty"`
	BroadcastedTxnCommonProperties
}

type InvokeTxnV0 struct {
	CommonTransaction
	FunctionCall
}

func (tx InvokeTxnV0) Hash() *felt.Felt {
	return tx.TransactionHash
}

type InvokeTxnV1 struct {
	CommonTransaction
	SenderAddress *felt.Felt `json:"sender_address"`
	// Calldata The parameters passed to the function
	Calldata []*felt.Felt `json:"calldata"`
}

func (tx InvokeTxnV1) Hash() *felt.Felt {
	return tx.TransactionHash
}

type InvokeTxn interface{}

type L1HandlerTxn struct {
	TransactionHash *felt.Felt      `json:"transaction_hash,omitempty"`
	Type            TransactionType `json:"type,omitempty"`
	// Version of the transaction scheme
	Version NumAsHex `json:"version"`
	// Nonce
	Nonce string `json:"nonce,omitempty"`
	FunctionCall
}

func (tx L1HandlerTxn) Hash() *felt.Felt {
	return tx.TransactionHash
}

type DeclareTxnV1 struct {
	CommonTransaction

	// ClassHash the hash of the declared class
	ClassHash *felt.Felt `json:"class_hash"`

	// SenderAddress the address of the account contract sending the declaration transaction
	SenderAddress *felt.Felt `json:"sender_address"`
}

type DeclareTxnV2 struct {
	CommonTransaction

	ClassHash *felt.Felt `json:"class_hash,omitempty"`

	// SenderAddress the address of the account contract sending the declaration transaction
	SenderAddress *felt.Felt `json:"sender_address"`

	CompiledClassHash *felt.Felt `json:"compiled_class_hash"`
}

func (tx DeclareTxnV1) Hash() *felt.Felt {
	return tx.TransactionHash
}
func (tx DeclareTxnV2) Hash() *felt.Felt {
	return tx.TransactionHash
}

type Transaction interface {
	Hash() *felt.Felt
}

// DeployTxn The structure of a deploy transaction. Note that this transaction type is deprecated and will no longer be supported in future versions
type DeployTxn struct {
	TransactionHash *felt.Felt `json:"transaction_hash,omitempty"`
	// ClassHash The hash of the deployed contract's class
	ClassHash *felt.Felt `json:"class_hash"`

	DeployTransactionProperties
}

func (tx DeployTxn) Hash() *felt.Felt {
	return tx.TransactionHash
}

type DeployAccountTransactionProperties struct {
	// ClassHash The hash of the deployed contract's class
	ClassHash *felt.Felt `json:"class_hash"`

	// ContractAddressSalt The salt for the address of the deployed contract
	ContractAddressSalt *felt.Felt `json:"contract_address_salt"`

	// ConstructorCalldata The parameters passed to the constructor
	ConstructorCalldata []*felt.Felt `json:"constructor_calldata"`
}

// DeployAccountTxn The structure of a deployAccount transaction.
type DeployAccountTxn struct {
	CommonTransaction
	DeployAccountTransactionProperties
}

func (tx DeployAccountTxn) Hash() *felt.Felt {
	return tx.TransactionHash
}

type Transactions []Transaction

func (txns *Transactions) UnmarshalJSON(data []byte) error {
	var dec []interface{}
	if err := json.Unmarshal(data, &dec); err != nil {
		return err
	}

	unmarshalled := make([]Transaction, len(dec))
	for i, t := range dec {
		var err error
		unmarshalled[i], err = unmarshalTxn(t)
		if err != nil {
			return err
		}
	}

	*txns = unmarshalled
	return nil
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
	case string:
		var txn InvokeTransactionReceipt
		txhash, err := utils.HexToFelt(casted)
		if err != nil {
			return txn, err
		}
		return TransactionHash{txhash}, nil
	case map[string]interface{}:
		switch TransactionType(casted["type"].(string)) {
		case TransactionType_Declare:

			switch TransactionType(casted["version"].(string)) {
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

type BroadcastedTransaction interface{}

type BroadcastedTxnCommonProperties struct {
	MaxFee *felt.Felt `json:"max_fee"`
	// Version of the transaction scheme, should be set to 0 or 1
	Version TransactionVersion `json:"version"`
	// Signature
	Signature []*felt.Felt    `json:"signature"`
	Nonce     *felt.Felt      `json:"nonce"`
	Type      TransactionType `json:"type"`
}

// BroadcastedInvokeV1Transaction is BROADCASTED_INVOKE_TXN
// since we only support InvokeV1 transactions
type BroadcastedInvokeV1Transaction struct {
	BroadcastedTxnCommonProperties
	SenderAddress *felt.Felt   `json:"sender_address"`
	Calldata      []*felt.Felt `json:"calldata"`
}

func (b BroadcastedInvokeV1Transaction) MarshalJSON() ([]byte, error) {
	output := map[string]interface{}{}
	output["type"] = b.Type
	if b.MaxFee != nil {
		output["max_fee"] = fmt.Sprintf("0x%x", b.MaxFee)
	}
	if b.Nonce != nil {
		output["nonce"] = fmt.Sprintf("0x%x", b.Nonce)
	}
	output["version"] = b.Version
	signature := b.Signature
	output["signature"] = signature
	output["sender_address"] = b.SenderAddress
	output["calldata"] = b.Calldata
	return json.Marshal(output)
}

type BroadcastedDeclareTransaction interface{}

var _ BroadcastedDeclareTransaction = BroadcastedDeclareTransactionV1{}
var _ BroadcastedDeclareTransaction = BroadcastedDeclareTransactionV2{}

type BroadcastedDeclareTransactionV1 struct {
	BroadcastedTxnCommonProperties
	ContractClass DeprecatedContractClass `json:"contract_class"`
	SenderAddress *felt.Felt              `json:"sender_address"`
}

func (b BroadcastedDeclareTransactionV1) MarshalJSON() ([]byte, error) {
	output := map[string]interface{}{}
	output["type"] = "DECLARE"
	if b.MaxFee != nil {
		output["max_fee"] = fmt.Sprintf("0x%x", b.MaxFee)
	}
	if b.Nonce != nil {
		output["nonce"] = fmt.Sprintf("0x%x", b.Nonce)
	}
	output["version"] = b.Version
	signature := b.Signature
	output["signature"] = signature
	output["sender_address"] = b.SenderAddress.String()
	output["contract_class"] = b.ContractClass
	return json.Marshal(output)
}

type BroadcastedDeclareTransactionV2 struct {
	BroadcastedTxnCommonProperties
	ContractClass     ContractClass `json:"contract_class"`
	SenderAddress     *felt.Felt    `json:"sender_address"`
	CompiledClassHash *felt.Felt    `json:"compiled_class_hash"`
}

func (b BroadcastedDeclareTransactionV2) MarshalJSON() ([]byte, error) {
	output := map[string]interface{}{}
	output["type"] = "DECLARE"
	if b.MaxFee != nil {
		output["max_fee"] = fmt.Sprintf("0x%x", b.MaxFee)
	}
	if b.Nonce != nil {
		output["nonce"] = fmt.Sprintf("0x%x", b.Nonce)
	}
	output["version"] = b.Version
	signature := b.Signature
	output["signature"] = signature
	output["sender_address"] = b.SenderAddress.String()
	output["contract_class"] = b.ContractClass
	output["compiled_class_hash"] = b.CompiledClassHash
	return json.Marshal(output)
}

type DeployTransactionProperties struct {
	Version             TransactionVersion `json:"version"`
	Type                TransactionType    `json:"type"`
	ContractAddressSalt *felt.Felt         `json:"contract_address_salt"`
	ConstructorCalldata []*felt.Felt       `json:"constructor_calldata"`
}

type BroadcastedDeployAccountTransaction struct {
	BroadcastedTxnCommonProperties
	ContractAddressSalt *felt.Felt   `json:"contract_address_salt"`
	ConstructorCalldata []*felt.Felt `json:"constructor_calldata"`
	ClassHash           *felt.Felt   `json:"class_hash"`
}

func (b BroadcastedDeployAccountTransaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(b)
}
