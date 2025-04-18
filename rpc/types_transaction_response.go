package rpc

import "github.com/NethermindEth/juno/core/felt"

// AddDeclareTransactionResponse provides the output for AddDeclareTransaction.
type AddDeclareTransactionResponse struct {
	Hash      *felt.Felt `json:"transaction_hash"`
	ClassHash *felt.Felt `json:"class_hash"`
}

// AddDeployTransactionResponse provides the output for AddDeployTransaction.
type AddDeployAccountTransactionResponse struct {
	Hash            *felt.Felt `json:"transaction_hash"`
	ContractAddress *felt.Felt `json:"contract_address"`
}

// AddInvokeTransactionResponse provides the output for AddInvokeTransaction.
type AddInvokeTransactionResponse struct {
	Hash *felt.Felt `json:"transaction_hash"`
}

// TransactionResponse is a generic response for all transaction types sent to the network.
type TransactionResponse struct {
	// Present for all transaction types
	Hash *felt.Felt `json:"transaction_hash"`
	// Present only for declare transactions
	ClassHash *felt.Felt `json:"class_hash,omitempty"`
	// Present only for deploy transactions
	ContractAddress *felt.Felt `json:"contract_address,omitempty"`
}
