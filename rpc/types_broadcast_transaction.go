package rpc

import "github.com/NethermindEth/juno/core/felt"

type BroadcastTxn interface{}

var (
	_ BroadcastTxn = BroadcastInvokev0Txn{}
	_ BroadcastTxn = BroadcastInvokev1Txn{}
	_ BroadcastTxn = BroadcastDeclareTxnV1{}
	_ BroadcastTxn = BroadcastDeclareTxnV2{}
	_ BroadcastTxn = BroadcastDeclareTxnV3{}
	_ BroadcastTxn = BroadcastDeployAccountTxn{}
)

type BroadcastInvokeTxnType interface {
	GetCalldata() []*felt.Felt
}

var (
	_ BroadcastInvokeTxnType = BroadcastInvokev0Txn{}
	_ BroadcastInvokeTxnType = BroadcastInvokev1Txn{}
	_ BroadcastInvokeTxnType = BroadcastInvokev3Txn{}
)

type BroadcastDeclareTxnType interface {
	GetContractClass() interface{}
}

var (
	_ BroadcastDeclareTxnType = BroadcastDeclareTxnV1{}
	_ BroadcastDeclareTxnType = BroadcastDeclareTxnV2{}
	_ BroadcastDeclareTxnType = BroadcastDeclareTxnV3{}
)

type BroadcastAddDeployTxnType interface {
	GetConstructorCalldata() []*felt.Felt
}

var (
	_ BroadcastAddDeployTxnType = BroadcastDeployAccountTxn{}
	_ BroadcastAddDeployTxnType = BroadcastDeployAccountTxnV3{}
)

type BroadcastInvokev0Txn struct {
	InvokeTxnV0
}

func (tx BroadcastInvokev0Txn) GetCalldata() []*felt.Felt {
	return tx.Calldata
}

type BroadcastInvokev1Txn struct {
	InvokeTxnV1
}

func (tx BroadcastInvokev1Txn) GetCalldata() []*felt.Felt {
	return tx.Calldata
}

type BroadcastInvokev3Txn struct {
	InvokeTxnV3
}

func (tx BroadcastInvokev3Txn) GetCalldata() []*felt.Felt {
	return tx.Calldata
}

type BroadcastDeclareTxnV1 struct {
	Type TransactionType `json:"type"`
	// SenderAddress the address of the account contract sending the declaration transaction
	SenderAddress *felt.Felt              `json:"sender_address"`
	MaxFee        *felt.Felt              `json:"max_fee"`
	Version       TransactionVersion      `json:"version"`
	Signature     []*felt.Felt            `json:"signature"`
	Nonce         *felt.Felt              `json:"nonce"`
	ContractClass DeprecatedContractClass `json:"contract_class"`
}

func (tx BroadcastDeclareTxnV1) GetContractClass() interface{} {
	return tx.ContractClass
}

type BroadcastDeclareTxnV2 struct {
	Type TransactionType `json:"type"`
	// SenderAddress the address of the account contract sending the declaration transaction
	SenderAddress     *felt.Felt         `json:"sender_address"`
	CompiledClassHash *felt.Felt         `json:"compiled_class_hash"`
	MaxFee            *felt.Felt         `json:"max_fee"`
	Version           TransactionVersion `json:"version"`
	Signature         []*felt.Felt       `json:"signature"`
	Nonce             *felt.Felt         `json:"nonce"`
	ContractClass     ContractClass      `json:"contract_class"`
}

func (tx BroadcastDeclareTxnV2) GetContractClass() interface{} {
	return tx.ContractClass
}

type BroadcastDeclareTxnV3 struct {
	Type              TransactionType       `json:"type"`
	SenderAddress     *felt.Felt            `json:"sender_address"`
	CompiledClassHash *felt.Felt            `json:"compiled_class_hash"`
	Version           TransactionVersion    `json:"version"`
	Signature         []*felt.Felt          `json:"signature"`
	Nonce             *felt.Felt            `json:"nonce"`
	ContractClass     *ContractClass        `json:"contract_class"`
	ResourceBounds    ResourceBoundsMapping `json:"resource_bounds"`
	Tip               U64                   `json:"tip"`
	// The data needed to allow the paymaster to pay for the transaction in native tokens
	PayMasterData []*felt.Felt `json:"paymaster_data"`
	// The data needed to deploy the account contract from which this tx will be initiated
	AccountDeploymentData []*felt.Felt `json:"account_deployment_data"`
	// The storage domain of the account's nonce (an account has a nonce per DA mode)
	NonceDataMode DataAvailabilityMode `json:"nonce_data_availability_mode"`
	// The storage domain of the account's balance from which fee will be charged
	FeeMode DataAvailabilityMode `json:"fee_data_availability_mode"`
}

func (tx BroadcastDeclareTxnV3) GetContractClass() interface{} {
	return *tx.ContractClass
}

type BroadcastDeployAccountTxn struct {
	DeployAccountTxn
}

func (tx BroadcastDeployAccountTxn) GetConstructorCalldata() []*felt.Felt {
	return tx.ConstructorCalldata
}

type BroadcastDeployAccountTxnV3 struct {
	DeployAccountTxnV3
}

func (tx BroadcastDeployAccountTxnV3) GetConstructorCalldata() []*felt.Felt {
	return tx.ConstructorCalldata
}
