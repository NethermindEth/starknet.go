package account

import (
	"github.com/NethermindEth/juno/core/felt"
	starknetgo "github.com/NethermindEth/starknet.go"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

// computeHashOnElementsFelt hashes the array of felts provided as input
func computeHashOnElementsFelt(feltArr []*felt.Felt) (*felt.Felt, error) {
	bigIntArr := utils.FeltArrToBigIntArr(feltArr)
	hash, err := starknetgo.Curve.ComputeHashOnElements(bigIntArr)
	if err != nil {
		return nil, err
	}
	return utils.BigIntToFelt(hash), nil
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

/*
Formats the multicall transactions in a format which can be signed and verified by the network and OpenZeppelin account contracts
*/
func FmtCalldata(fnCalls []rpc.FunctionCall) []*felt.Felt {
	callArray := []*felt.Felt{}
	callData := []*felt.Felt{new(felt.Felt).SetUint64(uint64(len(fnCalls)))}

	for _, tx := range fnCalls {
		callData = append(callData, tx.ContractAddress, tx.EntryPointSelector)

		if len(tx.Calldata) == 0 {
			callData = append(callData, &felt.Zero, &felt.Zero)
			continue
		}

		callData = append(callData, new(felt.Felt).SetUint64(uint64(len(callArray))), new(felt.Felt).SetUint64(uint64(len(tx.Calldata))+1))
		callArray = append(callArray, tx.Calldata...)
	}
	callData = append(callData, new(felt.Felt).SetUint64(uint64(len(callArray)+1)))
	callData = append(callData, callArray...)
	callData = append(callData, new(felt.Felt).SetUint64(0))
	return callData
}
