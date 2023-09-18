package account

import (
	"github.com/NethermindEth/juno/core/felt"
	starknetgo "github.com/NethermindEth/starknet.go"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

// computeHashOnElementsFelt hashes the array of felts provided as input
func computeHashOnElementsFelt(feltArr []*felt.Felt) (*felt.Felt, error) {
	bigIntArr, err := utils.FeltArrToBigIntArr(feltArr)
	if err != nil {
		return nil, err
	}
	hash, err := starknetgo.Curve.ComputeHashOnElements(*bigIntArr)
	if err != nil {
		return nil, err
	}
	return utils.BigIntToFelt(hash)
}

// calculateTransactionHashCommon [specification] calculates the transaction hash in the StarkNet network - a unique identifier of the transaction.
// [specification]: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/starknet/core/os/transaction_hash/transaction_hash.py#L27C5-L27C38
func calculateTransactionHashCommon(
	txHashPrefix *felt.Felt,
	version *felt.Felt,
	contractAddress *felt.Felt,
	entryPointSelector *felt.Felt,
	calldata *felt.Felt,
	maxFee *felt.Felt,
	chainId *felt.Felt,
	additionalData []*felt.Felt) (*felt.Felt, error) {

	dataToHash := []*felt.Felt{
		txHashPrefix,
		version,
		contractAddress,
		entryPointSelector,
		calldata,
		maxFee,
		chainId,
	}
	dataToHash = append(dataToHash, additionalData...)
	return computeHashOnElementsFelt(dataToHash)
}

// calculateDeployAccountTransactionHash computes the transaction hash for deployAccount transactions
func calculateDeployAccountTransactionHash(tx rpc.BroadcastedDeployAccountTransaction, contractAddress *felt.Felt, chainID string) (*felt.Felt, error) {
	Prefix_DEPLOY_ACCOUNT := new(felt.Felt).SetBytes([]byte("deploy_account"))
	chainIdFelt := new(felt.Felt).SetBytes([]byte(chainID))

	calldata := []*felt.Felt{tx.ClassHash, tx.ContractAddressSalt}
	calldata = append(calldata, tx.ConstructorCalldata...)
	calldataHash, err := computeHashOnElementsFelt(calldata)
	if err != nil {
		return nil, err
	}

	versionFelt, err := new(felt.Felt).SetString(string(tx.Version))
	if err != nil {
		return nil, err
	}

	return calculateTransactionHashCommon(
		Prefix_DEPLOY_ACCOUNT,
		versionFelt,
		contractAddress,
		&felt.Zero,
		calldataHash,
		tx.MaxFee,
		chainIdFelt,
		[]*felt.Felt{tx.Nonce},
	)
}
