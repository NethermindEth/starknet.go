package rpc

import (
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
)

type BroadcastTxn interface{}

// Note: this allow all types to pass, but are to help users of starknet.go
// understand which types are allowed where.
var (
	_ BroadcastTxn = BroadcastInvokeTxnV3{}
	_ BroadcastTxn = BroadcastDeclareTxnV3{}
	_ BroadcastTxn = BroadcastDeployAccountTxnV3{}
)

type BroadcastInvokeTxnV3 = InvokeTxnV3

type BroadcastDeployAccountTxnV3 = DeployAccountTxnV3

type BroadcastDeclareTxnV3 struct {
	Type              TransactionType          `json:"type"`
	SenderAddress     *felt.Felt               `json:"sender_address"`
	CompiledClassHash *felt.Felt               `json:"compiled_class_hash"`
	Version           TransactionVersion       `json:"version"`
	Signature         []*felt.Felt             `json:"signature"`
	Nonce             *felt.Felt               `json:"nonce"`
	ContractClass     *contracts.ContractClass `json:"contract_class"`
	ResourceBounds    ResourceBoundsMapping    `json:"resource_bounds"`
	Tip               U64                      `json:"tip"`
	// The data needed to allow the paymaster to pay for the transaction in native tokens
	PayMasterData []*felt.Felt `json:"paymaster_data"`
	// The data needed to deploy the account contract from which this tx will be initiated
	AccountDeploymentData []*felt.Felt `json:"account_deployment_data"`
	// The storage domain of the account's nonce (an account has a nonce per DA mode)
	NonceDataMode DataAvailabilityMode `json:"nonce_data_availability_mode"`
	// The storage domain of the account's balance from which fee will be charged
	FeeMode DataAvailabilityMode `json:"fee_data_availability_mode"`
}
