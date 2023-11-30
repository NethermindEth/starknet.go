package rpc

import "github.com/NethermindEth/juno/core/felt"

type BroadcastTxn interface{}

var (
	_ BroadcastTxn = BroadcastInvokev0Txn{}
	_ BroadcastTxn = BroadcastInvokev1Txn{}
	_ BroadcastTxn = BroadcastDeclareV1Txn{}
	_ BroadcastTxn = BroadcastDeclareV2Txn{}
	_ BroadcastTxn = BroadcastDeclareTxnV3{}
	_ BroadcastTxn = BroadcastDeployAccountTxn{}
)

type BroadcastInvokeTxnType interface{}

var (
	_ BroadcastInvokeTxnType = BroadcastInvokev0Txn{}
	_ BroadcastInvokeTxnType = BroadcastInvokev1Txn{}
	_ BroadcastInvokeTxnType = BroadcastInvokev3Txn{}
)

type BroadcastDeclareTxn interface{}

var (
	_ BroadcastDeclareTxn = BroadcastDeclareV1Txn{}
	_ BroadcastDeclareTxn = BroadcastDeclareV2Txn{}
)

type BroadcastInvokev0Txn struct {
	InvokeTxnV0
}

type BroadcastInvokev1Txn struct {
	InvokeTxnV1
}

type BroadcastInvokev3Txn struct {
	InvokeTxnV3
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

type BroadcastDeclareTxnV3 struct {
	Type              TransactionType       `json:"type"`
	SenderAddress     *felt.Felt            `json:"sender_address"`
	CompiledClassHash *felt.Felt            `json:"compiled_class_hash"`
	Version           NumAsHex              `json:"version"`
	Signature         []*felt.Felt          `json:"signature"`
	Nonce             *felt.Felt            `json:"nonce"`
	ContractClass     *ContractClass        `json:"contract_class"`
	ResourceBounds    ResourceBoundsMapping `json:"resource_bounds"`
	Tip               *felt.Felt            `json:"tip"`
	// The data needed to allow the paymaster to pay for the transaction in native tokens
	PayMasterData []*felt.Felt `json:"paymaster_data"`
	// The data needed to deploy the account contract from which this tx will be initiated
	AccountDeploymentData *felt.Felt `json:"account_deployment_data"`
	// The storage domain of the account's nonce (an account has a nonce per DA mode)
	NonceDataMode DataAvailabilityMode `json:"nonce_data_availability_mode"`
	// The storage domain of the account's balance from which fee will be charged
	FeeMode DataAvailabilityMode `json:"fee_data_availability_mode"`
}

type BroadcastDeployAccountTxn struct {
	DeployAccountTxn
}
