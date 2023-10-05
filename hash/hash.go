package hash

import (
	"github.com/NethermindEth/juno/core/felt"
	starknetgo "github.com/NethermindEth/starknet.go"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

// ComputeHashOnElementsFelt hashes the array of felts provided as input
func ComputeHashOnElementsFelt(feltArr []*felt.Felt) (*felt.Felt, error) {
	bigIntArr := utils.FeltArrToBigIntArr(feltArr)

	hash, err := starknetgo.Curve.ComputeHashOnElements(bigIntArr)
	if err != nil {
		return nil, err
	}
	return utils.BigIntToFelt(hash), nil
}

// CalculateTransactionHashCommon [specification] calculates the transaction hash in the StarkNet network - a unique identifier of the transaction.
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

	Version := "CONTRACT_CLASS_V" + contract.ContractClassVersion
	ContractClassVersionHash := new(felt.Felt).SetBytes([]byte(Version))
	ConstructorHash := hashEntryPointByType(contract.EntryPointsByType.Constructor)
	ExternalHash := hashEntryPointByType(contract.EntryPointsByType.External)
	L1HandleHash := hashEntryPointByType(contract.EntryPointsByType.L1Handler)
	SierraProgamHash := starknetgo.Curve.PoseidonArray(contract.SierraProgram...)
	ABIHash, err := starknetgo.Curve.StarknetKeccak([]byte(contract.ABI))
	if err != nil {
		return nil, err
	}

	// https://docs.starknet.io/documentation/architecture_and_concepts/Network_Architecture/transactions/#deploy_account_hash_calculation
	return starknetgo.Curve.PoseidonArray(ContractClassVersionHash, ExternalHash, L1HandleHash, ConstructorHash, ABIHash, SierraProgamHash), nil
}

func hashEntryPointByType(entryPoint []rpc.SierraEntryPoint) *felt.Felt {
	flattened := make([]*felt.Felt, 0, len(entryPoint))
	for _, elt := range entryPoint {
		flattened = append(flattened, elt.Selector, new(felt.Felt).SetUint64(uint64(elt.FunctionIdx)))
	}
	return starknetgo.Curve.PoseidonArray(flattened...)
}

func CompiledClassHash(casmClass contracts.CasmClass) *felt.Felt {
	ContractClassVersionHash := new(felt.Felt).SetBytes([]byte("COMPILED_CLASS_V1"))
	ExternalHash := hashCasmClassEntryPointByType(casmClass.EntryPointByType.External)
	L1HandleHash := hashCasmClassEntryPointByType(casmClass.EntryPointByType.L1Handler)
	ConstructorHash := hashCasmClassEntryPointByType(casmClass.EntryPointByType.Constructor)
	ByteCodeHasH := starknetgo.Curve.PoseidonArray(casmClass.ByteCode...)

	// https://github.com/software-mansion/starknet.py/blob/development/starknet_py/hash/casm_class_hash.py#L10
	return starknetgo.Curve.PoseidonArray(ContractClassVersionHash, ExternalHash, L1HandleHash, ConstructorHash, ByteCodeHasH)
}

func hashCasmClassEntryPointByType(entryPoint []contracts.CasmClassEntryPoint) *felt.Felt {
	flattened := make([]*felt.Felt, 0, len(entryPoint))
	for _, elt := range entryPoint {
		builtInFlat := []*felt.Felt{}
		for _, builtIn := range elt.Builtins {
			builtInFlat = append(builtInFlat, new(felt.Felt).SetBytes([]byte(builtIn)))
		}
		builtInHash := starknetgo.Curve.PoseidonArray(builtInFlat...)
		flattened = append(flattened, elt.Selector, new(felt.Felt).SetUint64(uint64(elt.Offset)), builtInHash)
	}
	return starknetgo.Curve.PoseidonArray(flattened...)
}
