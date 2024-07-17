package rpc

import "github.com/NethermindEth/juno/core/felt"

// AddDeclareTransactionResponse provides the output for AddDeclareTransaction.
type AddDeclareTransactionResponse struct {
	TransactionHash *felt.Felt `json:"transaction_hash"`
	ClassHash       *felt.Felt `json:"class_hash"`
}

// AddDeployTransactionResponse provides the output for AddDeployTransaction.
type AddDeployAccountTransactionResponse struct {
	TransactionHash *felt.Felt `json:"transaction_hash"`
	ContractAddress *felt.Felt `json:"contract_address"`
}

// AddInvokeTransactionResponse provides the output for AddInvokeTransaction.
type AddInvokeTransactionResponse struct {
	TransactionHash *felt.Felt `json:"transaction_hash"`
}

type TransactionResponse struct {
	TransactionHash *felt.Felt `json:"transaction_hash"`
	ClassHash       *felt.Felt `json:"class_hash,omitempty"`
	ContractAddress *felt.Felt `json:"contract_address,omitempty"`
}

// func ConvertToTransactionResponse(resp interface{}) *TransactionResponse {
// 	switch r := resp.(type) {
// 	case *AddInvokeTransactionResponse:
// 		return &TransactionResponse{
// 			TransactionHash: r.TransactionHash,
// 		}
// 	case *AddDeclareTransactionResponse:
// 		return &TransactionResponse{
// 			TransactionHash: r.TransactionHash,
// 			ClassHash:       r.ClassHash,
// 		}
// 	case *AddDeployAccountTransactionResponse:
// 		return &TransactionResponse{
// 			TransactionHash: r.TransactionHash,
// 			ContractAddress: r.ContractAddress,
// 		}
// 	default:
// 		return nil
// 	}
// }
