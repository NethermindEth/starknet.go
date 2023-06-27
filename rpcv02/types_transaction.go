package rpcv02

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	types "github.com/NethermindEth/caigo/types"
	"github.com/NethermindEth/caigo/types/felt"
	"github.com/NethermindEth/caigo/utils"
)

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
	TransactionHash *felt.Felt      `json:"transaction_hash,omitempty"`
	Type            TransactionType `json:"type,omitempty"`
	// MaxFee maximal fee that can be charged for including the transaction
	MaxFee string `json:"max_fee,omitempty"`
	// Version of the transaction scheme
	Version types.NumAsHex `json:"version"`
	// Signature
	Signature []string `json:"signature,omitempty"`
	// Nonce
	Nonce string `json:"nonce,omitempty"`
}

type InvokeTxnV0 struct {
	CommonTransaction
	ContractAddress    *felt.Felt `json:"contract_address"`
	EntryPointSelector string     `json:"entry_point_selector"`

	// Calldata The parameters passed to the function
	Calldata []string `json:"calldata"`
}

func (tx InvokeTxnV0) Hash() *felt.Felt {
	return tx.TransactionHash
}

type InvokeTxnV1 struct {
	CommonTransaction
	SenderAddress *felt.Felt `json:"sender_address"`
	// Calldata The parameters passed to the function
	Calldata []string `json:"calldata"`
}

func (tx InvokeTxnV1) Hash() *felt.Felt {
	return tx.TransactionHash
}

type InvokeTxn interface{}

type L1HandlerTxn struct {
	TransactionHash *felt.Felt      `json:"transaction_hash,omitempty"`
	Type            TransactionType `json:"type,omitempty"`
	// Version of the transaction scheme
	Version types.NumAsHex `json:"version"`
	// Nonce
	Nonce              string     `json:"nonce,omitempty"`
	ContractAddress    *felt.Felt `json:"contract_address"`
	EntryPointSelector string     `json:"entry_point_selector"`

	// Calldata The parameters passed to the function
	Calldata []string `json:"calldata"`
}

func (tx L1HandlerTxn) Hash() *felt.Felt {
	return tx.TransactionHash
}

type DeclareTxn struct {
	CommonTransaction

	// ClassHash the hash of the declared class
	ClassHash string `json:"class_hash"`

	// SenderAddress the address of the account contract sending the declaration transaction
	SenderAddress string `json:"sender_address"`
}

func (tx DeclareTxn) Hash() *felt.Felt {
	return tx.TransactionHash
}

// DeployTxn The structure of a deploy transaction. Note that this transaction type is deprecated and will no longer be supported in future versions
type DeployTxn struct {
	TransactionHash *felt.Felt `json:"transaction_hash,omitempty"`
	// ClassHash The hash of the deployed contract's class
	ClassHash string          `json:"class_hash"`
	Type      TransactionType `json:"type,omitempty"`
	// Version of the transaction scheme
	Version types.NumAsHex `json:"version"`
	// ContractAddressSalt The salt for the address of the deployed contract
	ContractAddressSalt string `json:"contract_address_salt"`
	// ConstructorCalldata The parameters passed to the constructor
	ConstructorCalldata []string `json:"constructor_calldata"`
}

func (tx DeployTxn) Hash() *felt.Felt {
	return tx.TransactionHash
}

type Transaction interface {
	Hash() *felt.Felt
}

// DeployTxn The structure of a deploy transaction. Note that this transaction type is deprecated and will no longer be supported in future versions
type DeployAccountTxn struct {
	CommonTransaction
	// ClassHash The hash of the deployed contract's class
	ClassHash string `json:"class_hash"`

	// ContractAddressSalt The salt for the address of the deployed contract
	ContractAddressSalt string `json:"contract_address_salt"`

	// ConstructorCalldata The parameters passed to the constructor
	ConstructorCalldata []string `json:"constructor_calldata"`
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
			var txn DeclareTxn
			remarshal(casted, &txn)
			return txn, nil
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

type TransactionVersion string

const (
	TransactionV1 TransactionVersion = "0x1"
)

type BroadcastedTransaction interface{}

type BroadcastedTxnCommonProperties struct {
	MaxFee *big.Int `json:"max_fee"`
	// Version of the transaction scheme, should be set to 0 or 1
	Version TransactionVersion `json:"version"`
	// Signature
	Signature []string `json:"signature"`
	Nonce     *big.Int `json:"nonce"`
	Type      string   `json:"type"`
}

type BroadcastedInvokeV1Transaction struct {
	BroadcastedTxnCommonProperties
	SenderAddress *felt.Felt `json:"sender_address"`
	Calldata      []string   `json:"calldata"`
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
	signature := []string{}
	for _, v := range b.Signature {
		s, _ := big.NewInt(0).SetString(v, 0)
		signature = append(signature, fmt.Sprintf("0x%x", s))
	}
	output["signature"] = signature
	output["sender_address"] = b.SenderAddress
	output["calldata"] = b.Calldata
	return json.Marshal(output)
}

type BroadcastedDeclareTransaction struct {
	BroadcastedTxnCommonProperties
	ContractClass types.ContractClass `json:"contract_class"`
	SenderAddress *felt.Felt          `json:"sender_address"`
}

func (b BroadcastedDeclareTransaction) MarshalJSON() ([]byte, error) {
	output := map[string]interface{}{}
	output["type"] = "DECLARE"
	if b.MaxFee != nil {
		output["max_fee"] = fmt.Sprintf("0x%x", b.MaxFee)
	}
	if b.Nonce != nil {
		output["nonce"] = fmt.Sprintf("0x%x", b.Nonce)
	}
	output["version"] = b.Version
	signature := []string{}
	for _, v := range b.Signature {
		s, _ := big.NewInt(0).SetString(v, 0)
		signature = append(signature, fmt.Sprintf("0x%x", s))
	}
	output["signature"] = signature
	output["sender_address"] = b.SenderAddress.String()
	output["contract_class"] = b.ContractClass
	return json.Marshal(output)
}

type BroadcastedDeployTransaction struct {
	Type string `json:"type"`
	// Version of the transaction scheme, should be set to 0 or 1
	Version             *big.Int            `json:"version"`
	ContractAddressSalt string              `json:"contract_address_salt"`
	ConstructorCalldata []string            `json:"constructor_calldata"`
	ContractClass       types.ContractClass `json:"contract_class"`
}

func (b BroadcastedDeployTransaction) MarshalJSON() ([]byte, error) {
	output := map[string]interface{}{}
	output["type"] = "DEPLOY"
	output["version"] = fmt.Sprintf("0x%x", b.Version)
	contractAddressSalt, ok := big.NewInt(0).SetString(b.ContractAddressSalt, 0)
	if !ok {
		return nil, errors.New("wrong salt")
	}
	output["contract_address_salt"] = fmt.Sprintf("0x%x", contractAddressSalt)
	constructorCalldata := []string{}
	for _, v := range b.ConstructorCalldata {
		constructorCall, ok := big.NewInt(0).SetString(v, 0)
		if !ok {
			return nil, errors.New("wrong constructor call data")
		}
		constructorCalldata = append(constructorCalldata, fmt.Sprintf("0x%x", constructorCall))
	}
	output["constructor_calldata"] = constructorCalldata
	output["contract_class"] = b.ContractClass
	return json.Marshal(output)
}

type BroadcastedDeployAccountTransaction struct {
	BroadcastedTxnCommonProperties
	ContractAddressSalt string     `json:"contract_address_salt"`
	ConstructorCalldata []string   `json:"constructor_calldata"`
	ClassHash           *felt.Felt `json:"class_hash"`
}

func (b BroadcastedDeployAccountTransaction) MarshalJSON() ([]byte, error) {
	output := map[string]interface{}{}
	output["type"] = "DEPLOY_ACCOUNT"
	if b.MaxFee != nil {
		output["max_fee"] = fmt.Sprintf("0x%x", b.MaxFee)
	}
	if b.Nonce != nil {
		output["nonce"] = fmt.Sprintf("0x%x", b.Nonce)
	}
	output["version"] = b.Version
	signature := []string{}
	for _, v := range b.Signature {
		s, _ := big.NewInt(0).SetString(v, 0)
		signature = append(signature, fmt.Sprintf("0x%x", s))
	}
	output["signature"] = signature
	contractAddressSalt, ok := big.NewInt(0).SetString(b.ContractAddressSalt, 0)
	if !ok {
		return nil, errors.New("wrong salt")
	}
	output["contract_address_salt"] = fmt.Sprintf("0x%x", contractAddressSalt)
	constructorCalldata := []string{}
	for _, v := range b.ConstructorCalldata {
		constructorCall, ok := big.NewInt(0).SetString(v, 0)
		if !ok {
			return nil, errors.New("wrong constructor call data")
		}
		constructorCalldata = append(constructorCalldata, fmt.Sprintf("0x%x", constructorCall))
	}
	output["constructor_calldata"] = constructorCalldata
	output["class_hash"] = b.ClassHash.String()
	return json.Marshal(output)
}
