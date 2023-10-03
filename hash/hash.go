package hash

import (
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	starknetgo "github.com/NethermindEth/starknet.go"
	newcontract "github.com/NethermindEth/starknet.go/newcontracts"
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
	// https://docs.starknet.io/documentation/architecture_and_concepts/Smart_Contracts/class-hash/
	// https://github.com/starkware-libs/cairo-lang/blob/7712b21fc3b1cb02321a58d0c0579f5370147a8b/src/starkware/starknet/core/os/contracts.cairo#L47

	Version := "CONTRACT_CLASS_V" + contract.ContractClassVersion
	ContractClassVersionHash := new(felt.Felt).SetBytes([]byte(Version))
	ConstructorHash := hashEntryPointByType(contract.EntryPointsByType.Constructor)
	ExternalHash := hashEntryPointByType(contract.EntryPointsByType.External)
	L1HandleHash := hashEntryPointByType(contract.EntryPointsByType.L1Handler)

	// The ABI Bytes seem to match, but the hash does not
	ABIHash, err := starknetgo.Curve.StarknetKeccak([]byte(contract.ABI))
	if err != nil {
		return nil, err
	}
	SierraProgamHash, err := ComputeHashOnElementsFelt(contract.SierraProgram)
	if err != nil {
		return nil, err
	}

	fmt.Println("ContractClassVersionHash", ContractClassVersionHash) // Correct
	fmt.Println("ExternalHash", ExternalHash)                         // Correct
	fmt.Println("L1HandleHash", L1HandleHash)                         // Correct
	fmt.Println("ConstructorHash", ConstructorHash)                   // Correct
	fmt.Println("newABIHash", ABIHash)                                // Correct
	fmt.Println("SierraProgamHash", SierraProgamHash)                 // Incorrect

	// https://docs.starknet.io/documentation/architecture_and_concepts/Network_Architecture/transactions/#deploy_account_hash_calculation
	return ComputeHashOnElementsFelt(
		[]*felt.Felt{
			ContractClassVersionHash,
			ExternalHash,
			L1HandleHash,
			ConstructorHash,
			ABIHash,
			SierraProgamHash},
	)
}

func hashEntryPointByType(entryPoint []rpc.SierraEntryPoint) *felt.Felt {
	flattened := []*felt.Felt{}
	for _, elt := range entryPoint {
		flattened = append(flattened, elt.Selector, new(felt.Felt).SetUint64(uint64(elt.FunctionIdx)))
	}
	return starknetgo.Curve.PoseidonArray(flattened...)
}

func CompiledClassHash(casmClass newcontract.CasmClass) (*felt.Felt, error) {
	ContractClassVersionHash := new(felt.Felt).SetBytes([]byte(casmClass.Version))
	ExternalHash, err := hashCasmClassEntryPointByType(casmClass.EntryPointByType.External)
	if err != nil {
		return nil, err
	}
	L1HandleHash, err := hashCasmClassEntryPointByType(casmClass.EntryPointByType.L1Handler)
	if err != nil {
		return nil, err
	}
	ConstructorHash, err := hashCasmClassEntryPointByType(casmClass.EntryPointByType.Constructor)
	if err != nil {
		return nil, err
	}
	ByteCodeHasH, err := ComputeHashOnElementsFelt(casmClass.ByteCode)
	if err != nil {
		return nil, err
	}
	// https://github.com/software-mansion/starknet.py/blob/development/starknet_py/hash/casm_class_hash.py#L10
	return ComputeHashOnElementsFelt(
		[]*felt.Felt{
			ContractClassVersionHash,
			ExternalHash,
			L1HandleHash,
			ConstructorHash,
			ByteCodeHasH},
	)
}
func hashCasmClassEntryPointByType(entryPoint []newcontract.CasmClassEntryPoint) (*felt.Felt, error) {
	flattened := []*felt.Felt{}
	for _, elt := range entryPoint {
		builtInFlat := []*felt.Felt{}
		for _, builtIn := range elt.Builtins {
			builtInFlat = append(builtInFlat, new(felt.Felt).SetBytes([]byte(builtIn)))
		}
		builtInHash, err := ComputeHashOnElementsFelt(builtInFlat)
		if err != nil {
			return nil, err
		}
		flattened = append(flattened, elt.Selector, new(felt.Felt).SetUint64(uint64(elt.Offset)), builtInHash)
	}
	return ComputeHashOnElementsFelt(flattened)
}
