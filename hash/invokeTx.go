package hash

import (
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/NethermindEth/starknet.go/rpc/rpcv9"
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
func TransactionHashInvokeV0[T invokeV0](tx T, chainID *felt.Felt) (*felt.Felt, error) {
	// https://docs.starknet.io/learn/cheatsheets/transactions-reference#invoke-v0
	if tx == nil {
		return nil, ErrTransactionNil
	}

	switch typedTx := any(tx).(type) {
	case *rpcv9.InvokeTxnV0:
		return calculateDeprecatedTransactionHashCommon(
			prefixInvoke,
			string(typedTx.Version),
			typedTx.ContractAddress,
			typedTx.EntryPointSelector,
			curve.PedersenArray(typedTx.Calldata...),
			typedTx.MaxFee,
			chainID,
			[]*felt.Felt{},
		)
	case *rpcv10.InvokeTxnV0:
		return calculateDeprecatedTransactionHashCommon(
			prefixInvoke,
			string(typedTx.Version),
			typedTx.ContractAddress,
			typedTx.EntryPointSelector,
			curve.PedersenArray(typedTx.Calldata...),
			typedTx.MaxFee,
			chainID,
			[]*felt.Felt{},
		)
	default:
		// Should never happen due to the generic type constraint
		return nil, errTxTypeNotSupported
	}
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
func TransactionHashInvokeV1[T invokeV1](tx T, chainID *felt.Felt) (*felt.Felt, error) {
	// https://docs.starknet.io/learn/cheatsheets/transactions-reference#invoke-v1
	if tx == nil {
		return nil, ErrTransactionNil
	}

	switch typedTx := any(tx).(type) {
	case *rpcv9.InvokeTxnV1:
		return calculateDeprecatedTransactionHashCommon(
			prefixInvoke,
			string(typedTx.Version),
			typedTx.SenderAddress,
			&felt.Zero,
			curve.PedersenArray(typedTx.Calldata...),
			typedTx.MaxFee,
			chainID,
			[]*felt.Felt{typedTx.Nonce},
		)
	case *rpcv10.InvokeTxnV1:
		return calculateDeprecatedTransactionHashCommon(
			prefixInvoke,
			string(typedTx.Version),
			typedTx.SenderAddress,
			&felt.Zero,
			curve.PedersenArray(typedTx.Calldata...),
			typedTx.MaxFee,
			chainID,
			[]*felt.Felt{typedTx.Nonce},
		)
	default:
		// Should never happen due to the generic type constraint
		return nil, errTxTypeNotSupported
	}
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
func TransactionHashInvokeV3[T invokeV3](tx T, chainID *felt.Felt) (*felt.Felt, error) {
	// https://docs.starknet.io/learn/cheatsheets/transactions-reference#invoke-v3
	if tx == nil {
		return nil, ErrTransactionNil
	}

	switch typedTx := any(tx).(type) {
	case *rpcv9.InvokeTxnV3:
		if isOrContainsNil(typedTx.AccountDeploymentData, typedTx.Calldata) {
			return nil, ErrNotAllParametersSet
		}
		return calculateV3TransactionHash(
			prefixInvoke,
			string(typedTx.Version),
			typedTx.SenderAddress,
			typedTx.Tip,
			typedTx.ResourceBounds,
			typedTx.PayMasterData,
			chainID,
			typedTx.Nonce,
			typedTx.FeeMode,
			typedTx.NonceDataMode,
			[]*felt.Felt{
				curve.PoseidonArray(typedTx.AccountDeploymentData...),
				curve.PoseidonArray(typedTx.Calldata...),
			},
		)
	case *rpcv10.InvokeTxnV3:
		if isOrContainsNil(typedTx.AccountDeploymentData, typedTx.Calldata) {
			return nil, ErrNotAllParametersSet
		}
		return calculateV3TransactionHash(
			prefixInvoke,
			string(typedTx.Version),
			typedTx.SenderAddress,
			typedTx.Tip,
			typedTx.ResourceBounds,
			typedTx.PayMasterData,
			chainID,
			typedTx.Nonce,
			typedTx.FeeMode,
			typedTx.NonceDataMode,
			[]*felt.Felt{
				curve.PoseidonArray(typedTx.AccountDeploymentData...),
				curve.PoseidonArray(typedTx.Calldata...),
			},
		)
	default:
		// Should never happen due to the generic type constraint
		return nil, errTxTypeNotSupported
	}
}
