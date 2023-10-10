package rpc

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
)

type BlockTransactions []BlockTransaction

type BlockTransaction interface {
	Hash() *felt.Felt
}

var _ BlockTransaction = BlockInvokeTxnV0{}
var _ BlockTransaction = BlockInvokeTxnV1{}
var _ BlockTransaction = BlockDeclareTxnV0{}
var _ BlockTransaction = BlockDeclareTxnV1{}
var _ BlockTransaction = BlockDeclareTxnV2{}
var _ BlockTransaction = BlockDeployTxn{}
var _ BlockTransaction = BlockDeployAccountTxn{}
var _ BlockTransaction = BlockL1HandlerTxn{}

// Hash returns the transaction hash of the BlockInvokeTxnV0.
//
// It returns a *felt.Felt, which represents the transaction hash.
func (tx BlockInvokeTxnV0) Hash() *felt.Felt {
	return tx.TransactionHash
}

// Hash returns the hash of the BlockInvokeTxnV1 transaction.
//
// It returns a pointer to a felt.Felt value.
func (tx BlockInvokeTxnV1) Hash() *felt.Felt {
	return tx.TransactionHash
}
// Hash returns the transaction hash of the BlockDeclareTxnV0.
//
// Returns a pointer to felt.Felt.
func (tx BlockDeclareTxnV0) Hash() *felt.Felt {
	return tx.TransactionHash
}

// Hash returns the transaction hash of the BlockDeclareTxnV1.
//
// It does not take any parameters.
// It returns a pointer to a felt.Felt object.
func (tx BlockDeclareTxnV1) Hash() *felt.Felt {
	return tx.TransactionHash
}

// Hash returns the transaction hash of the BlockDeclareTxnV2.
//
// No parameters.
// Returns *felt.Felt.
func (tx BlockDeclareTxnV2) Hash() *felt.Felt {
	return tx.TransactionHash
}

// Hash returns the hash of the BlockDeployTxn.
//
// No parameters.
// Return type: *felt.Felt.
func (tx BlockDeployTxn) Hash() *felt.Felt {
	return tx.TransactionHash
}

// Hash returns the Felt hash of the BlockDeployAccountTxn.
//
// It returns a pointer to a felt.Felt type.
func (tx BlockDeployAccountTxn) Hash() *felt.Felt {
	return tx.TransactionHash
}

// Hash returns the hash of the BlockL1HandlerTxn.
//
// No parameters.
// Returns a pointer to a felt.Felt.
func (tx BlockL1HandlerTxn) Hash() *felt.Felt {
	return tx.TransactionHash
}

type BlockInvokeTxnV0 struct {
	TransactionHash *felt.Felt `json:"transaction_hash"`
	InvokeTxnV0
}

type BlockInvokeTxnV1 struct {
	TransactionHash *felt.Felt `json:"transaction_hash"`
	InvokeTxnV1
}

type BlockL1HandlerTxn struct {
	TransactionHash *felt.Felt `json:"transaction_hash"`
	L1HandlerTxn
}

type BlockDeclareTxnV0 struct {
	TransactionHash *felt.Felt `json:"transaction_hash"`
	DeclareTxnV0
}

type BlockDeclareTxnV1 struct {
	TransactionHash *felt.Felt `json:"transaction_hash"`
	DeclareTxnV1
}

type BlockDeclareTxnV2 struct {
	TransactionHash *felt.Felt `json:"transaction_hash"`
	DeclareTxnV2
}

type BlockDeployTxn struct {
	TransactionHash *felt.Felt `json:"transaction_hash"`
	DeployTxn
}

type BlockDeployAccountTxn struct {
	TransactionHash *felt.Felt `json:"transaction_hash"`
	DeployAccountTxn
}

// UnmarshalJSON unmarshals the data into a BlockTransactions object.
//
// It takes a byte slice as the parameter, representing the JSON data to be unmarshalled.
// The function returns an error if the unmarshalling process fails.
func (txns *BlockTransactions) UnmarshalJSON(data []byte) error {
	var dec []interface{}
	if err := json.Unmarshal(data, &dec); err != nil {
		return err
	}

	unmarshalled := make([]BlockTransaction, len(dec))
	for i, t := range dec {
		var err error
		unmarshalled[i], err = unmarshalBlockTxn(t)
		if err != nil {
			return err
		}
	}

	*txns = unmarshalled
	return nil
}

// unmarshalBlockTxn unmarshals a given interface and returns a BlockTransaction.
//
// The function takes an interface `t` as a parameter, which is used to determine the type of transaction.
// It returns a BlockTransaction and an error.
func unmarshalBlockTxn(t interface{}) (BlockTransaction, error) {
	switch casted := t.(type) {
	case map[string]interface{}:
		switch TransactionType(casted["type"].(string)) {
		case TransactionType_Declare:

			switch TransactionType(casted["version"].(string)) {
			case "0x0":
				var txn BlockDeclareTxnV0
				remarshal(casted, &txn)
				return txn, nil
			case "0x1":
				var txn BlockDeclareTxnV1
				remarshal(casted, &txn)
				return txn, nil
			case "0x2":
				var txn BlockDeclareTxnV2
				remarshal(casted, &txn)
				return txn, nil
			default:
				return nil, errors.New("Internal error with Declare transaction version and unmarshalTxn()")
			}
		case TransactionType_Deploy:
			var txn BlockDeployTxn
			remarshal(casted, &txn)
			return txn, nil
		case TransactionType_DeployAccount:
			var txn BlockDeployAccountTxn
			remarshal(casted, &txn)
			return txn, nil
		case TransactionType_Invoke:
			if casted["version"].(string) == "0x0" {
				var txn BlockInvokeTxnV0
				remarshal(casted, &txn)
				return txn, nil
			} else {
				var txn BlockInvokeTxnV1
				remarshal(casted, &txn)
				return txn, nil
			}
		case TransactionType_L1Handler:
			var txn BlockL1HandlerTxn
			remarshal(casted, &txn)
			return txn, nil
		}
	}

	return nil, fmt.Errorf("unknown transaction type: %v", t)
}
