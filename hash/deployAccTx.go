package hash

import (
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/NethermindEth/starknet.go/rpc/rpcv9"
)

// deployAccountV1 represents a pointer to a deploy account V1
// transaction from all supported RPC versions.
type deployAccountV1 interface {
	*rpcv9.DeployAccountTxnV1 | *rpcv10.DeployAccountTxnV1
}

// deployAccountV3 represents a pointer to a deploy account V3
// transaction from all supported RPC versions.
type deployAccountV3 interface {
	*rpcv9.DeployAccountTxnV3 | *rpcv10.DeployAccountTxnV3
}

// deployAccountTx represents a pointer to a deploy account transaction
// from all supported RPC versions.
type deployAccountTx interface {
	deployAccountV1 | deployAccountV3
}

// @new
// TransactionHashDeployAccount calculates the transaction hash for a deploy account transaction.
//
// Parameters:
//   - txn: The deploy account transaction to calculate the hash for
//   - chainID: The chain ID as a *felt.Felt
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//   - error: an error if any
func TransactionHashDeployAccount[T deployAccountTx](tx T, contractAddress, chainID *felt.Felt) (*felt.Felt, error) {
	switch typedTx := any(tx).(type) {
	case *rpcv9.DeployAccountTxnV1:
		return TransactionHashDeployAccountV1(typedTx, contractAddress, chainID)
	case *rpcv9.DeployAccountTxnV3:
		return TransactionHashDeployAccountV3(typedTx, contractAddress, chainID)
	case *rpcv10.DeployAccountTxnV1:
		return TransactionHashDeployAccountV1(typedTx, contractAddress, chainID)
	case *rpcv10.DeployAccountTxnV3:
		return TransactionHashDeployAccountV3(typedTx, contractAddress, chainID)
	default:
		// Should never happen due to the generic type constraint
		return nil, errTxTypeNotSupported
	}
}

// TransactionHashDeployAccountV1 calculates the transaction hash for a deploy account V1 transaction.
//
// Parameters:
//   - txn: The deploy account V1 transaction to calculate the hash for
//   - chainID: The chain ID as a *felt.Felt
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//   - error: an error if any
func TransactionHashDeployAccountV1[T deployAccountV1](tx T, contractAddress, chainID *felt.Felt) (*felt.Felt, error) {
	// https://docs.starknet.io/learn/cheatsheets/transactions-reference#deploy-account-v1
	if tx == nil {
		return nil, ErrTransactionNil
	}

	switch typedTx := any(tx).(type) {
	case *rpcv9.DeployAccountTxnV1:
		calldata := []*felt.Felt{typedTx.ClassHash, typedTx.ContractAddressSalt}
		calldata = append(calldata, typedTx.ConstructorCalldata...)
		calldataHash := curve.PedersenArray(calldata...)

		return calculateDeprecatedTransactionHashCommon(
			prefixDeployAccount,
			string(typedTx.Version),
			contractAddress,
			&felt.Zero,
			calldataHash,
			typedTx.MaxFee,
			chainID,
			[]*felt.Felt{typedTx.Nonce},
		)
	case *rpcv10.DeployAccountTxnV1:
		calldata := []*felt.Felt{typedTx.ClassHash, typedTx.ContractAddressSalt}
		calldata = append(calldata, typedTx.ConstructorCalldata...)
		calldataHash := curve.PedersenArray(calldata...)

		return calculateDeprecatedTransactionHashCommon(
			prefixDeployAccount,
			string(typedTx.Version),
			contractAddress,
			&felt.Zero,
			calldataHash,
			typedTx.MaxFee,
			chainID,
			[]*felt.Felt{typedTx.Nonce},
		)
	default:
		// Should never happen due to the generic type constraint
		return nil, errTxTypeNotSupported
	}
}

// TransactionHashDeployAccountV3 calculates the transaction hash for a deploy account V3 transaction.
//
// Parameters:
//   - txn: The deploy account V3 transaction to calculate the hash for
//   - chainID: The chain ID as a *felt.Felt
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//   - error: an error if any
func TransactionHashDeployAccountV3[T deployAccountV3](tx T, contractAddress, chainID *felt.Felt) (*felt.Felt, error) {
	// https://docs.starknet.io/learn/cheatsheets/transactions-reference#deploy-account-v3
	if tx == nil {
		return nil, ErrTransactionNil
	}

	switch typedTx := any(tx).(type) {
	case *rpcv9.DeployAccountTxnV3:
		if isOrContainsNil(typedTx.ConstructorCalldata) {
			return nil, ErrNotAllParametersSet
		}
		return calculateV3TransactionHash(
			prefixDeployAccount,
			string(typedTx.Version),
			contractAddress,
			typedTx.Tip,
			typedTx.ResourceBounds,
			typedTx.PayMasterData,
			chainID,
			typedTx.Nonce,
			typedTx.FeeMode,
			typedTx.NonceDataMode,
			[]*felt.Felt{
				curve.PoseidonArray(typedTx.ConstructorCalldata...),
				typedTx.ClassHash,
				typedTx.ContractAddressSalt,
			},
		)
	case *rpcv10.DeployAccountTxnV3:
		if isOrContainsNil(typedTx.ConstructorCalldata) {
			return nil, ErrNotAllParametersSet
		}
		return calculateV3TransactionHash(
			prefixDeployAccount,
			string(typedTx.Version),
			contractAddress,
			typedTx.Tip,
			typedTx.ResourceBounds,
			typedTx.PayMasterData,
			chainID,
			typedTx.Nonce,
			typedTx.FeeMode,
			typedTx.NonceDataMode,
			[]*felt.Felt{
				curve.PoseidonArray(typedTx.ConstructorCalldata...),
				typedTx.ClassHash,
				typedTx.ContractAddressSalt,
			},
		)
	default:
		// Should never happen due to the generic type constraint
		return nil, errTxTypeNotSupported
	}
}

// @removed TransactionHashBroadcastDeclareV3. The logic was included in the
// TransactionHashDeclareV3 function.
