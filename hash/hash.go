package hash

import (
	"github.com/NethermindEth/juno/core/felt"
	starknetgo "github.com/NethermindEth/starknet.go"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

// computeHashOnElementsFelt hashes the array of felts provided as input
func ComputeHashOnElementsFelt(feltArr []*felt.Felt) (*felt.Felt, error) {
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
func CalculateTransactionHashCommon(
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
	return ComputeHashOnElementsFelt(dataToHash)
}

func ClassHash(contract rpc.ContractClass) (*felt.Felt, error) {
	//https://docs.starknet.io/documentation/architecture_and_concepts/Smart_Contracts/class-hash/

	SierraProgamHash, err := ComputeHashOnElementsFelt(contract.SierraProgram)
	if err != nil {
		return nil, err
	}
	ContractClassVersionHash := new(felt.Felt).SetBytes([]byte(contract.ContractClassVersion))
	ABHIHash := new(felt.Felt).SetBytes([]byte(contract.ABI))
	ExternalHash, err := hashEntryPointByType(contract.EntryPointsByType.External)
	if err != nil {
		return nil, err
	}
	L1HandleHash, err := hashEntryPointByType(contract.EntryPointsByType.L1Handler)
	if err != nil {
		return nil, err
	}
	ConstructorHash, err := hashEntryPointByType(contract.EntryPointsByType.Constructor)
	if err != nil {
		return nil, err
	}

	// https://docs.starknet.io/documentation/architecture_and_concepts/Network_Architecture/transactions/#deploy_account_hash_calculation
	return ComputeHashOnElementsFelt(
		[]*felt.Felt{
			ContractClassVersionHash,
			ExternalHash,
			L1HandleHash,
			ConstructorHash,
			ABHIHash,
			SierraProgamHash},
	)
}

func hashEntryPointByType(qwe []rpc.SierraEntryPoint) (*felt.Felt, error) {
	flattened := []*felt.Felt{}
	for _, elt := range qwe {
		flattened = append(flattened, elt.Selector, new(felt.Felt).SetUint64(uint64(elt.FunctionIdx)))
	}
	return ComputeHashOnElementsFelt(flattened)
}

// func CompiledClassHash() (*felt.Felt, error) {
// 	// https://github.com/software-mansion/starknet.py/blob/development/starknet_py/hash/casm_class_hash.py#L10
// 	panic("CompiledClassHash needs implementing")
// 	ContractClassVersionHash := new(felt.Felt).SetBytes([]byte(contract.ContractClassVersion))
// 	ExternalHash, err := hashEntryPointByType(contract.EntryPointsByType.External)
// 	if err != nil {
// 		return nil, err
// 	}
// 	L1HandleHash, err := hashEntryPointByType(contract.EntryPointsByType.L1Handler)
// 	if err != nil {
// 		return nil, err
// 	}
// 	ConstructorHash, err := hashEntryPointByType(contract.EntryPointsByType.Constructor)
// 	if err != nil {
// 		return nil, err
// 	}
// 	ByteCodeHasH, err := hashEntryPointByType(contract.ByteCode)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return ComputeHashOnElementsFelt(
// 		[]*felt.Felt{
// 			ContractClassVersionHash,
// 			ExternalHash,
// 			L1HandleHash,
// 			ConstructorHash,
// 			ByteCodeHasH},
// 	)
// }
