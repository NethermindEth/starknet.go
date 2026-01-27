package hash

import (
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/types/constraints"
)

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

	return CalculateDeprecatedTransactionHashCommon(
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

	return CalculateDeprecatedTransactionHashCommon(
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
			curve.PoseidonArray(tx.GetCalldata()...),
		},
	)
}
