package rpc

import "github.com/NethermindEth/juno/core/felt"

// Transaction is an interface that represents a Starknet transaction.
// It is used to provide a common interface for all transaction types.
// The 'Type' and 'Version' fields are present in all transaction types.
type Transaction interface {
	GetType() TransactionType
	GetVersion() TransactionVersion
}

var _ Transaction = InvokeTxnV0{}
var _ Transaction = InvokeTxnV1{}
var _ Transaction = InvokeTxnV3{}
var _ Transaction = DeclareTxnV1{}
var _ Transaction = DeclareTxnV2{}
var _ Transaction = DeclareTxnV3{}
var _ Transaction = DeployTxn{}
var _ Transaction = DeployAccountTxn{}
var _ Transaction = DeployAccountTxnV3{}
var _ Transaction = L1HandlerTxn{}

// Invoke transactions
func (tx InvokeTxnV0) GetType() TransactionType {
	return tx.Type
}
func (tx InvokeTxnV0) GetVersion() TransactionVersion {
	return tx.Version
}
func (tx InvokeTxnV1) GetType() TransactionType {
	return tx.Type
}
func (tx InvokeTxnV1) GetVersion() TransactionVersion {
	return tx.Version
}
func (tx InvokeTxnV3) GetType() TransactionType {
	return tx.Type
}
func (tx InvokeTxnV3) GetVersion() TransactionVersion {
	return tx.Version
}

// Declare transactions
func (tx DeclareTxnV0) GetType() TransactionType {
	return tx.Type
}
func (tx DeclareTxnV0) GetVersion() TransactionVersion {
	return tx.Version
}
func (tx DeclareTxnV1) GetType() TransactionType {
	return tx.Type
}
func (tx DeclareTxnV1) GetVersion() TransactionVersion {
	return tx.Version
}
func (tx DeclareTxnV2) GetType() TransactionType {
	return tx.Type
}
func (tx DeclareTxnV2) GetVersion() TransactionVersion {
	return tx.Version
}
func (tx DeclareTxnV3) GetType() TransactionType {
	return tx.Type
}
func (tx DeclareTxnV3) GetVersion() TransactionVersion {
	return tx.Version
}

// Deploy transaction
func (tx DeployTxn) GetType() TransactionType {
	return tx.Type
}
func (tx DeployTxn) GetVersion() TransactionVersion {
	return tx.Version
}

// DeployAccount transactions
func (tx DeployAccountTxn) GetType() TransactionType {
	return tx.Type
}
func (tx DeployAccountTxn) GetVersion() TransactionVersion {
	return tx.Version
}
func (tx DeployAccountTxnV3) GetType() TransactionType {
	return tx.Type
}
func (tx DeployAccountTxnV3) GetVersion() TransactionVersion {
	return tx.Version
}

// L1Handler transaction
func (tx L1HandlerTxn) GetType() TransactionType {
	return tx.Type
}
func (tx L1HandlerTxn) GetVersion() TransactionVersion {
	return tx.Version
}

// InvokeTxnType is an interface that represents a Starknet invoke transaction.
// It is used to provide a common interface for all invoke transaction types.
// The 'Calldata' field is present in all invoke transaction types.
type InvokeTxnType interface {
	GetCalldata() []*felt.Felt
}

func (tx InvokeTxnV0) GetCalldata() []*felt.Felt {
	return tx.Calldata
}
func (tx InvokeTxnV1) GetCalldata() []*felt.Felt {
	return tx.Calldata
}
func (tx InvokeTxnV3) GetCalldata() []*felt.Felt {
	return tx.Calldata
}
func (tx BroadcastInvokeTxnV3) GetCalldata() []*felt.Felt {
	return tx.Calldata
}

var _ InvokeTxnType = InvokeTxnV0{}
var _ InvokeTxnType = InvokeTxnV1{}
var _ InvokeTxnType = InvokeTxnV3{}
var _ InvokeTxnType = BroadcastInvokeTxnV3{}

// DeclareTxnType is an interface that represents a Starknet declare transaction.
// It is used to provide a common interface for all declare transaction types.
// The 'SenderAddress' field is present in all declare transaction types.
type DeclareTxnType interface {
	GetSenderAddress() *felt.Felt
}

func (tx DeclareTxnV0) GetSenderAddress() *felt.Felt {
	return tx.SenderAddress
}
func (tx DeclareTxnV1) GetSenderAddress() *felt.Felt {
	return tx.SenderAddress
}
func (tx DeclareTxnV2) GetSenderAddress() *felt.Felt {
	return tx.SenderAddress
}
func (tx DeclareTxnV3) GetSenderAddress() *felt.Felt {
	return tx.SenderAddress
}
func (tx BroadcastDeclareTxnV3) GetSenderAddress() *felt.Felt {
	return tx.SenderAddress
}

var _ DeclareTxnType = DeclareTxnV0{}
var _ DeclareTxnType = DeclareTxnV1{}
var _ DeclareTxnType = DeclareTxnV2{}
var _ DeclareTxnType = DeclareTxnV3{}
var _ DeclareTxnType = BroadcastDeclareTxnV3{}

// DeployAccountType is an interface that represents a Starknet deploy account transaction.
// It is used to provide a common interface for all deploy account transaction types.
// The 'ConstructorCalldata' field is present in all deploy account transaction types.
type DeployAccountType interface {
	GetConstructorCalldata() []*felt.Felt
}

func (tx DeployAccountTxn) GetConstructorCalldata() []*felt.Felt {
	return tx.ConstructorCalldata
}
func (tx DeployAccountTxnV3) GetConstructorCalldata() []*felt.Felt {
	return tx.ConstructorCalldata
}
func (tx BroadcastDeployAccountTxnV3) GetConstructorCalldata() []*felt.Felt {
	return tx.ConstructorCalldata
}

var _ DeployAccountType = DeployAccountTxn{}
var _ DeployAccountType = DeployAccountTxnV3{}
var _ DeployAccountType = BroadcastDeployAccountTxnV3{}
