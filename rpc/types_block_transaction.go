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
var _ BlockTransaction = BlockInvokeTxnV3{}
var _ BlockTransaction = BlockDeclareTxnV0{}
var _ BlockTransaction = BlockDeclareTxnV1{}
var _ BlockTransaction = BlockDeclareTxnV2{}
var _ BlockTransaction = BlockDeclareTxnV3{}
var _ BlockTransaction = BlockDeployTxn{}
var _ BlockTransaction = BlockDeployAccountTxn{}
var _ BlockTransaction = BlockL1HandlerTxn{}

// Hash returns the transaction hash of the BlockInvokeTxnV0.
func (tx BlockInvokeTxnV0) Hash() *felt.Felt {
	return tx.TransactionHash
}

// Hash returns the hash of the BlockInvokeTxnV1 transaction.
func (tx BlockInvokeTxnV1) Hash() *felt.Felt {
	return tx.TransactionHash
}

// Hash returns the hash of the BlockInvokeTxnV3 transaction.
func (tx BlockInvokeTxnV3) Hash() *felt.Felt {
	return tx.TransactionHash
}

// Hash returns the transaction hash of the BlockDeclareTxnV0.
func (tx BlockDeclareTxnV0) Hash() *felt.Felt {
	return tx.TransactionHash
}

// Hash returns the transaction hash of the BlockDeclareTxnV1.
func (tx BlockDeclareTxnV1) Hash() *felt.Felt {
	return tx.TransactionHash
}

// Hash returns the transaction hash of the BlockDeclareTxnV2.
func (tx BlockDeclareTxnV2) Hash() *felt.Felt {
	return tx.TransactionHash
}

// Hash returns the transaction hash of the BlockDeclareTxnV3.
func (tx BlockDeclareTxnV3) Hash() *felt.Felt {
	return tx.TransactionHash
}

// Hash returns the hash of the BlockDeployTxn.
func (tx BlockDeployTxn) Hash() *felt.Felt {
	return tx.TransactionHash
}

// Hash returns the Felt hash of the BlockDeployAccountTxn.
func (tx BlockDeployAccountTxn) Hash() *felt.Felt {
	return tx.TransactionHash
}

// Hash returns the hash of the BlockL1HandlerTxn.
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

type BlockInvokeTxnV3 struct {
	TransactionHash *felt.Felt `json:"transaction_hash"`
	InvokeTxnV3
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

type BlockDeclareTxnV3 struct {
	TransactionHash *felt.Felt `json:"transaction_hash"`
	DeclareTxnV3
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
//
// Parameters:
// - data: The JSON data to be unmarshalled
// Returns:
// - error: An error if the unmarshalling process fails
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
// Parameter:
// - t: The interface{} to be unmarshalled
// Returns:
// - BlockTransaction: a BlockTransaction
// - error: an error if the unmarshaling process fails
func unmarshalBlockTxn(t interface{}) (BlockTransaction, error) {
	switch casted := t.(type) {
	case map[string]interface{}:
		switch TransactionType(casted["type"].(string)) {
		case TransactionType_Declare:

			switch TransactionType(casted["version"].(string)) {
			case "0x0":
				var txn BlockDeclareTxnV0
				err := remarshal(casted, &txn)
				return txn, err
			case "0x1":
				var txn BlockDeclareTxnV1
				err := remarshal(casted, &txn)
				return txn, err
			case "0x2":
				var txn BlockDeclareTxnV2
				err := remarshal(casted, &txn)
				return txn, err
			case "0x3":
				var txn BlockDeclareTxnV3
				err := remarshal(casted, &txn)
				return txn, err
			default:
				return nil, errors.New("internal error with Declare transaction version and unmarshalTxn()")
			}
		case TransactionType_Deploy:
			var txn BlockDeployTxn
			err := remarshal(casted, &txn)
			return txn, err
		case TransactionType_DeployAccount:
			var txn BlockDeployAccountTxn
			err := remarshal(casted, &txn)
			return txn, err
		case TransactionType_Invoke:
			if casted["version"].(string) == "0x0" {
				var txn BlockInvokeTxnV0
				err := remarshal(casted, &txn)
				return txn, err
			} else if casted["version"].(string) == "0x1" {
				var txn BlockInvokeTxnV1
				err := remarshal(casted, &txn)
				return txn, err
			} else {
				var txn BlockInvokeTxnV3
				err := remarshal(casted, &txn)
				return txn, err
			}
		case TransactionType_L1Handler:
			var txn BlockL1HandlerTxn
			err := remarshal(casted, &txn)
			return txn, err
		}
	}

	return nil, fmt.Errorf("unknown transaction type: %v", t)
}
