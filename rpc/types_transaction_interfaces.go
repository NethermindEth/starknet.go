package rpc

type AddDeclareTxnInput interface{}

var _ AddDeclareTxnInput = DeclareTxnV1{}
var _ AddDeclareTxnInput = DeclareTxnV2{}

type Transaction interface {
	GetType() TransactionType
}

var _ Transaction = InvokeTxnV0{}
var _ Transaction = InvokeTxnV1{}
var _ Transaction = DeclareTxnV1{}
var _ Transaction = DeclareTxnV2{}
var _ Transaction = DeployTxn{}
var _ Transaction = DeployAccountTxn{}
var _ Transaction = L1HandlerTxn{}

// GetType returns the transaction type of the InvokeTxnV0 struct.
//
// Parameters:
//  none
// Returns:
// - TransactionType: the transaction type
func (tx InvokeTxnV0) GetType() TransactionType {
	return tx.Type
}

// GetType returns the type of the InvokeTxnV1 transaction.
//
// Parameters:
//  none
// Returns:
// - TransactionType: the transaction type
func (tx InvokeTxnV1) GetType() TransactionType {
	return tx.Type
}

// GetType returns the TransactionType of the DeclareTxnV0 object.
//
// Parameters:
//  none
// Returns:
// - TransactionType: the transaction type
func (tx DeclareTxnV0) GetType() TransactionType {
	return tx.Type
}

// GetType returns the transaction type of DeclareTxnV1.
//
// Parameters:
//  none
// Returns:
// - TransactionType: the transaction type
func (tx DeclareTxnV1) GetType() TransactionType {
	return tx.Type
}

// GetType returns the type of the transaction.
//
// Parameters:
//  none
// Returns:
// - TransactionType: the transaction type
func (tx DeclareTxnV2) GetType() TransactionType {
	return tx.Type
}

// GetType returns the type of the DeployTxn.
//
// Parameters:
//  none
// Returns:
// - TransactionType: the transaction type
func (tx DeployTxn) GetType() TransactionType {
	return tx.Type
}

// GetType returns the transaction type of the DeployAccountTxn.
//
// Parameters:
//  none
// Returns:
// - TransactionType: the transaction type
func (tx DeployAccountTxn) GetType() TransactionType {
	return tx.Type
}

// GetType returns the transaction type of the L1HandlerTxn.
//
// Parameters:
//  none
// Returns:
// - TransactionType: the transaction type
func (tx L1HandlerTxn) GetType() TransactionType {
	return tx.Type
}

// Note: these allow all types to pass, but are to help users of starknet.go
// understand which types are allowed where.

type InvokeTxnType interface{}

var _ InvokeTxnType = InvokeTxnV0{}
var _ InvokeTxnType = InvokeTxnV1{}

type DeclareTxnType interface{}

var _ DeclareTxnType = DeclareTxnV0{}
var _ DeclareTxnType = DeclareTxnV1{}
var _ DeclareTxnType = DeclareTxnV2{}

type DeployAccountType interface{}

var _ DeployAccountType = DeployAccountTxn{}
