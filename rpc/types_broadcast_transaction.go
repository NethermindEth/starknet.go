package rpc

import "github.com/NethermindEth/juno/core/felt"

type BroadcastTxn interface{}

var (
	_ BroadcastTxn = BroadcastInvokev0Txn{}
	_ BroadcastTxn = BroadcastInvokev1Txn{}
	_ BroadcastTxn = BroadcastDeclareV1Txn{}
	_ BroadcastTxn = BroadcastDeclareV2Txn{}
	_ BroadcastTxn = BroadcastDeployAccountTxn{}
)

type BroadcastInvokeTxn interface{}

var (
	_ BroadcastInvokeTxn = BroadcastInvokev0Txn{}
	_ BroadcastInvokeTxn = BroadcastInvokev1Txn{}
)

type BroadcastDeclareTxn interface{}

var (
	_ BroadcastInvokeTxn = BroadcastDeclareV1Txn{}
	_ BroadcastInvokeTxn = BroadcastDeclareV2Txn{}
)

type BroadcastInvokev0Txn struct {
	InvokeTxnV0
}

type BroadcastInvokev1Txn struct {
	InvokeTxnV1
}

type BroadcastDeclareV1Txn struct {
	Type TransactionType `json:"type"`
	// SenderAddress the address of the account contract sending the declaration transaction
	SenderAddress *felt.Felt              `json:"sender_address"`
	MaxFee        *felt.Felt              `json:"max_fee"`
	Version       NumAsHex                `json:"version"`
	Signature     []*felt.Felt            `json:"signature"`
	Nonce         *felt.Felt              `json:"nonce"`
	ContractClass DeprecatedContractClass `json:"contract_class"`
}
type BroadcastDeclareV2Txn struct {
	Type TransactionType `json:"type"`
	// SenderAddress the address of the account contract sending the declaration transaction
	SenderAddress     *felt.Felt    `json:"sender_address"`
	CompiledClassHash *felt.Felt    `json:"compiled_class_hash"`
	MaxFee            *felt.Felt    `json:"max_fee"`
	Version           NumAsHex      `json:"version"`
	Signature         []*felt.Felt  `json:"signature"`
	Nonce             *felt.Felt    `json:"nonce"`
	ContractClass     ContractClass `json:"contract_class"`
}

type BroadcastDeployAccountTxn struct {
	DeployAccountTxn
}
