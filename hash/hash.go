package hash

import (
	"errors"
	"slices"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

// CalculateTransactionHashCommon calculates the transaction hash common to be used in the StarkNet network - a unique identifier of the transaction.
// [specification]: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/starknet/core/os/transaction_hash/transaction_hash.py#L27C5-L27C38
//
// Parameters:
// - txHashPrefix: The prefix of the transaction hash
// - version: The version of the transaction
// - contractAddress: The address of the contract
// - entryPointSelector: The selector of the entry point
// - calldata: The data of the transaction
// - maxFee: The maximum fee for the transaction
// - chainId: The ID of the blockchain
// - additionalData: Additional data to be included in the hash
// Returns:
// - *felt.Felt: the calculated transaction hash
func CalculateTransactionHashCommon(
	txHashPrefix *felt.Felt,
	version *felt.Felt,
	contractAddress *felt.Felt,
	entryPointSelector *felt.Felt,
	calldata *felt.Felt,
	maxFee *felt.Felt,
	chainId *felt.Felt,
	additionalData []*felt.Felt) *felt.Felt {

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
	return curve.PedersenArray(dataToHash...)
}

// ClassHash calculates the hash of a contract class.
//
// It takes a contract class as input and calculates the hash by combining various elements of the class.
// The hash is calculated using the PoseidonArray function from the Curve package.
// The elements used in the hash calculation include the contract class version, constructor entry point, external entry point, L1 handler entry point, ABI, and Sierra program.
// The ABI is converted to bytes and then hashed using the StarknetKeccak function from the Curve package.
// Finally, the ContractClassVersionHash, ExternalHash, L1HandleHash, ConstructorHash, ABIHash, and SierraProgamHash are combined using the PoseidonArray function from the Curve package.
//
// Parameters:
// - contract: A contract class object of type rpc.ContractClass.
// Returns:
// - *felt.Felt: a pointer to a felt.Felt object that represents the calculated hash.
// - error: an error object if there was an error during the hash calculation.
func ClassHash(contract rpc.ContractClass) *felt.Felt {
	// https://docs.starknet.io/architecture-and-concepts/smart-contracts/class-hash/

	Version := "CONTRACT_CLASS_V" + contract.ContractClassVersion
	ContractClassVersionHash := new(felt.Felt).SetBytes([]byte(Version))
	ConstructorHash := hashEntryPointByType(contract.EntryPointsByType.Constructor)
	ExternalHash := hashEntryPointByType(contract.EntryPointsByType.External)
	L1HandleHash := hashEntryPointByType(contract.EntryPointsByType.L1Handler)
	SierraProgamHash := curve.PoseidonArray(contract.SierraProgram...)
	ABIHash := curve.StarknetKeccak([]byte(contract.ABI))

	// https://docs.starknet.io/architecture-and-concepts/smart-contracts/class-hash/#computing_the_cairo_1_class_hash
	return curve.PoseidonArray(ContractClassVersionHash, ExternalHash, L1HandleHash, ConstructorHash, ABIHash, SierraProgamHash)
}

// hashEntryPointByType calculates the hash of an entry point by type.
//
// Parameters:
// - entryPoint: A slice of rpc.SierraEntryPoint objects
// Returns:
// - *felt.Felt: a pointer to a felt.Felt object that represents the calculated hash.
func hashEntryPointByType(entryPoint []rpc.SierraEntryPoint) *felt.Felt {
	flattened := make([]*felt.Felt, 0, len(entryPoint))
	for _, elt := range entryPoint {
		flattened = append(flattened, elt.Selector, new(felt.Felt).SetUint64(uint64(elt.FunctionIdx)))
	}
	return curve.PoseidonArray(flattened...)
}

type hasherFunc = func() *felt.Felt

type bytecodeSegment struct {
	Hash hasherFunc
	Size uint64
}

// getByteCodeSegmentHasher calculates hasher function for byte code array from casm file
// this code is adaptation of:
// https://github.com/starkware-libs/cairo-lang/blob/v0.13.1/src/starkware/starknet/core/os/contract_class/compiled_class_hash.py
//
// Parameters:
// - bytecode: Array of compiled bytecode values from casm file
// - bytecodeSegmentLengths: Nested datastructure of bytecode_segment_lengths values from casm file
// - visitedPcs: array pointer for tracking which bytecode bits were already processed, needed for recursive processing
// - bytecodeOffset: pointer at current offset in bytecode array, needed for recursive processing organisation
// Returns:
// - hasherFunc: closure that calculates hash for given bytecode array, or nil in case of error
// - uint64: size of the current processed bytecode array, or nil in case of error
// - error: error if any happened or nil if everything fine
func getByteCodeSegmentHasher(
	bytecode []*felt.Felt,
	bytecodeSegmentLengths contracts.NestedUInts,
	visitedPcs *[]uint64,
	bytecodeOffset uint64,
) (hasherFunc, uint64, error) {
	if !bytecodeSegmentLengths.IsArray {
		segmentValue := *bytecodeSegmentLengths.Value
		segmentEnd := bytecodeOffset + segmentValue

		for {
			visitedPcsData := *visitedPcs

			if len(visitedPcsData) <= 0 {
				break
			}

			lastVisitedPcs := visitedPcsData[len(visitedPcsData)-1]

			if (bytecodeOffset > lastVisitedPcs) || (lastVisitedPcs >= segmentEnd) {
				break
			}

			*visitedPcs = visitedPcsData[:len(visitedPcsData)-1]
		}

		bytecodePart := bytecode[bytecodeOffset:segmentEnd]

		return func() *felt.Felt {
			return curve.PoseidonArray(bytecodePart...)
		}, segmentValue, nil
	}

	segments := []bytecodeSegment{}
	totalLen := uint64(0)

	for _, item := range bytecodeSegmentLengths.Values {
		visitedPcsData := *visitedPcs
		var visitedPcBefore *uint64 = nil

		if len(visitedPcsData) > 0 {
			visitedPcBefore = &visitedPcsData[len(visitedPcsData)-1]
		}

		segmentHash, segmentLen, err := getByteCodeSegmentHasher(
			bytecode,
			item,
			visitedPcs,
			bytecodeOffset,
		)

		if err != nil {
			return nil, 0, err
		}

		var visitedPcAfter *uint64 = nil
		if len(visitedPcsData) > 0 {
			visitedPcAfter = &visitedPcsData[len(visitedPcsData)-1]
		}

		is_used := visitedPcAfter != visitedPcBefore

		if is_used && *visitedPcBefore != bytecodeOffset {
			return nil, 0, errors.New(
				"invalid segment structure: PC {visited_pc_before} was visited, " +
					"but the beginning of the segment ({bytecode_offset}) was not",
			)
		}

		segments = append(segments, bytecodeSegment{
			Hash: segmentHash,
			Size: segmentLen,
		})
		bytecodeOffset += segmentLen
		totalLen += segmentLen
	}

	return func() *felt.Felt {
		components := make([]*felt.Felt, len(segments)*2)

		for i, val := range segments {
			components[i*2] = utils.Uint64ToFelt(val.Size)
			components[i*2+1] = val.Hash()
		}

		return felt.Zero.Add(
			utils.Uint64ToFelt(1),
			curve.PoseidonArray(components...),
		)
	}, totalLen, nil
}

// getByteCodeSegmentHasher calculates hash for byte code array from casm file
//
// Parameters:
// - bytecode: Array of compiled bytecode values from casm file
// - bytecodeSegmentLengths: Nested datastructure of bytecode_segment_lengths values from casm file
// Returns:
// - *felt.Felt: Hash value
// - error: Error message
func hashCasmClassByteCode(
	bytecode []*felt.Felt,
	bytecodeSegmentLengths contracts.NestedUInts,
) (*felt.Felt, error) {
	visited := []uint64{}

	for i := range len(bytecode) {
		visited = append(visited, uint64(i))
	}

	slices.Reverse(visited)

	hasher, _, err := getByteCodeSegmentHasher(
		bytecode, bytecodeSegmentLengths, &visited, uint64(0),
	)

	if err != nil {
		return nil, err
	}

	return hasher(), nil
}

// CompiledClassHash calculates the hash of a compiled class in the Casm format.
//
// Parameters:
// - casmClass: A `contracts.CasmClass` object
// Returns:
// - *felt.Felt: a pointer to a felt.Felt object that represents the calculated hash.
func CompiledClassHash(casmClass contracts.CasmClass) *felt.Felt {
	ContractClassVersionHash := new(felt.Felt).SetBytes([]byte("COMPILED_CLASS_V1"))
	ExternalHash := hashCasmClassEntryPointByType(casmClass.EntryPointByType.External)
	L1HandleHash := hashCasmClassEntryPointByType(casmClass.EntryPointByType.L1Handler)
	ConstructorHash := hashCasmClassEntryPointByType(casmClass.EntryPointByType.Constructor)

	var ByteCodeHasH *felt.Felt

	if casmClass.BytecodeSegmentLengths != nil {
		ByteCodeHasH, _ = hashCasmClassByteCode(casmClass.ByteCode, *casmClass.BytecodeSegmentLengths)
	} else {
		ByteCodeHasH = curve.PoseidonArray(casmClass.ByteCode...)
	}

	// https://github.com/software-mansion/starknet.py/blob/development/starknet_py/hash/casm_class_hash.py#L10
	return curve.PoseidonArray(ContractClassVersionHash, ExternalHash, L1HandleHash, ConstructorHash, ByteCodeHasH)
}

// hashCasmClassEntryPointByType calculates the hash of a CasmClassEntryPoint array.
//
// Parameters:
// - entryPoint: An array of CasmClassEntryPoint objects
// Returns:
// - *felt.Felt: a pointer to a Felt type
func hashCasmClassEntryPointByType(entryPoint []contracts.CasmClassEntryPoint) *felt.Felt {
	flattened := make([]*felt.Felt, 0, len(entryPoint))
	for _, elt := range entryPoint {
		builtInFlat := []*felt.Felt{}
		for _, builtIn := range elt.Builtins {
			builtInFlat = append(builtInFlat, new(felt.Felt).SetBytes([]byte(builtIn)))
		}
		builtInHash := curve.PoseidonArray(builtInFlat...)
		flattened = append(flattened, elt.Selector, new(felt.Felt).SetUint64(uint64(elt.Offset)), builtInHash)
	}
	return curve.PoseidonArray(flattened...)
}
