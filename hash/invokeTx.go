package hash

import (
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/NethermindEth/starknet.go/rpc/rpcv9"
	"github.com/NethermindEth/starknet.go/types/constraints"
)

// invokeV0 represents a pointer to an invoke V0
// transaction from all supported RPC versions.
type invokeV0 interface {
	*rpcv9.InvokeTxnV0 | *rpcv10.InvokeTxnV0
}

// invokeV1 represents a pointer to an invoke V1
// transaction from all supported RPC versions.
type invokeV1 interface {
	*rpcv9.InvokeTxnV1 | *rpcv10.InvokeTxnV1
}

// invokeV3 represents a pointer to an invoke V3
// transaction from all supported RPC versions.
type invokeV3 interface {
	*rpcv9.InvokeTxnV3 | *rpcv10.InvokeTxnV3
}

// invokeTx represents a pointer to an invoke transaction
// from all supported RPC versions.
type invokeTx interface {
	invokeV0 | invokeV1 | invokeV3
}

// @new
// TransactionHashInvoke calculates the transaction hash for an invoke transaction.
//
// Parameters:
//   - txn: The invoke transaction to calculate the hash for
//   - chainID: The chain ID as a *felt.Felt
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//   - error: an error if any
func TransactionHashInvoke[T invokeTx](tx T, chainID *felt.Felt) (*felt.Felt, error) {
	switch typedTx := any(tx).(type) {
	// **** v9 ****
	case *rpcv9.InvokeTxnV0:
		return TransactionHashInvokeV0(typedTx, chainID)
	case *rpcv9.InvokeTxnV1:
		return TransactionHashInvokeV1(typedTx, chainID)
	case *rpcv9.InvokeTxnV3:
		return TransactionHashInvokeV3(typedTx, chainID)
	// **** v10 ****
	case *rpcv10.InvokeTxnV0:
		return TransactionHashInvokeV0(typedTx, chainID)
	case *rpcv10.InvokeTxnV1:
		return TransactionHashInvokeV1(typedTx, chainID)
	case *rpcv10.InvokeTxnV3:
		return TransactionHashInvokeV3(typedTx, chainID)
	default:
		// Should never happen due to the generic type constraint
		return nil, errTxTypeNotSupported
	}
}

// @changed this and all other functions are now generic
// TransactionHashInvokeV0 calculates the transaction hash for a invoke V0 transaction.
//
// Parameters:
//   - txn: The invoke V0 transaction to calculate the hash for
//   - chainID: The chain ID as a *felt.Felt
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//   - error: an error if any
func TransactionHashInvokeV0[
	TxType, TxVersion ~string,
](tx constraints.InvokeTxnV0Interface[TxType, TxVersion],
	chainID *felt.Felt,
) (*felt.Felt, error) {
	// https://docs.starknet.io/learn/cheatsheets/transactions-reference#invoke-v0

	return calculateDeprecatedTransactionHashCommon(
		prefixInvoke,
		string(tx.GetVersion()),
		tx.GetContractAddress(),
		tx.GetEntryPointSelector(),
		curve.PedersenArray(tx.GetCalldata()...),
		tx.GetMaxFee(),
		chainID,
		[]*felt.Felt{},
	)
}

// TransactionHashInvokeV1 calculates the transaction hash for a invoke V1 transaction.
//
// Parameters:
//   - txn: The invoke V1 transaction to calculate the hash for
//   - chainID: The chain ID as a *felt.Felt
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//   - error: an error if any
func TransactionHashInvokeV1[
	TxType, TxVersion ~string,
](tx constraints.InvokeTxnV1Interface[TxType, TxVersion],
	chainID *felt.Felt,
) (*felt.Felt, error) {
	// https://docs.starknet.io/learn/cheatsheets/transactions-reference#invoke-v1

	return calculateDeprecatedTransactionHashCommon(
		prefixInvoke,
		string(tx.GetVersion()),
		tx.GetSenderAddress(),
		&felt.Zero,
		curve.PedersenArray(tx.GetCalldata()...),
		tx.GetMaxFee(),
		chainID,
		[]*felt.Felt{tx.GetNonce()},
	)
}

// TransactionHashInvokeV3 calculates the transaction hash for a invoke V3 transaction.
//
// Parameters:
//   - txn: The invoke V3 transaction to calculate the hash for
//   - chainID: The chain ID as a *felt.Felt
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//   - error: an error if any
func TransactionHashInvokeV3[
	TxType, TxVersion ~string,
	u64 constraints.U64,
	u128 constraints.U128,
	RB constraints.ResourceBounds[u64, u128],
	RBM constraints.ResourceBoundsMapping[u64, u128, RB],
	DA constraints.DataAvailabilityMode,
](tx constraints.InvokeTxnV3Interface[TxType, TxVersion, u64, u128, RB, RBM, DA],
	chainID *felt.Felt,
) (*felt.Felt, error) {
	// https://docs.starknet.io/learn/cheatsheets/transactions-reference#invoke-v3

	return calculateV3TransactionHash(
		prefixInvoke,
		string(tx.GetVersion()),
		tx.GetSenderAddress(),
		tx.GetTip(),
		tx.GetResourceBounds(),
		tx.GetPayMasterData(),
		chainID,
		tx.GetNonce(),
		tx.GetFeeMode(),
		tx.GetNonceDataMode(),
		[]*felt.Felt{
			curve.PoseidonArray(tx.GetAccountDeploymentData()...),
			curve.PoseidonArray(tx.GetCalldata()...),
		},
	)
}
