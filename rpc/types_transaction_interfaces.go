package rpc

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

// Note: these allow all types to pass, but are to help users of starknet.go
// understand which types are allowed where.

type InvokeTxnType interface{}

var _ InvokeTxnType = InvokeTxnV0{}
var _ InvokeTxnType = InvokeTxnV1{}
var _ InvokeTxnType = InvokeTxnV3{}

type DeclareTxnType interface{}

var _ DeclareTxnType = DeclareTxnV0{}
var _ DeclareTxnType = DeclareTxnV1{}
var _ DeclareTxnType = DeclareTxnV2{}
var _ DeclareTxnType = DeclareTxnV3{}

type DeployAccountType interface{}

var _ DeployAccountType = DeployAccountTxn{}
var _ DeployAccountType = DeployAccountTxnV3{}
