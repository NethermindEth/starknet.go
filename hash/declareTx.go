package hash

import (
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/NethermindEth/starknet.go/rpc/rpcv9"
)

// declareV0 represents a pointer to a declare V0
// transaction from all supported RPC versions.
type declareV0 interface {
	*rpcv9.DeclareTxnV0 | *rpcv10.DeclareTxnV0
}

// declareV1 represents a pointer to a declare V1
// transaction from all supported RPC versions.
type declareV1 interface {
	*rpcv9.DeclareTxnV1 | *rpcv10.DeclareTxnV1
}

// declareV2 represents a pointer to a declare V2
// transaction from all supported RPC versions.
type declareV2 interface {
	*rpcv9.DeclareTxnV2 | *rpcv10.DeclareTxnV2
}

// declareV3 represents a pointer to a declare V3
// transaction from all supported RPC versions.
type declareV3 interface {
	*rpcv9.DeclareTxnV3 | *rpcv9.BroadcastDeclareTxnV3 |
		*rpcv10.DeclareTxnV3 | *rpcv10.BroadcastDeclareTxnV3
}

// declareTx represents a pointer to a declare transaction
// from all supported RPC versions.
type declareTx interface {
	declareV0 | declareV1 | declareV2 | declareV3
}

// @new
// TransactionHashDeclare calculates the transaction hash for a declare transaction.
//
// Parameters:
//   - txn: The declare transaction to calculate the hash for
//   - chainID: The chain ID as a *felt.Felt
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//   - error: an error if any
func TransactionHashDeclare[T declareTx](tx T, chainID *felt.Felt) (*felt.Felt, error) {
	switch typedTx := any(tx).(type) {
	// **** v9 ****
	case *rpcv9.DeclareTxnV0:
		return TransactionHashDeclareV0(typedTx, chainID)
	case *rpcv9.DeclareTxnV1:
		return TransactionHashDeclareV1(typedTx, chainID)
	case *rpcv9.DeclareTxnV2:
		return TransactionHashDeclareV2(typedTx, chainID)
	case *rpcv9.DeclareTxnV3:
		return TransactionHashDeclareV3(typedTx, chainID)
	case *rpcv9.BroadcastDeclareTxnV3:
		return TransactionHashDeclareV3(typedTx, chainID)
	// **** v10 ****
	case *rpcv10.DeclareTxnV0:
		return TransactionHashDeclareV0(typedTx, chainID)
	case *rpcv10.DeclareTxnV1:
		return TransactionHashDeclareV1(typedTx, chainID)
	case *rpcv10.DeclareTxnV2:
		return TransactionHashDeclareV2(typedTx, chainID)
	case *rpcv10.DeclareTxnV3:
		return TransactionHashDeclareV3(typedTx, chainID)
	case *rpcv10.BroadcastDeclareTxnV3:
		return TransactionHashDeclareV3(typedTx, chainID)
	default:
		// Should never happen due to the generic type constraint
		return nil, errTxTypeNotSupported
	}
}

// @new
// TransactionHashDeclareV0 calculates the transaction hash for a declare V0 transaction.
//
// Parameters:
//   - txn: The declare V0 transaction to calculate the hash for
//   - chainID: The chain ID as a *felt.Felt
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//   - error: an error if any
func TransactionHashDeclareV0[T declareV0](tx T, chainID *felt.Felt) (*felt.Felt, error) {
	// https://docs.starknet.io/learn/cheatsheets/transactions-reference#declare-v0
	if tx == nil {
		return nil, ErrTransactionNil
	}

	switch typedTx := any(tx).(type) {
	case *rpcv9.DeclareTxnV0:
		return CalculateDeprecatedTransactionHashCommon(
			prefixDeclare,
			string(typedTx.Version),
			typedTx.SenderAddress,
			&felt.Zero,
			curve.PedersenArray(),
			typedTx.MaxFee,
			chainID,
			[]*felt.Felt{typedTx.ClassHash},
		)
	case *rpcv10.DeclareTxnV0:
		return CalculateDeprecatedTransactionHashCommon(
			prefixDeclare,
			string(typedTx.Version),
			typedTx.SenderAddress,
			&felt.Zero,
			curve.PedersenArray(),
			typedTx.MaxFee,
			chainID,
			[]*felt.Felt{typedTx.ClassHash},
		)
	default:
		// Should never happen due to the generic type constraint
		return nil, errTxTypeNotSupported
	}
}

// TransactionHashDeclareV1 calculates the transaction hash for a declare V1 transaction.
//
// Parameters:
//   - txn: The declare V1 transaction to calculate the hash for
//   - chainID: The chain ID as a *felt.Felt
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//   - error: an error if any
func TransactionHashDeclareV1[T declareV1](tx T, chainID *felt.Felt) (*felt.Felt, error) {
	// https://docs.starknet.io/learn/cheatsheets/transactions-reference#declare-v1
	if tx == nil {
		return nil, ErrTransactionNil
	}

	switch typedTx := any(tx).(type) {
	case *rpcv9.DeclareTxnV1:
		return CalculateDeprecatedTransactionHashCommon(
			prefixDeclare,
			string(typedTx.Version),
			typedTx.SenderAddress,
			&felt.Zero,
			curve.PedersenArray(typedTx.ClassHash),
			typedTx.MaxFee,
			chainID,
			[]*felt.Felt{typedTx.Nonce},
		)
	case *rpcv10.DeclareTxnV1:
		return CalculateDeprecatedTransactionHashCommon(
			prefixDeclare,
			string(typedTx.Version),
			typedTx.SenderAddress,
			&felt.Zero,
			curve.PedersenArray(typedTx.ClassHash),
			typedTx.MaxFee,
			chainID,
			[]*felt.Felt{typedTx.Nonce},
		)
	default:
		// Should never happen due to the generic type constraint
		return nil, errTxTypeNotSupported
	}
}

// TransactionHashDeclareV2 calculates the transaction hash for a declare V2 transaction.
//
// Parameters:
//   - txn: The declare V2 transaction to calculate the hash for
//   - chainID: The chain ID as a *felt.Felt
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//   - error: an error if any
func TransactionHashDeclareV2[T declareV2](tx T, chainID *felt.Felt) (*felt.Felt, error) {
	// https://docs.starknet.io/learn/cheatsheets/transactions-reference#declare-v2
	if tx == nil {
		return nil, ErrTransactionNil
	}

	switch typedTx := any(tx).(type) {
	case *rpcv9.DeclareTxnV2:
		return CalculateDeprecatedTransactionHashCommon(
			prefixDeclare,
			string(typedTx.Version),
			typedTx.SenderAddress,
			&felt.Zero,
			curve.PedersenArray(typedTx.ClassHash),
			typedTx.MaxFee,
			chainID,
			[]*felt.Felt{typedTx.Nonce, typedTx.CompiledClassHash},
		)
	case *rpcv10.DeclareTxnV2:
		return CalculateDeprecatedTransactionHashCommon(
			prefixDeclare,
			string(typedTx.Version),
			typedTx.SenderAddress,
			&felt.Zero,
			curve.PedersenArray(typedTx.ClassHash),
			typedTx.MaxFee,
			chainID,
			[]*felt.Felt{typedTx.Nonce, typedTx.CompiledClassHash},
		)
	default:
		// Should never happen due to the generic type constraint
		return nil, errTxTypeNotSupported
	}
}

// TransactionHashDeclareV3 calculates the transaction hash for a declare V3 transaction.
//
// Parameters:
//   - txn: The declare V3 transaction to calculate the hash for
//   - chainID: The chain ID as a *felt.Felt
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//   - error: an error if any
func TransactionHashDeclareV3[T declareV3](tx T, chainID *felt.Felt) (*felt.Felt, error) {
	// https://docs.starknet.io/learn/cheatsheets/transactions-reference#declare-v3
	if tx == nil {
		return nil, ErrTransactionNil
	}

	switch typedTx := any(tx).(type) {
	// **** v9 ****
	case *rpcv9.DeclareTxnV3:
		if isOrContainsNil(typedTx.AccountDeploymentData) {
			return nil, ErrNotAllParametersSet
		}
		return CalculateV3TransactionHash(
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
				typedTx.ClassHash,
				typedTx.CompiledClassHash,
			},
		)
	case *rpcv9.BroadcastDeclareTxnV3:
		if isOrContainsNil(typedTx.AccountDeploymentData, typedTx.ContractClass) {
			return nil, ErrNotAllParametersSet
		}
		return CalculateV3TransactionHash(
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
				ClassHash(typedTx.ContractClass),
				typedTx.CompiledClassHash,
			},
		)
	// **** v10 ****
	case *rpcv10.DeclareTxnV3:
		if isOrContainsNil(typedTx.AccountDeploymentData) {
			return nil, ErrNotAllParametersSet
		}
		return CalculateV3TransactionHash(
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
				typedTx.ClassHash,
				typedTx.CompiledClassHash,
			},
		)
	case *rpcv10.BroadcastDeclareTxnV3:
		if isOrContainsNil(typedTx.AccountDeploymentData, typedTx.ContractClass) {
			return nil, ErrNotAllParametersSet
		}
		return CalculateV3TransactionHash(
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
				ClassHash(typedTx.ContractClass),
				typedTx.CompiledClassHash,
			},
		)
	default:
		// Should never happen due to the generic type constraint
		return nil, errTxTypeNotSupported
	}
}
