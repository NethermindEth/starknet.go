package hash

import (
	"errors"
	"slices"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/curve"
)

// ClassHash calculates the hash of a contract class.
//
// Parameters:
//   - contract: A contract class object of type [*contracts.ContractClass].
//
// Returns:
//   - *felt.Felt: a pointer to a felt.Felt object that represents the calculated hash.
//   - error: an error object if there was an error during the hash calculation.
func ClassHash(contract *contracts.ContractClass) *felt.Felt {
	// https://docs.starknet.io/architecture-and-concepts/smart-contracts/class-hash/

	Version := "CONTRACT_CLASS_V" + contract.ContractClassVersion
	ContractClassVersionHash := new(felt.Felt).SetBytes([]byte(Version))
	ConstructorHash := hashEntryPointByType(contract.EntryPointsByType.Constructor)
	ExternalHash := hashEntryPointByType(contract.EntryPointsByType.External)
	L1HandleHash := hashEntryPointByType(contract.EntryPointsByType.L1Handler)
	SierraProgamHash := curve.PoseidonArray(contract.SierraProgram...)
	ABIHash := curve.StarknetKeccak([]byte(contract.ABI))

	// https://docs.starknet.io/architecture-and-concepts/smart-contracts/class-hash/#computing_the_cairo_1_class_hash
	//nolint:lll // The link would be unclickable if we break the line.
	return curve.PoseidonArray(
		ContractClassVersionHash,
		ExternalHash,
		L1HandleHash,
		ConstructorHash,
		ABIHash,
		SierraProgamHash,
	)
}

// CompiledClassHash calculates the hash of a compiled class in the Casm format
// using the Poseidon hash function.
// This function will be deprecated in Starknet v0.14.1 onwards.
//
// Parameters:
//   - casmClass: A `contracts.CasmClass` object
//
// Returns:
//   - *felt.Felt: a pointer to a felt.Felt object that represents the calculated hash.
func CompiledClassHash(casmClass *contracts.CasmClass) (*felt.Felt, error) {
	return compiledClassHash(casmClass, curve.PoseidonArray)
}

// CompiledClassHashV2 calculates the hash of a compiled class in the Casm format
// using the Blake2s hash function. This is correct hash function to calculate the
// compiled class hash for Starknet v0.14.1 onwards.
//
// Parameters:
//   - casmClass: A `contracts.CasmClass` object
//
// Returns:
//   - *felt.Felt: a pointer to a felt.Felt object that represents the calculated hash.
func CompiledClassHashV2(casmClass *contracts.CasmClass) (*felt.Felt, error) {
	return compiledClassHash(casmClass, curve.Blake2sArray)
}

// compiledClassHash calculates the hash of a compiled class in the Casm format
// using the provided hash function.
//
// Parameters:
//   - casmClass: A `contracts.CasmClass` object
//
// Returns:
//   - *felt.Felt: a pointer to a felt.Felt object that represents the calculated hash.
func compiledClassHash(
	casmClass *contracts.CasmClass,
	hashFunc func(...*felt.Felt) *felt.Felt,
) (*felt.Felt, error) {
	ContractClassVersionHash := new(felt.Felt).SetBytes([]byte("COMPILED_CLASS_V1"))
	ExternalHash := hashCasmEntryPoints(casmClass.EntryPointsByType.External, hashFunc)
	L1HandleHash := hashCasmEntryPoints(casmClass.EntryPointsByType.L1Handler, hashFunc)
	ConstructorHash := hashCasmEntryPoints(casmClass.EntryPointsByType.Constructor, hashFunc)

	var ByteCodeHasH *felt.Felt
	var err error

	if casmClass.BytecodeSegmentLengths != nil {
		ByteCodeHasH, err = hashCasmClassByteCode(
			casmClass.ByteCode,
			*casmClass.BytecodeSegmentLengths,
			hashFunc,
		)
		if err != nil {
			return nil, err
		}
	} else {
		ByteCodeHasH = hashFunc(casmClass.ByteCode...)
	}

	//nolint:lll // The link would be unclickable if we break the line.
	// https://github.com/software-mansion/starknet.py/blob/39af414389984efbc6edc48b0fe1f914ea5b9a77/starknet_py/hash/casm_class_hash.py#L18
	return hashFunc(
		ContractClassVersionHash,
		ExternalHash,
		L1HandleHash,
		ConstructorHash,
		ByteCodeHasH,
	), nil
}

// hashEntryPointByType calculates the hash of an entry point by type.
func hashEntryPointByType(entryPoint []contracts.SierraEntryPoint) *felt.Felt {
	flattened := make([]*felt.Felt, 0, len(entryPoint))
	for _, elt := range entryPoint {
		flattened = append(
			flattened,
			elt.Selector,
			new(felt.Felt).SetUint64(uint64(elt.FunctionIdx)),
		)
	}

	return curve.PoseidonArray(flattened...)
}

// getByteCodeSegmentHasher calculates hasher function for byte code array from
// casm file. This code is adaptation of:
// https://github.com/starkware-libs/cairo-lang/blob/efa9648f57568aad8f8a13fbf027d2de7c63c2c0/src/starkware/starknet/core/os/contract_class/compiled_class_hash.py
//
// Parameters:
//   - bytecode: Array of compiled bytecode values from casm file
//   - bytecodeSegmentLengths: Nested datastructure of bytecode_segment_lengths
//     values from casm file
//   - visitedPcs: array pointer for tracking which bytecode bits were already
//     processed, needed for recursive processing
//   - bytecodeOffset: pointer at current offset in bytecode array, needed for
//     recursive processing organisation
//
// Returns:
//   - hasherFunc: closure that calculates hash for given bytecode array, or nil
//     in case of error
//   - uint64: size of the current processed bytecode array, or nil in case of error
//   - error: error if any happened or nil if everything fine
//
//nolint:lll // The link would be unclickable if we break the line.
func getByteCodeSegmentHasher(
	bytecode []*felt.Felt,
	bytecodeSegmentLengths contracts.NestedUints,
	visitedPcs *[]uint64,
	bytecodeOffset uint64,
	hashFunc func(...*felt.Felt) *felt.Felt,
) (hash *felt.Felt, size uint64, err error) {
	if !bytecodeSegmentLengths.IsArray {
		segmentValue := *bytecodeSegmentLengths.Value
		segmentEnd := bytecodeOffset + segmentValue

		for {
			visitedPcsData := *visitedPcs

			if len(visitedPcsData) == 0 {
				break
			}

			lastVisitedPcs := visitedPcsData[len(visitedPcsData)-1]

			if (bytecodeOffset > lastVisitedPcs) || (lastVisitedPcs >= segmentEnd) {
				break
			}

			*visitedPcs = visitedPcsData[:len(visitedPcsData)-1]
		}

		bytecodePart := bytecode[bytecodeOffset:segmentEnd]

		return hashFunc(bytecodePart...), segmentValue, nil
	}

	type bytecodeSegment struct {
		Value *felt.Felt
		Size  uint64
	}

	segments := []bytecodeSegment{}
	totalLen := uint64(0)

	for _, item := range bytecodeSegmentLengths.Values {
		visitedPcsData := *visitedPcs
		var visitedPcBefore *uint64

		if len(visitedPcsData) > 0 {
			visitedPcBefore = &visitedPcsData[len(visitedPcsData)-1]
		}

		segmentHash, segmentLen, err := getByteCodeSegmentHasher(
			bytecode,
			item,
			visitedPcs,
			bytecodeOffset,
			hashFunc,
		)
		if err != nil {
			return nil, 0, err
		}

		var visitedPcAfter *uint64
		if len(visitedPcsData) > 0 {
			visitedPcAfter = &visitedPcsData[len(visitedPcsData)-1]
		}

		isUsed := visitedPcAfter != visitedPcBefore

		if isUsed && *visitedPcBefore != bytecodeOffset {
			return nil, 0, errors.New(
				"invalid segment structure: PC {visited_pc_before} was visited, " +
					"but the beginning of the segment ({bytecode_offset}) was not",
			)
		}

		segments = append(segments, bytecodeSegment{
			Value: segmentHash,
			Size:  segmentLen,
		})
		bytecodeOffset += segmentLen
		totalLen += segmentLen
	}

	components := make([]*felt.Felt, len(segments)*2)

	for i, val := range segments {
		components[i*2] = felt.NewFromUint64[felt.Felt](val.Size)
		components[i*2+1] = val.Value
	}

	return new(felt.Felt).Add(
		felt.NewFromUint64[felt.Felt](1),
		hashFunc(components...),
	), totalLen, nil
}

// getByteCodeSegmentHasher calculates hash for byte code array from casm file
//
// Parameters:
//   - bytecode: Array of compiled bytecode values from casm file
//   - bytecodeSegmentLengths: Nested datastructure of bytecode_segment_lengths
//     values from casm file
//
// Returns:
//   - *felt.Felt: Hash value
//   - error: Error message
func hashCasmClassByteCode(
	bytecode []*felt.Felt,
	bytecodeSegmentLengths contracts.NestedUints,
	hashFunc func(...*felt.Felt) *felt.Felt,
) (*felt.Felt, error) {
	visited := make([]uint64, len(bytecode))

	for i := range bytecode {
		visited[i] = uint64(i) //nolint:gosec // Never overflows
	}

	slices.Reverse(visited)

	hash, _, err := getByteCodeSegmentHasher(
		bytecode,
		bytecodeSegmentLengths,
		&visited,
		uint64(0),
		hashFunc,
	)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

// hashCasmEntryPoints calculates the hash of a CasmClassEntryPoint array
// using the provided hash function.
func hashCasmEntryPoints(
	entryPoint []contracts.CasmEntryPoint,
	hashFunc func(...*felt.Felt) *felt.Felt,
) *felt.Felt {
	flattened := make([]*felt.Felt, 0, len(entryPoint))
	for _, elt := range entryPoint {
		builtInFlat := []*felt.Felt{}
		for _, builtIn := range elt.Builtins {
			builtInFlat = append(builtInFlat, new(felt.Felt).SetBytes([]byte(builtIn)))
		}
		builtInHash := hashFunc(builtInFlat...)
		flattened = append(
			flattened,
			elt.Selector,
			new(felt.Felt).SetUint64(uint64(elt.Offset)),
			builtInHash,
		)
	}

	return hashFunc(flattened...)

}
