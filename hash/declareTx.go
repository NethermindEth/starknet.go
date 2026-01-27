package hash

import (
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/types/constraints"
)

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
func TransactionHashDeclareV0[
	TxType, TxVersion ~string,
](tx constraints.DeclareTxnV0Interface[TxType, TxVersion],
	chainID *felt.Felt,
) (*felt.Felt, error) {
	// https://docs.starknet.io/learn/cheatsheets/transactions-reference#declare-v0

	return CalculateDeprecatedTransactionHashCommon(
		prefixDeclare,
		string(tx.GetVersion()),
		tx.GetSenderAddress(),
		&felt.Zero,
		curve.PedersenArray(),
		tx.GetMaxFee(),
		chainID,
		[]*felt.Felt{tx.GetClassHash()},
	)
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
func TransactionHashDeclareV1[
	TxType, TxVersion ~string,
](tx constraints.DeclareTxnV1Interface[TxType, TxVersion],
	chainID *felt.Felt,
) (*felt.Felt, error) {
	// https://docs.starknet.io/learn/cheatsheets/transactions-reference#declare-v1

	return CalculateDeprecatedTransactionHashCommon(
		prefixDeclare,
		string(tx.GetVersion()),
		tx.GetSenderAddress(),
		&felt.Zero,
		curve.PedersenArray(tx.GetClassHash()),
		tx.GetMaxFee(),
		chainID,
		[]*felt.Felt{tx.GetNonce()},
	)
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
func TransactionHashDeclareV2[
	TxType, TxVersion ~string,
](tx constraints.DeclareTxnV2Interface[TxType, TxVersion],
	chainID *felt.Felt,
) (*felt.Felt, error) {
	// https://docs.starknet.io/learn/cheatsheets/transactions-reference#declare-v2

	return CalculateDeprecatedTransactionHashCommon(
		prefixDeclare,
		string(tx.GetVersion()),
		tx.GetSenderAddress(),
		&felt.Zero,
		curve.PedersenArray(tx.GetClassHash()),
		tx.GetMaxFee(),
		chainID,
		[]*felt.Felt{tx.GetNonce(), tx.GetCompiledClassHash()},
	)
}

// @changed generics + accepts a class hash as a parameter
// TransactionHashDeclareV3 calculates the transaction hash for a declare V3 transaction.
//
// Parameters:
//   - txn: The declare V3 transaction to calculate the hash for
//   - chainID: The chain ID as a *felt.Felt
//   - classHash: The class hash as a *felt.Felt
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//   - error: an error if any
func TransactionHashDeclareV3[
	TxType, TxVersion ~string,
	u64 constraints.U64,
	u128 constraints.U128,
	RB constraints.ResourceBounds[u64, u128],
	RBM constraints.ResourceBoundsMapping[u64, u128, RB],
	DA constraints.DataAvailabilityMode,
](tx constraints.DeclareTxnV3Interface[TxType, TxVersion, u64, u128, RB, RBM, DA],
	chainID, classHash *felt.Felt,
) (*felt.Felt, error) {
	// https://docs.starknet.io/learn/cheatsheets/transactions-reference#declare-v3

	if isOrContainsNil(tx.GetAccountDeploymentData(), classHash) {
		return nil, ErrNotAllParametersSet
	}
	return CalculateV3TransactionHash(
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
			classHash,
			tx.GetCompiledClassHash(),
		},
	)
}

// @removed TransactionHashBroadcastDeclareV3. The logic was included in the
// TransactionHashDeclareV3 function.
