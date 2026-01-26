package rpcv10

import (
	"github.com/NethermindEth/juno/core/felt"
)

// Transaction is an interface that represents a Starknet transaction.
// It is used to provide a common interface for all transaction types.
// The 'Type' and 'Version' fields are present in all transaction types.
type Transaction interface {
	GetType() TransactionType
	GetVersion() TransactionVersion
}

var (
	_ Transaction = InvokeTxnV0{}
	_ Transaction = InvokeTxnV1{}
	_ Transaction = InvokeTxnV3{}
	_ Transaction = DeclareTxnV1{}
	_ Transaction = DeclareTxnV2{}
	_ Transaction = DeclareTxnV3{}
	_ Transaction = DeployTxn{}
	_ Transaction = DeployAccountTxnV1{}
	_ Transaction = DeployAccountTxnV3{}
	_ Transaction = L1HandlerTxn{}
)

// @changed this and all other interfaces now include the Transaction interface
// InvokeTxnType is an interface that represents a Starknet invoke transaction.
// It is used to provide a common interface for all invoke transaction types.
// The 'Calldata' field is present in all invoke transaction types.
type InvokeTxnType interface {
	Transaction
	GetCalldata() []*felt.Felt
}

var (
	_ InvokeTxnType = InvokeTxnV0{}
	_ InvokeTxnType = InvokeTxnV1{}
	_ InvokeTxnType = InvokeTxnV3{}
	_ InvokeTxnType = BroadcastInvokeTxnV3{}
)

// DeclareTxnType is an interface that represents a Starknet declare transaction.
// It is used to provide a common interface for all declare transaction types.
// The 'SenderAddress' field is present in all declare transaction types.
type DeclareTxnType interface {
	Transaction
	GetSenderAddress() *felt.Felt
}

var (
	_ DeclareTxnType = DeclareTxnV0{}
	_ DeclareTxnType = DeclareTxnV1{}
	_ DeclareTxnType = DeclareTxnV2{}
	_ DeclareTxnType = DeclareTxnV3{}
	_ DeclareTxnType = BroadcastDeclareTxnV3{}
)

// DeployAccountType is an interface that represents a Starknet deploy account transaction.
// It is used to provide a common interface for all deploy account transaction types.
// The 'ConstructorCalldata' field is present in all deploy account transaction types.
type DeployAccountType interface {
	Transaction
	GetConstructorCalldata() []*felt.Felt
}

var (
	_ DeployAccountType = DeployAccountTxnV1{}
	_ DeployAccountType = DeployAccountTxnV3{}
	_ DeployAccountType = BroadcastDeployAccountTxnV3{}
)
