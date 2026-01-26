package rpcv10

import (
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
)

type BroadcastTxn interface{}

// Note: this allow all types to pass, but are to help users of starknet.go
// understand which types are allowed where.
var (
	_ BroadcastTxn = (*BroadcastInvokeTxnV3)(nil)
	_ BroadcastTxn = (*BroadcastDeclareTxnV3)(nil)
	_ BroadcastTxn = (*BroadcastDeployAccountTxnV3)(nil)
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
	ResourceBounds    *ResourceBoundsMapping   `json:"resource_bounds"`
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

func (tx BroadcastDeclareTxnV3) GetType() TransactionType                   { return tx.Type }
func (tx BroadcastDeclareTxnV3) GetSenderAddress() *felt.Felt               { return tx.SenderAddress }
func (tx BroadcastDeclareTxnV3) GetCompiledClassHash() *felt.Felt           { return tx.CompiledClassHash }
func (tx BroadcastDeclareTxnV3) GetVersion() TransactionVersion             { return tx.Version }
func (tx BroadcastDeclareTxnV3) GetSignature() []*felt.Felt                 { return tx.Signature }
func (tx BroadcastDeclareTxnV3) GetNonce() *felt.Felt                       { return tx.Nonce }
func (tx BroadcastDeclareTxnV3) GetContractClass() *contracts.ContractClass { return tx.ContractClass }
func (tx BroadcastDeclareTxnV3) GetResourceBounds() *ResourceBoundsMapping  { return tx.ResourceBounds }
func (tx BroadcastDeclareTxnV3) GetTip() U64                                { return tx.Tip }
func (tx BroadcastDeclareTxnV3) GetPayMasterData() []*felt.Felt             { return tx.PayMasterData }
func (tx BroadcastDeclareTxnV3) GetAccountDeploymentData() []*felt.Felt {
	return tx.AccountDeploymentData
}
func (tx BroadcastDeclareTxnV3) GetNonceDataMode() DataAvailabilityMode { return tx.NonceDataMode }
func (tx BroadcastDeclareTxnV3) GetFeeMode() DataAvailabilityMode       { return tx.FeeMode }
