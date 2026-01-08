package rpcv10

import (
	"encoding/json"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc/types"
)

type BlockTransaction struct {
	Hash *felt.Felt `json:"transaction_hash"`
	types.Transaction
}

// SubPendingTxnsInput is the optional input of the
// starknet_subscribeNewTransactionReceipts subscription.
type SubNewTxnReceiptsInput struct {
	// Optional: A vector of finality statuses to receive updates for.
	// Only `PRE_CONFIRMED` and `ACCEPTED_ON_L2` are supported. Default is
	// `ACCEPTED_ON_L2`.
	FinalityStatus []TxnFinalityStatus `json:"finality_status,omitempty"`
	// Optional: Filter transaction receipts to only include transactions
	// sent by the specified addresses
	SenderAddress []*felt.Felt `json:"sender_address,omitempty"`
}

// SubNewTxnsInput is the optional input of the
// starknet_subscribeNewTransactions subscription.
type SubNewTxnsInput struct {
	// Optional: A vector of finality statuses to receive updates for.
	// Support all transaction statuses, except `ACCEPTED_ON_L1`. Default is
	// `ACCEPTED_ON_L2`.
	FinalityStatus []TxnStatus `json:"finality_status,omitempty"`
	// Optional: Filter transaction receipts to only include transactions sent
	// by the specified addresses
	SenderAddress []*felt.Felt `json:"sender_address,omitempty"`
}

// TxnWithHashAndStatus is the response of the
// starknet_subscribeNewTransactions subscription.
type TxnWithHashAndStatus struct {
	// Transaction with hash and status
	BlockTransaction
	// Finality status of the transaction, except `ACCEPTED_ON_L1`.
	FinalityStatus TxnStatus `json:"finality_status"`
}

func (txn *TxnWithHashAndStatus) UnmarshalJSON(data []byte) error {
	// type alias TxnWithHashAndStatus
	var aux BlockTransaction

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	var aux2 struct {
		FinalityStatus TxnStatus `json:"finality_status"`
	}

	err = json.Unmarshal(data, &aux2)
	if err != nil {
		return err
	}

	txn.BlockTransaction = aux
	txn.FinalityStatus = aux2.FinalityStatus

	return nil
}

// MarshalJSON marshals the TxnWithHashAndStatus object into a JSON byte slice.
func (txn *TxnWithHashAndStatus) MarshalJSON() ([]byte, error) {
	blockTxnData, err := json.Marshal(&txn.BlockTransaction)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(blockTxnData, &result); err != nil {
		return nil, err
	}

	result["finality_status"] = txn.FinalityStatus

	return json.Marshal(result)
}

// UnmarshalJSON unmarshals the data into a BlockTransaction object.
//
// It takes a byte slice as the parameter, representing the JSON data to be
// unmarshalled.
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
// The function returns a byte slice representing the JSON data and an error if
// the marshalling process fails.
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
