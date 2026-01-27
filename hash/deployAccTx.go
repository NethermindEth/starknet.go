package hash

import (
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/types/constraints"
)

// TransactionHashDeployAccountV1 calculates the transaction hash for a deploy account V1 transaction.
//
// Parameters:
//   - txn: The deploy account V1 transaction to calculate the hash for
//   - chainID: The chain ID as a *felt.Felt
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//   - error: an error if any
func TransactionHashDeployAccountV1[
	TxType, TxVersion ~string,
](tx constraints.DeployAccountTxnV1Interface[TxType, TxVersion],
	contractAddress, chainID *felt.Felt,
) (*felt.Felt, error) {
	// https://docs.starknet.io/learn/cheatsheets/transactions-reference#deploy-account-v1

	calldata := []*felt.Felt{tx.GetClassHash(), tx.GetContractAddressSalt()}
	calldata = append(calldata, tx.GetConstructorCalldata()...)
	calldataHash := curve.PedersenArray(calldata...)

	return CalculateDeprecatedTransactionHashCommon(
		prefixDeployAccount,
		string(tx.GetVersion()),
		contractAddress,
		&felt.Zero,
		calldataHash,
		tx.GetMaxFee(),
		chainID,
		[]*felt.Felt{tx.GetNonce()},
	)
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
func TransactionHashDeployAccountV3[
	TxType, TxVersion ~string,
	u64 constraints.U64,
	u128 constraints.U128,
	RB constraints.ResourceBounds[u64, u128],
	RBM constraints.ResourceBoundsMapping[u64, u128, RB],
	DA constraints.DataAvailabilityMode,
](tx constraints.DeployAccountTxnV3Interface[TxType, TxVersion, u64, u128, RB, RBM, DA],
	contractAddress, chainID *felt.Felt,
) (*felt.Felt, error) {
	// https://docs.starknet.io/learn/cheatsheets/transactions-reference#deploy-account-v3

	if isOrContainsNil(tx.GetConstructorCalldata()) {
		return nil, ErrNotAllParametersSet
	}
	return CalculateV3TransactionHash(
		prefixDeployAccount,
		string(tx.GetVersion()),
		contractAddress,
		tx.GetTip(),
		tx.GetResourceBounds(),
		tx.GetPayMasterData(),
		chainID,
		tx.GetNonce(),
		tx.GetFeeMode(),
		tx.GetNonceDataMode(),
		[]*felt.Felt{
			curve.PoseidonArray(tx.GetConstructorCalldata()...),
			tx.GetClassHash(),
			tx.GetContractAddressSalt(),
		},
	)
}
