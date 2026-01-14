package rpcv10

import (
	"encoding/json"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc/types"
)

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
	types.BlockTransaction
	// Finality status of the transaction, except `ACCEPTED_ON_L1`.
	FinalityStatus TxnStatus `json:"finality_status"`
}

func (txn *TxnWithHashAndStatus) UnmarshalJSON(data []byte) error {
	// type alias TxnWithHashAndStatus
	var aux types.BlockTransaction

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
